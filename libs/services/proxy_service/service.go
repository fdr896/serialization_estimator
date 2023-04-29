package proxy_service

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"serialization_estimator/libs/estimator"
	"serialization_estimator/libs/services/protocol"
	"serialization_estimator/libs/support"
	"strconv"
	"strings"
	"syscall"
	"time"

	zlog "github.com/rs/zerolog/log"
)

type Service struct {
    port string
	connByMethod map[string]*net.UDPConn

    multicastConn *net.UDPConn
    multicastRespListener *net.UDPAddr
    multicastRespChan chan string
    multicastErrChain chan error
}

func New(port string) (*Service, error) {
    service := &Service{
        port: port,
    }

	err := service.establishConnections()
	if err != nil {
		return nil, err
	} else {
    	return service, nil
	}
}

func (s *Service) Start() error {
    zlog.Info().
        Str("port", s.port).
		Interface("available methods", s.connByMethod).
        Msg("Running proxy service")

    addr, err := net.ResolveUDPAddr("udp", net.JoinHostPort("proxy", s.port))
    if err != nil {
        return err
    }

    multicastAddr := support.GetMulticastGroupAddrFromEnv()
    if multicastAddr != nil {
        {
            addr, err := net.ResolveUDPAddr("udp", net.JoinHostPort("proxy", "2030"))
            if (err != nil) {
                return err
            }
            go s.serveMulticastResponses(addr)
        }

        conn, err := net.DialUDP("udp", nil, multicastAddr)
        if err != nil {
            return err
        }
        s.multicastConn = conn
        zlog.Info().
            Str("local addr", conn.LocalAddr().String()).
            Str("remote addr", conn.RemoteAddr().String()).
            Msg("Establised connection to multicast group")

        s.multicastRespChan = make(chan string)
        s.multicastErrChain = make(chan error)
    }

    conn, err := net.ListenUDP("udp", addr)
    if err != nil {
        return err
    }
    zlog.Info().
        Str("addr", addr.String()).
        Msg("Proxy service listens")

    for {
        buf := make([]byte, protocol.MAX_DATAGRAM_SIZE)
        n, sender, err := conn.ReadFrom(buf)
        if err != nil {
            zlog.Error().Err(err).Msg("Failed to handle request")
            continue
        }

        var req protocol.ProxyRequest
        if err := json.Unmarshal(buf[:n], &req); err != nil {
            zlog.Error().Err(err).Msg("Failed to parse request")
            continue
        }
        if !req.Valid() {
            zlog.Error().
                Str("method", req.Method).
                Msg("Incorrect request method or param")
            continue
        }

        zlog.Info().
            Int("bytes", n).
            Str("sender", sender.String()).
            Str("method", req.Method).
            Str("param", req.Param).
            Msg("Proxy accepted request")
        
        switch req.Method {
        case "get_result":
            resp, err := s.handleGetResult(&req);
            if err != nil {
                zlog.Error().Err(err).
                    Msg("Failed to handle get result")
                continue
            }
            _, err = conn.WriteTo([]byte(*resp), sender)
            if err != nil {
                zlog.Error().Err(err).Msg("Failed to response")
            }
        case "get_result_all":
            if err := s.handleGetResultAll(); err != nil {
                zlog.Error().Err(err).
                    Msg("Failed to handle get result all")
                continue
            }

            for i := 0; i < len(s.connByMethod); i++ {
                select {
                case resp := <-s.multicastRespChan:
                    _, err = conn.WriteTo([]byte(resp), sender)
                    if err != nil {
                        zlog.Error().Err(err).Msg("Failed to respond")
                    }
                case err := <-s.multicastErrChain:
                    _, err = conn.WriteTo([]byte(err.Error()), sender)
                    if err != nil {
                        zlog.Error().Err(err).Msg("Failed to respond")
                    }
                }
            }
        }
    }
}

func (s *Service) handleGetResult(req *protocol.ProxyRequest) (*protocol.ProxyResponse, error) {
    conn := s.connByMethod[req.Param]

    // send request to estimator
    estimatorReq := protocol.NewEstimatorGetResultRequest(nil)
    buf, err := json.Marshal(estimatorReq)
    if err != nil {
        return nil, err
    }
    _, err = conn.Write(buf)
    if err != nil {
        return nil, err
    }

    // receive estimator's response
    respBuf := make([]byte, protocol.MAX_DATAGRAM_SIZE)
    conn.SetReadDeadline(time.Now().Add(3 * time.Second))
    n, _, err := conn.ReadFromUDP(respBuf)
    if err != nil {
        return nil, err
    }

    var resp protocol.EstimatorResponse
    if err := json.Unmarshal(respBuf[:n], &resp); err != nil {
        return nil, err
    }

    return protocol.NewProxyResponse(&resp), nil

}

func (s *Service) handleGetResultAll() error {
    // multicast request to all estimators
    estimatorReq := protocol.NewEstimatorGetResultRequest(s.multicastRespListener)
    buf, err := json.Marshal(estimatorReq)
    if err != nil {
        return err
    }

    s.multicastConn.SetWriteDeadline(time.Now().Add(3 * time.Second))
    _, err = s.multicastConn.Write(buf)
    if err != nil {
        return err
    }

    return nil
}

func (s *Service) serveMulticastResponses(addr *net.UDPAddr) {
    s.multicastRespListener = addr

    conn, err := net.ListenUDP("udp", addr)
    if err != nil {
        panic(err)
    }
    file, _:= conn.File()
    syscall.SetsockoptInt(int(file.Fd()), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)

    zlog.Info().
        Str("addr", addr.String()).
        Msg("Listening multicast responses")

    for {
        buf := make([]byte, protocol.MAX_DATAGRAM_SIZE)
        n, sender, err := conn.ReadFromUDP(buf)
        if err != nil {
            zlog.Error().Err(err).Msg("Failed to handle response")
            s.multicastErrChain <- err
            continue
        }

        zlog.Info().
            Int("bytes", n).
            Str("sender", sender.String()).
            Msg("Accepted response on multicast")

        var resp protocol.EstimatorResponse
        if err := json.Unmarshal(buf[:n], &resp); err != nil {
            s.multicastErrChain <- err
            continue
        }

        s.multicastRespChan <- *protocol.NewProxyResponse(&resp)
    }
}

func (s *Service) establishConnections() error {
    s.connByMethod = make(map[string]*net.UDPConn)
	for _, method := range estimator.ESTIMATION_METHODS.Keys() {
		envName := makeEnvMethodVariable(method)
		port, set := os.LookupEnv(envName)
		if !set {
			return fmt.Errorf("%s not set", envName)
		}

		addr, err := resolveAddr(method, port)
		if err != nil {
			return err
		}

		conn, err := net.DialUDP("udp", nil, addr)
		if err != nil {
			return err
		}
		s.connByMethod[method] = conn
	}

	return nil
}

func makeEnvMethodVariable(method string) string {
	return strings.ToUpper(method) + "_PORT"
}

func resolveAddr(method, port string) (*net.UDPAddr, error) {
	if _, err := strconv.Atoi(port); err != nil {
		return nil, fmt.Errorf("incorrect port: %s", port)
	}
	return net.ResolveUDPAddr("udp", net.JoinHostPort(method, port))
}

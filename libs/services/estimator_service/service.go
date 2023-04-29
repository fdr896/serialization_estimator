package estimator_service

import (
	"encoding/json"
	"net"
	"sync"
	"syscall"
	"time"

	"serialization_estimator/libs/estimator"
	"serialization_estimator/libs/services/protocol"
	"serialization_estimator/libs/support"

	zlog "github.com/rs/zerolog/log"
)

type Service struct {
    port string
    multicastAddr *net.UDPAddr
    estimator estimator.Estimator

    mutex sync.Mutex
}

func New(port, method string) *Service {
    service := &Service{
        port: port,
        multicastAddr: support.GetMulticastGroupAddrFromEnv(),
        estimator: estimator.New(method),
    }

    return service
}

func (s *Service) Start() error {
    var wg sync.WaitGroup

    unicastAddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(s.estimator.Method(), s.port))
    if err != nil {
        return err
    }

    var unicastConn, multicastConn *net.UDPConn

    unicastConn, err = net.ListenUDP("udp", unicastAddr)
    if err != nil {
        return err
    }
    file, _:= unicastConn.File()
    syscall.SetsockoptInt(int(file.Fd()), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
    if s.multicastAddr != nil {
        multicastConn, err = net.ListenMulticastUDP("udp", nil, s.multicastAddr)
        if err != nil {
            return err
        }
        file, _ := multicastConn.File()
        syscall.SetsockoptInt(int(file.Fd()), syscall.IPPROTO_IP, syscall.IP_MULTICAST_LOOP, 1)
    }

    // Start handling unicast requests
    wg.Add(1)
    go func() {
        zlog.Info().
            Str("addr", unicastAddr.String()).
            Str("method", s.estimator.Method()).
            Msg("Estimator service serves unicast")

        s.serve(unicastConn)
        wg.Done()
    }()

    if (multicastConn != nil) {
        // Start listening multicast address
        wg.Add(1)
        go func() {
            zlog.Info().
                IPAddr("group ip", s.multicastAddr.IP).
                Int16("port", int16(s.multicastAddr.Port)).
                Str("method", s.estimator.Method()).
                Msg("Estimator service serves multicast")

            s.serve(multicastConn)
            wg.Done()
        }()
    }

    wg.Wait()

    return nil
}

func (s *Service) serve(conn *net.UDPConn) {
    for {
        buf := make([]byte, protocol.MAX_DATAGRAM_SIZE)
        n, sender, err := conn.ReadFromUDP(buf)
        if err != nil {
            zlog.Error().Err(err).Msg("Failed to handle request")
            continue
        }

        var req protocol.EstimatorRequest
        if err := json.Unmarshal(buf[:n], &req); err != nil {
            zlog.Error().Err(err).Msg("Failed to parse request")
            continue
        }
        if !req.Valid() {
            zlog.Error().Str("method", req.Method).Msg("Unacceptable method")
            continue
        }

        zlog.Info().
            Int("bytes", n).
            Str("sender", sender.String()).
            Str("method", req.Method).
            Msg("Estimator accepted request")
            
        // Serailization method estimation
        s.mutex.Lock()
        serializedSize, serializeDuration, deserializeDuration := estimator.Estimate(s.estimator)
        s.mutex.Unlock()

        resp := protocol.NewEstimatorResponse(s.estimator.Method(), serializedSize, serializeDuration, deserializeDuration)
        respBytes, err := json.Marshal(resp)
        if err != nil {
            zlog.Error().Err(err).Msg("Failed to marshal response")
            continue
        }

        if req.RespPort != nil {
            respAddr := net.JoinHostPort(req.RespPort.IP, req.RespPort.Port)
            zlog.Debug().Str("addr", respAddr).Msg("Responding to")

            sender, err = net.ResolveUDPAddr("udp", respAddr)
            if err != nil {
                zlog.Error().Err(err).Msg("Incorrect multicast addr")
                continue
            }
        }

        conn.SetWriteDeadline(time.Now().Add(3 * time.Second))
        _, err = conn.WriteToUDP(respBytes, sender)
        if err != nil {
            zlog.Error().Err(err).Msg("Failed to response")
        }
    }
}

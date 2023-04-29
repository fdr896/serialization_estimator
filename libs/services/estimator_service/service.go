package estimator_service

import (
	"encoding/json"
	"net"
	"sync"
	"syscall"

	"serialization_estimator/libs/estimator"
	"serialization_estimator/libs/services/protocol"
	"serialization_estimator/libs/support"

	zlog "github.com/rs/zerolog/log"
)

type Service struct {
    port string
    multicastAddr *net.UDPAddr
    estimator estimator.Estimator
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

    unicastAddr, err := net.ResolveUDPAddr("udp", ":" + s.port)
    if err != nil {
        return err
    }

    var unicastConn, multicastConn *net.UDPConn

    unicastConn, err = net.ListenUDP("udp", unicastAddr)
    if err != nil {
        return err
    }
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
            Str("port", s.port).
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
        serializedSize, serializeDuration, deserializeDuration := estimator.Estimate(s.estimator)

        resp := protocol.NewEstimatorResponse(s.estimator.Method(), serializedSize, serializeDuration, deserializeDuration)
        respBytes, err := json.Marshal(resp)
        if err != nil {
            zlog.Error().Err(err).Msg("Failed to marshal response")
            continue
        }

        if req.RespPort != nil {
            zlog.Debug().Str("Port", req.RespPort.Port).Msg("Responding to")
            sender, err = net.ResolveUDPAddr("udp", ":" + req.RespPort.Port)
            if err != nil {
                zlog.Error().Err(err).Msg("Incorrect multicast addr")
                continue
            }
        }

        _, err = conn.WriteToUDP(respBytes, sender)
        if err != nil {
            zlog.Error().Err(err).Msg("Failed to response")
        }
    }
}

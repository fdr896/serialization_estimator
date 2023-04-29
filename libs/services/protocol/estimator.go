package protocol

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

// REQUEST

type MulticastRespPort struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
}

type EstimatorRequest struct {
	Method   string             `json:"method"`
	RespPort *MulticastRespPort `json:"multicast_resp_port,omitempty"`
}

func NewEstimatorGetResultRequest(addr *net.UDPAddr) *EstimatorRequest {
	er := &EstimatorRequest{
		Method: "get_result",
	}
	if addr != nil {
		er.RespPort = &MulticastRespPort{
			IP: addr.IP.String(),
			Port: strconv.Itoa(addr.Port),
		}
	}

	return er
}

func (er *EstimatorRequest) Valid() bool {
	return er.Method == "get_result"
}

// RESPONSE

type EstimatorResponse struct {
	Method					  string `json:"method"`
	SerializedSize            uint32 `json:"serialized_size"`
	SerializationDurationNs   int64  `json:"serialization_duration_ns"`
	DeserializationDurationNs int64  `json:"deserialization_duration_ns"`
}

func NewEstimatorResponse(method string, size uint32, serializeDur time.Duration, deserializeDur time.Duration) *EstimatorResponse {
	return &EstimatorResponse{
		Method: method,
		SerializedSize: size,
		SerializationDurationNs: serializeDur.Nanoseconds(),
		DeserializationDurationNs: deserializeDur.Nanoseconds(),
	}
}

func (er *EstimatorResponse) format() string {
	return fmt.Sprintf("%s - %db - %dns - %dns\n",
			er.Method, er.SerializedSize,
			er.SerializationDurationNs, er.DeserializationDurationNs)

}

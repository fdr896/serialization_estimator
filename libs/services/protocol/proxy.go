package protocol

import (
	"fmt"
	"serialization_estimator/libs/estimator"
)

// REQUEST

type ProxyRequest struct {
	Method string `json:"method"`
	Param  string `json:"param,omitempty"`
}

func (pr *ProxyRequest) Valid() bool {
	switch pr.Method {
	case "get_result":
		return estimator.ESTIMATION_METHODS.Contains(pr.Param)
	case "get_result_all":
		return len(pr.Param) == 0
	default:
		return true
	}
}

// RESPONSE

type ProxyResponse = string

func NewProxyResponse(er *EstimatorResponse) *ProxyResponse {
	resp := ProxyResponse(er.format())
	return &resp
}

func NewEmptyProxyResponse() *ProxyResponse {
	resp := ProxyResponse("")
	return &resp
}

func AppendEstimatorResponse(pr *ProxyResponse, er *EstimatorResponse) *ProxyResponse {
	resp := ProxyResponse(fmt.Sprintf("%s\n%s", *pr, er.format()))
	return &resp
}

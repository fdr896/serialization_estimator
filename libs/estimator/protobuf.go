package estimator

import (
	"serialization_estimator/libs/estimator/impl"
	protobuf "serialization_estimator/libs/estimator/proto"
	"serialization_estimator/libs/estimator/templates"
	"time"
)

type protoEstimator struct {
	obj protobuf.Object
    data []byte
}

func newProtobufEstimator() *protoEstimator {
    return &protoEstimator{}
}

func (e *protoEstimator) Method() string {
    return "protobuf"
}

func (e *protoEstimator) loadObject() {
	e.obj = templates.GetProtobufObjectForEstimation()
}

func (e *protoEstimator) serialize() (uint32, time.Duration) {
	return impl.SerializeProtoImpl(&e.obj, &e.data)
}

func (e *protoEstimator) deserialize() time.Duration {
	return impl.DeseializeProtoImpl(&e.data)
}

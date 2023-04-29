package estimator

import (
	"serialization_estimator/libs/estimator/impl"
	"serialization_estimator/libs/estimator/templates"
	"time"

	"github.com/vmihailenco/msgpack"
)

type msgpackEstimator struct {
	obj templates.Object
    data []byte
}

func newMsgpackEstimator() *msgpackEstimator {
    return &msgpackEstimator{}
}

func (e *msgpackEstimator) Method() string {
    return "msgpack"
}

func (e *msgpackEstimator) loadObject() {
    e.obj = templates.GetMarshiableObjectForEstimation()
}

func (e *msgpackEstimator) serialize() (uint32, time.Duration) {
	return impl.SerializeImpl(&e.obj, &e.data, msgpack.Marshal)
}

func (e *msgpackEstimator) deserialize() time.Duration {
	return impl.DeseializeImpl(&e.data, msgpack.Unmarshal)
}

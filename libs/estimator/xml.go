package estimator

import (
	"encoding/xml"
	"serialization_estimator/libs/estimator/impl"
	"serialization_estimator/libs/estimator/templates"
	"time"
)

type xmlEstimator struct {
	obj templates.Object
    data []byte
}

func newxmlEstimator() *xmlEstimator {
    return &xmlEstimator{}
}

func (e *xmlEstimator) Method() string {
    return "xml"
}

func (e *xmlEstimator) loadObject() {
    e.obj = templates.GetMarshiableObjectForEstimation()
}

func (e *xmlEstimator) serialize() (uint32, time.Duration) {
	return impl.SerializeImpl(&e.obj, &e.data, xml.Marshal)
}

func (e *xmlEstimator) deserialize() time.Duration {
	return impl.DeseializeImpl(&e.data, xml.Unmarshal)
}

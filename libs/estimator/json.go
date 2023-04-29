package estimator

import (
	"encoding/json"
	"serialization_estimator/libs/estimator/impl"
	"serialization_estimator/libs/estimator/templates"
	"time"
)

type jsonEstimator struct {
	obj templates.Object
    data []byte
}

func newJsonEstimator() *jsonEstimator {
    return &jsonEstimator{}
}

func (e *jsonEstimator) Method() string {
    return "json"
}

func (e *jsonEstimator) loadObject() {
    e.obj = templates.GetMarshiableObjectForEstimation()
}

func (e *jsonEstimator) serialize() (uint32, time.Duration) {
	return impl.SerializeImpl(&e.obj, &e.data, json.Marshal)
}

func (e *jsonEstimator) deserialize() time.Duration {
	return impl.DeseializeImpl(&e.data, json.Unmarshal)
}

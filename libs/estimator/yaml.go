package estimator

import (
	"serialization_estimator/libs/estimator/impl"
	"serialization_estimator/libs/estimator/templates"
	"time"

	"gopkg.in/yaml.v2"
)

type yamlEstimator struct {
	obj templates.Object
    data []byte
}

func newYamlEstimator() *yamlEstimator {
    return &yamlEstimator{}
}

func (e *yamlEstimator) Method() string {
    return "yaml"
}

func (e *yamlEstimator) loadObject() {
    e.obj = templates.GetMarshiableObjectForEstimation()
}

func (e *yamlEstimator) serialize() (uint32, time.Duration) {
	return impl.SerializeImpl(&e.obj, &e.data, yaml.Marshal)
}

func (e *yamlEstimator) deserialize() time.Duration {
	return impl.DeseializeImpl(&e.data, yaml.Unmarshal)
}

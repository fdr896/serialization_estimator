package estimator

import (
	avro_scheme "serialization_estimator/libs/estimator/avro"
	"serialization_estimator/libs/estimator/impl"
	"serialization_estimator/libs/estimator/templates"
	"time"
)

var SCHEMA = avro_scheme.GetAvroSchema()

type avroEstimator struct {
	obj templates.Object
    data []byte
}

func newAvrobufEstimator() *avroEstimator {
    return &avroEstimator{}
}

func (e *avroEstimator) Method() string {
    return "avro"
}

func (e *avroEstimator) loadObject() {
	e.obj = templates.GetAvroObjectForEstimation()
}

func (e *avroEstimator) serialize() (uint32, time.Duration) {
	return impl.SerializeAvroImpl(&e.obj, &e.data, SCHEMA)
}

func (e *avroEstimator) deserialize() time.Duration {
	return impl.DeseializeAvroImpl(&e.data, SCHEMA)
}

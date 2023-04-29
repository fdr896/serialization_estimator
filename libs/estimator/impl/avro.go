package impl

import (
	"serialization_estimator/libs/estimator/templates"
	"time"

	"github.com/hamba/avro/v2"
)

func SerializeAvroImpl(obj *templates.Object, data *[]byte, schema *avro.Schema) (uint32, time.Duration) {
    start := time.Now()
    {
		var err error
		*data, err = avro.Marshal(*schema, obj)
		if err != nil {
			panic(err)
		}
    }
    return uint32(len(*data)), time.Since(start)
}

func DeseializeAvroImpl(data *[]byte, schema *avro.Schema) time.Duration {
    start := time.Now()
    {
		var obj templates.Object
		if err := avro.Unmarshal(*schema, *data, &obj); err != nil {
			panic(err)
		}
    }
    return time.Since(start)
}

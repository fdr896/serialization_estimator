package impl

import (
	"serialization_estimator/libs/estimator/templates"
	"time"
)

type Marshal func(v any) ([]byte, error)
type Unmarshal func(data []byte, v any) error

func SerializeImpl(obj *templates.Object, data *[]byte, marshal Marshal) (uint32, time.Duration) {
    start := time.Now()
    {
		var err error
		*data, err = marshal(obj)
		if err != nil {
			panic(err)
		}
    }
    return uint32(len(*data)), time.Since(start)
}

func DeseializeImpl(data *[]byte, unmarshal Unmarshal) time.Duration {
    start := time.Now()
    {
		var obj templates.Object
		if err := unmarshal(*data, &obj); err != nil {
			panic(err)
		}
    }
    return time.Since(start)
}

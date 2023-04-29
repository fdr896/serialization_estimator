package estimator

import (
	"bytes"
	"encoding/gob"
	"serialization_estimator/libs/estimator/templates"
	"time"
)

type nativeEstimator struct {
    obj templates.Object
    data []byte
}

func newNativeEstimator() *nativeEstimator {
    return &nativeEstimator{}
}

func (e *nativeEstimator) Method() string {
    return "native"
}

func (e *nativeEstimator) loadObject() {
    e.obj = templates.GetMarshiableObjectForEstimation()
}

func (e *nativeEstimator) serialize() (uint32, time.Duration) {
    data := bytes.Buffer{}
    encoder := gob.NewEncoder(&data)

    start := time.Now()
    {
        if err := encoder.Encode(e.obj); err != nil {
            panic(err)
        }
        e.data = data.Bytes()
    }
    return uint32(len(e.data)), time.Since(start)
}

func (e *nativeEstimator) deserialize() time.Duration {
    data := make([]byte, len(e.data))

    start := time.Now()
    {
        copy(data, e.data)
    }
    return time.Since(start)
}

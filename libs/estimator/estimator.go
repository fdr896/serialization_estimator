package estimator

import (
	"fmt"
	"serialization_estimator/libs/support"
	"time"
)

const ESTIMATION_ITERS = 1000
var ESTIMATION_METHODS = support.NewStringSet(
    []string{"native", "json", "xml", "msgpack", "yaml", "protobuf", "avro"})

type Estimator interface {
    Method() string

    loadObject()
    serialize() (uint32, time.Duration)
    deserialize() time.Duration
}

func New(method string) Estimator {
    switch method {
    case "native": return newNativeEstimator()
    case "json": return newJsonEstimator()
    case "xml": return newxmlEstimator()
    case "msgpack": return newMsgpackEstimator()
    case "yaml": return newYamlEstimator()
    case "protobuf": return newProtobufEstimator()
    case "avro": return newAvrobufEstimator()

    default: {
        panic(fmt.Errorf("invalid serialization method: %s", method))
    }
    }
}

func Estimate(e Estimator) (uint32, time.Duration, time.Duration) {
    e.loadObject()

    var avgSize uint32
    var avgSerializeDuration, avgDeserializeDuration time.Duration

    for i := 0; i < ESTIMATION_ITERS; i++ {
        size, serializeDuration := e.serialize()
        desirializeDuration := e.deserialize()

        avgSize = size
        avgSerializeDuration += serializeDuration
        avgDeserializeDuration += desirializeDuration
    }

    avgSerializeDuration /= ESTIMATION_ITERS
    avgDeserializeDuration /= ESTIMATION_ITERS

    return avgSize, avgSerializeDuration, avgDeserializeDuration
}

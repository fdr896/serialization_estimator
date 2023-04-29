package impl

import (
	protobuf "serialization_estimator/libs/estimator/proto"
	"time"

	"google.golang.org/protobuf/proto"
)

func SerializeProtoImpl(obj *protobuf.Object, data *[]byte) (uint32, time.Duration) {
    start := time.Now()
    {
		var err error
		*data, err = proto.Marshal(obj)
		if err != nil {
			panic(err)
		}
    }
    return uint32(len(*data)), time.Since(start)
}

func DeseializeProtoImpl(data *[]byte) time.Duration {
    start := time.Now()
    {
		var obj protobuf.Object
		if err := proto.Unmarshal(*data, &obj); err != nil {
			panic(err)
		}
    }
    return time.Since(start)
}

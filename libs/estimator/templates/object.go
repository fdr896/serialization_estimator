package templates

import (
	protobuf "serialization_estimator/libs/estimator/proto"
	"serialization_estimator/libs/support"
)

type Object struct {
	I int64                    `json:"i" xml:"i" yaml:"i" msgpack:"i" avro:"i"`
	F float32                  `json:"f" xml:"f" yaml:"f" msgpack:"f" avro:"f"`
	S string                   `json:"s" xml:"s" yaml:"s" msgpack:"s" avro:"s"`
	A []string                 `json:"a" xml:"a" yaml:"a" msgpack:"a" avro:"a"`
	M support.SerializiableMap `json:"m" xml:"m" yaml:"m" msgpack:"m" avro:"m"`
}

func GetMarshiableObjectForEstimation() Object {
	return Object{
		I: 1337,
		F: 55.6849494,
		S: "aBacOba7695839..1)%%%33894##010",
		A: []string{"13893", "9595afKfk", "....", "0049)))2838", "8484949"},
		M: support.SerializiableMap{
			"abc": "24443",
			"rnn": "..z9z7505///.a,1",
			"dkd058": ".48jszjGF8rJ4m1",
			"a500": "//////////",
			"s900": "100500",
		},
	}
}

func GetProtobufObjectForEstimation() protobuf.Object {
	obj := GetMarshiableObjectForEstimation()
	return protobuf.Object{
		I: obj.I,
		F: obj.F,
		S: obj.S,
		A: obj.A,
		M: obj.M,
	}
}

func GetAvroObjectForEstimation() Object {
	return GetMarshiableObjectForEstimation()
}

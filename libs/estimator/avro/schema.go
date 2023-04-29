package avro_scheme

import "github.com/hamba/avro/v2"

func GetAvroSchema() *avro.Schema {
	schema, _ := avro.Parse(`{
		"type": "record",
		"name": "object",
		"fields" : [
			{"name": "i", "type": "long"},
			{"name": "f", "type": "float"},
			{"name": "s", "type": "string"},
			{"name": "a", "type": {
				"type": "array",
				"items": "string"
			}},
			{"name": "m", "type": {
				"type": "map",
				"values": "string"
			}}
		]
	}`)

	return &schema
}

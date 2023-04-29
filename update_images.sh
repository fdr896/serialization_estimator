#!/bin/bash

port=2000
for method in "native" "json" "xml" "msgpack" "yaml" "protobuf" "avro"; do
    port=$((port+1))
    
    docker image tag serialization_estimator:$method fdr400/serialization_estimator:$method
    docker image push fdr400/serialization_estimator:$method
done

docker image push fdr400/serialization_estimator_proxy

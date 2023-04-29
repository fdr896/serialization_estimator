#!/bin/bash

port=2000
for method in "native" "json" "xml" "msgpack" "yaml" "protobuf" "avro"; do
    port=$((port+1))

    docker build -t serialization_estimator:$method -f bin/estimator/Dockerfile \
        --build-arg PORT=$port --build-arg METHOD=$method .
done

docker build -t fdr400/serialization_estimator_proxy -f bin/proxy/Dockerfile \
    --build-arg PORT=2000 \
    --build-arg NATIVE_PORT=2001 --build-arg JSON_PORT=2002 \
    --build-arg XML_PORT=2003 --build-arg MSGPACK_PORT=2004 \
    --build-arg YAML_PORT=2005 --build-arg PROTOBUF_PORT=2006 \
    --build-arg AVRO_PORT=2007 .

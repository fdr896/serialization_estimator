#!/bin/bash

BIN_PATH="./bin"
ESTIMATOR_PATH="$BIN_PATH/estimator"
PROXY_PATH="$BIN_PATH/proxy"

MULTICAST_IP=224.0.0.7

# build
go build -o $ESTIMATOR_PATH/estimator $ESTIMATOR_PATH/main.go
go build -o $PROXY_PATH/proxy $PROXY_PATH/main.go

port=2000
for method in "native" "json" "xml" "msgpack" "yaml" "protobuf" "avro"; do
    port=$((port+1))

    MULTICAST_ADDR=$MULTICAST_IP \
    MULTICAST_PORT=2010 \
    $ESTIMATOR_PATH/estimator --port $port --method $method &
done

MULTICAST_ADDR=$MULTICAST_IP \
MULTICAST_PORT=2010 \
NATIVE_PORT=2001 \
JSON_PORT=2002 \
XML_PORT=2003 \
MSGPACK_PORT=2004 \
YAML_PORT=2005 \
PROTOBUF_PORT=2006 \
AVRO_PORT=2007 \
$PROXY_PATH/proxy --port 2000

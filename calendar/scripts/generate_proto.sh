#!/usr/bin/env bash

cd api/protobuf-spec

protoc ./events.proto --go_out=plugins=grpc:../../internal/pkg/events

cd -

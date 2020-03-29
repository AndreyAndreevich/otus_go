#!/usr/bin/env bash

cd api

protoc ./events.proto --go_out=plugins=grpc:../external/pkg/events

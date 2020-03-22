#!/usr/bin/env bash
node opc/opc.js
go run checkerr/main.go
go run unloader/main.go
go run pusher/main.go
go run api/main.go
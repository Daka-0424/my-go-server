#!/usr/bin/env bash
set -euo pipefail

go mod download

go install github.com/ramya-rao-a/go-outline@latest
go install github.com/cweill/gotests/gotests@v1.6.0
go install github.com/fatih/gomodifytags@latest
go install github.com/josharian/impl@v1.4.0
go install github.com/haya14busa/goplay/cmd/goplay@latest
go install github.com/go-delve/delve/cmd/dlv@v1.27.2
go install honnef.co/go/tools/cmd/staticcheck@v0.8.2
go install golang.org/x/tools/gopls@v0.24.0
go install github.com/golang/mock/mockgen@v1.7.0-rc.1
go install github.com/swaggo/swag/cmd/swag@v1.8.12

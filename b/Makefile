VERSION := 0.2

.PHONY: build-BMainFunction
build-BMainFunction:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=${VERSION}" -o ${ARTIFACTS_DIR}/main

VERSION := 0.1

.PHONY: build-MainFunction
build-MainFunction:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=${VERSION}" -o ${ARTIFACTS_DIR}/main

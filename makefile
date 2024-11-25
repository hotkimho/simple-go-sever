build:
	go mod tidy
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build
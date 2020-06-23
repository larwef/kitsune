APP_NAME=kitsune
VERSION=v0.0.1

TARGET=target
SRC=cmd/$(APP_NAME)/main.go

all: test build-linux build-mac build-windows

run:
	go run cmd/kitsune/main.go

# PHONY used to mitigate conflict with dir name test
.PHONY: test
test:
	go mod tidy
	go fmt ./...
	go vet ./...
	golint ./...
	go test ./...

integration:
	go test ./... -tags=integration

coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out

build-linux:
	GOOS=linux go build -ldflags "-X main.Version=$(VERSION)" -o $(TARGET)/linux/$(APP_NAME) $(SRC)
	zip -j $(TARGET)/deployment-linux.zip $(TARGET)/linux/$(APP_NAME)

build-mac:
	GOOS=darwin go build -ldflags "-X main.Version=$(VERSION)" -o $(TARGET)/mac/$(APP_NAME) $(SRC)
	zip -j $(TARGET)/deployment-mac.zip $(TARGET)/mac/$(APP_NAME)

build-windows:
	GOOS=windows go build -ldflags "-X main.Version=$(VERSION)" -o $(TARGET)/windows/$(APP_NAME) $(SRC)
	zip -j $(TARGET)/deployment-windows.zip $(TARGET)/windows/$(APP_NAME)

release: test build-linux build-mac build-windows
	git tag -a $(VERSION) -m "Release $(VERSION)"
	git push origin $(VERSION)

clean:
	rm -rf $(TARGET)

rebuild:
	clean all

doc:
	godoc -http=":6060"

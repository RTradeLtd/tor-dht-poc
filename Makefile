
GO111MODULE=on

#GOPATH=$(PWD)/.go

VERSION=0.0.01
USER_GH=RTradeLtd

echo:
	@echo $(GOPATH)
	@echo "type make version to do release $(VERSION)"


version:
	gothub release -s $(GITHUB_TOKEN) -u $(USER_GH) -r go-garlic-tcp-transport -t v$(VERSION) -d "version $(VERSION)"

del:
	gothub delete -s $(GITHUB_TOKEN) -u $(USER_GH) -r go-garlic-tcp-transport -t v$(VERSION)


build:
	go build -o i2p-dht-poc ./go-i2p-dht-poc

fmt:
	find . -name '*.go' -exec gofmt -w {} \;

set:
	find ./go-i2p-dht-poc -name '*.go' -exec sed -i 's|cretz\/tor-dht-poc\/go-i2p-dht-poc|RTradeLtd\/tor-dht-poc\/go-i2p-dht-poc|g' {} \;

reset:
	find ./go-i2p-dht-poc -name '*.go' -exec sed -i 's|RTradeLtd\/tor-dht-poc\/go-i2p-dht-poc|cretz\/tor-dht-poc\/go-i2p-dht-poc|g' {} \;

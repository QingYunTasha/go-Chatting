HUB?=myhub
IMAGE=${HUB}/go-chatting

all:
	build push

gobuild:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cmd/main.exe cmd/main.go

build:
	docker build -t $(IMAGE):latest .
	
push:
	docker push $(IMAGE):latest 
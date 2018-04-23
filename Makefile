GOPATH = $(shell go env GOPATH)
GOCMD = go
GOBUILD = $(GOCMD) build -o
GOCLEAN = $(GOCMD) clean
GOPARM  = CGO_ENABLED=0 GOOS=linux GOARCH=amd64
DOCKERFILE = ./docker/Dockerfile
APP = ./docker/app
DOCKERIMAGE = registry.gitlab.com/iisidotnetgroup/ifttt

clean:
	$(GOCLEAN) -r

run:
	$(GOCMD) run *.go

build:
	$(GOPARM) $(GOBUILD) main
	mkdir -p $(APP)
	cp ./main $(APP)/main
	cp -a ./config $(APP)/config

dockerbuild: build
	docker build --file $(DOCKERFILE) --build-arg APP=$(APP) -t $(DOCKERIMAGE) .
	rm -r $(APP)

dockerpush:
	docker login registry.gitlab.com
	docker push $(DOCKERIMAGE)

dockerrun:
	docker run --rm -it -p 8888:8888 $(DOCKERIMAGE)

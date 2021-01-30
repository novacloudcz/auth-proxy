OWNER=novacloud
IMAGE_NAME=auth-proxy
QNAME=$(OWNER)/$(IMAGE_NAME)

GIT_TAG=$(QNAME):$(GITHUB_SHA)
BUILD_TAG=$(QNAME):$(GITHUB_RUN_ID).$(GITHUB_SHA)
TAG=$(QNAME):`echo $(GITHUB_REF) | sed 's/refs\/heads\///' | sed 's/master/latest/;s/develop/unstable/'`

lint:
	docker run -it --rm -v "$(PWD)/Dockerfile:/Dockerfile:ro" redcoolbeans/dockerlint

build:
	# go get ./...
	# gox -osarch="linux/amd64" -output="bin/devops-alpine"
	# CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/binary .
	docker build -t $(GIT_TAG) .
	
tag:
	docker tag $(GIT_TAG) $(BUILD_TAG)
	docker tag $(GIT_TAG) $(TAG)
	
login:
	@docker login -u "$(DOCKER_USER)" -p "$(DOCKER_PASSWORD)"
push: login
	# docker push $(GIT_TAG)
	# docker push $(BUILD_TAG)
	docker push $(TAG)

generate:
	go run github.com/99designs/gqlgen
	go generate ./...

build-local:
	# go get ./...
	# go build -o $(IMAGE_NAME) ./server/server.go
	go build -o app

deploy-local:
	make build-local
	mv app /usr/local/bin/${IMAGE_NAME}

run:
	make build-local && REQUIRED_JWT_SCOPES="test aa" PROXY_URL=http://example.com/ JWKS_PROVIDER_URL=https://id.novacloud.cz/.well-known/jwks.json PORT=8080 ./app

# test:
# 	DATABASE_URL=sqlite3://test.db $(IMAGE_NAME) server -p 8005
	# DATABASE_URL="mysql://root:root@tcp(localhost:3306)/test?parseTime=true" go run *.go server -p 8000

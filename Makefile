API = api/v1
IMAGE_NAME = quake-grcp-server
CONTAINER_NAME= quake-grpc-server

.PHONY: all
.DEFAULT: all
all: api-stubs docker-image

api-stubs:
	$(MAKE) -C $(API) all

.PHONY: docker-image
docker-image:
	docker build -f mapi.Dockerfile . -t $(IMAGE_NAME)

.PHONY: swagger
swagger: docker-image
	docker run -t --rm -d --name $(CONTAINER_NAME)-tmp $(IMAGE_NAME)
	docker cp $(CONTAINER_NAME)-tmp:/opt/quake/api/v1/quake_api.swagger.json .
	docker rm -f $(CONTAINER_NAME)-tmp

.PHONY: run
run: swagger
	docker rm -f $(CONTAINER_NAME) || true
	docker run -it --rm -d -p 8081:8081 --name $(CONTAINER_NAME) $(IMAGE_NAME)

.PHONY: fuzz
fuzz: run
	-mapi target create forallsecure/abrewer-quake-grpc-example http://localhost:8081
	mapi run forallsecure/abrewer-quake-grpc-example 30 quake_api.swagger.json

.PHONY: stop
stop:
	docker rm -f $(CONTAINER_NAME)

clean:
	-rm -f quake_api.swagger.json
	$(MAKE) -C $(API) clean

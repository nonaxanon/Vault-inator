#!make
include .env

SHELL=/bin/bash
BUILD_ID=$(shell date +%Y%m%d%H%M)
IMAGE_NAME?=gameserver
IMAGE_TAG?=vlabs
REGISTRY?=jorvelazquez3
DEPLOYMENT_FILE?=./deployment.yaml

.PHONY: api
api:
	cd api && make build-api

.PHONY: clean-api
clean-api:
	cd api && make clear-gen-directories

.PHONY: remove-validate-js-ts
remove-validate-js-ts:
	cd api/gen/web && find . -name "*.d.ts" | xargs sed -i 's/import \* as validate_validate_pb from '\''\.\.\/\.\.\/\.\.\/validate\/validate_pb'\''\;//g'; \
			find . -name "*_pb.js" | xargs sed -i '/validate/d'; \

.PHONY: remove-validate
remove-validate:
	cd api/gen/csharp && find . -name "*.cs" | xargs sed -e s/global::Validate.ValidateReflection.Descriptor,//g -i *;

.PHONY: add-serializable
add-serializable:
	cd api/gen/csharp && find . -name "*.cs" | xargs sed -e '/public sealed partial class PlayerCards /i [Serializable]' -i
	cd api/gen/csharp && find . -name "*.cs" | xargs sed -e '/using pb = global::Google.Protobuf/i using System;' -i;

.PHONY: sql
sql:
	cd tools/gorm2sql && go run main.go postgresql --f=../../internal/migrator/aggregates/aggregates.pb.gorm.go --s=ShopCardORM  --o="../../internal/db.sql"

# delete images
.PHONY: delete-image
delete-image $(image):
	docker image remove $(image)

.PHONY: delete
delete $(SVC):
	kubectl delete -f internal/$(SVC)/deployment.yaml

.PHONY: myproxy
myproxy:
	docker build -f ./internal/proxy_image/Dockerfile -t jorvelazquez3/proxy:latest .
	docker push  projectaresdemoacr.azurecr.io/proxy:0.0.1

.PHONY: mysaltapi
mysaltapi:
	docker build -f ./internal/salt_images/salt-api/Dockerfile -t jorvelazquez3/salt-api:latest .
	docker push  projectaresdemoacr.azurecr.io/salt-api:latest

.PHONY: basic
basic:
	docker build -f ./internal/taskvalidator/templates/basic/Dockerfile -t jorvelazquez3/basic:latest .
	docker push jorvelazquez3/basic:latest

# Postgres db for local dev, change the values as needed
.PHONY: postgres postgres-clean
postgres:
	docker run --name dev-sql-db \
		-e POSTGRES_PASSWORD=admin \
		-e POSTGRES_USER=admin \
		-e POSTGRES_DB=dev \
		-d --rm -it -p 5432:5432 postgres

postgres-clean:
	docker kill dev-sql-db

.PHONY: run-dev
run-dev:
	go run -tags dev ./cmd/server/...

.PHONY: run-prod
run-prod:
	go run ./cmd/server/...

.PHONY: build-dev
build-dev:
	go build -tags dev -o bin/server-dev ./cmd/server/...

.PHONY: build-prod
build-prod:
	go build -o bin/server ./cmd/server/...

.PHONY: build-basic
build-basic:
	docker build -f ./container_templates/basic/Dockerfile -t jorvelazquez3/basic:latest .
	docker push jorvelazquez3/basic:latest

.PHONY: build-vlabs
build-vlabs:
	docker build --build-arg PAT=$(PAT) -f ./Dockerfile -t jorvelazquez3/gameserver:vlabs .
	docker push jorvelazquez3/gameserver:vlabs
	kubectl delete -f ./deployment.yaml
	kubectl apply -f ./deployment.yaml

.PHONY: deploy
deploy:
	@echo "Building Docker image $(REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)..."
	docker build -t $(REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG) .

	@echo "Pushing Docker image to registry..."
	docker push $(REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)

	@echo "Deleting existing deployment..."
	kubectl delete -f $(DEPLOYMENT_FILE)

	@echo "Applying new deployment..."
	kubectl apply -f $(DEPLOYMENT_FILE)

	@echo "Deployment complete!"


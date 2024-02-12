DOCKER_REGISTRY ?= docker.io
REPOSITORY ?= sujaykumarsuman
IMAGE_NAME ?= kubepod
IMAGE_VERSION ?= v1.0.0-im-$(shell git rev-parse --short HEAD)
CHART_VERSION ?= v1.0.0-hc-$(shell git rev-parse --short HEAD)


docker-image-exists:
	@if docker images $(REPOSITORY)/$(IMAGE_NAME):$(IMAGE_VERSION) | awk 'NR>1 { exit 1 }'; then \
		echo "Docker image does not exist. Building..."; \
		make docker-build; \
	else \
		echo "Docker image already exists. Skipping build."; \
	fi
.PHONY: docker-image-exists

docker-build:
	@docker build -t $(REPOSITORY)/$(IMAGE_NAME):$(IMAGE_VERSION) .
.PHONY: docker-build

docker-push: docker-image-exists
	@docker push $(REPOSITORY)/$(IMAGE_NAME):$(IMAGE_VERSION)
.PHONY: docker-push

build-dir-exist:
	@mkdir -p build
.PHONY: build-dir-exist

go-build: build-dir-exist
	@go build -o build/$(IMAGE_NAME) main.go
.PHONY: go-build

update-tag-in-values:
	@sed -i.bak 's@tag:.*@tag: $(IMAGE_VERSION)@g' ./helm/$(IMAGE_NAME)/values.yaml
	@rm -f ./helm/$(IMAGE_NAME)/values.yaml.bak
.PHONY: update-tag-in-values

helm-package: build-dir-exist update-tag-in-values
	@helm package ./helm/$(IMAGE_NAME) --destination ./build --version $(CHART_VERSION)
.PHONE: helm-package

helm-push: helm-package
	@helm push ./build/$(IMAGE_NAME)-$(CHART_VERSION).tgz oci://$(DOCKER_REGISTRY)/$(REPOSITORY)
.PHONE: helm-push

helm-install:
	@helm install $(IMAGE_NAME) ./helm/$(IMAGE_NAME)
.PHONE: helm-install

helm-uninstall:
	@helm uninstall $(IMAGE_NAME)
.PHONE: helm-uninstall

cleanup:
	@rm -rf build
	@docker rmi -f $(REPOSITORY)/$(IMAGE_NAME):$(IMAGE_VERSION)
.PHONY: cleanup

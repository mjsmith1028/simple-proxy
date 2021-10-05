PWD := ${CURDIR}
IMAGE ?= mjsmith1028/simple-proxy
TAG ?= latest
HOST_PORT ?= 80
APP_HOST ?= 169.254.169.254
APP_PORT ?= 8080
APP_PROTOCOL ?= http://
K8S_NAMESPACE ?= default
TEST_URL ?= http://$(APP_HOST)/metadata/identity/oauth2/token/?resource=https://management.core.windows.net/

image:
	docker build -t $(IMAGE):$(TAG) .

publish: image
	docker push $(IMAGE):$(TAG)

add-alias:
	sudo ifconfig lo0 $(APP_HOST) alias

remove-alias:
	sudo ifconfig lo0 $(APP_HOST) -alias

port-forward:
	sudo kubectl port-forward -n $(K8S_NAMESPACE) deployment/simple-proxy $(HOST_PORT):$(APP_PORT) --address=$(APP_HOST)

install:
	helm install simple-proxy charts/simple-proxy --namespace $(K8S_NAMESPACE)

uninstall:
	helm uninstall simple-proxy --namespace $(K8S_NAMESPACE)

test-connection:
	curl $(TEST_URL)

.PHONY: image publish add-alias remove-alias port-forward install uninstall test-connection

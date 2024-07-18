GO_PATH=$(shell go env GOPATH)

.PHONY: tools-install
tools-install:
	@echo "installing dev tools"
	@go install k8s.io/code-generator/cmd/deepcopy-gen@v0.29.0

# .PHONY: crd-gen
# crd-gen:
# 	@${GO_PATH}/bin/controller-gen crd paths=github.com/kong/dataplane-controller/api/v1alpha1 output:crd:dir=deployments/chart/crds/ output:stdout

.PHONY: deepcopy-gen
deepcopy-gen:
	go run sigs.k8s.io/controller-tools/cmd/controller-gen object:headerFile="internal/model/boilerplate.go.txt" paths="./..."

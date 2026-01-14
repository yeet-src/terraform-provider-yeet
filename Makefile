.PHONY: build install clean test fmt help init plan apply destroy reinstall

BINARY_NAME=terraform-provider-yeet
BUILD_DIR=.
TERRAFORM_PLUGINS_DIR=~/.terraform.d/plugins
PROVIDER_NAMESPACE=local/yeet-src/yeet
PROVIDER_VERSION=1.0.0
OS_ARCH=$$(go env GOOS)_$$(go env GOARCH)
INSTALL_PATH=$(TERRAFORM_PLUGINS_DIR)/$(PROVIDER_NAMESPACE)/$(PROVIDER_VERSION)/$(OS_ARCH)

help:
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build:
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME)
	@echo "Build complete: $(BINARY_NAME)"

install: build
	@echo "Installing provider to $(INSTALL_PATH)..."
	@mkdir -p $(INSTALL_PATH)
	@cp $(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "Provider installed successfully!"
	@echo ""
	@echo "Add this to your ~/.terraformrc to use the local provider:"
	@echo ""
	@echo "provider_installation {"
	@echo "  dev_overrides {"
	@echo "    \"$(PROVIDER_NAMESPACE)\" = \"$(shell pwd)\""
	@echo "  }"
	@echo "  direct {}"
	@echo "}"

reinstall: clean install

clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)
	rm -rf .terraform
	rm -f .terraform.lock.hcl
	rm -f terraform.tfstate
	rm -f terraform.tfstate.backup
	rm -f crash.log
	@echo "Clean complete"

fmt:
	go fmt ./...

test:
	go test -v ./...

tidy:
	go mod tidy

init:
	@echo "Note: 'terraform init' is not needed when using dev_overrides"
	@echo "Just run 'make plan' or 'make apply' directly"

plan: build
	terraform plan

apply: build
	terraform apply -auto-approve

destroy:
	terraform destroy -auto-approve

show:
	terraform show

output:
	terraform output

dev-setup:
	@echo "Setting up development environment..."
	@$(MAKE) build
	@echo ""
	@echo "Add this to your ~/.terraformrc:"
	@echo ""
	@echo "provider_installation {"
	@echo "  dev_overrides {"
	@echo "    \"$(PROVIDER_NAMESPACE)\" = \"$(shell pwd)\""
	@echo "  }"
	@echo "  direct {}"
	@echo "}"
	@echo ""
	@echo "Or run: make terraformrc >> ~/.terraformrc"

terraformrc:
	@echo ""
	@echo "provider_installation {"
	@echo "  dev_overrides {"
	@echo "    \"$(PROVIDER_NAMESPACE)\" = \"$(shell pwd)\""
	@echo "  }"
	@echo "  direct {}"
	@echo "}"
	@echo ""

quick-test: build
	@terraform apply -auto-approve

quick-destroy:
	@terraform destroy -auto-approve

cycle: clean build apply

all: clean build install

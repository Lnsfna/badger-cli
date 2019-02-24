BUILD_IMAGE ?= teletraan/golang:1.11-ubuntu18.04

container-bin:
	@echo "build container binary"
	@docker run -it --rm \
		-v "${GOPATH}/pkg/mod/:/go/pkg/mod" \
		-v "$$(pwd):/app" \
		-w /app \
		$(BUILD_IMAGE) \
		go build -v
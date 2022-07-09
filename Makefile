.PHONY: help
## help: prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: clean
## clean: remove all generated files
clean:
	@echo "Cleaning generated files..."
	@rm -r -f output/*

.PHONY: build
## build: run genblog, copy static files
build:
	@genblog; cp -R _static/ output/

.PHONY: build-podman
## build-podman: build image with podman
build-podman:
	@podman build -t micro:local .

.PHONY: run-podman
## run-podman: run image with podman
run-podman:
	@podman stop micro || true;
	@podman rm micro || true;
	@podman run -d --name micro -p 8080:80 micro:local

.PHONY: rebuild-podman
## rebuild-podman: alias for buid, build-podman, run-podman
rebuild-podman: build build-podman run-podman

.PHONY: run-docker
## run-docker: run the docker container
run-docker:
	@docker-compose up -d nginx

.PHONY: watch
## watch: watch for changes and rebuild
watch:
	@echo "Watching for file changes..."
	@fswatch -or -e "output" -e ".git" . | xargs -n1 sh -c "genblog; cp -R _static/ output/"

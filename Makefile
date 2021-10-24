.PHONY: help
## help: prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: build
## build: run genblog, copy static files
build:
	@genblog; cp -R static/ output/

.PHONY: run-docker
## run-docker: run the docker container on http://127.0.0.1:8080
run-docker:
	@docker-compose up -d nginx

.PHONY: watch
## watch: watch for changes and rebuild
watch:
	@echo "Watching for file changes..."
	@fswatch -or -e "output" -e ".git" . | xargs -n1 sh -c "genblog; cp -R static/ output/"

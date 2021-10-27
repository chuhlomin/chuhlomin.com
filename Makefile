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
	@genblog; cp -R static/ output/

.PHONY: run-docker
## run-docker: run the docker container
run-docker:
	@docker-compose up -d nginx

.PHONY: watch
## watch: watch for changes and rebuild
watch:
	@echo "Watching for file changes..."
	@fswatch -or -e "output" -e ".git" . | xargs -n1 sh -c "genblog; cp -R static/ output/"

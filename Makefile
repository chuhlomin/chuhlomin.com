.PHONY: help
## help: prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: clean
## clean: remove all generated files
clean:
	@echo "Cleaning generated files..."
	@rm -r -f output search_index search_index_temp

.PHONY: run-genblog
## run-genblog: run genblog, copy static files
run-genblog:
	@genblog

.PHONY: build-generator
## build-generator: build the generator
build-generator:
	@echo "Building generator..."
	@cd generator && go build -o ../bin/generator

.PHONY: run-generator
## run-generator: run the generator
run-generator:
	@echo "Running generator..."
	@./bin/generator

.PHONY: build
## build: build & run generator
build: build-generator run-generator

.PHONY: watch
## watch: watch for changes and rebuild
watch:
	@echo "Watching for file changes..."
	@fswatch -or -e "output" -e ".git" . | xargs -n1 sh -c "genblog; cp -R _static/ output/"

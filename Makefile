.PHONY: help
## help: prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: clean
## clean: remove all generated files
clean:
	@echo "Cleaning generated files..."
	@rm -r -f output search_index

.PHONY: build
## build: run genblog, copy static files
build:
	@genblog; cp -R _static/ output/

.PHONY: watch
## watch: watch for changes and rebuild
watch:
	@echo "Watching for file changes..."
	@fswatch -or -e "output" -e ".git" . | xargs -n1 sh -c "genblog; cp -R _static/ output/"

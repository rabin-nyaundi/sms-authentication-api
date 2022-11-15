.PHONY: help
help:
	@echo 'Usage'
	@echo -n "s/^##//p" ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'


.PHONY: run
run/api:
	@echo 
	@go run cmd/api/*.go
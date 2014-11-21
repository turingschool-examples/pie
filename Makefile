default: generate

generate:
	@ego templates
	@go fmt ego.go

.PHONY: default generate

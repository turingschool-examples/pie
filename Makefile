default: generate

generate: assets templates

assets:
	@go-bindata -pkg assets -prefix assets -o assets/bindata.go -ignore bindata.go assets

templates:
	@ego templates
	@go fmt ego.go

.PHONY: assets default generate templates

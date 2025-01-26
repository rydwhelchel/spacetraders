# May be necessary to build with these tags for sqlite
build:
	@go build -tags "darwin arm64" -o ./tmp/spacetraders

run: build
	@./tmp/spacetraders

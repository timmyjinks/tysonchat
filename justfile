build:
    pnpm run build
    go build .

run: build
    go run .

fmt:
    gofumpt -w .

install: build
    go install -v .

build:
	CGO_ENABLED=0 GOARCH=adm64 GOOS=linux go build ./cmd/goshortener

dev:
	modd

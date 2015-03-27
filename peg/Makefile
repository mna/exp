pigeon.go:
	go run ./bootstrap/cmd/bootstrap-build/main.go grammar/bootstrap.peg | goimports > bootstrap/cmd/bootstrap-pigeon/bootstrap-pigeon.go
	go build -o bootstrap/cmd/bootstap-pigeon/bootstrap-pigeon ./bootstrap/cmd/bootstrap-pigeon/
	bootstrap/cmd/bootstrap-pigeon/bootstrap-pigeon grammar/pigeon.peg | goimports > $@

.PHONY: pigeon.go


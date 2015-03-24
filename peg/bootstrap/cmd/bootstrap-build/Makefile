GRAMMAR ?= "../../bootstrap-grammar.peg"

../bootstrap-pigeon/bootstrap-pigeon.go:
	go run main.go $(GRAMMAR) | goimports > $@

.PHONY: ../bootstrap-pigeon/bootstrap-pigeon.go

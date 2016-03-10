cmdnames = server client callee
cmds = $(addprefix juggler-, $(cmdnames))

all: $(cmds)

$(cmds):
	go build $(flags) ./cmd/$@ 

.PHONY: all $(cmds)


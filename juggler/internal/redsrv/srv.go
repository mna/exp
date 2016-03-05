package redsrv

import (
	"github.com/PuerkitoBio/exp/juggler/msg"
	"github.com/pborman/uuid"
)

// TODO : move elsewhere
type Caller interface {
	// push call onto the list
	Call(msg.Call) error

	// loops and BRPOPs over result list
	ProcessResult(connUUID uuid.UUID) error
}

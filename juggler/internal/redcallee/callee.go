package redcallee

// TODO : to move elsewhere, but just thinking out loud here...
type Callee interface {
	ProcessURI(uri string) error // loops and BRPOPs calls, invokes, stores results
}

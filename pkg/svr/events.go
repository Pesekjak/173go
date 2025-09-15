package svr

const (
	StartEvt string = "server.start"
	StopEvt         = "server.stop"
)

type StartEventData struct {
	Server *Server
}

type StopEventData struct {
	Server *Server
}

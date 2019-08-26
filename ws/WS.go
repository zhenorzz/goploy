package ws

// State
const (
	Fail    = 0
	Success = 1
)

// Init websocket
func Init() {
	go GetSyncHub().Run()
}

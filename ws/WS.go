package ws

// Init websocket
func Init() {
	go GetSyncHub().Run()
}

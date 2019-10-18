package ws

// Init websocket
func Init() {
	go GetUnicastHub().Run()
	go GetBroadcastHub().Run()
}

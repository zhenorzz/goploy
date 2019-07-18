package cron

import (
	"github.com/jasonlvhit/gocron"
)

// Run cron
func Run() {
	// Do jobs without params
	gocron.Every(1).Second().Do(updatePublishState)
	<-gocron.Start()
}

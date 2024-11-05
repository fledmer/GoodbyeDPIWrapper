package dpi

import (
	"dpiWrapper/internal/config"
	"dpiWrapper/internal/proc"
	"fmt"
	"log"
)

var singleProc proc.Async
var Сfg *config.Config

func Stop() {
	err := singleProc.Stop()
	if err != nil {
		fmt.Println("failed to stop proc: %w", err)
		return
	}
	log.Println("process is stopped")
}

func RunYoutube() {
	_, err := singleProc.Start(Сfg.DPIYoutubePath, Сfg.DpiYoutubeCommands...)
	if err != nil {
		log.Println("failed to start youtube dpi: ", err)
		return
	}
	log.Println("youtube dpi started!")
}

package main

import (
	"bufio"
	"dpiWrapper/internal/commands"
	"dpiWrapper/internal/config"
	"dpiWrapper/internal/dpi"
	"dpiWrapper/internal/vpn"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	// procManager *manager.Manager
	cfg     *config.Config
	scanner *bufio.Scanner
)

func main() {
	bytes, err := os.ReadFile("./.cfg")
	if err != nil {
		log.Fatalln("failed to read config", err)
	}
	cfg = &config.Config{}
	err = json.Unmarshal(bytes, cfg)
	if err != nil {
		log.Fatalln("failed to parse config", err)
	}
	vpn.Cfg = cfg.VPN
	dpi.Сfg = cfg
	scanner = bufio.NewScanner(os.Stdin)
	printStartMessage()
	inputLoop()
}

func printStartMessage() {
	fmt.Println("GoodbyeDPI wrapped")
}

func inputLoop() {
	fmt.Printf(">")
	for scanner.Scan() {
		switch parseCommand(scanner.Text()) {
		case commands.StartYoutube:
			dpi.RunYoutube()
		case commands.StartRussiaBlackList:
		case commands.StartVPN:
			vpn.Run()
		case commands.StopVPN:
			vpn.Stop()
		case commands.StopDpi:
			dpi.Stop()
		case commands.Help:
			fmt.Println(commands.CommandMap.String())
		case commands.Exit:
			os.Exit(0)
		default:
			log.Println("unknown command, input: 'h' for help")
		}
		fmt.Printf(">")
	}
}

func parseCommand(text string) commands.Command {
	return commands.CommandMap[strings.ToLower(text)]
}

package commands

import (
	"fmt"
	"maps"
	"slices"
)

func init() {
	CommandMap = map[string]Command{
		"start_youtube": StartYoutube,
		"start_y":       StartYoutube,
		"s_y":           StartYoutube,
		"sy":            StartYoutube,
		"youtube":       StartYoutube,
		"y":             StartYoutube,

		"start_blacklist": StartRussiaBlackList,
		"start_b":         StartRussiaBlackList,
		"s_b":             StartRussiaBlackList,
		"sb":              StartRussiaBlackList,
		"blacklist":       StartRussiaBlackList,

		"!d":       StopDpi,
		"!y":       StopDpi,
		"stop_dpi": StopDpi,

		"v":         StartVPN,
		"start_vpn": StartVPN,
		"sv":        StartVPN,

		"!v":       StopVPN,
		"stop_vpn": StopVPN,

		"help": Help,
		"h":    Help,
		"exit": Exit,
		"e":    Exit,
	}
}

type Command string

const (
	BadRequest           Command = ""
	StartYoutube         Command = "Start youtube"
	StartRussiaBlackList Command = "Start russia blackList"
	StopDpi              Command = "Stop any dpi"
	StartVPN             Command = "Start vpn"
	StopVPN              Command = "Stop vpn"
	Help                 Command = "Help"
	Exit                 Command = "Exit"
)

type commandMap map[string]Command

func (c commandMap) String() (str string) {
	keyMap := map[Command][]string{}
	commands := []Command{}
	for key := range maps.Keys(c) {
		sl := keyMap[c[key]]
		keyMap[c[key]] = append(sl, key)
	}
	for command := range maps.Keys(keyMap) {
		slices.Sort(keyMap[command])
		commands = append(commands, command)
	}
	slices.Sort(commands)
	for _, command := range commands {
		str += fmt.Sprintf("For %s: use %s \n", command, keyMap[command])
	}
	return str
}

var CommandMap commandMap

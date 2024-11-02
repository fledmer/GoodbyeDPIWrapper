package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"maps"
	"os"
	"os/exec"
	"slices"
	"strings"
)

func init() {
	commandMap = map[string]Command{
		"start_youtube": StartYoutube,
		"start_y":       StartYoutube,
		"s_y":           StartYoutube,
		"sy":            StartYoutube,
		"youtube":       StartYoutube,

		"start_blacklist": StartRussiaBlackList,
		"start_b":         StartRussiaBlackList,
		"s_b":             StartRussiaBlackList,
		"sb":              StartRussiaBlackList,
		"blacklist":       StartRussiaBlackList,

		"stop": StopProc,
		"s":    StopProc,

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
	StopProc             Command = "Stop process"
	Help                 Command = "Help"
	Exit                 Command = "Exit"
)

type CommandMap map[string]Command

func (c CommandMap) String() (str string) {
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

var commandMap CommandMap

type Config struct {
	DPIYoutubePath     string   `json:"dpi_youtube_path,omitempty"`
	DpiYoutubeCommands []string `json:"dpi_youtube_commands"`
	RussiaBlackListDPI string   `json:"russia_black_list_dpi,omitempty"`
}

var (
	cfg     *Config
	scanner *bufio.Scanner
	proc    *exec.Cmd
)

func main() {
	bytes, err := os.ReadFile("./.cfg")
	if err != nil {
		log.Fatalln("failed to read config", err)
	}
	cfg = &Config{}
	err = json.Unmarshal(bytes, cfg)
	if err != nil {
		log.Fatalln("failed to parse config", err)
	}
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
		case StartYoutube:
			startYoutube()
		case StartRussiaBlackList:
		case StopProc:
			stop()
		case Help:
			fmt.Println(commandMap.String())
		case Exit:
			os.Exit(0)
		default:
			log.Println("unknown command, input: 'h' for help")
		}
		fmt.Printf(">")
	}
}

func parseCommand(text string) Command {
	return commandMap[strings.ToLower(text)]
}

func startYoutube() {
	localProc, err := startSingleProc(cfg.DPIYoutubePath, cfg.DpiYoutubeCommands...)
	if err != nil {
		log.Println("failed to start youtube: ", err)
		return
	}
	outReader, err := localProc.StdoutPipe()
	if err != nil {
		log.Println("failed to read stdout pipe: ", err)
		return
	}
	if err = localProc.Start(); err != nil {
		log.Println("failed to start process: ", err)
		return
	} else {
		log.Println("youtube dpi started")
		go readOutput(outReader)
		proc = localProc
	}
}

func stop() {
	if proc == nil {
		log.Println("no running process")
		return
	}
	if proc != nil {
		proc.Process.Kill()
		state, err := proc.Process.Wait()
		if err != nil {
			log.Println("failed to stop proc: ", err)
			return
		}
		if state.Exited() {
			log.Println("process is stopped")
			proc = nil
		} else {
			log.Println("process is not exited")
		}
	}
}

func startSingleProc(path string, args ...string) (*exec.Cmd, error) {
	if err := singleProcCheck(); err != nil {
		return nil, fmt.Errorf("failed to start new proc single check failed, %w", err)
	}

	proc := exec.Command(path, args...)
	if proc.Err != nil {
		return nil, fmt.Errorf("failed to start new proc %w", proc.Err)
	}
	return proc, nil
}

func singleProcCheck() (err error) {
	if proc != nil {
		return fmt.Errorf("dpi proc with pid: %v already run", proc.String())
	}
	return
}

func readOutput(reader io.ReadCloser) {
	bytes := make([]byte, 4096)
	var err error
	var n int
	for err == nil {
		n, err = reader.Read(bytes)
		os.Stdout.Write(bytes[:n])
	}
}

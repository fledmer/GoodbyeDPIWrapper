package vpn

import (
	"dpiWrapper/internal/proc"
	"log"
)

type Config struct {
	WireguardConfPath string `json:"wireguard_conf_path,omitempty"`
	TunelName         string `json:"tunel_name,omitempty"`
}

var Cfg Config

func Run() {
	if err := proc.StartSync("wireguard.exe", "/installtunnelservice", Cfg.WireguardConfPath); err != nil {
		log.Println("failed to run vpn, error: ", err)
		return
	}
	log.Println("vpn running")
}

func Stop() {
	if err := proc.StartSync("wireguard.exe", "/uninstalltunnelservice", Cfg.TunelName); err != nil {
		log.Println("failed to stop vpn, error: ", err)
		return
	}
	log.Println("vpn stoped")
}

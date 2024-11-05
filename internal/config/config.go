package config

import "dpiWrapper/internal/vpn"

type Config struct {
	DPIYoutubePath     string     `json:"dpi_youtube_path,omitempty"`
	DpiYoutubeCommands []string   `json:"dpi_youtube_commands,omitempty"`
	RussiaBlackListDPI string     `json:"russia_black_list_dpi,omitempty"`
	VPN                vpn.Config `json:"vpn,omitempty"`
}

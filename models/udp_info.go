package models

// it's Opcode 2 Ready payload but called it that is quite bad lol
type UDPInfo struct {
	Ssrc  uint32   `json:"ssrc" mapstructure:"ssrc"`
	IP    string   `json:"ip" mapstructure:"ip"`
	Port  uint16   `json:"port" mapstructure:"port"`
	Modes []string `json:"modes" mapstructure:"modes"`
}

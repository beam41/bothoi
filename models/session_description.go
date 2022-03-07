package models

type SessionDescription struct {
	Mode      string `json:"mode" mapstructure:"mode"`
	SecretKey [32]byte `json:"secret_key" mapstructure:"secret_key"`
}

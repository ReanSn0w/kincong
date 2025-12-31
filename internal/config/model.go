package config

import "github.com/ReanSn0w/kincong/internal/utils"

const (
	RequestTypeBat       RequestType = "bat"
	RequestTypeWireguard RequestType = "wireguard"
	RequestTypeAmnezia   RequestType = "amnezia"
)

type RequestType string

type Request struct {
	Type   string `json:"type"`
	Rules  []Rule `json:"rules"`
	Config string `json:"config"`
}

type Rule struct {
	Value utils.Value `json:"value"`
	Type  string      `json:"type"`
}

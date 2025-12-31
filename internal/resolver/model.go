package resolver

import "github.com/ReanSn0w/kincong/internal/utils"

const (
	resolverTypeEmpty ResolverType = iota
	ResolverTypeIP
	ResolverTypeDNS
	ResolverTypeASN
)

type ResolverType int

type ResolverItem interface {
	Type() ResolverType
	Resolve(string) ([]utils.ResolvedSubnet, error)
}

type IPData struct {
	Asns   []string `json:"asns"`
	Prefix string   `json:"prefix"`
}

package ip

import (
	"errors"
	"net"
	"strings"

	"github.com/ReanSn0w/kincong/internal/resolver"
	"github.com/ReanSn0w/kincong/internal/utils"
)

func New() *Resolver {
	return &Resolver{}
}

type Resolver struct{}

func (r *Resolver) Type() resolver.ResolverType {
	return resolver.ResolverTypeIP
}

func (r *Resolver) Resolve(value string) ([]utils.ResolvedSubnet, error) {
	parts := strings.Split(value, "/")
	switch len(parts) {
	case 1:
		parts = append(parts, "32")
	case 2:
	default:
		return nil, errors.New("invalid IP address")
	}

	parsedIP := net.ParseIP(parts[0])
	if parsedIP == nil {
		return nil, errors.New("invalid IP address")
	}

	_, _, err := net.ParseCIDR(parts[0] + "/" + parts[1])
	if err != nil {
		return nil, err
	}

	return []utils.ResolvedSubnet{utils.ResolvedSubnet(strings.Join(parts, "/"))}, nil
}

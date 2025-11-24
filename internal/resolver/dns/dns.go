package dns

import (
	"context"
	"errors"
	"net"
	"time"

	"math/rand/v2"

	"github.com/ReanSn0w/kincong/internal/resolver"
	"github.com/ReanSn0w/kincong/internal/utils"
)

func New(upstream ...string) *Resolver {
	return &Resolver{
		upstream: upstream,
	}
}

type Resolver struct {
	upstream []string
}

func (r *Resolver) Type() resolver.ResolverType {
	return resolver.ResolverTypeDNS
}

func (r *Resolver) Resolve(domain string) ([]utils.ResolvedSubnet, error) {
	// Select a random DNS server from upstream
	randomDNS := r.upstream[rand.IntN(len(r.upstream))]

	// Resolve the domain to an IP address using the selected DNS server
	resolver := net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Second * 5,
			}
			return d.DialContext(ctx, network, randomDNS+":53")
		},
	}

	ipRecords, err := resolver.LookupIP(context.Background(), "ip4", domain)
	if err != nil {
		return nil, err
	}

	if len(ipRecords) == 0 {
		return nil, errors.New("no A records found")
	}

	result := make([]utils.ResolvedSubnet, len(ipRecords))
	for i, ip := range ipRecords {
		result[i] = utils.ResolvedSubnet(ip.String() + "/32")
	}

	return result, nil
}

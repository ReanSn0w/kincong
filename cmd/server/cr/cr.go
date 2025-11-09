package cr

import (
	"context"
	"fmt"
	"strings"

	"github.com/ReanSn0w/kincong/internal/resolver"
	"github.com/ReanSn0w/kincong/internal/resolver/asn"
	"github.com/ReanSn0w/kincong/internal/resolver/dns"
	"github.com/ReanSn0w/kincong/internal/resolver/ip"
	"github.com/ReanSn0w/kincong/internal/utils"
)

func New(ctx context.Context) *CR {
	asnResolver := asn.New()
	domainResolver := dns.New("1.1.1.1")

	return &CR{
		cache:      newCache(ctx),
		domainInfo: domainResolver,
		ipInfo:     asnResolver,
		resolver: resolver.NewResolver(
			ip.New(),
			domainResolver,
			asnResolver,
		),
	}
}

type CR struct {
	cache      *Cache
	resolver   *resolver.Resolver
	ipInfo     *asn.Resolver
	domainInfo *dns.Resolver
}

type DomainInfo struct {
	IPs      []utils.Value `json:"ips"`
	Networks []utils.Value `json:"networks"`
	ASNs     []utils.Value `json:"asns"`
}

func (c *CR) Resolve(values []utils.Value) ([]utils.ResolvedSubnet, error) {
	var (
		result    = make([]utils.ResolvedSubnet, 0)
		newValues = make([]utils.Value, 0)
	)

	// Получение данных из кеша
	{
		for _, value := range values {
			subnets, ok := c.cache.GetSubnet(value)
			if ok {
				result = append(result, subnets...)
			} else {
				newValues = append(newValues, value)
			}
		}
	}

	// Резолвинг новых значений
	{
		valuesResolvedMap, err := c.resolver.ResolveByValue(newValues...)
		if err != nil {
			return nil, err
		}

		for value, subnets := range valuesResolvedMap {
			result = append(result, subnets...)
			c.cache.SetSubnet(value, subnets)
		}
	}

	return result, nil
}

func (c *CR) DomainInfo(domain utils.Value) (*DomainInfo, error) {
	if !domain.IsDomain() {
		return nil, fmt.Errorf("invalid domain")
	}

	if data, ok := c.cache.GetDomainInfo(domain); ok {
		return data, nil
	}

	var (
		ipMap      = make(map[utils.Value]struct{})
		networkMap = make(map[utils.Value]struct{})
		asnMap     = make(map[utils.Value]struct{})
	)

	{
		domainData, err := c.domainInfo.Resolve(string(domain))
		if err != nil {
			return nil, err
		}

		for _, ip := range domainData {
			ipValue := strings.TrimSuffix(string(ip), "/32")
			ipMap[utils.Value(ipValue)] = struct{}{}
		}

		for ip, _ := range ipMap {
			data, err := c.ipInfo.InfoByIP(string(ip))
			if err != nil {
				return nil, err
			}

			networkMap[utils.Value(data.Prefix)] = struct{}{}

			for _, asn := range data.Asns {
				asnMap[utils.Value("AS"+asn)] = struct{}{}
			}
		}
	}

	result := DomainInfo{
		IPs:      make([]utils.Value, 0),
		Networks: make([]utils.Value, 0),
		ASNs:     make([]utils.Value, 0),
	}

	{
		for ip := range ipMap {
			result.IPs = append(result.IPs, ip)
		}

		for network := range networkMap {
			result.Networks = append(result.Networks, network)
		}

		for asn := range asnMap {
			result.ASNs = append(result.ASNs, asn)
		}
	}

	c.cache.SetDomainInfo(domain, result)
	return &result, nil
}

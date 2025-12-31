package resolver

import (
	"context"
	"errors"
	"sort"
	"time"

	"github.com/ReanSn0w/kincong/internal/utils"
)

var (
	ErrResolverNotFound = errors.New("resolver not found")
)

func NewResolver(ctx context.Context, items ...ResolverItem) *Resolver {
	resolvers := make(map[ResolverType]ResolverItem)
	for _, item := range items {
		resolvers[item.Type()] = item
	}
	return &Resolver{
		resolvers:    resolvers,
		cacheSubnets: utils.NewCache[[]utils.ResolvedSubnet](ctx, time.Hour*12),
		cacheInfo:    utils.NewCache[InfoResult](ctx, time.Hour*24),
		cacheASNInfo: utils.NewCache[IPData](ctx, time.Hour*24),
	}
}

type Resolver struct {
	resolvers    map[ResolverType]ResolverItem
	cacheSubnets *utils.Cache[[]utils.ResolvedSubnet]
	cacheInfo    *utils.Cache[InfoResult]
	cacheASNInfo *utils.Cache[IPData]
}

// ResolveMany преобразует массив items в карту массивов подсетей полученных из item по item
//
// items - набор значений ip, доменов и asn
//
// в ответе карта посетей
// map[<item>][](<ip сеть>)
// и карта ошибок
// map[<item>]error
//
// при этом выполнение не прерывается при получении ошибки
// решение о том, как поступать с данными следует принимать
// после выполнения
func (r *Resolver) ResolveMany(items ...utils.Value) (map[utils.Value]utils.ResolvedSubnet, error) {
	var (
		result = make(map[utils.Value]utils.ResolvedSubnet)
		errMap = make(utils.ErrorsMap)
	)

	for _, item := range items {
		resolved, err := r.Resolve(item)
		if err != nil {
			errMap[string(item)] = err
		} else {
			for _, subnet := range resolved {
				result[item] = subnet
			}
		}
	}

	return result, errMap.HasError()
}

// Resolve преобразует item в массив подсетей полученных из item
//
// items - набор значений ip, доменов и asn
//
// фунция возвращает массив подсетей или ошибку
func (r *Resolver) Resolve(item utils.Value) ([]utils.ResolvedSubnet, error) {
	t := r.detectValueType(item)
	if t == resolverTypeEmpty {
		return nil, ErrResolverNotFound
	}

	value, err := r.cacheSubnets.Must(string(item), func() (*[]utils.ResolvedSubnet, error) {
		resolver, ok := r.resolvers[t]
		if !ok {
			return nil, ErrResolverNotFound
		}

		subnets, err := resolver.Resolve(string(item))
		if err != nil {
			return nil, err
		}

		return &subnets, err
	})

	if err != nil {
		return nil, err
	}

	return *value, nil
}

type InfoResult struct {
	ASN      []utils.Value `json:"asn,omitempty"`
	Domains  []utils.Value `json:"domains,omitempty"`
	Networks []utils.Value `json:"networks,omitempty"`
}

// InfoByValue извлекает данные о домене, ANS и связанных сетях по value
func (r *Resolver) InfoByValue(item utils.Value) (*InfoResult, error) {
	return r.cacheInfo.Must(string(item), func() (*InfoResult, error) {
		var (
			t   = r.detectValueType(item)
			res = &InfoResult{}

			networksMap = map[utils.Value]bool{}
			domainsMap  = map[utils.Value]bool{}
			asnsMap     = map[utils.Value]bool{}
		)

		switch t {
		case ResolverTypeIP:
			if ipResolver, ok := r.resolvers[ResolverTypeIP]; ok {
				subnets, _ := ipResolver.Resolve(string(item))
				for _, subnet := range subnets {
					networksMap[utils.Value(subnet)] = true
				}

				for _, subnet := range subnets {
					ip, ok := subnet.IP()
					if !ok {
						continue
					}

					info, err := r.resolveASNInfo(ip)
					if err != nil {
						continue
					}

					for _, asn := range info.Asns {
						asnsMap[utils.Value(asn)] = true
					}
				}
			}

			if asnResolver, ok := r.resolvers[ResolverTypeASN]; ok {
				for asn := range asnsMap {
					items, err := asnResolver.Resolve(string(asn))
					if err != nil {
						continue
					}

					for _, item := range items {
						networksMap[utils.Value(item)] = true
					}
				}
			}
		case ResolverTypeDNS:
			domainsMap[item] = true

			if dnsResolver, ok := r.resolvers[ResolverTypeDNS]; ok {
				subnets, err := dnsResolver.Resolve(string(item))
				if err != nil {
					return nil, err
				}

				for _, subnet := range subnets {
					ip, ok := subnet.IP()
					if !ok {
						continue
					}

					info, err := r.resolveASNInfo(ip)
					if err != nil {
						continue
					}

					for _, asn := range info.Asns {
						asnsMap[utils.Value(asn)] = true
					}
				}
			}

			if asnResolver, ok := r.resolvers[ResolverTypeASN]; ok {
				for asn := range asnsMap {
					items, err := asnResolver.Resolve(string(asn))
					if err != nil {
						continue
					}

					for _, item := range items {
						networksMap[utils.Value(item)] = true
					}
				}
			}
		case ResolverTypeASN:
			asnsMap[item] = true

			if asnResolver, ok := r.resolvers[ResolverTypeASN]; ok {
				for asn := range asnsMap {
					items, err := asnResolver.Resolve(string(asn))
					if err != nil {
						continue
					}

					for _, item := range items {
						networksMap[utils.Value(item)] = true
					}
				}
			}
		default:
			return nil, ErrResolverNotFound
		}

		{
			for asn := range asnsMap {
				res.ASN = append(res.ASN, asn)
			}

			for domain := range domainsMap {
				res.Domains = append(res.Domains, domain)
			}

			for network := range networksMap {
				res.Networks = append(res.Networks, network)
			}

			sort.Slice(res.ASN, func(i, j int) bool {
				return res.ASN[i] < res.ASN[j]
			})

			sort.Slice(res.Domains, func(i, j int) bool {
				return res.Domains[i] < res.Domains[j]
			})

			sort.Slice(res.Networks, func(i, j int) bool {
				return res.Networks[i] < res.Networks[j]
			})
		}

		return res, nil
	})
}

type ipInfoResolver interface {
	InfoByIP(ip string) (*IPData, error)
}

func (r *Resolver) resolveASNInfo(ip string) (*IPData, error) {
	return r.cacheASNInfo.Must(ip, func() (*IPData, error) {
		_, ok := r.resolvers[ResolverTypeASN]
		if !ok {
			return nil, errors.New("asn resolver unavaliable")
		}

		resolver := any(r.resolvers[ResolverTypeASN]).(ipInfoResolver)
		return resolver.InfoByIP(ip)
	})
}

func (r *Resolver) detectValueType(item utils.Value) ResolverType {
	if item.IsASN() {
		return ResolverTypeASN
	}
	if item.IsDomain() {
		return ResolverTypeDNS
	}
	if item.IsNetwork() {
		return ResolverTypeIP
	}
	if item.IsIP() {
		return ResolverTypeIP
	}
	return resolverTypeEmpty
}

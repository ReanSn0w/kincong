package resolver

import (
	"errors"

	"github.com/ReanSn0w/kincong/internal/utils"
)

var (
	ErrResolverNotFound = errors.New("resolver not found")
)

func NewResolver(items ...ResolverItem) *Resolver {
	resolvers := make(map[ResolverType]ResolverItem)
	for _, item := range items {
		resolvers[item.Type()] = item
	}
	return &Resolver{resolvers: resolvers}
}

type Resolver struct {
	resolvers map[ResolverType]ResolverItem
}

func (r *Resolver) Resolve(items ...utils.Value) (map[string]utils.ResolvedSubnet, error) {
	var (
		result = make(map[string]utils.ResolvedSubnet)
		errMap = make(utils.ErrorsMap)
	)

	for _, item := range items {
		resolved, err := r.resolveItem(item)
		if err != nil {
			errMap[string(item)] = err
		} else {
			for _, subnet := range resolved {
				result[subnet.Hash()] = subnet
			}
		}
	}

	return result, errMap.HasError()
}

func (r *Resolver) ResolveByValue(items ...utils.Value) (map[utils.Value][]utils.ResolvedSubnet, error) {
	var (
		result = make(map[utils.Value][]utils.ResolvedSubnet)
		errMap = make(utils.ErrorsMap)
	)

	for _, item := range items {
		resolved, err := r.resolveItem(item)
		if err != nil {
			errMap[string(item)] = err
		} else {
			result[item] = resolved
		}
	}

	return result, errMap.HasError()
}

func (r *Resolver) resolveItem(item utils.Value) ([]utils.ResolvedSubnet, error) {
	t := r.detectValueType(item)
	if t == resolverTypeEmpty {
		return nil, nil
	}

	resolver, ok := r.resolvers[t]
	if !ok {
		return nil, ErrResolverNotFound
	}
	return resolver.Resolve(string(item))
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

package asn

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ReanSn0w/kincong/internal/resolver"
	"github.com/ReanSn0w/kincong/internal/utils"
)

type (
	Resolver struct{}

	resolveResponse struct {
		Version        string `json:"version"`
		DataCallName   string `json:"data_call_name"`
		DataCallStatus string `json:"data_call_status"`
		Cached         bool   `json:"cached"`
		Data           data   `json:"data"`
		QueryID        string `json:"query_id"`
		ProcessTime    int64  `json:"process_time"`
		ServerID       string `json:"server_id"`
		BuildVersion   string `json:"build_version"`
		Status         string `json:"status"`
		StatusCode     int64  `json:"status_code"`
		Time           string `json:"time"`
	}

	data struct {
		Prefixes []prefix `json:"prefixes"`
	}

	prefix struct {
		InBGP   bool   `json:"in_bgp"`
		InWhois bool   `json:"in_whois"`
		Prefix  string `json:"prefix"`
	}

	ipInfoResponse struct {
		Data       IPData `json:"data"`
		Status     string `json:"status"`
		StatusCode int64  `json:"status_code"`
	}

	IPData struct {
		Asns   []string `json:"asns"`
		Prefix string   `json:"prefix"`
	}
)

func New() *Resolver {
	return &Resolver{}
}

func (r *Resolver) Type() resolver.ResolverType {
	return resolver.ResolverTypeASN
}

func (r *Resolver) InfoByIP(ip string) (*IPData, error) {
	var response ipInfoResponse
	err := r.request(fmt.Sprintf("data/network-info/data.json?resource=%s", ip), &response)
	if err != nil {
		return nil, err
	}

	return &response.Data, nil
}

func (r *Resolver) Resolve(value string) ([]utils.ResolvedSubnet, error) {
	var responseResult resolveResponse
	err := r.request(fmt.Sprintf("data/as-routing-consistency/data.json?resource=%s", value), &responseResult)
	if err != nil {
		return nil, err
	}

	var result = make([]utils.ResolvedSubnet, 0, len(responseResult.Data.Prefixes))
	for _, prefix := range responseResult.Data.Prefixes {
		if !prefix.InBGP {
			continue
		}

		if !prefix.InWhois {
			continue
		}

		result = append(result, utils.ResolvedSubnet(prefix.Prefix))
	}

	return result, nil
}

func (r *Resolver) request(path string, respBody any) error {
	resp, err := http.Get(fmt.Sprintf("https://stat.ripe.net/%s", path))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(respBody)
	if err != nil {
		return err
	}

	return nil
}

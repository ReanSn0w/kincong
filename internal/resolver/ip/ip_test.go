package ip_test

import (
	"testing"

	"github.com/ReanSn0w/kincong/internal/resolver/ip"
	"github.com/ReanSn0w/kincong/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestIP_Resolve(t *testing.T) {
	cases := []struct {
		Name   string
		Input  string
		Output []utils.ResolvedSubnet
	}{
		{
			Name:   "Valid IP",
			Input:  "192.168.1.1",
			Output: []utils.ResolvedSubnet{"192.168.1.1/32"},
		},
		{
			Name:   "Valid Subnet",
			Input:  "192.168.0.0/16",
			Output: []utils.ResolvedSubnet{"192.168.0.0/16"},
		},
		{
			Name:   "Empty Input",
			Input:  "",
			Output: nil,
		},
	}

	resolver := ip.New()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			output, _ := resolver.Resolve(tc.Input)
			assert.Equal(t, tc.Output, output)
		})
	}
}

package utils_test

import (
	"testing"

	"github.com/ReanSn0w/kincong/internal/utils"
	"github.com/stretchr/testify/assert"
)

func Test_Value(t *testing.T) {
	cases := []struct {
		Name      string
		Value     string
		IsIP      bool
		IsNetwork bool
		IsDomain  bool
		IsASN     bool
	}{
		{
			Name:  "Empty Value",
			Value: "",
		},
		{
			Name:  "Invalid Value",
			Value: "invalid",
		},
		{
			Name:     "Normal Domain 1",
			Value:    "animevost.org",
			IsDomain: true,
		},
		{
			Name:     "Normal Domain 2",
			Value:    "rutracker.org",
			IsDomain: true,
		},
		{
			Name:     "Normal Domain 3",
			Value:    "coder.papkovda.ru",
			IsDomain: true,
		},
		{
			Name:     "Normal Domain 4",
			Value:    "site.dev.aaabbbcc.ru",
			IsDomain: true,
		},
		{
			Name:  "Normal ASN 1",
			Value: "AS15169",
			IsASN: true,
		},
		{
			Name:  "Normal ASN 2",
			Value: "AS15169",
			IsASN: true,
		},
		{
			Name:  "Invalid ASN 1",
			Value: "ASAABBCC",
		},
		{
			Name:  "Invalid ASN 2",
			Value: "AS",
		},
		{
			Name:      "Normal Network 1",
			Value:     "103.31.4.0/22",
			IsNetwork: true,
		},
		{
			Name:      "Normal Network 2",
			Value:     "104.16.0.0/13",
			IsNetwork: true,
		},
		{
			Name:      "Normal Network 3",
			Value:     "104.24.0.0/14",
			IsNetwork: true,
		},
		{
			Name:  "Invalid Network",
			Value: "256.24.0.0/14",
		},
		{
			Name:  "Invalid Network",
			Value: "104.24.0.0/33",
		},
		{
			Name:  "Normal IP 1",
			Value: "104.24.0.0",
			IsIP:  true,
		},
		{
			Name:  "Normal IP 2",
			Value: "104.24.0.1",
			IsIP:  true,
		},
		{
			Name:  "Normal IP 3",
			Value: "104.24.0.2",
			IsIP:  true,
		},
		{
			Name:  "Invalid IP 1",
			Value: "104.24.0.256",
		},
		{
			Name:  "Invalid IP 2",
			Value: "255.255.255.256",
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			value := utils.Value(tc.Value)
			assert.Equal(t, tc.IsIP, value.IsIP(), "ip check")
			assert.Equal(t, tc.IsNetwork, value.IsNetwork(), "network check")
			assert.Equal(t, tc.IsDomain, value.IsDomain(), "domain check")
			assert.Equal(t, tc.IsASN, value.IsASN(), "asn check")
		})
	}
}

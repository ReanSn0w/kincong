package server

import (
	"github.com/ReanSn0w/gokit/pkg/web"
	"github.com/ReanSn0w/kincong/internal/utils"
)

// ValueSearchRequest - Запрос на поиск значений
type ValueSearchRequest struct {
	// Value - Значение для поиска
	//
	// Домен сайта, например "example.com"
	Value utils.Value `json:"value" example:"example.com"`
}

// ValueSearchResponse - Ответ на запрос поиска значений
type ValueSearchResponse web.Response[[]ValueSearchItem]

// ValueSearchItem - Элемент результата поиска значений
type ValueSearchItem struct {
	// From - Источник данных, например "example.com"
	From string `json:"from" example:"example.com"`

	// Value - Значение, например
	//
	// Для типа:
	// asn - AS12345
	// cidr - 192.168.0.0/24
	// ip - 192.168.0.1
	Value string `json:"value" example:"AS12345"`

	Type ValueSearchItemType `json:"type" example:"asn"`
}

const (
	ASN  ValueSearchItemType = "asn"
	CIDR ValueSearchItemType = "cidr"
	IP   ValueSearchItemType = "ip"
)

// ValueSearchItemType - Тип для value
type ValueSearchItemType string

package utils

import (
	"crypto/md5"
	"fmt"
	"net"
	"regexp"
	"strings"
)

var (
	domainPattern = regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*\.[a-zA-Z]+$`)
	asnPattern    = regexp.MustCompile(`^AS\d+(\.\d+)?$`)
)

// Value - Значение введенное пользоватем
type Value string

// IsDomain - Проверка на то, что значение является сайтом
func (v Value) IsDomain() bool {
	return domainPattern.MatchString(string(v)) && strings.Contains(string(v), ".")
}

// IsASN - Проверка на то, что значение является ASN
func (v Value) IsASN() bool {
	return asnPattern.MatchString(string(v))
}

// IsIP - Проверка на то, что значение является IP
func (v Value) IsIP() bool {
	return net.ParseIP(string(v)) != nil
}

// IsNetwork - Проверка на то, что значение является сетью
func (v Value) IsNetwork() bool {
	_, _, err := net.ParseCIDR(string(v))
	return err == nil
}

// ResolvedSubnet - Значение добавляемое в конфигурацию роутера
type ResolvedSubnet string

// IP возвращает строку с IP адресом из ResolvedSubnet
func (r ResolvedSubnet) IP() (string, bool) {
	ipValue, _, err := net.ParseCIDR(string(r))
	if err != nil {
		return "", false
	}

	return ipValue.String(), true
}

// Mask возвращает маску в формате IP адреса
func (r ResolvedSubnet) Mask() (string, bool) {
	_, ipnet, err := net.ParseCIDR(string(r))
	if err != nil {
		return "", false
	}

	mask := net.IP(ipnet.Mask).String()
	return mask, true
}

// Hash - Возвращает хеш значения
//
// Планируется использование его в качестве части описания маршрута на роутере
// для определения маршрутов, добавленных с помощью конфигурации, на случае если
// добавленные ранее маршруты придется удалить
func (r ResolvedSubnet) Hash() string {
	return fmt.Sprintf("%x", md5.Sum([]byte(r)))
}

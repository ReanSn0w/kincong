package main

import (
	"html/template"
	"net/url"
	"time"

	"github.com/ReanSn0w/gokit/pkg/app"
	"github.com/ReanSn0w/kincong/internal/resolver"
	"github.com/ReanSn0w/kincong/internal/resolver/asn"
	"github.com/ReanSn0w/kincong/internal/resolver/dns"
	"github.com/ReanSn0w/kincong/internal/resolver/ip"
	"github.com/ReanSn0w/kincong/internal/server"
)

var (
	revision = "debug"
	opts     = struct {
		app.Debug

		Port    int    `long:"port" env:"PORT" default:"8080" description:"Port to listen on"`
		BaseURL string `long:"base-url" env:"BASE_URL" default:"http://localhost:8080" description:"Base URL for the server"`

		DNS []string `long:"dns" env:"DNS" env-delim:"," default:"1.1.1.1" description:"dns resolver"`
	}{}
)

func main() {
	app := app.New("Network Data Resolver", revision, &opts)

	baseURL, err := url.Parse(opts.BaseURL)
	if err != nil {
		panic(err.Error())
	}

	tmpl, err := template.ParseFiles("static/app.html")
	if err != nil {
		panic(err.Error())
	}

	server := server.New(revision, baseURL, resolver.NewResolver(
		app.Context(),
		ip.New(),
		dns.New(opts.DNS...),
		asn.New(),
	), tmpl)
	app.Add(server.Stop)
	server.Start(app.CancelCause(), opts.Port)

	app.GracefulShutdown(time.Second * 10)
}

package server

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"sort"
	"time"

	"github.com/ReanSn0w/gokit/pkg/web"
	"github.com/ReanSn0w/kincong/internal/resolver"
	"github.com/ReanSn0w/kincong/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-pkgz/lgr"
)

func New(revision string, baseURL *url.URL, resolver *resolver.Resolver, tmpl *template.Template) *Server {
	return &Server{
		revision: revision,
		tmpl:     tmpl,
		server:   web.New(lgr.Default()),
		api:      newAPI(revision, baseURL, resolver),
	}
}

type Server struct {
	revision string
	tmpl     *template.Template
	server   *web.Server
	api      *API
}

type ResponseError web.Response[string]

func (s *Server) Start(cancel context.CancelCauseFunc, port int) {
	s.server.Run(cancel, port, s.handler())
}

func (s *Server) Stop(ctx context.Context) {
	s.server.Shutdown(ctx)
}

func (s *Server) handler() http.Handler {
	r := chi.NewRouter()

	{
		r.Handle("/docs/*", utils.SwaggerHandler(s.api.baseURL.JoinPath("/v1/api"), s.revision))
		r.Handle("/static/*", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))
		r.Mount("/v1", s.api.handler())
	}

	r.Get("/", s.mainPageHandler)

	r.Route("/htmx", func(r chi.Router) {
		r.Post("/update-query", s.htmxUpdateQueryHandler)
		r.Post("/add-config-item", s.htmxAddConfigItemHandler)
		r.Post("/remove-config-item", s.htmxRemoveConfigItemHandler)
		r.Post("/generate-config", s.htmxGenerateConfigHandler)
	})

	return r
}

func (s *Server) mainPageHandler(w http.ResponseWriter, r *http.Request) {
	s.tmpl.ExecuteTemplate(w, "app", nil)
}

type Value struct {
	Input  string `json:"input"`
	Type   string `json:"type"`
	Output string `json:"output"`
}

func (s *Server) htmxUpdateQueryHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	query := r.FormValue("query")

	result, err := s.api.resolver.InfoByValue(utils.Value(query))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get result: %s", err), http.StatusInternalServerError)
		return
	}

	values := []Value{}

	for _, value := range result.ASN {
		values = append(values, Value{
			Input:  query,
			Type:   "asn",
			Output: "AS" + string(value),
		})
	}

	for _, value := range result.Networks {
		values = append(values, Value{
			Input:  query,
			Type:   "network",
			Output: string(value),
		})
	}

	s.tmpl.ExecuteTemplate(w, "add-item-array", values)
}

func (s *Server) htmxAddConfigItemHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	input := r.FormValue("input")
	t := r.FormValue("type")
	output := r.FormValue("output")

	s.tmpl.ExecuteTemplate(w, "delete-item", Value{
		Input:  input,
		Type:   t,
		Output: output,
	})
}

func (s *Server) htmxRemoveConfigItemHandler(w http.ResponseWriter, r *http.Request) {
}

func (s *Server) htmxGenerateConfigHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	values := []utils.Value{}

	for _, ip := range r.Form["ip"] {
		values = append(values, utils.Value(ip))
	}

	for _, network := range r.Form["network"] {
		values = append(values, utils.Value(network))
	}

	for _, asn := range r.Form["asn"] {
		values = append(values, utils.Value(asn))
	}

	var subnets []utils.ResolvedSubnet

	for _, value := range values {
		ss, err := s.api.resolver.Resolve(value)
		if err != nil {
			lgr.Default().Logf("[ERROR] resolve subnets err: %s", err)
			continue
		}

		subnets = append(subnets, ss...)
	}

	sort.Slice(subnets, func(i, j int) bool {
		return subnets[i] > subnets[j]
	})

	switch r.FormValue("route_option") {
	case "bat":
		buffer := new(bytes.Buffer)

		for _, value := range subnets {
			ip, ok := value.IP()
			if !ok {
				continue
			}

			mask, ok := value.Mask()
			if !ok {
				continue
			}

			fmt.Fprintf(buffer, "route ADD %s MASK %s\n", ip, mask)
		}

		w.Header().Set("Content-Disposition", `attachment; filename="result.bat"`)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		http.ServeContent(w, r, "result.bat", time.Now(), bytes.NewReader(buffer.Bytes()))
	default:
		buffer := new(bytes.Buffer)

		for _, value := range subnets {
			fmt.Fprintf(buffer, "%s\n", value)
		}

		w.Header().Set("Content-Disposition", `attachment; filename="result.txt"`)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		http.ServeContent(w, r, "result.txt", time.Now(), bytes.NewReader(buffer.Bytes()))
	}
}

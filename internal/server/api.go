package server

import (
	"net/http"
	"net/url"

	"github.com/ReanSn0w/gokit/pkg/web"
	"github.com/ReanSn0w/gokit/pkg/web/mv/json"
	"github.com/ReanSn0w/kincong/internal/config"
	"github.com/ReanSn0w/kincong/internal/resolver"
	"github.com/ReanSn0w/kincong/internal/utils"
	"github.com/go-chi/chi/v5"
)

func newAPI(revision string, baseURL *url.URL, resolver *resolver.Resolver) *API {
	return &API{
		revision: revision,
		resolver: resolver,
		baseURL:  baseURL,
	}
}

// @title       Kincong Resolver Rest Server
// @version	    debug
// @description	API для получения подсетей на основе IP адресов, доменов и ASN
// @host        localhost
// @schemes     http
// @BasePath    /
type API struct {
	revision string
	baseURL  *url.URL
	resolver *resolver.Resolver
}

func (a *API) handler() http.Handler {
	r := chi.NewRouter()

	{
		r.MethodNotAllowed(web.JSON_MethodNotAllowedHandlerFunc)
		r.NotFound(web.JSON_NotFoundHandlerFunc)
	}

	r.Route("/api", func(r chi.Router) {
		r.
			With(json.Decoder[ResolveRequest]).
			Post("/resolve", a.resolveHandler)

		r.
			With(json.Decoder[ValueRequest]).
			Post("/value", a.valueHandler)

		r.
			With(json.Decoder[config.Request]).
			Post("/upload", a.modifyConfigHandler)
	})

	return r
}

// MARK: - Value Resolver

type ValueRequest struct {
	Value utils.Value `json:"value" example:"google.com"`
}

type ValueResponse web.Response[ValueResponseData]

type ValueResponseData resolver.InfoResult

// @Summary		Получение информации о значении
// @Description	Метод возвращает массивы IP, Подсетей и ASN записей
// @Description связанных с указанным значением
// @Accept		json
// @Produce		json
// @Param       body body ValueRequest true "Параметры запроса"
// @Success		200	    {object} ValueResponse
// @Failure		default	{object} ResponseError
// @Router		/value [post]
func (a *API) valueHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := json.Get[ValueRequest](ctx)

	info, err := a.resolver.InfoByValue(req.Value)
	if err != nil {
		web.NewResponse(err).Write(http.StatusInternalServerError, w)
		return
	}

	web.NewResponse(info).Write(http.StatusOK, w)
}

// MARK: - Networks resolver

type ResolveRequest struct {
	// Список IP адресов для разрешения
	IPs []string `json:"ip,omitempty" example:"8.8.8.8"`
	// Список доменных имен для разрешения
	Domains []string `json:"domain,omitempty" example:"google.com"`
	// Список ASN для разрешения
	ASNs []string `json:"asn,omitempty" example:"AS12345"`
}

func (r *ResolveRequest) Values() []utils.Value {
	result := make([]utils.Value, 0, len(r.IPs)+len(r.Domains)+len(r.ASNs))
	for _, ip := range r.IPs {
		result = append(result, utils.Value(ip))
	}
	for _, domain := range r.Domains {
		result = append(result, utils.Value(domain))
	}
	for _, asn := range r.ASNs {
		result = append(result, utils.Value(asn))
	}
	return result
}

type ResolveResponse web.Response[ResolveResponseData]

type ResolveResponseData map[utils.Value][]utils.ResolvedSubnet

// @Summary		Получение списка подсетей
// @Description	Метод возвращает пользователю список доступных в системе документов
// @Accept		json
// @Produce		json
// @Param       body body ResolveRequest true "Параметры запроса"
// @Success		200	    {object} ResolveResponse
// @Failure		default	{object} ResponseError
// @Router		/resolve [post]
func (a *API) resolveHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := json.Get[ResolveRequest](ctx)

	networks, err := a.resolver.ResolveMany(req.Values()...)
	if err != nil {
		web.NewResponse(err).Write(http.StatusInternalServerError, w)
		return
	}

	web.NewResponse(networks).Write(http.StatusOK, w)
}

// MARK: - Modify Config

// @Summary		Модификация конфигурации
// @Description	Метод принимает текущую конфигурацию и модифицирует
// @Description ее на основе переданных правил
// @Accept		json
// @Produce		application/octet-stream
// @Param       body body config.Request true "Параметры запроса"
// @Success		200	    {object} []byte
// @Failure		default	{object} ResponseError
// @Router		/modify-config [post]
func (a *API) modifyConfigHandler(w http.ResponseWriter, r *http.Request) {

}

package router

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

const (
	SecretHeaderName = "X-App-Secret"
)

type ProducerRouter interface {
	Server() *http.Server
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

type producerRouter struct {
	addr        string
	marketUrl   string
	profilePath string

	client *http.Client
}

func NewProducerRouter(addr, marketUrl, profilePath string) ProducerRouter {
	return &producerRouter{
		addr:        addr,
		marketUrl:   marketUrl,
		profilePath: profilePath,
		client:      createClient(),
	}
}

func (p *producerRouter) Server() *http.Server {
	return nil
}

func (p *producerRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	secret := req.Header.Get(SecretHeaderName)
	if secret == "" {
		onNoSecret(w)
		return
	}

	profile, err := fetchProfileBySecretKey(p.client, p.marketUrl, p.profilePath, secret)
	if err != nil {
		onTradeError(w, err)
		return
	}

	if profile == nil {
		onNoProfile(w)
		return
	}

	srvName, _ := profile.MetaServiceName(KeyProducer)
	srv, err := consulService(profile.ConsulHost, srvName)
	if err != nil {
		onConsulError(w, err)
		return
	}

	url := &url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%d", srv.ServiceAddress, srv.ServicePort),
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ServeHTTP(w, req)
}

package gateway

import (
	"context"
	"fmt"
	"github.com/csh0820/gateway/config"
	"github.com/csh0820/gateway/pkg/etcd"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

	"github.com/csh0820/gateway/pkg/registry"

	"github.com/gin-gonic/gin"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type GatewayHandler struct {
	registry etcd.Registry
	// registry *discovery.EtcdRegistry
	// routes   map[string]*Route
	// proxies  map[string]*httputil.ReverseProxy

	client    *clientv3.Client
	Instances map[string]*registry.ServiceInstance
	prefix    string
	mu        sync.RWMutex

	ctx context.Context
}

func NewGatewayHandler(registry etcd.Registry) *GatewayHandler {
	ctx := context.Background()
	registry.GetService(ctx, config.GetConfig().GatewayAddress)
	watch, err := registry.Watch(ctx, config.GetConfig().GatewayAddress)
	if err != nil {
		log.Fatal(err)
	}

	return &GatewayHandler{
		registry: registry,
	}
}

func HandleRequest(c *gin.Context) {
	path := c.Request.URL.Path

	targetURL := fmt.Sprintf("http://%s:%d", "localhost", 9090)
	target, err := url.Parse(targetURL)
	if err != nil {
		log.Println(err)
	}

	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			// if c.Request.URL.Path == "" {
			// 	c.Request.URL.Path = "/"
			// }

			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host
			req.URL.Path = path
			req.Host = target.Host
		},
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}

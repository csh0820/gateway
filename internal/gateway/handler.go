package gateway

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

	"github.com/gin-gonic/gin"
)

type GatewayHandler struct {
	mu sync.RWMutex
}

func NewGatewayHandler() *GatewayHandler {
	return &GatewayHandler{}
}

func (h *GatewayHandler) HandleRequest(c *gin.Context) {
	path := c.Request.URL.Path

	targetURL := fmt.Sprintf("http://%s:%d", "localhost", 9090)
	target, err := url.Parse(targetURL)
	if err != nil {
		log.Println(err)
	}
	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host
			req.URL.Path = path
			req.Host = target.Host
		},
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}

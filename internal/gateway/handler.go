package gateway

import (
	"fmt"
	"net/http/httputil"
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
	fmt.Println(path)

	proxy := &httputil.ReverseProxy{}

	proxy.ServeHTTP(c.Writer, c.Request)
}

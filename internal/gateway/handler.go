package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
)

type GatewayHandler struct {
	// registry *discovery.EtcdRegistry
	// routes   map[string]*Route
	// proxies  map[string]*httputil.ReverseProxy

	client    *clientv3.Client
	Instances map[string]*ServiceInstance
	prefix    string
	mu        sync.RWMutex
}

type ServiceInstance struct {
	ID          string              `json:"id"`
	ServiceName string              `json:"service_name"`
	Address     string              `json:"address"`
	Port        int                 `json:"port"`
	Metadata    map[string][]string `json:"metadata"`
}

func NewGatewayHandler(client *clientv3.Client) *GatewayHandler {
	gh := &GatewayHandler{
		client:    client,
		Instances: make(map[string]*ServiceInstance),
		prefix:    "/gateway",
	}

	gh.initAllInstances()

	go gh.watch()

	return gh
}

func (gh *GatewayHandler) initAllInstances() {
	resp, err := gh.client.Get(context.Background(), gh.prefix, clientv3.WithPrefix())
	if err != nil {
		log.Println(err)
	}

	for _, kv := range resp.Kvs {
		instance := &ServiceInstance{}
		err = json.Unmarshal(kv.Value, instance)
		if err != nil {
			log.Println(err)
			continue
		}

		gh.Instances[strings.TrimPrefix(string(kv.Key), gh.prefix)] = instance
	}
}

func (gh *GatewayHandler) watch() {
	// WithPrefix() 用于监听该路径下所有子 key 的变化
	watchChan := gh.client.Watch(context.Background(), gh.prefix, clientv3.WithPrefix())

	for watchResp := range watchChan {
		// 检查 Watch 是否出错（如连接丢失、认证过期）
		if watchResp.Err() != nil {
			log.Printf("Watch 响应错误: %v，正在重试...", watchResp.Err())
			continue
		}

		for _, ev := range watchResp.Events {
			gh.mu.Lock()
			switch ev.Type {
			case clientv3.EventTypePut:
				instance := &ServiceInstance{}
				err := json.Unmarshal(ev.Kv.Value, instance)
				if err != nil {
					log.Println(err)
					continue
				}
				gh.Instances[strings.TrimPrefix(string(ev.Kv.Key), gh.prefix)] = instance
			case clientv3.EventTypeDelete:
				delete(gh.Instances, strings.TrimPrefix(string(ev.Kv.Key), gh.prefix))
			}
			gh.mu.Unlock()
		}
	}

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

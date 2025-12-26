package discovery

import (
	"context"
	"fmt"
	"sync"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdRegistry struct {
	client *clientv3.Client
	prefix string
	ctx    context.Context
	cancel context.CancelFunc

	mu sync.RWMutex
}

func NewEtcdRegistry(client *clientv3.Client) *EtcdRegistry {
	ctx, cancel := context.WithCancel(context.Background())

	registry := &EtcdRegistry{
		client: client,
		prefix: "/gateway",
		ctx:    ctx,
		cancel: cancel,
	}

	return registry
}

func (r *EtcdRegistry) watchAllServices() {
	watchChan := r.client.Watch(r.ctx, r.prefix, clientv3.WithPrefix())

	for {
		select {
		case <-r.ctx.Done():
			return
		case resp := <-watchChan:
			if resp.Err() != nil {
				continue
			}

			for _, event := range resp.Events {
				key := string(event.Kv.Key)

				switch event.Type {
				case clientv3.EventTypePut:

				case clientv3.EventTypeDelete:
				}

				fmt.Print(key)
			}
		}
	}
}

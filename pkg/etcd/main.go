package etcd

import (
	"log"

	"github.com/csh0820/gateway/config"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func NewEtcd() *clientv3.Client {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: config.GetConfig().Etcd.Endpoints,
	})
	if err != nil {
		log.Fatalln("init etcd failed:", err)
	}

	return cli
}

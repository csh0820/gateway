# gateway
使用 gin 框架实现 gateway

### 区分服务
每个服务的唯一ID
ID = service_name + address + port

确认ID所在的服务集群
map[ID]service_name

etcd 存储

key: /prefix/ID

value(json) :
{
    "service_name": "order",
    "address": "localhost",
    "port": 8080,
    "metadata": map[string][]string
}

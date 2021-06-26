package etcd

type Config struct {
	Endpoints []string `json:"endpoints" toml:"endpoints"`
}

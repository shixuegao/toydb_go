package config

import "time"

type Config struct {
	Raft `yaml:"raft"`
}

type Raft struct {
	// 当前节点ID
	ID string `yaml:"id"`
	// rpc请求超时时间
	RpcRequestTimeout time.Duration `yaml:"rpcRequestTimeout"`
	// <id, addr>
	Peers map[string]string `yaml:"peers"`
}

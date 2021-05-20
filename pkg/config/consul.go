package config

import (
	"github.com/ilooky/go-layout/pkg/guava"
	"errors"
	"fmt"
	consul "github.com/hashicorp/consul/api"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type cloud struct {
	host   string
	port   string
	consul *consul.Client
}
type Cloud interface {
	ReadConfig(serverName string) (*Config, error)
	Register(serverName string, serverIp string, serverPort int) error
	UnRegister(serverName string, serverIp string, serverPort int) error
	GetServerUri(name string) (string, error)
}

var Client Cloud

func InitCloud() (Cloud, error) {
	if client, err := newConsulClient(); err != nil {
		return nil, err
	} else {
		Client = client
		return client, err
	}
}

func newConsulClient() (Cloud, error) {
	host := guava.GetEnv("CONSUL_HOST", "192.168.1.2")
	port := guava.GetEnv("CONSUL_PORT", "18500")
	config := consul.DefaultConfig()
	config.Address = host + ":" + port
	if c, err := consul.NewClient(config); err != nil {
		return nil, err
	} else {
		return &cloud{consul: c, host: host, port: port}, nil
	}
}

func (c *cloud) ReadConfig(serverName string) (*Config, error) {
	if strings.HasSuffix(serverName, "local") {
		dir, _ := os.Getwd()
		return NewConfig(dir + string(filepath.Separator) + serverName + ".yaml")
	}
	kv := c.consul.KV()
	kvps, _, err := kv.List(serverName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch k/v pairs from consul: %+v", err)
	}
	for _, kvp := range kvps {
		if val := kvp.Value; val != nil {
			k := kvp.Key
			if k == serverName {
				return ParseConfig(string(val[:]))
			}
		}
	}
	return nil, errors.New("not find config")
}

func (c *cloud) Register(serverName string, serverIp string, serverPort int) error {
	healthUrl := fmt.Sprintf("http://%s:%d%s", serverIp, serverPort, "/health")
	reg := &consul.AgentServiceRegistration{
		ID:      ServerId(serverName, serverIp, serverPort),
		Address: serverIp,
		Name:    serverName,
		Port:    serverPort,
		Check: &consul.AgentServiceCheck{
			Status:   consul.HealthPassing,
			Interval: "30s",
			Timeout:  "20s",
			HTTP:     healthUrl,
			Method:   "GET",
		},
	}
	return c.consul.Agent().ServiceRegister(reg)
}

func (c *cloud) UnRegister(serverName string, serverIp string, serverPort int) error {
	return c.consul.Agent().ServiceDeregister(ServerId(serverName, serverIp, serverPort))
}

func (c *cloud) GetServerUri(serviceName string) (uri string, err error) {
	addr, _, err := c.consul.Health().Service(serviceName, "", true, nil)
	if len(addr) == 0 && err == nil {
		return "", fmt.Errorf("service ( %s ) was not found", serviceName)
	}
	if err != nil {
		return "", err
	}
	var urls []string
	for _, entry := range addr {
		if entry.Checks.AggregatedStatus() == consul.HealthPassing {
			urls = append(urls, fmt.Sprintf("http://%s:%d/", entry.Service.Address, entry.Service.Port))
		}
	}
	if len(urls) == 0 {
		panic("No service was found available：" + serviceName)
	}
	i := rand.Intn(len(urls))
	return urls[i], err
}

func ServerId(serverName string, serverIp string, serverPort int) string {
	return serverName + "-" + serverIp + "-" + strconv.Itoa(serverPort)
}

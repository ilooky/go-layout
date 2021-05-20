package config

import (
	"fmt"
	"github.com/ilooky/go-layout/pkg/guava"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

type Config struct {
	Host  string
	Port  string
	Name  string
	Tag   []string
	Mysql Mysql
	DM    DM
	Redis Redis
	Mq    Mq
	Feign Feign
	Log   Log
}
type Log struct {
	Level   string
	Path    string
	Release bool
	Style   string
}

type Feign struct {
	Da2       string
	Da3       string
	Diagram   string
	Manage    string
	Equipment string
}

type Mysql struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	ShowSql  bool
}
type DM struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

type Redis struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

type Mq struct {
	Host        string
	Port        string
	Username    string
	Password    string
	VirtualHost string `yaml:"virtual-host"`
	Queues      []string
	Exchange    string
}

func ParseConfig(value string) (setting *Config, err error) {
	set := &Config{}
	err = yaml.Unmarshal([]byte(value), set)
	if err != nil {
		panic("read config happen err:" + err.Error())
	}
	if set.Name == "" {
		set.Name = "us-diagram"
	}
	if set.Host == "" {
		set.Host = guava.GetEnv("SERVER_HOST", "127.0.0.1")
	}
	prefix := guava.GetEnv("DB_PREFIX", "")
	if set.Mysql.Host == "" {
		set.Mysql.Host = guava.GetEnv("MYSQL_HOST", "127.0.0.1")
	}
	if set.Mysql.Port == "" {
		set.Mysql.Port = guava.GetEnv("MYSQL_PORT", "3306")
	}
	if set.Mysql.Username == "" {
		set.Mysql.Username = guava.GetEnv("MYSQL_USERNAME", "root")
	}
	if set.Mysql.Password == "" {
		set.Mysql.Password = guava.GetEnv("MYSQL_PASSWD", "1234rewq!")
	}
	if set.Mysql.Database == "" {
		set.Mysql.Database = "us_diagram"
	}
	set.Mysql.Database = prefix + set.Mysql.Database

	if set.DM.Host == "" {
		set.DM.Host = guava.GetEnv("DM_HOST", "127.0.0.1")
	}
	if set.DM.Port == "" {
		set.DM.Port = guava.GetEnv("DM_PORT", "5236")
	}
	if set.DM.Username == "" {
		set.DM.Username = guava.GetEnv("DM_USERNAME", "SYSDBA")
	}
	if set.DM.Password == "" {
		set.DM.Password = guava.GetEnv("DM_PASSWD", "SYSDBA!")
	}
	if set.DM.Database == "" {
		set.DM.Database = "us_diagram"
	}
	set.DM.Database = prefix + set.DM.Database

	if set.Redis.Host == "" {
		set.Redis.Host = guava.GetEnv("REDIS_HOST", "127.0.0.1")
	}
	if set.Redis.Port == "" {
		set.Redis.Port = guava.GetEnv("REDIS_PORT", "18160")
	}
	if set.Redis.Username == "" {
		set.Redis.Username = guava.GetEnv("REDIS_USERNAME", "")
	}
	if set.Redis.Password == "" {
		set.Redis.Password = guava.GetEnv("REDIS_PASSWD", "")
	}
	if set.Redis.Database == "" {
		set.Redis.Database = guava.GetEnv("REDIS_DATABASE", "0")
	}
	if set.Mq.Host == "" {
		set.Mq.Host = guava.GetEnv("RABBIT_HOST", "127.0.0.1")
	}
	if set.Mq.Port == "" {
		set.Mq.Port = guava.GetEnv("RABBIT_PORT", "5672")
	}
	if set.Mq.Username == "" {
		set.Mq.Username = guava.GetEnv("RABBIT_USERNAME", "us")
	}
	if set.Mq.Password == "" {
		set.Mq.Password = guava.GetEnv("RABBIT_PASSWD", "1234rewq!")
	}
	if set.Mq.VirtualHost == "" {
		set.Mq.VirtualHost = guava.GetEnv("RABBIT_VHOST", "us")
	}
	set.Mq.VirtualHost = prefix + set.Mq.VirtualHost
	if set.Mq.Exchange == "" {
		set.Mq.Exchange = guava.GetEnv("RABBIT_EXCHANGE", "push")
	}
	if set.Mq.Queues == nil || len(set.Mq.Queues) == 0 {
		set.Mq.Queues = []string{"line", "diagram", "global"}
	}
	if set.Feign.Da2 == "" {
		set.Feign.Da2 = "us-da-v2"
	}
	if set.Feign.Da3 == "" {
		set.Feign.Da3 = "us-da-v3"
	}
	if set.Feign.Manage == "" {
		set.Feign.Manage = "us-manage"
	}
	if set.Feign.Equipment == "" {
		set.Feign.Equipment = "us-equipment"
	}
	if set.Feign.Diagram == "" {
		set.Feign.Diagram = "us-diagram"
	}
	if set.Log.Level == "" {
		set.Log.Level = "info"
	}
	if set.Log.Path == "" {
		set.Log.Path = "/root/projects/us-edps.log"
	}
	return set, nil
}

func NewConfig(configPath string) (*Config, error) {
	if err := validateConfigPath(configPath); err != nil {
		return nil, err
	}
	config := Config{}
	file, err := os.Open(filepath.Clean(configPath))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	d := yaml.NewDecoder(file)
	if err := d.Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

func validateConfigPath(path string) error {
	s, err := os.Stat(filepath.Clean(path))
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory", path)
	}
	return nil
}

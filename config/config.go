package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

type DB struct {
	DSN string `yaml:"dsn" toml:"dsn" json:"dsn"`
}

func (d *DB) valid() []string {
	var msg []string
	if len(d.DSN) == 0 {
		msg = append(msg, "db.dsn is required")
	}

	return msg
}

func newDB() DB {
	return DB{
		DSN: "user=user password=pwd dbname=dbname host=127.0.0.0 port=5432 sslmode=disable TimeZone=Asia/Shanghai",
	}
}

type TSDBType int8

const (
	TSDBTypeInfluxDB TSDBType = iota + 1
	TSDBTypeGreptime
)

var _tsdbMap = map[string]TSDBType{"influxdb": TSDBTypeInfluxDB, "influxDB": TSDBTypeInfluxDB, "greptime": TSDBTypeGreptime}

type TSDB struct {
	Type     string    `yaml:"type" toml:"type" json:"type"`
	TSDBType TSDBType  `yaml:"-" toml:"-" json:"-"`
	InfluxDB *InfluxDB `yaml:"influxDB" toml:"influxDB" json:"influxDB"`
	Greptime *Greptime `yaml:"greptime" toml:"greptime" json:"greptime"`
}

func (t *TSDB) valid() []string {
	var msg []string
	if len(t.Type) == 0 {
		msg = append(msg, "tsdb.type is required")
	}
	tsdbType, ok := _tsdbMap[t.Type]
	if !ok {
		msg = append(msg, fmt.Sprintf("tsdb.type '%s' is not supported", t.Type))
	}
	switch tsdbType {
	case TSDBTypeInfluxDB:
		if t.InfluxDB == nil {
			msg = append(msg, fmt.Sprintf("when tsdb type is '%s' tsdb.influxDB is required", t.Type))
		} else {
			msg = append(msg, t.InfluxDB.valid()...)
		}
	case TSDBTypeGreptime:
		if t.Greptime == nil {
			msg = append(msg, fmt.Sprintf("when tsdb type is '%s' tsdb.greptime is required", t.Type))
		} else {
			msg = append(msg, t.Greptime.valid()...)
		}
	}

	return msg
}

func newTSDB() TSDB {
	return TSDB{
		Type:     "influxdb/greptime",
		InfluxDB: newInfluxDB(),
		Greptime: newGreptime(),
	}
}

type InfluxDB struct {
	URL    string `yaml:"url" toml:"url" json:"url"`
	Token  string `yaml:"token" toml:"token" json:"token"`
	Org    string `yaml:"org" toml:"org" json:"org"`
	Bucket string `yaml:"bucket" toml:"bucket" json:"bucket"`
}

func (i *InfluxDB) valid() []string {
	var msg []string
	if len(i.URL) == 0 {
		msg = append(msg, "tsdb.influxDB.url is required")
	}
	if len(i.Bucket) == 0 {
		msg = append(msg, "tsdb.influxDB.bucket is required")
	}

	return msg
}

func newInfluxDB() *InfluxDB {
	return &InfluxDB{
		URL:    "http://influxdb:8086",
		Token:  "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		Org:    "org",
		Bucket: "test",
	}
}

type Greptime struct {
	Host     string `yaml:"host" toml:"host" json:"host"`
	Port     int    `yaml:"port" toml:"port" json:"port"`
	Database string `yaml:"database" toml:"database" json:"database"`
	Username string `yaml:"username" toml:"username" json:"username"`
	Password string `yaml:"password" toml:"password" json:"password"`
}

func (g Greptime) valid() []string {
	var msg []string
	if len(g.Host) == 0 {
		msg = append(msg, "tsdb.greptime.host is required")
	}
	if g.Port == 0 {
		msg = append(msg, "tsdb.greptime.port is required")
	}
	if len(g.Database) == 0 {
		msg = append(msg, "tsdb.greptime.database is required")
	}
	return msg
}

func newGreptime() *Greptime {
	return &Greptime{
		Host:     "127.0.0.1",
		Port:     4001,
		Username: "username",
		Password: "password",
	}
}

type Server struct {
	Addr string `yaml:"addr" toml:"addr" json:"addr"`
}

func (s *Server) valid() []string {
	var msg []string
	if len(s.Addr) == 0 {
		msg = append(msg, "server.addr is required")
	}

	return msg
}

func newServer() Server {
	return Server{
		Addr: ":5001",
	}
}

var _c *Config

func Init(path string) error {
	c, err := Load(path)
	if err != nil {
		return fmt.Errorf("load config error: %v", err)
	}

	_c = c

	return nil
}

func C() *Config {
	return _c
}

type Config struct {
	DB     DB     `yaml:"db" toml:"db" json:"db"`
	TSDB   TSDB   `yaml:"tsdb" toml:"tsdb" json:"tsdb"`
	Server Server `yaml:"server" toml:"server" json:"server"`
}

func (c *Config) valid() []string {
	var msg []string
	msg = append(msg, c.DB.valid()...)
	msg = append(msg, c.TSDB.valid()...)
	msg = append(msg, c.Server.valid()...)

	return msg
}

func NewConfig() *Config {
	return &Config{
		DB:     newDB(),
		TSDB:   newTSDB(),
		Server: newServer(),
	}
}

func Load(path string) (*Config, error) {
	var cfg Config
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config file %s error: %v", path, err)
	}
	defer f.Close()

	decoder := toml.NewDecoder(f)
	if _, err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("decode config file %s error: %v", path, err)
	}

	msg := cfg.valid()
	if len(msg) != 0 {
		return nil, fmt.Errorf("config file %s error: %v", path, strings.Join(msg, "\n"))
	}

	return &cfg, nil
}

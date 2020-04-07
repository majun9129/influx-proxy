// Copyright 2016 Eleme. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package backend

import (
    "encoding/json"
    "log"
    "os"
)

const (
    VERSION = "1.4.2"
)

// listen_addr: proxy listen addr
// db: proxy db, client's db must be same with it
// username: proxy username
// password: proxy password
// data_dir: data dir to save .dat .rec, default is data
// log_path: log file path, default "" for stdout
// idle_timeout: keep-alives wait time, default is 10000ms
// stat_interval: interval to collect statistics, default is 10000ms
// write_tracing: enable logging for the write, default is false
// query_tracing: enable logging for the query, default is false
// https_enabled: enable https, default is false
// https_cert: the ssl certificate to use when https is enabled
// https_key: use a separate private key location
type NodeConfig struct {
    ListenAddr   string         `json:"listen_addr"`
    DB           string         `json:"db"`
    Username     string         `json:"username"`
    Password     string         `json:"password"`
    DataDir      string         `json:"data_dir"`
    LogPath      string         `json:"log_path"`
    IdleTimeout  int            `json:"idle_timeout"`
    StatInterval int            `json:"stat_interval"`
    WriteTracing bool           `json:"write_tracing"`
    QueryTracing bool           `json:"query_tracing"`
    HTTPSEnabled bool           `json:"https_enabled"`
    HTTPSCert    string         `json:"https_cert"`
    HTTPSKey     string         `json:"https_key"`
}

// url: influxdb addr or other http backend which supports influxdb line protocol
// db: influxdb db
// username: influxdb username
// password: influxdb password
// flush_size: default config is 10000, wait 10000 points write
// flush_time: default config is 1000ms, wait 1 second write whether point count has bigger than flush_size config
// timeout: default config is 10000ms, write timeout until 10 seconds
// check_interval: default config is 1000ms, check backend active every 1 second
// rewrite_interval: default config is 10000ms, rewrite every 10 seconds
// write_only: default is false
type BackendConfig struct {
    URL             string      `json:"url"`
    DB              string      `json:"db"`
    Username        string      `json:"username"`
    Password        string      `json:"password"`
    FlushSize       int         `json:"flush_size"`
    FlushTime       int         `json:"flush_time"`
    Timeout         int         `json:"timeout"`
    CheckInterval   int         `json:"check_interval"`
    RewriteInterval int         `json:"rewrite_interval"`
    WriteOnly       bool        `json:"write_only"`
}

// KEYMAPS
// measurement: [BACKENDS keys], the key must be in the BACKENDS
// http request with the measurement matching the following principles will be forwarded to the backends
// exact match first:
//      for instance, we use cpu.load for measurement's name
//      The KEYMAPS has cpu and cpu.load keys, it will use the cpu.load corresponding backends
// prefix match then:
//      for instance, we use cpu.load for measurement's name
//      the KEYMAPS only has cpu key, it will use the cpu corresponding backends
// _default_ match finally:
//      if the measurement don't match the above principles, it will use the _default_ corresponding backends

// BACKENDS, KEYMAPS, NODE
type FileConfigSource struct {
    BACKENDS     map[string]BackendConfig
    KEYMAPS      map[string][]string
    NODE         NodeConfig
}

func NewFileConfigSource(cfgfile string) (fcs *FileConfigSource, err error) {
    fcs = &FileConfigSource{}
    file, err := os.Open(cfgfile)
    if err != nil {
        return
    }
    defer file.Close()
    dec := json.NewDecoder(file)
    err = dec.Decode(fcs)
    return
}

func (fcs *FileConfigSource) LoadNode() (nodecfg NodeConfig) {
    nodecfg = fcs.NODE
    if nodecfg.ListenAddr == "" {
        nodecfg.ListenAddr = ":7076"
    }
    if nodecfg.DataDir == "" {
        nodecfg.DataDir = "data"
    }
    if nodecfg.IdleTimeout == 0 {
        nodecfg.IdleTimeout = 10000
    }
    if nodecfg.StatInterval == 0 {
        nodecfg.StatInterval = 10000
    }
    return
}

func (fcs *FileConfigSource) LoadBackends() (backends map[string]*BackendConfig, err error) {
    backends = make(map[string]*BackendConfig)
    for name, val := range fcs.BACKENDS {
        cfg := &BackendConfig{
            URL: val.URL,
            DB: val.DB,
            Username: val.Username,
            Password: val.Password,
            FlushSize: val.FlushSize,
            FlushTime: val.FlushTime,
            Timeout: val.Timeout,
            CheckInterval: val.CheckInterval,
            RewriteInterval: val.RewriteInterval,
            WriteOnly: val.WriteOnly,
        }
        if cfg.FlushSize == 0 {
            cfg.FlushSize = 10000
        }
        if cfg.FlushTime == 0 {
            cfg.FlushTime = 1000
        }
        if cfg.Timeout == 0 {
            cfg.Timeout = 10000
        }
        if cfg.CheckInterval == 0 {
            cfg.CheckInterval = 1000
        }
        if cfg.RewriteInterval == 0 {
            cfg.RewriteInterval = 10000
        }
        backends[name] = cfg
    }
    log.Printf("%d backends loaded from file", len(backends))
    return
}

func (fcs *FileConfigSource) LoadMeasurements() (m_map map[string][]string, err error) {
    m_map = fcs.KEYMAPS
    log.Printf("%d measurements loaded from file", len(m_map))
    return
}

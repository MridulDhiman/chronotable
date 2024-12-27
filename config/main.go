package config

import (
	"strings"
	"sync"
	"github.com/spf13/viper"
)

type ConfigHandler struct {
	mtx *sync.Mutex;
	viper *viper.Viper;
}


func NewConfigHandler() *ConfigHandler {
	newViperInst := viper.New()
	newViperInst.AddConfigPath(CHRONO_MAIN_DIR)
	configFileSeg:= strings.Split(CONFIG_FILE, ".")
    newViperInst.SetConfigName(configFileSeg[0]) 
    newViperInst.SetConfigType(configFileSeg[1])   
	return &ConfigHandler{
		viper: newViperInst,
	}
}

func (c *ConfigHandler) Read() error {
	return c.read()
}

func (c *ConfigHandler) Get(key string) (interface{}, error) {
	if err := c.read(); err != nil {
		return nil, err
	}
	return c.viper.Get(key), nil
}

func (c *ConfigHandler) Set(key string, value any) {
	c.mtx.Lock()
	defer c.mtx.Unlock();
	c.viper.Set(key, value);
}

func (c *ConfigHandler) read() error {
	return viper.ReadInConfig()
}


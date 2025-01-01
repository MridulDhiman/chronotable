package config

import (
	"fmt"
	"os"
	"path"
	"strings"
	"sync"
	"github.com/MridulDhiman/chronotable/internal/utils"
	"github.com/spf13/viper"
)

var (
	currVersion int
	latestVersion int
)

type ConfigHandler struct {
	mtx sync.Mutex;
	viper *viper.Viper;
}


func NewConfigHandler() *ConfigHandler {
	newViperInst := viper.New()
	newViperInst.AddConfigPath(CHRONO_MAIN_DIR)
	configFileFullPath := path.Join(CHRONO_MAIN_DIR, CONFIG_FILE)
	yes, err := utils.Exists(configFileFullPath)
	if err != nil {
		panic("Could not find config file")
	}

	if !yes {
		fmt.Println("config file does not exist.")
		fmt.Println("Creating new file...");
		_, err := os.Create(configFileFullPath)
		if err != nil {
			panic("could not create config file")
		}

		fmt.Println("Config file created successful at path: ", configFileFullPath)
		}

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
	if err := c.viper.WriteConfig(); err != nil {
		fmt.Println("(error) *ConfigHandler.Set()", err)
	}
}

func (c *ConfigHandler) read() error {
	return viper.ReadInConfig()
}

func (c* ConfigHandler) UpdateConfigFile(version int) {
	fmt.Println("updating config file...")
	c.Set(ConfigKeyCurrVersion, version);
	c.Set(ConfigKeyLatestVersion, version);
}


func (c* ConfigHandler) FetchLatestVersion() (int,int, error) {
	tempCurr, errCurr := c.Get(ConfigKeyCurrVersion)
	if errCurr != nil {
		return -1, -1, errCurr
	}

	var ok bool;
	if currVersion, ok = tempCurr.(int); !ok {
		return -1, -1, fmt.Errorf("could not convert curr version into integer");
}

	tempLatest, errLatest := c.Get(ConfigKeyLatestVersion)
	if latestVersion, ok = tempLatest.(int); !ok {
		return -1, -1, fmt.Errorf("could not convert curr version into integer");
}
	if errLatest != nil {
		return -1, -1, errCurr
	}
	
	return currVersion, latestVersion, nil
}
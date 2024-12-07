package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"github.com/MridulDhiman/chronotable/config"
)

func UpdateConfigFile(version int) {
	file, err := os.OpenFile(filepath.Join("./", config.CHRONO_MAIN_DIR, config.CONFIG_FILE), os.O_CREATE|os.O_RDWR, os.FileMode(0644))
	if err != nil {
		fmt.Println("(error) could not open file: ", err)
	}
	defer file.Close()

	file.WriteString(`CURR=`+ strconv.Itoa(version));
}

func FetchLatestVersion() (int, error) {
	file, err := os.OpenFile(filepath.Join("./", config.CHRONO_MAIN_DIR, config.CONFIG_FILE), os.O_CREATE|os.O_RDWR, os.FileMode(0644))
	if err != nil {
		return -1,fmt.Errorf("(error) could not open file: %v", err)
	}
	defer file.Close()

	var buf []byte = make([]byte, 100)
	n, err := file.Read(buf)

	if err != nil {
		return -1, fmt.Errorf("(error) could not read file: %v", err)
	}
	data := buf[:n]
	if strings.HasPrefix(string(data),"CURR=") {
		version, err:= strconv.Atoi(strings.TrimPrefix(string(data), "CURR="));
		if err != nil {
			return -1, err
		}
		return version, nil
	} else {
		return -1, errors.New("(error) CURR prefix not found")
	}
}


func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return false, err
}
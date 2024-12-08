package utils

import (
	"bufio"
	
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
	defer file.Close();
	

	file.WriteString(`CURR=`+ strconv.Itoa(version)+"\n");
	file.WriteString(`LATEST=`+ strconv.Itoa(version));
}

func FetchLatestVersion() (int,int, error) {
	file, err := os.OpenFile(filepath.Join("./", config.CHRONO_MAIN_DIR, config.CONFIG_FILE), os.O_CREATE|os.O_RDWR, os.FileMode(0644))
	if err != nil {
		return -1,-1, fmt.Errorf("(error) could not open file: %v", err)
	}
	defer file.Close();
	
	scanner:= bufio.NewScanner(file)
	var CurrentVersion , LatestVersion int = -1,-1

	for scanner.Scan() {
		data := scanner.Text()
		if strings.HasPrefix(data, "CURR=") {
			currentVersion, err := strconv.Atoi(strings.TrimPrefix(data, "CURR="));
			if err != nil {
				return -1, -1, fmt.Errorf("(error) could not convert CURR to integer: %v", err)
			}
			CurrentVersion = currentVersion
			
		} else if(strings.HasPrefix(data, "LATEST=")) {
			latestVersion, err := strconv.Atoi(strings.TrimPrefix(data, "LATEST="));
			if err != nil {
				return -1, -1, fmt.Errorf("(error) could not convert LATEST to integer: %v", err)
			}
			

			LatestVersion = latestVersion
		}
	}

	return CurrentVersion, LatestVersion, nil
}


func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return false, err
}
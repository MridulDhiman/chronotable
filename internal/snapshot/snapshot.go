package snapshot

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/MridulDhiman/chronotable/internal/decoder"
	"github.com/MridulDhiman/chronotable/internal/encoder"
	"github.com/MridulDhiman/chronotable/config"
)

type SnapShot struct {
	LatestVersion   int64
	CurrentVersion  int64
}

type Version struct {
	Id        int64
	Timestamp time.Time
	Data      map[string]interface{}
	Path      string
	Prev int64
}

func New() *SnapShot {
	return &SnapShot{
		LatestVersion:   0,
		CurrentVersion:  0,
	}
}

// TODO: Make the configHandler loosely coupled 
func (snapshot *SnapShot) Create(m map[string]any,  configHandler *config.ConfigHandler) (*Version, error) {
	if snapshot.LatestVersion != 0 {
		latestVersion, err := decodeVersionBinary(getVersionFilePath(snapshot.LatestVersion))
		if err != nil {
			return nil, err
		}
		if same := compareWithLatestVersion(latestVersion, m); same {
			fmt.Println("latest version: ", latestVersion)
			return nil, errors.New("no change since last snapshot")
		}
	}
	newSnapshot, err := createSnapshot(snapshot.LatestVersion, m)
	if err != nil {
		return nil, err
	}

	newSnapshot.Prev = snapshot.CurrentVersion
	snapshot.CurrentVersion = newSnapshot.Id
	configHandler.UpdateConfigFile(newSnapshot.Id)
	snapshot.LatestVersion = newSnapshot.Id
	return newSnapshot, nil
}

func (snapshot *SnapShot) GetVersion(version int64) (*Version, bool) {
	desiredVersion, err := decodeVersionBinary(getVersionFilePath(version))
	if err != nil {
		fmt.Println("(error) GetVersion() could not decode version binary", err)
		return nil, false
	}
	return desiredVersion, true
}

func (snapshot *SnapShot) GetLatestVersion() (*Version, bool) {
	if snapshot.LatestVersion == 0 {
		return nil, true
	}
	return snapshot.GetVersion(snapshot.LatestVersion)
}

func (snapshot *SnapShot) SetLatestVersion(version int64)  {
	snapshot.LatestVersion = version;
}

func (snapshot *SnapShot) SetCurrentVersion(version int64) {
	snapshot.CurrentVersion = version;
}



func decodeVersionBinary(versionFile string) (*Version, error) {
	file, err := os.OpenFile(versionFile, os.O_RDONLY, os.FileMode(0644))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	binDecoder := decoder.NewDecoder(file)
	desiredVersion := new(Version)
	if err := binDecoder.Decode(desiredVersion); err != nil {
		return nil, err
	}
	return desiredVersion, nil
}

func createSnapshot(currentVersion int64, m map[string]interface{}) (*Version, error) {
	_path := path.Join("./", config.CHRONO_MAIN_DIR, fmt.Sprintf("%d", currentVersion+1)+config.SNAPSHOT_EXT)
	newVersion := &Version{
		Timestamp: time.Now(),
		Data:      deepCopy(m),
		Id:        currentVersion + int64(1),
		Path:      _path,
	}
	file, err := os.OpenFile(_path, os.O_WRONLY|os.O_CREATE, os.FileMode(0644))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	binEncoder := encoder.NewEncoder(file)
	if err := binEncoder.Encode(newVersion); err != nil {
		return nil, err
	}
	return newVersion, nil
}

func compareWithLatestVersion(latestVersion *Version, m2 map[string]any) bool {
	m1 := latestVersion.Data
	if len(m1) != len(m2) {
		return false
	}
	for key, val1 := range m1 {
		val2, exists := m2[key]
		if !exists || val1 != val2 {
			return false
		}
	}
	return true
}

func deepCopy(src map[string]interface{}) map[string]interface{} {
	dst := make(map[string]interface{})
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func getVersionFilePath(version int64) string {
	return filepath.Join("./", config.CHRONO_MAIN_DIR, strconv.FormatInt(version, 10) + config.SNAPSHOT_EXT)
}

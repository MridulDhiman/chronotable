package snapshot

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/MridulDhiman/chronotable/config"
	"github.com/MridulDhiman/chronotable/internal/decoder"
	"github.com/MridulDhiman/chronotable/internal/encoder"
	"github.com/MridulDhiman/chronotable/internal/utils"
)

type SnapShot struct {
	LatestVersion   int
	CurrentVersion  int
}

type Version struct {
	Id        int
	Timestamp time.Time
	Data      map[string]interface{}
	Path      string
	AOFStart  int64
	AOFEnd    int64
}

func New(initialized bool) *SnapShot {
	if !initialized {
		utils.UpdateConfigFile(0)
	}

	return &SnapShot{
		LatestVersion:   0,
		CurrentVersion:  0,
	}
}

func (snapshot *SnapShot) Create(m map[string]any, start, end int64) (*Version, error) {
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
	newSnapshot, err := createSnapshot(snapshot.LatestVersion, m, start, end)
	if err != nil {
		return nil, err
	}
	snapshot.CurrentVersion = newSnapshot.Id
	go utils.UpdateConfigFile(newSnapshot.Id)
	snapshot.LatestVersion = newSnapshot.Id
	return newSnapshot, nil
}

func (snapshot *SnapShot) GetVersion(version int) (*Version, bool) {
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

func (snapshot *SnapShot) SetLatestVersion(version int)  {
	snapshot.LatestVersion = version;
}

func (snapshot *SnapShot) SetCurrentVersion(version int) {
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

func createSnapshot(currentVersion int, m map[string]interface{}, start, end int64) (*Version, error) {
	_path := path.Join("./", config.CHRONO_MAIN_DIR, fmt.Sprintf("%d", currentVersion+1)+config.SNAPSHOT_EXT)
	newVersion := &Version{
		Timestamp: time.Now(),
		Data:      deepCopy(m),
		Id:        currentVersion + 1,
		Path:      _path,
		AOFStart:  start,
		AOFEnd:    end,
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

func getVersionFilePath(version int) string {
	return filepath.Join("./", config.CHRONO_MAIN_DIR, strconv.Itoa(version) + config.SNAPSHOT_EXT)
}

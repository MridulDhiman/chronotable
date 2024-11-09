package snapshot

import (
	"errors"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/MridulDhiman/chronotable/config"
	"github.com/MridulDhiman/chronotable/internal/encoder"
)

type SnapShot struct {
	LatestVersion   *Version
	VersionRegistry map[int]*Version
}

type Version struct {
	Id               int
	Timestamp        time.Time
	Data             map[string]interface{}
}

func New() *SnapShot {
	return &SnapShot{
		VersionRegistry: make(map[int]*Version),
		LatestVersion:   &Version{
			Id: 0,
			Timestamp: time.Now(),
			Data: make(map[string]interface{}),
		},
	}
}

func (snapshot *SnapShot) Create(m map[string]any) (*Version,error) {
	if snapshot.LatestVersion.Id != 0 {
		if same := compareWithLatestVersion(snapshot.LatestVersion, m); same {
			return nil, errors.New("no change since last snapshot")
		}
	}
	 newSnapshot, err := createSnapshot(snapshot.LatestVersion.Id, m)
	 if err != nil {
		return nil, err
	 }
	snapshot.LatestVersion = newSnapshot
	snapshot.VersionRegistry[newSnapshot.Id] = newSnapshot
	return newSnapshot, nil
}

func createSnapshot(currentVersion int, m map[string]interface{}) (*Version, error) {
	newVersion := &Version{
		Timestamp: time.Now(),
		Data:      deepCopy(m),
		Id:        currentVersion +1,
	}
	_path := path.Join("./", fmt.Sprintf("%d", currentVersion+ 1) + config.SNAPSHOT_EXT)
	file, err := os.OpenFile(_path, os.O_WRONLY|os.O_CREATE, os.FileMode(0644))
	if err != nil {
		return nil, err
	}
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


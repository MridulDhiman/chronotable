package chronotable

import "time"

type Options struct {
	EnableAOF        bool
	AOFPath          string
	EnableSnapShot   bool
	SnapShotPath     string
	SnapShotInterval time.Duration
}
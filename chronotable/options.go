package chronotable

import "log"

type Options struct {
	EnableAOF      bool
	AOFPath        string
	EnableSnapshot bool
	Initialized    bool
	Logger         *log.Logger
}

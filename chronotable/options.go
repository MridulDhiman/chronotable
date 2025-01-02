package chronotable

import "log"

type Options struct {
	EnableAOF      bool
	AOFPath        string
	EnableSnapshot bool
	Initialized    bool
	Logger         *log.Logger
	Mode           string
}

type InputOpts struct {
	IsReplayed bool
}
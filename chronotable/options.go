package chronotable

type Options struct {
	EnableAOF      bool
	AOFPath        string
	EnableSnapshot bool
	Initialized bool
}

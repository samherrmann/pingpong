package about

import "log"

type About struct {
	Name    string
	Version string
}

func New(version string) *About {
	a := new(About)
	a.Name = "pingpong"
	a.Version = version
	return a
}

func (a *About) Log() {
	log.Println("Running " + a.Name + " version " + a.Version)
}

package version

import "fmt"

var (
	RainyqtVersion = Version{
		Major: 0, Minor: 1, Patch: 0, Build: "clean",
	}
)

type Version struct {
	Major int
	Minor int
	Patch int
	Build string
}

func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d-%s", v.Major, v.Minor, v.Patch, v.Build)
}

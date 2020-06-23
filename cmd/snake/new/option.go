package new

// ProjectInfo ...
type ProjectInfo struct {
	// project dir
	Path string
	// project name
	Name      string
	ModPrefix string
}

var (
	project            ProjectInfo
	DefaultProjectName = "snake-demo"
)

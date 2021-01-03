package new

var (
	Project            ProjectInfo
	DefaultProjectName = "snake-demo"
)

// ProjectInfo ...
type ProjectInfo struct {
	// project dir
	Path string
	// project name
	Name      string
	ModPrefix string
}

package git

import (
	"os"

	"github.com/wx-chevalier/go-utils/commandutils"
)

// Git represents a Git project.
type Git struct {
	dir string
}

// New creates a new git project.
func New(dir string) (Git, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return Git{}, err
	}
	return Git{dir: dir}, nil
}

func (g *Git) command(args ...string) *command.Model {
	cmd := command.New("git", args...)
	cmd.SetDir(g.dir)
	cmd.SetEnvs(append(os.Environ(), "GIT_ASKPASS=echo")...)
	return cmd
}

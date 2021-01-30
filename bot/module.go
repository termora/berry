package bot

import (
	"github.com/starshine-sys/bcr"
)

type botModule struct {
	name     string
	commands bcr.Commands
}

// String returns the module's name
func (b botModule) String() string {
	return b.name
}

// Commands returns a list of commands
func (b *botModule) Commands() []*bcr.Command {
	return b.commands
}

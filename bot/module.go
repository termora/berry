package bot

import (
	"github.com/starshine-sys/bcr"
)

type botModule struct {
	name     string
	commands bcr.Commands
}

func (b botModule) String() string {
	return b.name
}

func (b *botModule) Commands() []*bcr.Command {
	return b.commands
}

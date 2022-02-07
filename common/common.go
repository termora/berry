package common

import (
	"os/exec"
	"strings"

	"github.com/termora/berry/common/log"
)

var Version = "[unknown]"

func init() {
	if Version == "[unknown]" {
		log.Info("Version not set, falling back to checking current directory")

		git := exec.Command("git", "rev-parse", "--short", "HEAD")
		// ignoring errors *should* be fine? if there's no output we just fall back to "unknown"
		b, _ := git.Output()
		Version = strings.TrimSpace(string(b))
		if Version == "" {
			Version = "[unknown]"
		}
	}
}

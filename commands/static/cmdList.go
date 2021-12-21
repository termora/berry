package static

import (
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/bcr/bot"
	"github.com/termora/berry/db"
)

func (bot *Bot) commandList(ctx *bcr.Context) (err error) {
	embeds := make([]discord.Embed, 0)

	// get an accurate page count, modules with 0 non-hidden commands don't show up at all
	var modCount int
	for _, m := range bot.Modules {
		modCount += commandCount(m)
	}

	// create a list of commands per module
	for i, mod := range bot.Modules {
		cmds := make([]string, 0)
		for _, cmd := range mod.Commands() {
			if cmd.Hidden {
				// skip hidden commands
				continue
			}
			cmds = append(cmds, fmt.Sprintf("`%v%v`: %v",
				ifThing(
					cmd.CustomPermissions == nil && cmd.Permissions == 0 && !cmd.OwnerOnly,
					"", "[!] ",
				), cmd.Name, cmd.Summary,
			))
		}

		// if the module has no commands, skip this embed
		if len(cmds) == 0 {
			continue
		}

		embeds = append(embeds, discord.Embed{
			Title:       fmt.Sprintf("%v (%v)", mod.String(), len(mod.Commands())),
			Description: strings.Join(cmds, "\n"),
			Color:       db.EmbedColour,

			Footer: &discord.EmbedFooter{
				Text: fmt.Sprintf("Commands marked with [!] need extra permissions. Page %v/%v", i+1, modCount),
			},
		})
	}

	_, err = ctx.PagedEmbed(embeds, false)
	return err
}

func commandCount(m bot.Module) int {
	var c int
	for _, i := range m.Commands() {
		if !i.Hidden {
			c++
		}
	}

	if c == 0 {
		return 0
	}
	return 1
}

func ifThing(b bool, t, f string) string {
	if b {
		return t
	}
	return f
}

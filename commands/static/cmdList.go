package static

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Starshine113/bcr"
	"github.com/Starshine113/berry/db"
	"github.com/diamondburned/arikawa/v2/discord"
)

type cmdList []*bcr.Command

func (c cmdList) Len() int      { return len(c) }
func (c cmdList) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c cmdList) Less(i, j int) bool {
	return sort.StringsAreSorted([]string{c[i].Name, c[j].Name})
}

func (c *Commands) commandList(ctx *bcr.Context) (err error) {
	var cmds cmdList = ctx.Router.Commands()
	sort.Sort(cmds)

	cmdSlices := make([][]*bcr.Command, 0)

	for i := 0; i < len(cmds); i += 10 {
		end := i + 10

		if end > len(cmds) {
			end = len(cmds)
		}

		cmdSlices = append(cmdSlices, cmds[i:end])
	}

	embeds := make([]discord.Embed, 0)

	for i, slice := range cmdSlices {
		var s strings.Builder
		for _, c := range slice {
			s.WriteString(fmt.Sprintf("`%v%v`: %v\n",
				ifThing(
					c.CustomPermissions == nil && c.Permissions == 0 && !c.OwnerOnly,
					"", "[!] ",
				), c.Name, c.Summary,
			))
		}

		embeds = append(embeds, discord.Embed{
			Title:       fmt.Sprintf("List of commands (%v)", len(cmds)),
			Description: s.String(),
			Color:       db.EmbedColour,

			Footer: &discord.EmbedFooter{
				Text: fmt.Sprintf("Commands marked with [!] need extra permissions. Page %v/%v", i+1, len(cmdSlices)),
			},
		})
	}

	_, err = ctx.PagedEmbed(embeds, false)
	return err
}

func ifThing(b bool, t, f string) string {
	if b {
		return t
	}
	return f
}

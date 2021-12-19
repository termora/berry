package pronouns

import "github.com/starshine-sys/bcr"

func (bot *Bot) random(ctx *bcr.Context) (err error) {
	// we don't wanna repeat code so just call c.use with a random set
	// get a random pronoun set
	set, err := bot.DB.RandomPronouns()
	if err != nil {
		return bot.DB.InternalError(ctx, err)
	}
	// set the arguments to that set
	ctx.RawArgs = set.String()
	ctx.Args = []string{set.String()}
	// return c.use
	return bot.use(ctx)
}

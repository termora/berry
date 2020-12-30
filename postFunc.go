package main

import "github.com/Starshine113/crouter"

// Log when a command is used
func postFunc(ctx *crouter.Ctx) {
	sugar.Debugf("Command executed: `%v` (arguments %v) by %v (channel %v, guild %v)", ctx.Cmd.Name, ctx.Args, ctx.Author.ID, ctx.Channel.ID, ctx.Message.GuildID)
}

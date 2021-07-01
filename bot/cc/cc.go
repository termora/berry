// Package cc contains types and functions for generating static commands from JSON files.
package cc

import (
	"encoding/json"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
)

// Command is a single command, able to be serialized as JSON.
type Command struct {
	// The first name is used as the bcr Command name, the others are used as aliases.
	Names []string `json:"names"`

	Summary     string `json:"summary"`
	Description string `json:"description"`

	Hidden bool `json:"hidden"`

	Response struct {
		Content string         `json:"content"`
		Embed   *discord.Embed `json:"embed"`
	} `json:"response"`
}

// ToBcrCommand converts a Command into a bcr Command.
func (c Command) ToBcrCommand() *bcr.Command {
	if len(c.Names) == 0 {
		return nil
	}

	name := c.Names[0]
	var aliases []string

	if len(c.Names) > 1 {
		aliases = c.Names[1:]
	}

	cmd := &bcr.Command{
		Name:    name,
		Aliases: aliases,

		Summary:     c.Summary,
		Description: c.Description,

		Hidden: c.Hidden,

		Command: func(ctx *bcr.Context) (err error) {
			_, err = ctx.Send(c.Response.Content, c.Response.Embed)
			return err
		},
	}

	return cmd
}

// ParseBytes parses JSON bytes into a slice of bcr Commands.
func ParseBytes(b []byte) (out []*bcr.Command, err error) {
	var s []Command

	err = json.Unmarshal(b, &s)
	if err != nil {
		return nil, err
	}

	for _, c := range s {
		c := c

		out = append(out, c.ToBcrCommand())
	}

	return out, nil
}

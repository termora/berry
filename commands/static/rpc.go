package static

import (
	"context"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/termora/berry/commands/static/rpc"
)

// SendGuildMemberChunk ...
func (c *Commands) SendGuildMemberChunk(ctx context.Context, data *rpc.SendGuildMemberChunkData) (resp *rpc.SendGuildMemberChunkResponse, err error) {
	c.memberMu.Lock()
	defer c.memberMu.Unlock()
	for _, m := range data.GetMembers() {
		c.Sugar.Debugf("Received member %v#%v (%v)", m.Username, m.Discriminator, m.UserId)

		var roles []discord.RoleID
		for _, r := range m.RoleIds {
			roles = append(roles, discord.RoleID(r))
		}

		c.SupportServerMembers[discord.UserID(m.UserId)] = discord.Member{
			User: discord.User{
				ID:            discord.UserID(m.UserId),
				Username:      m.Username,
				Discriminator: m.Discriminator,
			},
			Nick:    m.Nickname,
			RoleIDs: roles,
		}
	}
	c.guildMembersChunked = true

	return &rpc.SendGuildMemberChunkResponse{}, nil
}

// UpdateGuildMember ...
func (c *Commands) UpdateGuildMember(ctx context.Context, data *rpc.UpdateGuildMemberData) (*rpc.UpdateGuildMemberResponse, error) {
	c.Sugar.Debugf("Received member %v#%v (%v)", data.Member.Username, data.Member.Discriminator, data.Member.UserId)

	c.memberMu.Lock()
	defer c.memberMu.Unlock()
	m := discord.Member{
		User: discord.User{
			ID:            discord.UserID(data.Member.UserId),
			Username:      data.Member.Username,
			Discriminator: data.Member.Discriminator,
		},
		Nick: data.Member.Nickname,
	}
	var roles []discord.RoleID
	for _, r := range data.Member.RoleIds {
		roles = append(roles, discord.RoleID(r))
	}
	m.RoleIDs = roles

	c.SupportServerMembers[discord.UserID(data.Member.UserId)] = m

	return &rpc.UpdateGuildMemberResponse{
		ChunkedSent: c.guildMembersChunked,
	}, nil
}

// RemoveGuildMember ...
func (c *Commands) RemoveGuildMember(ctx context.Context, data *rpc.RemoveGuildMemberData) (*rpc.RemoveGuildMemberResponse, error) {
	c.Sugar.Debugf("Removing user %v", data.UserId)

	c.memberMu.Lock()
	defer c.memberMu.Unlock()
	delete(c.SupportServerMembers, discord.UserID(data.UserId))

	return &rpc.RemoveGuildMemberResponse{
		ChunkedSent: c.guildMembersChunked,
	}, nil
}

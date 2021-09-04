package main

import (
	"context"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/termora/berry/commands/static/rpc"
)

func guildCreate(ev *gateway.GuildCreateEvent) {
	if ev.ID != guildID {
		return
	}

	requestMemberChunks()
}

func requestMemberChunks() {
	err := state.RequestGuildMembers(gateway.RequestGuildMembersData{
		GuildIDs: []discord.GuildID{guildID},
	})
	if err != nil {
		log.Errorf("Error requesting guild members: %v", err)
	}
}

func guildMemberAdd(ev *gateway.GuildMemberAddEvent) {
	if ev.GuildID != guildID {
		return
	}

	userMu.Lock()
	users[ev.User.ID] = ev.Member
	userMu.Unlock()

	var roles []uint64
	for _, r := range ev.RoleIDs {
		roles = append(roles, uint64(r))
	}

	resp, err := client.UpdateGuildMember(context.Background(), &rpc.UpdateGuildMemberData{
		Member: &rpc.Member{
			UserId:        uint64(ev.User.ID),
			Username:      ev.User.Username,
			Discriminator: ev.User.Discriminator,
			RoleIds:       roles,
			Nickname:      ev.Nick,
		},
	})
	if err != nil {
		log.Errorf("Error updating guild member: %v", err)
		return
	}
	if !resp.GetChunkedSent() {
		requestMemberChunks()
	}
}

func guildMemberUpdate(ev *gateway.GuildMemberUpdateEvent) {
	if ev.GuildID != guildID {
		return
	}

	userMu.Lock()
	defer userMu.Unlock()
	m, ok := users[ev.User.ID]
	if !ok {
		member, err := state.Member(guildID, ev.User.ID)
		if err != nil {
			log.Errorf("Couldn't get member: %v", err)
			return
		}
		m = *member
	}
	ev.Update(&m)
	users[ev.User.ID] = m

	var roles []uint64
	for _, r := range ev.RoleIDs {
		roles = append(roles, uint64(r))
	}

	resp, err := client.UpdateGuildMember(context.Background(), &rpc.UpdateGuildMemberData{
		Member: &rpc.Member{
			UserId:        uint64(ev.User.ID),
			Username:      ev.User.Username,
			Discriminator: ev.User.Discriminator,
			RoleIds:       roles,
			Nickname:      m.Nick,
		},
	})
	if err != nil {
		log.Errorf("Error updating guild member: %v", err)
		return
	}
	if !resp.GetChunkedSent() {
		requestMemberChunks()
	}
}

func guildMemberRemove(ev *gateway.GuildMemberRemoveEvent) {
	if ev.GuildID != guildID {
		return
	}

	userMu.Lock()
	defer userMu.Unlock()
	delete(users, ev.User.ID)

	resp, err := client.RemoveGuildMember(context.Background(), &rpc.RemoveGuildMemberData{
		UserId: uint64(ev.User.ID),
	})
	if err != nil {
		log.Errorf("Error removing guild member: %v", err)
		return
	}
	if !resp.GetChunkedSent() {
		requestMemberChunks()
	}
}

func guildMemberChunk(ev *gateway.GuildMembersChunkEvent) {
	if ev.GuildID != guildID {
		return
	}

	members := []*rpc.Member{}
	for _, m := range ev.Members {
		var roles []uint64
		for _, r := range m.RoleIDs {
			roles = append(roles, uint64(r))
		}
		members = append(members, &rpc.Member{
			UserId:        uint64(m.User.ID),
			Username:      m.User.Username,
			Discriminator: m.User.Discriminator,
			RoleIds:       roles,
			Nickname:      m.Nick,
		})
	}

	_, err := client.SendGuildMemberChunk(context.Background(), &rpc.SendGuildMemberChunkData{
		Members: members,
	})
	if err != nil {
		log.Errorf("Error sending member chunk: %v", err)
	}
}

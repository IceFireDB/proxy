// Copyright 2014 Wandoujia Inc. All Rights Reserved.
// Licensed under the MIT (MIT-LICENSE.txt) license.

package group

import (
	"github.com/IceFireDB/kit/pkg/models"

	log "github.com/IceFireDB/kit/pkg/logger"
)

type Group struct {
	master       string
	redisServers map[string]models.Server
}

func (g *Group) Master() string {
	return g.master
}

func NewGroup(groupInfo models.ServerGroup) *Group {
	g := &Group{
		redisServers: make(map[string]models.Server),
	}

	// todo get leader correctly
	for _, server := range groupInfo.Servers {
		if server.Type == models.ServerTypeLeader {
			if len(g.master) > 0 {
				log.Fatalf("two master not allowed: %+v", groupInfo)
			}

			g.master = server.Addr
		}
		g.redisServers[server.Addr] = server
	}

	if len(g.master) == 0 {
		log.Fatalf("master not found: %+v", groupInfo)
	}

	return g
}

// todo watch leader state change
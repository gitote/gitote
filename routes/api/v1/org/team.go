package org

import (
	api "gitote.com/gitote/go-gitote-client"

	"gitote.com/gitote/gitote/pkg/context"
	"gitote.com/gitote/gitote/routes/api/v1/convert"
)

func ListTeams(c *context.APIContext) {
	org := c.Org.Organization
	if err := org.GetTeams(); err != nil {
		c.Error(500, "GetTeams", err)
		return
	}

	apiTeams := make([]*api.Team, len(org.Teams))
	for i := range org.Teams {
		apiTeams[i] = convert.ToTeam(org.Teams[i])
	}
	c.JSON(200, apiTeams)
}

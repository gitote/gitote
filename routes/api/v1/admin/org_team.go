package admin

import (
	api "gitote.com/gitote/go-gitote-client"

	"gitote.com/gitote/gitote/models"
	"gitote.com/gitote/gitote/pkg/context"
	"gitote.com/gitote/gitote/routes/api/v1/convert"
	"gitote.com/gitote/gitote/routes/api/v1/user"
)

func CreateTeam(c *context.APIContext, form api.CreateTeamOption) {
	team := &models.Team{
		OrgID:       c.Org.Organization.ID,
		Name:        form.Name,
		Description: form.Description,
		Authorize:   models.ParseAccessMode(form.Permission),
	}
	if err := models.NewTeam(team); err != nil {
		if models.IsErrTeamAlreadyExist(err) {
			c.Error(422, "", err)
		} else {
			c.Error(500, "NewTeam", err)
		}
		return
	}

	c.JSON(201, convert.ToTeam(team))
}

func AddTeamMember(c *context.APIContext) {
	u := user.GetUserByParams(c)
	if c.Written() {
		return
	}
	if err := c.Org.Team.AddMember(u.ID); err != nil {
		c.Error(500, "AddMember", err)
		return
	}

	c.Status(204)
}

func RemoveTeamMember(c *context.APIContext) {
	u := user.GetUserByParams(c)
	if c.Written() {
		return
	}

	if err := c.Org.Team.RemoveMember(u.ID); err != nil {
		c.Error(500, "RemoveMember", err)
		return
	}

	c.Status(204)
}

package admin

import (
	"gitote.com/gitote/gitote/models"
	"gitote.com/gitote/gitote/models/errors"
	"gitote.com/gitote/gitote/pkg/context"
)

func GetRepositoryByParams(c *context.APIContext) *models.Repository {
	repo, err := models.GetRepositoryByName(c.Org.Team.OrgID, c.Params(":reponame"))
	if err != nil {
		if errors.IsRepoNotExist(err) {
			c.Status(404)
		} else {
			c.Error(500, "GetRepositoryByName", err)
		}
		return nil
	}
	return repo
}

func AddTeamRepository(c *context.APIContext) {
	repo := GetRepositoryByParams(c)
	if c.Written() {
		return
	}
	if err := c.Org.Team.AddRepository(repo); err != nil {
		c.Error(500, "AddRepository", err)
		return
	}

	c.Status(204)
}

func RemoveTeamRepository(c *context.APIContext) {
	repo := GetRepositoryByParams(c)
	if c.Written() {
		return
	}
	if err := c.Org.Team.RemoveRepository(repo.ID); err != nil {
		c.Error(500, "RemoveRepository", err)
		return
	}

	c.Status(204)
}

package admin

import (
	api "gitote.com/gitote/go-gitote-client"

	"gitote.com/gitote/gitote/pkg/context"
	"gitote.com/gitote/gitote/routes/api/v1/repo"
	"gitote.com/gitote/gitote/routes/api/v1/user"
)

func CreateRepo(c *context.APIContext, form api.CreateRepoOption) {
	owner := user.GetUserByParams(c)
	if c.Written() {
		return
	}

	repo.CreateUserRepo(c, owner, form)
}

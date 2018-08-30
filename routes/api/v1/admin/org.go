package admin

import (
	api "gitote.com/gitote/go-gitote-client"

	"gitote.com/gitote/gitote/pkg/context"
	"gitote.com/gitote/gitote/routes/api/v1/org"
	"gitote.com/gitote/gitote/routes/api/v1/user"
)

func CreateOrg(c *context.APIContext, form api.CreateOrgOption) {
	org.CreateOrgForUser(c, form, user.GetUserByParams(c))
}

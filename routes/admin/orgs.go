package admin

import (
	"gitote.com/gitote/gitote/models"
	"gitote.com/gitote/gitote/pkg/context"
	"gitote.com/gitote/gitote/pkg/setting"
	"gitote.com/gitote/gitote/routes"
)

const (
	ORGS = "admin/org/list"
)

func Organizations(c *context.Context) {
	c.Data["Title"] = c.Tr("admin.organizations")
	c.Data["PageIsAdmin"] = true
	c.Data["PageIsAdminOrganizations"] = true

	routes.RenderUserSearch(c, &routes.UserSearchOptions{
		Type:     models.USER_TYPE_ORGANIZATION,
		Counter:  models.CountOrganizations,
		Ranger:   models.Organizations,
		PageSize: setting.UI.Admin.OrgPagingNum,
		OrderBy:  "id ASC",
		TplName:  ORGS,
	})
}

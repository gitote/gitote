package context

import (
	"gitote.com/gitote/gitote/models"
)

type APIOrganization struct {
	Organization *models.User
	Team         *models.Team
}

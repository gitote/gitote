package user

import (
	api "gitote.com/gitote/go-gitote-client"

	"gitote.com/gitote/gitote/models"
	"gitote.com/gitote/gitote/pkg/context"
)

func ListAccessTokens(c *context.APIContext) {
	tokens, err := models.ListAccessTokens(c.User.ID)
	if err != nil {
		c.Error(500, "ListAccessTokens", err)
		return
	}

	apiTokens := make([]*api.AccessToken, len(tokens))
	for i := range tokens {
		apiTokens[i] = &api.AccessToken{tokens[i].Name, tokens[i].Sha1}
	}
	c.JSON(200, &apiTokens)
}

func CreateAccessToken(c *context.APIContext, form api.CreateAccessTokenOption) {
	t := &models.AccessToken{
		UID:  c.User.ID,
		Name: form.Name,
	}
	if err := models.NewAccessToken(t); err != nil {
		c.Error(500, "NewAccessToken", err)
		return
	}
	c.JSON(201, &api.AccessToken{t.Name, t.Sha1})
}

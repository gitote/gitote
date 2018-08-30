package user

import (
	api "gitote.com/gitote/go-gitote-client"

	"gitote.com/gitote/gitote/models"
	"gitote.com/gitote/gitote/pkg/context"
	"gitote.com/gitote/gitote/pkg/setting"
	"gitote.com/gitote/gitote/routes/api/v1/convert"
)

func ListEmails(c *context.APIContext) {
	emails, err := models.GetEmailAddresses(c.User.ID)
	if err != nil {
		c.Error(500, "GetEmailAddresses", err)
		return
	}
	apiEmails := make([]*api.Email, len(emails))
	for i := range emails {
		apiEmails[i] = convert.ToEmail(emails[i])
	}
	c.JSON(200, &apiEmails)
}

func AddEmail(c *context.APIContext, form api.CreateEmailOption) {
	if len(form.Emails) == 0 {
		c.Status(422)
		return
	}

	emails := make([]*models.EmailAddress, len(form.Emails))
	for i := range form.Emails {
		emails[i] = &models.EmailAddress{
			UID:         c.User.ID,
			Email:       form.Emails[i],
			IsActivated: !setting.Service.RegisterEmailConfirm,
		}
	}

	if err := models.AddEmailAddresses(emails); err != nil {
		if models.IsErrEmailAlreadyUsed(err) {
			c.Error(422, "", "Email address has been used: "+err.(models.ErrEmailAlreadyUsed).Email)
		} else {
			c.Error(500, "AddEmailAddresses", err)
		}
		return
	}

	apiEmails := make([]*api.Email, len(emails))
	for i := range emails {
		apiEmails[i] = convert.ToEmail(emails[i])
	}
	c.JSON(201, &apiEmails)
}

func DeleteEmail(c *context.APIContext, form api.CreateEmailOption) {
	if len(form.Emails) == 0 {
		c.Status(204)
		return
	}

	emails := make([]*models.EmailAddress, len(form.Emails))
	for i := range form.Emails {
		emails[i] = &models.EmailAddress{
			UID:   c.User.ID,
			Email: form.Emails[i],
		}
	}

	if err := models.DeleteEmailAddresses(emails); err != nil {
		c.Error(500, "DeleteEmailAddresses", err)
		return
	}
	c.Status(204)
}

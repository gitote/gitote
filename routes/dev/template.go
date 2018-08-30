package dev

import (
	"gitote.com/gitote/gitote/models"
	"gitote.com/gitote/gitote/pkg/context"
	"gitote.com/gitote/gitote/pkg/setting"
)

func TemplatePreview(c *context.Context) {
	c.Data["User"] = models.User{Name: "Unknown"}
	c.Data["AppName"] = setting.AppName
	c.Data["AppVer"] = setting.AppVer
	c.Data["AppURL"] = setting.AppURL
	c.Data["Code"] = "2014031910370000009fff6782aadb2162b4a997acb69d4400888e0b9274657374"
	c.Data["ActiveCodeLives"] = setting.Service.ActiveCodeLives / 60
	c.Data["ResetPwdCodeLives"] = setting.Service.ResetPwdCodeLives / 60
	c.Data["CurDbValue"] = ""

	c.HTML(200, (c.Params("*")))
}

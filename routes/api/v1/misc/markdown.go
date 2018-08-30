package misc

import (
	api "gitote.com/gitote/go-gitote-client"

	"gitote.com/gitote/gitote/pkg/context"
	"gitote.com/gitote/gitote/pkg/markup"
)

func Markdown(c *context.APIContext, form api.MarkdownOption) {
	if c.HasApiError() {
		c.Error(422, "", c.GetErrMsg())
		return
	}

	if len(form.Text) == 0 {
		c.Write([]byte(""))
		return
	}

	switch form.Mode {
	case "gfm":
		c.Write(markup.Markdown([]byte(form.Text), form.Context, nil))
	default:
		c.Write(markup.RawMarkdown([]byte(form.Text), ""))
	}
}

func MarkdownRaw(c *context.APIContext) {
	body, err := c.Req.Body().Bytes()
	if err != nil {
		c.Error(422, "", err)
		return
	}
	c.Write(markup.RawMarkdown(body, ""))
}

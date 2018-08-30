package context

import (
	"fmt"
	"strings"

	"gitote.com/yoginth/paginater"
	log "gopkg.in/clog.v1"
	"gopkg.in/macaron.v1"

	"gitote.com/gitote/gitote/pkg/setting"
)

type APIContext struct {
	*Context
	Org *APIOrganization
}

// FIXME: move to gitlab.com/gitote/gitote-client
const DOC_URL = "https://gitlab.com/gitote/gitote-client/wiki"

// Error responses error message to client with given message.
// If status is 500, also it prints error to log.
func (c *APIContext) Error(status int, title string, obj interface{}) {
	var message string
	if err, ok := obj.(error); ok {
		message = err.Error()
	} else {
		message = obj.(string)
	}

	if status == 500 {
		log.Error(3, "%s: %s", title, message)
	}

	c.JSON(status, map[string]string{
		"message": message,
		"url":     DOC_URL,
	})
}

// SetLinkHeader sets pagination link header by given total number and page size.
func (c *APIContext) SetLinkHeader(total, pageSize int) {
	page := paginater.New(total, pageSize, c.QueryInt("page"), 0)
	links := make([]string, 0, 4)
	if page.HasNext() {
		links = append(links, fmt.Sprintf("<%s%s?page=%d>; rel=\"next\"", setting.AppURL, c.Req.URL.Path[1:], page.Next()))
	}
	if !page.IsLast() {
		links = append(links, fmt.Sprintf("<%s%s?page=%d>; rel=\"last\"", setting.AppURL, c.Req.URL.Path[1:], page.TotalPages()))
	}
	if !page.IsFirst() {
		links = append(links, fmt.Sprintf("<%s%s?page=1>; rel=\"first\"", setting.AppURL, c.Req.URL.Path[1:]))
	}
	if page.HasPrevious() {
		links = append(links, fmt.Sprintf("<%s%s?page=%d>; rel=\"prev\"", setting.AppURL, c.Req.URL.Path[1:], page.Previous()))
	}

	if len(links) > 0 {
		c.Header().Set("Link", strings.Join(links, ","))
	}
}

func APIContexter() macaron.Handler {
	return func(ctx *Context) {
		c := &APIContext{
			Context: ctx,
		}
		ctx.Map(c)
	}
}

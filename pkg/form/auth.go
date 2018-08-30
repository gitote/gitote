package form

import (
	"github.com/go-macaron/binding"
	"gopkg.in/macaron.v1"
)

type Authentication struct {
	ID                int64
	Type              int    `binding:"Range(2,5)"`
	Name              string `binding:"Required;MaxSize(30)"`
	Host              string
	Port              int
	BindDN            string
	BindPassword      string
	UserBase          string
	UserDN            string
	AttributeUsername string
	AttributeName     string
	AttributeSurname  string
	AttributeMail     string
	AttributesInBind  bool
	Filter            string
	AdminFilter       string
	GroupEnabled      bool
	GroupDN           string
	GroupFilter       string
	GroupMemberUID    string
	UserUID           string
	IsActive          bool
	SMTPAuth          string
	SMTPHost          string
	SMTPPort          int
	AllowedDomains    string
	SecurityProtocol  int `binding:"Range(0,2)"`
	TLS               bool
	SkipVerify        bool
	PAMServiceName    string
}

func (f *Authentication) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return validate(errs, ctx.Data, f, ctx.Locale)
}

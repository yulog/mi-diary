package presenter

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	cm "github.com/yulog/mi-diary/components"
	"github.com/yulog/mi-diary/logic"
)

func IndexPresentation(c echo.Context, o *logic.IndexOutput) templ.Component {
	return cm.IndexParams{
		Title:     o.Title,
		Profile:   o.Profile,
		Reactions: o.Reactions,
	}.Index()
}

func HashTagPresentation(c echo.Context, o *logic.HashTagOutput) templ.Component {
	return cm.HashTags(o.Profile, o.HashTags)
}

func UserPresentation(c echo.Context, o *logic.UserOutput) templ.Component {
	return cm.Users(o.Profile, o.Users)
}

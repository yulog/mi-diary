package presenter

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	cm "github.com/yulog/mi-diary/components"
	"github.com/yulog/mi-diary/logic"
)

func SelectProfilePresentation(c echo.Context, o *logic.SelectProfileOutput) templ.Component {
	return cm.SelectProfile(o.Title, o.Profiles)
}

func AddProfilePresentation(c echo.Context, o *logic.AddProfileOutput) templ.Component {
	return cm.AddProfile(o.Title)
}

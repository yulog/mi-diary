package presenter

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	cm "github.com/yulog/mi-diary/components"
	"github.com/yulog/mi-diary/logic"
)

func ArchivesPresentation(c echo.Context, o *logic.ArchivesOutput) templ.Component {
	return cm.ArchiveParams{
		Title:   o.Title,
		Profile: o.Profile,
		Items:   o.Items,
	}.Archive()
}

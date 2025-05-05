package presenter

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	cm "github.com/yulog/mi-diary/components"
	"github.com/yulog/mi-diary/logic"
)

func NoteWithPagesPresentation(c echo.Context, o *logic.NoteWithPages) templ.Component {
	return cm.Note{
		Title:      o.Note.Title,
		Profile:    o.Note.Profile,
		Host:       o.Note.Host,
		SearchPath: o.Note.SearchPath,
		Items:      o.Note.Items,
	}.WithPages(
		cm.Pages{
			Current: o.Pages.Current,
			Prev:    cm.Page{Index: o.Pages.Prev.Index, Has: o.Pages.Prev.Has},
			Next:    cm.Page{Index: o.Pages.Next.Index, Has: o.Pages.Next.Has},
			Last:    cm.Page{Index: o.Pages.Last.Index},
			QueryParams: cm.QueryParams{
				Page: o.Pages.QueryParams.Page,
				S:    o.Pages.QueryParams.S,
			},
		})
}

func FileWithPagesPresentation(c echo.Context, o *logic.FileWithPages) templ.Component {
	return cm.File{
		Title:          o.File.Title,
		Profile:        o.File.Profile,
		Host:           o.File.Host,
		FileFilterPath: o.File.FileFilterPath,
		Items:          o.File.Items,
	}.WithPages(
		cm.Pages{
			Current: o.Pages.Current,
			Prev:    cm.Page{Index: o.Pages.Prev.Index, Has: o.Pages.Prev.Has},
			Next:    cm.Page{Index: o.Pages.Next.Index, Has: o.Pages.Next.Has},
			Last:    cm.Page{Index: o.Pages.Last.Index},
			QueryParams: cm.QueryParams{
				Page: o.Pages.QueryParams.Page,
				S:    o.Pages.QueryParams.S,
			},
		})
}

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

func ManagePresentation(c echo.Context, o *logic.ManageOutput) templ.Component {
	if len(o.Profiles) > 0 {
		return cm.ManageInit(o.Title, o.Profiles)
	}
	return cm.ManageStart(o.Title)
}

func JobStartPresentation(c echo.Context, o *logic.JobStartOutput) templ.Component {
	return cm.Start(o.Placeholder, o.Button, o.Profile, o.JobType, o.JobID)
}

func JobProgressPresentation(c echo.Context, o *logic.JobProgressOutput) (int, bool, templ.Component) {
	return o.Progress, o.Completed, cm.Progress(o.ProgressMessage)
}

func JobFinishedPresentation(c echo.Context, o *logic.JobFinishedOutput) templ.Component {
	return cm.Job(o.Placeholder, o.Button, o.ProgressMessage, o.Profiles)
}

func SelectProfilePresentation(c echo.Context, o *logic.SelectProfileOutput) templ.Component {
	return cm.SelectProfile(o.Title, o.Profiles)
}

func AddProfilePresentation(c echo.Context, o *logic.AddProfileOutput) templ.Component {
	return cm.AddProfile(o.Title)
}

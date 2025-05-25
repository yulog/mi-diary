package presenter

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	cm "github.com/yulog/mi-diary/components"
	"github.com/yulog/mi-diary/logic"
)

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

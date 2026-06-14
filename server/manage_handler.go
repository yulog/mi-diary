package server

import (
	"net/http"

	"github.com/labstack/echo/v5"
	cm "github.com/yulog/mi-diary/components"
	"github.com/yulog/mi-diary/domain/model"
	"github.com/yulog/mi-diary/logic"
)

// type ManageHandler struct {
// 	logic       *logic.Logic
// 	manageLogic *logic.ManageLogic
// }

// func NewManageHandler(l *logic.Logic, m *logic.ManageLogic) *ManageHandler {
// 	return &ManageHandler{logic: l, manageLogic: m}
// }

// ManageHandler は /manage のハンドラ
func (srv *Server) ManageHandler(c *echo.Context) error {
	out := srv.logic.ManageLogic(c.Request().Context())
	if len(out.Profiles) > 0 {
		return renderer(c, cm.ManageInit(out.Title, out.Profiles))
	}
	return renderer(c, cm.ManageStart(out.Title))
}

// JobStartHandler は /job/start のハンドラ
func (srv *Server) JobStartHandler(c *echo.Context) error {
	j := new(Job)
	if err := c.Bind(j); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	if err := j.Validate(); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	job := logic.Job{
		Profile: j.Profile,
		Type:    model.JobType(j.Type),
		ID:      j.ID,
	}
	out := srv.logic.JobStartLogic(c.Request().Context(), job)
	return renderer(c, cm.Start(out.Placeholder, out.Button, out.Profile, out.JobType, out.JobID))
}

// JobProgressHandler は /job/progress のハンドラ
func (srv *Server) JobProgressHandler(c *echo.Context) error {
	out := srv.logic.JobProgressLogic(c.Request().Context())

	if out.Completed {
		c.Response().Header().Set("hx-trigger", "done")
	}

	return renderer(c, cm.Progress(out.ProgressMessage))
}

// JobHandler は /job のハンドラ
func (srv *Server) JobHandler(c *echo.Context) error {
	var params Params
	if err := c.Bind(&params); err != nil {
		return err
	}
	out := srv.logic.JobLogic(c.Request().Context(), params.Profile)

	return renderer(c, cm.Job(out.Placeholder, out.Button, out.ProgressMessage, out.Profiles))
}

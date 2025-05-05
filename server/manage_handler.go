package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yulog/mi-diary/domain/model"
	"github.com/yulog/mi-diary/logic"
	"github.com/yulog/mi-diary/presenter"
)

// type ManageHandler struct {
// 	logic       *logic.Logic
// 	manageLogic *logic.ManageLogic
// }

// func NewManageHandler(l *logic.Logic, m *logic.ManageLogic) *ManageHandler {
// 	return &ManageHandler{logic: l, manageLogic: m}
// }

// ManageHandler は /manage のハンドラ
func (srv *Server) ManageHandler(c echo.Context) error {

	return renderer(c, presenter.ManagePresentation(c, srv.logic.ManageLogic(c.Request().Context())))
}

// JobStartHandler は /job/start のハンドラ
func (srv *Server) JobStartHandler(c echo.Context) error {
	j := new(Job)
	if err := c.Bind(j); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	if err := c.Validate(j); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	job := logic.Job{
		Profile: j.Profile,
		Type:    model.JobType(j.Type),
		ID:      j.ID,
	}

	return renderer(c, presenter.JobStartPresentation(c, srv.logic.JobStartLogic(c.Request().Context(), job)))
}

// JobProgressHandler は /job/progress のハンドラ
func (srv *Server) JobProgressHandler(c echo.Context) error {
	_, d, out := presenter.JobProgressPresentation(c, srv.logic.JobProgressLogic(c.Request().Context()))

	if d {
		c.Response().Header().Set("hx-trigger", "done")
	}

	return renderer(c, out)
}

// JobHandler は /job のハンドラ
func (srv *Server) JobHandler(c echo.Context) error {
	var params Params
	if err := c.Bind(&params); err != nil {
		return err
	}
	out := srv.logic.JobLogic(c.Request().Context(), params.Profile)

	return renderer(c, presenter.JobFinishedPresentation(c, out))
}

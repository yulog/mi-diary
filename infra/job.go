package infra

import (
	"github.com/yulog/mi-diary/app"
	"github.com/yulog/mi-diary/logic"
)

type JobInfra struct {
	app *app.App
}

func NewJobInfra(a *app.App) logic.JobRepositorier {
	return &JobInfra{app: a}
}

func (infra *JobInfra) GetProgress() (int, int) {
	infra.app.Progress.RLock()
	defer infra.app.Progress.RUnlock()
	return infra.app.Progress.Progress, infra.app.Progress.Total
}

func (infra *JobInfra) SetProgress(p, t int) (int, int) {
	infra.app.Progress.Lock()
	defer infra.app.Progress.Unlock()
	infra.app.Progress.Progress = p
	infra.app.Progress.Total = t
	return p, t
}

func (infra *JobInfra) UpdateProgress(p, t int) (int, int) {
	cp, ct := infra.GetProgress()
	return infra.SetProgress(cp+p, ct+t)
}

func (infra *JobInfra) GetProgressDone() bool {
	infra.app.Progress.RLock()
	defer infra.app.Progress.RUnlock()
	return infra.app.Progress.Done
}

func (infra *JobInfra) SetProgressDone(d bool) bool {
	infra.app.Progress.Lock()
	defer infra.app.Progress.Unlock()
	infra.app.Progress.Done = d
	return d
}

func (infra *JobInfra) GetJob() chan app.Job {
	return infra.app.Job
}

func (infra *JobInfra) SetJob(j app.Job) {
	infra.app.Job <- j
}

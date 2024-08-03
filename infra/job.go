package infra

import "github.com/yulog/mi-diary/app"

func (infra *Infra) GetProgress() (int, int) {
	infra.app.Progress.RLock()
	defer infra.app.Progress.RUnlock()
	return infra.app.Progress.Progress, infra.app.Progress.Total
}

func (infra *Infra) SetProgress(p, t int) (int, int) {
	infra.app.Progress.Lock()
	defer infra.app.Progress.Unlock()
	infra.app.Progress.Progress = p
	infra.app.Progress.Total = t
	return p, t
}

func (infra *Infra) GetProgressDone() bool {
	infra.app.Progress.RLock()
	defer infra.app.Progress.RUnlock()
	return infra.app.Progress.Done
}

func (infra *Infra) SetProgressDone(d bool) bool {
	infra.app.Progress.Lock()
	defer infra.app.Progress.Unlock()
	infra.app.Progress.Done = d
	return d
}

func (infra *Infra) GetJob() chan app.Job {
	return infra.app.Job
}

func (infra *Infra) SetJob(j app.Job) {
	infra.app.Job <- j
}

package repository

import "github.com/yulog/mi-diary/app"

type JobRepositorier interface {
	GetJob() chan app.Job
	SetJob(j app.Job)

	GetProgress() (int, int)
	SetProgress(p, t int) (int, int)
	UpdateProgress(p, t int) (int, int)
	GetProgressDone() bool
	SetProgressDone(d bool) bool
}

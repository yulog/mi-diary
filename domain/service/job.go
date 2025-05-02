package service

import (
	"context"

	"github.com/yulog/mi-diary/domain/model"
)

type JobServicer interface {
	Execute(ctx context.Context, progressCallback func(int, int)) error
}

type JobWorker interface {
	CreateJob(task JobServicer)
	StartWorker(ctx context.Context)
	GetJobProgress() (int, int)
	GetJobStatus() model.JobStatus
}

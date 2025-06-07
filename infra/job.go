package infra

import (
	"context"
	"sync"

	"github.com/yulog/mi-diary/domain/model"
	"github.com/yulog/mi-diary/domain/service"
)

type JobWorker struct {
	// app      *app.App
	JobQueue chan *JobWrapper
	Progress *Progress
}

type XxxJob struct {
	Profile string
	// Type    JobType
	ID string
}

type Progress struct {
	sync.RWMutex
	Progress int
	Total    int
	Status   model.JobStatus
}

type JobWrapper struct {
	Task service.JobServicer
}

func NewJobWorker() service.JobWorker {
	return &JobWorker{
		JobQueue: make(chan *JobWrapper),
		Progress: &Progress{},
	}
}

func (jobWorker *JobWorker) setJobStatus(s model.JobStatus) {
	jobWorker.Progress.Lock()
	defer jobWorker.Progress.Unlock()
	jobWorker.Progress.Status = s
}

func (jobWorker *JobWorker) worker(ctx context.Context, jobWrapper *JobWrapper) {
	jobWorker.setJobStatus(model.Running)
	err := jobWrapper.Task.Execute(ctx, func(progress, total int) {
		jobWorker.Progress.Lock()
		defer jobWorker.Progress.Unlock()
		jobWorker.Progress.Progress = progress
		jobWorker.Progress.Total = total
	})
	if err != nil {
		jobWorker.setJobStatus(model.Failed)
	} else {
		jobWorker.setJobStatus(model.Completed)
	}
}

// func (jobWorker *JobWorker) CreateJob(task service.JobServicer) *JobWrapper {
func (jobWorker *JobWorker) CreateJob(task service.JobServicer) {
	jobWrapper := &JobWrapper{Task: task}
	jobWorker.JobQueue <- jobWrapper
}

func (jobWorker *JobWorker) StartWorker(ctx context.Context) {
	for jobWrapper := range jobWorker.JobQueue {
		jobWorker.worker(ctx, jobWrapper)
	}
}

func (jobWorker *JobWorker) GetJobProgress() (int, int) {
	jobWorker.Progress.RLock()
	defer jobWorker.Progress.RUnlock()
	return jobWorker.Progress.Progress, jobWorker.Progress.Total
}

func (jobWorker *JobWorker) GetJobStatus() model.JobStatus {
	jobWorker.Progress.RLock()
	defer jobWorker.Progress.RUnlock()
	return jobWorker.Progress.Status
}

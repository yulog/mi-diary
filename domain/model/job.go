package model

type JobStatus int

const (
	Pending JobStatus = iota + 1
	Running
	Completed
	Failed
)

func (s JobStatus) String() string {
	switch s {
	case Pending:
		return "Pending"
	case Running:
		return "Running"
	case Completed:
		return "Completed"
	case Failed:
		return "Failed"
	default:
		return "Unknown"
	}
}

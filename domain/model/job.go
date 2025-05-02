package model

type JobType int

const (
	Reaction JobType = iota + 1
	ReactionOne
	ReactionFull
	Emoji
	Color
)

func (j JobType) String() string {
	switch j {
	case Reaction:
		return "reaction"
	case ReactionOne:
		return "reaction(one)"
	case ReactionFull:
		return "reaction(full scan)"
	case Emoji:
		return "emoji"
	case Color:
		return "color"
	default:
		return "unknown"
	}
}

type JobStatus int

const (
	Pending JobStatus = iota
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

package logic

import "github.com/yulog/mi-diary/domain/model"

type ArchivesOutput struct {
	Title   string
	Profile string
	Items   []model.Month
}

type IndexOutput struct {
	Title     string
	Profile   string
	Reactions []model.ReactionEmoji
}

type HashTagOutput struct {
	Profile  string
	HashTags []model.HashTag
}

type UserOutput struct {
	Profile string
	Users   []model.User
}

type ManageOutput struct {
	Title    string
	Profiles []string
}

type JobStartOutput struct {
	Placeholder string
	Button      string
	Profile     string
	JobType     string
	JobID       string
}

type JobProgressOutput struct {
	Progress        int
	Completed       bool
	ProgressMessage string
}

type JobFinishedOutput struct {
	Placeholder     string
	Button          string
	ProgressMessage string
	Profiles        []string
}

type SelectProfileOutput struct {
	Title    string
	Profiles []string
}

type AddProfileOutput struct {
	Title string
}

type NoteWithPages struct {
	Note  Note
	Pages Pages
}

type FileWithPages struct {
	File  File
	Pages Pages
}

type Note struct {
	Title      string
	Profile    string
	Host       string
	SearchPath string
	Items      []model.Note
}

type File struct {
	Title          string
	Profile        string
	Host           string
	FileFilterPath string
	Items          []model.File
}

type Pages struct {
	Current int
	Prev    Page
	Next    Page
	Last    Page

	QueryParams Params
}

type Page struct {
	Index int
	Has   bool
}

type Params struct {
	Page  int
	S     string
	Color string
}

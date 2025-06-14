package logic

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"strings"
	"time"
	"unicode"

	"github.com/yulog/mi-diary/color"
	"github.com/yulog/mi-diary/domain/model"
	mi "github.com/yulog/miutil"
)

// type ManageLogic struct {
// 	Repo         Repositorier
// 	ProgressRepo ProgressRepositorier
// }

// func NewManageLogic(r *infra.Infra, p *infra.ProgressInfra) *ManageLogic {
// 	return &ManageLogic{Repo: r, ProgressRepo: p}
// }

type Job struct {
	Logic   *Logic
	Profile string
	Type    model.JobType
	ID      string
}

type DummyJob struct {
	Job
}

type ReactionJob struct {
	Job
}

type ReactionOneJob struct {
	Job
}

type ReactionFullJob struct {
	Job
}

type EmojiOneJob struct {
	Job
}

type EmojiFullJob struct {
	Job
}

type ColorOneJob struct {
	Job
}

type ColorFullJob struct {
	Job
}

func (l *Logic) ManageLogic(ctx context.Context) *ManageOutput {
	p, _ := l.JobWorkerService.GetJobProgress()
	// TODO: 進行中の判定これで良いの？
	if p > 0 {
		return &ManageOutput{Title: "Manage"}
	}
	return &ManageOutput{Title: "Manage", Profiles: l.ConfigRepo.GetProfilesSortedKey()}
}

func (l *Logic) JobStartLogic(ctx context.Context, job Job) *JobStartOutput {
	l.CreateJob(ctx, job)

	return &JobStartOutput{
		Placeholder: "",
		Button:      "Get",
		Profile:     job.Profile,
		JobType:     job.Type.String(),
		JobID:       job.ID,
	}
}

func (l *Logic) JobProgressLogic(ctx context.Context) *JobProgressOutput {
	p, t := l.JobWorkerService.GetJobProgress()
	completed := (l.JobWorkerService.GetJobStatus() == model.Completed || l.JobWorkerService.GetJobStatus() == model.Failed)

	return &JobProgressOutput{
		Progress:        p,
		Completed:       completed,
		ProgressMessage: fmt.Sprintf("%d / %d", p, t),
	}
}

func (l *Logic) JobLogic(ctx context.Context, profile string) *JobFinishedOutput {
	p, t := l.JobWorkerService.GetJobProgress()
	// TODO: 進捗、ステータスのリセットをする必要がある？

	return &JobFinishedOutput{
		Placeholder:     "",
		Button:          "Get",
		ProgressMessage: fmt.Sprintf("%d / %d", p, t),
		Profiles:        l.ConfigRepo.GetProfilesSortedKey(),
	}
}

func (l *Logic) CreateJob(ctx context.Context, job Job) {
	switch job.Type {
	case model.Reaction:
		j := &ReactionJob{Job{Logic: l, Profile: job.Profile, Type: job.Type, ID: job.ID}}
		l.JobWorkerService.CreateJob(j)
	case model.ReactionOne:
		j := &ReactionOneJob{Job{Logic: l, Profile: job.Profile, Type: job.Type, ID: job.ID}}
		l.JobWorkerService.CreateJob(j)
	case model.ReactionFull:
		j := &ReactionFullJob{Job{Logic: l, Profile: job.Profile, Type: job.Type, ID: job.ID}}
		l.JobWorkerService.CreateJob(j)
	case model.Emoji:
		if job.ID != "" {
			j := &EmojiOneJob{Job{Logic: l, Profile: job.Profile, Type: job.Type, ID: job.ID}}
			l.JobWorkerService.CreateJob(j)
		} else {
			j := &EmojiFullJob{Job{Logic: l, Profile: job.Profile, Type: job.Type, ID: job.ID}}
			l.JobWorkerService.CreateJob(j)
		}
	case model.Color:
		if job.ID != "" {
			j := &ColorOneJob{Job{Logic: l, Profile: job.Profile, Type: job.Type, ID: job.ID}}
			l.JobWorkerService.CreateJob(j)
		} else {
			j := &ColorFullJob{Job{Logic: l, Profile: job.Profile, Type: job.Type, ID: job.ID}}
			l.JobWorkerService.CreateJob(j)
		}
	default:
		// progressの動作確認用
		j := &DummyJob{Job{Logic: l, Profile: job.Profile, Type: job.Type, ID: job.ID}}
		l.JobWorkerService.CreateJob(j)
	}
}

func (l *Logic) InsertReactionTx(ctx context.Context, profile string, r *mi.Reactions) (rows int64) {
	if len(*r) == 0 {
		return 0
	}

	l.UOWRepo.RunInTx(ctx, profile, func(ctx context.Context) error {
		// JSONの中身をモデルへ移す
		var (
			users       []model.User
			notes       []model.Note
			reactions   []model.ReactionEmoji
			noteToTags  []model.NoteToTag
			files       []model.File
			noteToFiles []model.NoteToFile
		)

		for _, v := range *r {
			var dn string
			if v.Note.User.Name == nil {
				dn = v.Note.User.Username
			} else {
				dn = v.Note.User.Name.(string)
			}
			u := model.User{
				ID:          v.Note.User.ID,
				Name:        v.Note.User.Username,
				DisplayName: dn,
				AvatarURL:   v.Note.User.AvatarURL,
			}
			users = append(users, u)

			for _, tv := range v.Note.Tags {
				ht := model.HashTag{Text: tv}
				err := l.HashTagRepo.Insert(ctx, profile, &ht)
				if err != nil {
					return err
				}
				slog.Info("HashTag ID", slog.Int64("ID", ht.ID))
				// pp.Println(ht.ID)
				// id is scanned automatically
				noteToTags = append(noteToTags, model.NoteToTag{NoteID: v.Note.ID, HashTagID: ht.ID})
			}

			for _, fv := range v.Note.Files {
				f := model.File{
					ID:           fv.ID,
					Name:         fv.Name,
					URL:          fv.URL,
					ThumbnailURL: fv.ThumbnailURL,
					Type:         fv.Type,
					CreatedAt:    fv.CreatedAt,
				}
				files = append(files, f)
				noteToFiles = append(noteToFiles, model.NoteToFile{NoteID: v.Note.ID, FileID: f.ID})
			}

			reactionName := strings.TrimSuffix(strings.TrimPrefix(v.Note.MyReaction, ":"), "@.:")
			n := model.Note{
				ID:                v.Note.ID,
				ReactionID:        v.ID,
				UserID:            v.Note.User.ID,
				ReactionEmojiName: reactionName,
				Text:              v.Note.Text,
				CreatedAt:         v.Note.CreatedAt, // SQLite は日時をUTCで保持する
			}
			notes = append(notes, n)

			symbol := false
			for _, rune := range reactionName {
				if unicode.IsSymbol(rune) {
					symbol = true
					break
				}
			}
			r := model.ReactionEmoji{
				Name:     reactionName,
				IsSymbol: symbol,
			}
			reactions = append(reactions, r)
		}

		// 重複していたらアップデート
		err := l.UserRepo.Insert(ctx, profile, &users)
		if err != nil {
			return err
		}

		// 重複していたら登録しない(エラーにしない)
		rows, err = l.NoteRepo.Insert(ctx, profile, &notes)
		if err != nil {
			return err
		}
		slog.Info("Notes inserted", slog.Int64("count", rows))

		err = l.EmojiRepo.Insert(ctx, profile, &reactions)
		if err != nil {
			return err
		}

		// 0件の場合がある
		if len(noteToTags) > 0 {
			err = l.NoteRepo.InsertNoteToTags(ctx, profile, &noteToTags)
			if err != nil {
				return err
			}
		}

		if len(files) > 0 {
			err = l.FileRepo.Insert(ctx, profile, &files)
			if err != nil {
				return err
			}
		}

		if len(noteToFiles) > 0 {
			err = l.NoteRepo.InsertNoteToFiles(ctx, profile, &noteToFiles)
			if err != nil {
				return err
			}
		}

		err = l.EmojiRepo.UpdateCount(ctx, profile)
		if err != nil {
			return err
		}
		err = l.HashTagRepo.UpdateCount(ctx, profile)
		if err != nil {
			return err
		}
		err = l.UserRepo.UpdateCount(ctx, profile)
		if err != nil {
			return err
		}
		err = l.ArchiveRepo.UpdateCountMonthly(ctx, profile)
		if err != nil {
			return err
		}
		err = l.ArchiveRepo.UpdateCountDaily(ctx, profile)
		if err != nil {
			return err
		}

		return nil
	})
	return rows
}

func (l *Logic) getReactions(ctx context.Context, profile, id string, limit int) (int, *mi.Reactions, error) {
	prof, err := l.ConfigRepo.GetProfile(profile)
	if err != nil {
		return 0, &mi.Reactions{}, err
	}
	count, r, err := l.MisskeyService.Client(prof.Host, prof.I).GetUserReactions(prof.UserID, id, limit)
	if err != nil {
		return 0, &mi.Reactions{}, err
	}

	return count, r, nil
}

func (l *Logic) getEmoji(ctx context.Context, profile, name string) (*mi.Emoji, error) {
	prof, err := l.ConfigRepo.GetProfile(profile)
	if err != nil {
		return &mi.Emoji{}, err
	}
	emoji, err := l.MisskeyService.Client(prof.Host, prof.I).GetEmoji(name)
	if err != nil {
		return &mi.Emoji{}, err
	}

	return emoji, nil
}

func (j *DummyJob) Execute(ctx context.Context, progressCallback func(int, int)) error {
	// progressの動作確認用
	var progress int
	for i := 0; i < 10; i++ {
		progress += 10
		progressCallback(progress, 0)
		fmt.Println(j.Job, progress)
		time.Sleep(time.Second)
	}
	return nil
}

func (j *ReactionJob) Execute(ctx context.Context, progressCallback func(int, int)) error {
	var rid = j.ID
	var progress int
	var total int
	for {
		gc, r, err := j.Logic.getReactions(ctx, j.Profile, rid, 20)
		if err != nil {
			// TODO: エラー処理
			slog.Error(err.Error())
		}
		if gc == 0 || r == nil {
			break
		}

		ac := j.Logic.InsertReactionTx(ctx, j.Profile, r)
		slog.Info("Notes inserted(caller)", slog.Int64("count", ac))

		progress += int(ac)
		total += gc
		progressCallback(progress, total)

		slog.Info("reaction progress", slog.Int("progress", progress), slog.Int("total", total))
		if gc == 0 || ac == 0 {
			break
		}
		rid = (*r)[gc-1].ID
		time.Sleep(rand.N(time.Second))
	}
	return nil
}

func (j *ReactionOneJob) Execute(ctx context.Context, progressCallback func(int, int)) error {
	gc, r, err := j.Logic.getReactions(ctx, j.Profile, j.ID, 1)
	if err != nil {
		// TODO: エラー処理
		slog.Error(err.Error())
	}
	if gc == 0 || r == nil {
		return nil
	}

	ac := j.Logic.InsertReactionTx(ctx, j.Profile, r)
	slog.Info("Notes inserted(caller)", slog.Int64("count", ac))

	progressCallback(int(ac), gc)

	slog.Info("reaction progress", slog.Int("progress", int(ac)), slog.Int("total", gc))

	return nil
}

func (j *ReactionFullJob) Execute(ctx context.Context, progressCallback func(int, int)) error {
	var rid = j.ID
	var progress int
	var total int
	for {
		gc, r, err := j.Logic.getReactions(ctx, j.Profile, rid, 20)
		if err != nil {
			// TODO: エラー処理
			slog.Error(err.Error())
		}

		ac := j.Logic.InsertReactionTx(ctx, j.Profile, r)
		slog.Info("Notes inserted(caller)", slog.Int64("count", ac))

		progress += int(ac)
		total += gc
		progressCallback(progress, total)

		slog.Info("reaction progress", slog.Int("progress", progress), slog.Int("total", total))
		if gc == 0 {
			break
		}
		rid = (*r)[gc-1].ID
		time.Sleep(rand.N(time.Second))
	}
	return nil
}

func (j *EmojiOneJob) Execute(ctx context.Context, progressCallback func(int, int)) error {
	res, err := j.Logic.EmojiRepo.GetByName(ctx, j.Profile, j.ID)
	if err != nil {
		// TODO: エラー処理
		slog.Error(err.Error())
		return err
	}
	emoji, err := j.Logic.getEmoji(ctx, j.Profile, j.ID)
	if err != nil {
		// TODO: エラー処理
		slog.Error(err.Error())
		return err
	}
	j.Logic.EmojiRepo.UpdateByPKWithImage(ctx, j.Profile, res.ID, emoji)

	progressCallback(1, 1)
	slog.Info("emoji progress", slog.Int("progress", 1), slog.Int("total", 1))
	return nil
}

func (j *EmojiFullJob) Execute(ctx context.Context, progressCallback func(int, int)) error {
	r, err := j.Logic.EmojiRepo.GetByEmptyImage(ctx, j.Profile)
	if err != nil {
		// TODO: エラー処理
		slog.Error(err.Error())
		return err
	}

	var progress int
	for _, v := range r {
		// unicode emojiならスキップしたい
		symbol := false
		for _, rune := range v.Name {
			if unicode.IsSymbol(rune) {
				symbol = true
				break
			}
		}
		if symbol {
			continue
		}
		emoji, err := j.Logic.getEmoji(ctx, j.Profile, v.Name)
		if err != nil {
			// TODO: エラー処理
			slog.Error(err.Error())
			continue
		}
		j.Logic.EmojiRepo.UpdateByPKWithImage(ctx, j.Profile, v.ID, emoji)

		progress += 1
		progressCallback(progress, len(r))
		slog.Info("emoji progress", slog.Int("progress", progress), slog.Int("total", len(r)))

		time.Sleep(rand.N(time.Second))
	}
	return nil
}

func (j *ColorOneJob) Execute(ctx context.Context, progressCallback func(int, int)) error {
	r, err := j.Logic.FileRepo.GetByNoteID(ctx, j.Profile, j.ID)
	if err != nil {
		// TODO: エラー処理
		slog.Error(err.Error())
		return err
	}

	var progress int
	for _, v := range r {
		if !strings.HasPrefix(v.Type, "image") {
			continue
		}
		c1, c2, err := color.Color(v.ThumbnailURL)
		if err != nil {
			// TODO: エラー処理
			slog.Error(err.Error(), slog.String("file_id", v.ID), slog.String("url", v.ThumbnailURL), slog.String("dominant_color", c1), slog.String("group_color", c2))
			continue
		}
		slog.Info("get color", slog.String("file_id", v.ID), slog.String("url", v.ThumbnailURL), slog.String("dominant_color", c1), slog.String("group_color", c2))
		j.Logic.FileRepo.UpdateByPKWithColor(ctx, j.Profile, v.ID, c1, c2)

		progress += 1
		progressCallback(progress, len(r))
		slog.Info("color progress", slog.Int("progress", progress), slog.Int("total", len(r)))

		time.Sleep(rand.N(5 * time.Second))
	}
	return nil
}

func (j *ColorFullJob) Execute(ctx context.Context, progressCallback func(int, int)) error {
	r, err := j.Logic.FileRepo.GetByEmptyColor(ctx, j.Profile)
	if err != nil {
		// TODO: エラー処理
		slog.Error(err.Error())
		return err
	}

	var progress int
	for _, v := range r {
		if !strings.HasPrefix(v.Type, "image") {
			continue
		}
		c1, c2, err := color.Color(v.ThumbnailURL)
		if err != nil {
			// TODO: エラー処理
			slog.Error(err.Error(), slog.String("file_id", v.ID), slog.String("url", v.ThumbnailURL), slog.String("dominant_color", c1), slog.String("group_color", c2))
			// TODO: エラーだったとき、次回の処理対象にならないようにする
			continue
		}
		slog.Info("get color", slog.String("file_id", v.ID), slog.String("url", v.ThumbnailURL), slog.String("dominant_color", c1), slog.String("group_color", c2))
		j.Logic.FileRepo.UpdateByPKWithColor(ctx, j.Profile, v.ID, c1, c2)

		progress += 1
		progressCallback(progress, len(r))
		slog.Info("color progress", slog.Int("progress", progress), slog.Int("total", len(r)))

		time.Sleep(rand.N(5 * time.Second))
	}
	return nil
}

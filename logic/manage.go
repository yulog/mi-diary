package logic

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"strings"
	"time"
	"unicode"

	"github.com/a-h/templ"
	"github.com/yulog/mi-diary/app"
	"github.com/yulog/mi-diary/color"
	cm "github.com/yulog/mi-diary/components"
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

func (l *Logic) ManageLogic(ctx context.Context) templ.Component {
	p, _ := l.JobRepo.GetProgress()
	// TODO: 進行中の判定これで良いの？
	if p > 0 {
		return cm.ManageStart("Manage")
	}
	return cm.ManageInit("Manage", l.ConfigRepo.GetProfilesSortedKey())
}

func (l *Logic) JobStartLogic(ctx context.Context, job app.Job) templ.Component {
	l.JobRepo.SetJob(job)

	return cm.Start("", "Get", job.Profile, job.Type.String(), job.ID)
}

func (l *Logic) JobProgressLogic(ctx context.Context) (int, bool, templ.Component) {
	p, t := l.JobRepo.GetProgress()

	return p, l.JobRepo.GetProgressDone(), cm.Progress(fmt.Sprintf("%d / %d", p, t))
}

func (l *Logic) JobLogic(ctx context.Context, profile string) templ.Component {
	p, t := l.JobRepo.GetProgress()
	l.JobRepo.SetProgress(0, 0)
	l.JobRepo.SetProgressDone(false)

	return cm.Job("", "Get", fmt.Sprintf("%d / %d", p, t), l.ConfigRepo.GetProfilesSortedKey())
}

func (l *Logic) JobProcesser(ctx context.Context) {
	for j := range l.JobRepo.GetJob() {
		switch j.Type {
		case app.Reaction:
			l.reactionJob(ctx, j)
		case app.ReactionOne:
			l.reactionOneJob(ctx, j)
		case app.ReactionFull:
			l.reactionFullJob(ctx, j)
		case app.Emoji:
			if j.ID != "" {
				l.emojiOneJob(ctx, j)
			} else {
				l.emojiFullJob(ctx, j)
			}
		case app.Color:
			if j.ID != "" {
				l.colorOneJob(ctx, j)
			} else {
				l.colorFullJob(ctx, j)
			}
		default:
			// progressの動作確認用
			for i := 0; i < 10; i++ {
				p, _ := l.JobRepo.GetProgress()
				p, t := l.JobRepo.SetProgress(p+10, 0)
				fmt.Println(j, p, t)
				time.Sleep(time.Second)
			}
		}
		l.JobRepo.SetProgressDone(true)
	}
}

func (l *Logic) reactionJob(ctx context.Context, j app.Job) {
	var rid = j.ID
	for {
		gc, r, err := l.getReactions(ctx, j.Profile, rid, 20)
		if err != nil {
			// TODO: エラー処理
			slog.Error(err.Error())
		}
		if gc == 0 || r == nil {
			break
		}

		ac := l.InsertReactionTx(ctx, j.Profile, r)
		slog.Info("Notes inserted(caller)", slog.Int64("count", ac))

		p, t := l.JobRepo.UpdateProgress(int(ac), gc)

		slog.Info("reaction progress", slog.Int("progress", p), slog.Int("total", t))
		if gc == 0 || ac == 0 {
			break
		}
		rid = (*r)[gc-1].ID
		time.Sleep(rand.N(time.Second))
	}
}

func (l *Logic) reactionOneJob(ctx context.Context, j app.Job) {
	gc, r, err := l.getReactions(ctx, j.Profile, j.ID, 1)
	if err != nil {
		// TODO: エラー処理
		slog.Error(err.Error())
	}
	if gc == 0 || r == nil {
		return
	}

	ac := l.InsertReactionTx(ctx, j.Profile, r)
	slog.Info("Notes inserted(caller)", slog.Int64("count", ac))

	p, t := l.JobRepo.UpdateProgress(int(ac), gc)

	slog.Info("reaction progress", slog.Int("progress", p), slog.Int("total", t))
	if gc == 0 || ac == 0 {
		return
	}
}

func (l *Logic) reactionFullJob(ctx context.Context, j app.Job) {
	var rid = j.ID
	for {
		gc, r, err := l.getReactions(ctx, j.Profile, rid, 20)
		if err != nil {
			// TODO: エラー処理
			slog.Error(err.Error())
		}

		ac := l.InsertReactionTx(ctx, j.Profile, r)
		slog.Info("Notes inserted(caller)", slog.Int64("count", ac))

		p, t := l.JobRepo.UpdateProgress(int(ac), gc)

		slog.Info("reaction progress", slog.Int("progress", p), slog.Int("total", t))
		if gc == 0 {
			break
		}
		rid = (*r)[gc-1].ID
		time.Sleep(rand.N(time.Second))
	}
}

func (l *Logic) InsertReactionTx(ctx context.Context, profile string, r *mi.Reactions) (rows int64) {
	if len(*r) == 0 {
		return 0
	}

	l.Repo.RunInTx(ctx, profile, func(ctx context.Context) error {
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

			r := model.ReactionEmoji{
				Name: reactionName,
			}
			reactions = append(reactions, r)
		}

		// 重複していたらアップデート
		err := l.UserRepo.Insert(ctx, profile, &users)
		if err != nil {
			return err
		}

		// 重複していたら登録しない(エラーにしない)
		rows, err = l.Repo.InsertNotes(ctx, profile, &notes)
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
			err = l.Repo.InsertNoteToTags(ctx, profile, &noteToTags)
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
			err = l.Repo.InsertNoteToFiles(ctx, profile, &noteToFiles)
			if err != nil {
				return err
			}
		}

		err = l.Repo.Count(ctx, profile)
		// TODO: あってもなくても変わらない vs 統一感
		if err != nil {
			return err
		}
		return err
	})
	return rows
}

func (l *Logic) emojiOneJob(ctx context.Context, j app.Job) {
	res, err := l.EmojiRepo.GetByName(ctx, j.Profile, j.ID)
	if err != nil {
		// TODO: エラー処理
		slog.Error(err.Error())
	}
	emoji, err := l.getEmoji(ctx, j.Profile, j.ID)
	if err != nil {
		// TODO: エラー処理
		slog.Error(err.Error())
	}
	l.EmojiRepo.UpdateByPKWithImage(ctx, j.Profile, res.ID, emoji)

	p, _ := l.JobRepo.GetProgress()
	l.JobRepo.SetProgress(p+1, 1)
	slog.Info("emoji progress", slog.Int("progress", p+1), slog.Int("total", 1))
}

func (l *Logic) emojiFullJob(ctx context.Context, j app.Job) {
	r, err := l.EmojiRepo.GetByEmptyImage(ctx, j.Profile)
	if err != nil {
		// TODO: エラー処理
		slog.Error(err.Error())
	}

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
		emoji, err := l.getEmoji(ctx, j.Profile, v.Name)
		if err != nil {
			// TODO: エラー処理
			slog.Error(err.Error())
		}
		l.EmojiRepo.UpdateByPKWithImage(ctx, j.Profile, v.ID, emoji)

		p, _ := l.JobRepo.GetProgress()
		l.JobRepo.SetProgress(p+1, len(r))
		slog.Info("emoji progress", slog.Int("progress", p+1), slog.Int("total", len(r)))

		time.Sleep(rand.N(time.Second))
	}
}

func (l *Logic) colorOneJob(ctx context.Context, j app.Job) {
	r, err := l.FileRepo.GetByNoteID(ctx, j.Profile, j.ID)
	if err != nil {
		// TODO: エラー処理
		slog.Error(err.Error())
	}

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
		l.FileRepo.UpdateByPKWithColor(ctx, j.Profile, v.ID, c1, c2)

		p, _ := l.JobRepo.GetProgress()
		l.JobRepo.SetProgress(p+1, len(r))
		slog.Info("color progress", slog.Int("progress", p+1), slog.Int("total", len(r)))

		time.Sleep(rand.N(5 * time.Second))
	}
}

func (l *Logic) colorFullJob(ctx context.Context, j app.Job) {
	r, err := l.FileRepo.GetByEmptyColor(ctx, j.Profile)
	if err != nil {
		// TODO: エラー処理
		slog.Error(err.Error())
	}

	p, _ := l.JobRepo.GetProgress()
	l.JobRepo.SetProgress(p, len(r))
	slog.Info("color progress", slog.Int("progress", p), slog.Int("total", len(r)))

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
		l.FileRepo.UpdateByPKWithColor(ctx, j.Profile, v.ID, c1, c2)

		p, _ := l.JobRepo.GetProgress()
		l.JobRepo.SetProgress(p+1, len(r))
		slog.Info("color progress", slog.Int("progress", p+1), slog.Int("total", len(r)))

		time.Sleep(rand.N(5 * time.Second))
	}
}

func (l *Logic) getReactions(ctx context.Context, profile, id string, limit int) (int, *mi.Reactions, error) {
	count, r, err := l.MisskeyAPIRepo.GetUserReactions(profile, id, limit)
	if err != nil {
		return 0, &mi.Reactions{}, err
	}

	return count, r, nil
}

func (l *Logic) getEmoji(ctx context.Context, profile, name string) (*mi.Emoji, error) {
	emoji, err := l.MisskeyAPIRepo.GetEmoji(profile, name)
	if err != nil {
		return &mi.Emoji{}, err
	}

	return emoji, nil
}

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
		gc, r := l.getReactions(ctx, j.Profile, rid, 20)
		if gc == 0 || r == nil {
			break
		}
		ac := l.Repo.Insert(ctx, j.Profile, r)

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
	gc, r := l.getReactions(ctx, j.Profile, j.ID, 1)
	if gc == 0 || r == nil {
		return
	}
	ac := l.Repo.Insert(ctx, j.Profile, r)

	p, t := l.JobRepo.UpdateProgress(int(ac), gc)

	slog.Info("reaction progress", slog.Int("progress", p), slog.Int("total", t))
	if gc == 0 || ac == 0 {
		return
	}
}

func (l *Logic) reactionFullJob(ctx context.Context, j app.Job) {
	var rid = j.ID
	for {
		gc, r := l.getReactions(ctx, j.Profile, rid, 20)
		ac := l.Repo.Insert(ctx, j.Profile, r)

		p, t := l.JobRepo.UpdateProgress(int(ac), gc)

		slog.Info("reaction progress", slog.Int("progress", p), slog.Int("total", t))
		if gc == 0 {
			break
		}
		rid = (*r)[gc-1].ID
		time.Sleep(rand.N(time.Second))
	}
}

func (l *Logic) emojiOneJob(ctx context.Context, j app.Job) {
	res, err := l.Repo.ReactionOne(ctx, j.Profile, j.ID)
	if err != nil {
		// TODO: エラー処理
		slog.Error(err.Error())
	}
	emoji := l.getEmoji(ctx, j.Profile, j.ID)
	l.Repo.InsertEmoji(ctx, j.Profile, res.ID, emoji)

	p, _ := l.JobRepo.GetProgress()
	l.JobRepo.SetProgress(p+1, 1)
	slog.Info("emoji progress", slog.Int("progress", p+1), slog.Int("total", 1))
}

func (l *Logic) emojiFullJob(ctx context.Context, j app.Job) {
	r, err := l.Repo.ReactionImageEmpty(ctx, j.Profile)
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
		emoji := l.getEmoji(ctx, j.Profile, v.Name)
		l.Repo.InsertEmoji(ctx, j.Profile, v.ID, emoji)

		p, _ := l.JobRepo.GetProgress()
		l.JobRepo.SetProgress(p+1, len(r))
		slog.Info("emoji progress", slog.Int("progress", p+1), slog.Int("total", len(r)))

		time.Sleep(rand.N(time.Second))
	}
}

func (l *Logic) colorOneJob(ctx context.Context, j app.Job) {
	r, err := l.Repo.FilesByNoteID(ctx, j.Profile, j.ID)
	if err != nil {
		// TODO: エラー処理
		slog.Error(err.Error())
	}

	for _, v := range r {
		if strings.HasPrefix(v.Type, "image") {
			continue
		}
		c1, c2, err := color.Color(v.ThumbnailURL)
		if err != nil {
			// TODO: エラー処理
			slog.Error(err.Error(), slog.String("file_id", v.ID), slog.String("url", v.ThumbnailURL), slog.String("dominant_color", c1), slog.String("group_color", c2))
			continue
		}
		slog.Info("get color", slog.String("file_id", v.ID), slog.String("url", v.ThumbnailURL), slog.String("dominant_color", c1), slog.String("group_color", c2))
		l.Repo.InsertColor(ctx, j.Profile, v.ID, c1, c2)

		p, _ := l.JobRepo.GetProgress()
		l.JobRepo.SetProgress(p+1, len(r))
		slog.Info("color progress", slog.Int("progress", p+1), slog.Int("total", len(r)))

		time.Sleep(rand.N(5 * time.Second))
	}
}

func (l *Logic) colorFullJob(ctx context.Context, j app.Job) {
	r, err := l.Repo.FilesColorEmpty(ctx, j.Profile)
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
		l.Repo.InsertColor(ctx, j.Profile, v.ID, c1, c2)

		p, _ := l.JobRepo.GetProgress()
		l.JobRepo.SetProgress(p+1, len(r))
		slog.Info("color progress", slog.Int("progress", p+1), slog.Int("total", len(r)))

		time.Sleep(rand.N(5 * time.Second))
	}
}

func (l *Logic) getReactions(ctx context.Context, profile, id string, limit int) (int, *mi.Reactions) {
	prof, err := l.ConfigRepo.GetProfile(profile)
	if err != nil {
		return 0, &mi.Reactions{}
	}
	count, r, err := l.Repo.GetUserReactions(prof, id, limit)
	if err != nil {
		// TODO: エラー処理
		slog.Error(err.Error())
	}

	return count, r
}

func (l *Logic) getEmoji(ctx context.Context, profile, name string) *mi.Emoji {
	prof, err := l.ConfigRepo.GetProfile(profile)
	if err != nil {
		return &mi.Emoji{}
	}
	emoji, err := l.Repo.GetEmoji(prof, name)
	if err != nil {
		// TODO: エラー処理
		slog.Error(err.Error())
	}

	return emoji
}

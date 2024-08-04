package logic

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"time"
	"unicode"

	"github.com/a-h/templ"
	"github.com/yulog/mi-diary/app"
	cm "github.com/yulog/mi-diary/components"
	mi "github.com/yulog/miutil"
)

func (l *Logic) JobStartLogic(ctx context.Context, job app.Job) templ.Component {
	l.repo.SetJob(job)

	return cm.Start("", "Get", job.Profile, job.Type.String(), job.ID)
}

func (l *Logic) JobProgressLogic(ctx context.Context) (int, bool, templ.Component) {
	p, t := l.repo.GetProgress()

	return p, l.repo.GetProgressDone(), cm.Progress(fmt.Sprintf("%d / %d", p, t))
}

func (l *Logic) JobLogic(ctx context.Context, profile string) templ.Component {
	p, t := l.repo.GetProgress()
	l.repo.SetProgress(0, 0)
	l.repo.SetProgressDone(false)
	var ps []string
	for k := range *l.repo.GetProfiles() {
		ps = append(ps, k)
	}

	return cm.Job("", "Get", fmt.Sprintf("%d / %d", p, t), ps)
}

func (l *Logic) JobProcesser(ctx context.Context) {
	for j := range l.repo.GetJob() {
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
		default:
			// progressの動作確認用
			for i := 0; i < 10; i++ {
				p, _ := l.repo.GetProgress()
				p, t := l.repo.SetProgress(p+10, 0)
				fmt.Println(j, p, t)
				time.Sleep(time.Second)
			}
		}
		l.repo.SetProgressDone(true)
	}
}

func (l *Logic) reactionJob(ctx context.Context, j app.Job) {
	var rid = j.ID
	for {
		gc, r := l.getReactions(ctx, j.Profile, rid, 20)
		if gc == 0 || r == nil {
			break
		}
		ac := l.repo.Insert(ctx, j.Profile, r)

		p, t := l.repo.UpdateProgress(int(ac), gc)

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
	ac := l.repo.Insert(ctx, j.Profile, r)

	p, t := l.repo.UpdateProgress(int(ac), gc)

	slog.Info("reaction progress", slog.Int("progress", p), slog.Int("total", t))
	if gc == 0 || ac == 0 {
		return
	}
}

func (l *Logic) reactionFullJob(ctx context.Context, j app.Job) {
	var rid = j.ID
	for {
		gc, r := l.getReactions(ctx, j.Profile, rid, 20)
		ac := l.repo.Insert(ctx, j.Profile, r)

		p, t := l.repo.UpdateProgress(int(ac), gc)

		slog.Info("reaction progress", slog.Int("progress", p), slog.Int("total", t))
		if gc == 0 {
			break
		}
		rid = (*r)[gc-1].ID
		time.Sleep(rand.N(time.Second))
	}
}

func (l *Logic) emojiOneJob(ctx context.Context, j app.Job) {
	res, err := l.repo.ReactionOne(ctx, j.Profile, j.ID)
	if err != nil {
		// TODO: エラー処理
		slog.Error(err.Error())
	}
	emoji := l.getEmoji(ctx, j.Profile, j.ID)
	l.repo.InsertEmoji(ctx, j.Profile, res.ID, emoji)

	p, _ := l.repo.GetProgress()
	l.repo.SetProgress(p+1, 1)
	slog.Info("emoji progress", slog.Int("progress", p+1), slog.Int("total", 1))
}

func (l *Logic) emojiFullJob(ctx context.Context, j app.Job) {
	r, err := l.repo.ReactionImageEmpty(ctx, j.Profile)
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
		l.repo.InsertEmoji(ctx, j.Profile, v.ID, emoji)

		p, _ := l.repo.GetProgress()
		l.repo.SetProgress(p+1, len(r))
		slog.Info("emoji progress", slog.Int("progress", p+1), slog.Int("total", len(r)))

		time.Sleep(rand.N(time.Second))
	}
}

func (l *Logic) getReactions(ctx context.Context, profile, id string, limit int) (int, *mi.Reactions) {
	count, r, err := l.repo.GetUserReactions(profile, id, limit)
	if err != nil {
		// TODO: エラー処理
		slog.Error(err.Error())
	}

	return count, r
}

func (l *Logic) getEmoji(ctx context.Context, profile, name string) *mi.Emoji {
	emoji, err := l.repo.GetEmoji(profile, name)
	if err != nil {
		// TODO: エラー処理
		slog.Error(err.Error())
	}

	return emoji
}

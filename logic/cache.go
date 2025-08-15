package logic

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"

	"github.com/yulog/mi-diary/domain/model"
)

func (l *Logic) CacheLogic(ctx context.Context, profile, name string, emoji model.ReactionEmoji) (CacheOutput, error) {
	host, err := l.ConfigRepo.GetProfileHost(profile)
	if err != nil {
		return CacheOutput{}, err
	}

	val, err := l.CacheRepo.Get(host, name)
	if err != nil {
		log.Fatal(err)
		return CacheOutput{}, err
	}
	if val != nil {
		log.Println("cache hit")
		return CacheOutput{
			StatusCode:  http.StatusOK,
			ContentType: http.DetectContentType(val),
			Response:    io.NopCloser(bytes.NewBuffer(val)),
			DoCache:     func() {},
		}, nil
	}

	resp, err := http.Get(emoji.Image)
	if err != nil {
		log.Fatal(err)
		return CacheOutput{}, err
	}
	// ここでdeferせず、呼び出し元でする
	// defer resp.Body.Close()
	var dest bytes.Buffer
	tee := io.TeeReader(resp.Body, &dest) // teeが読みだされるとdestにも値が入る

	return CacheOutput{
		StatusCode:  resp.StatusCode,
		ContentType: resp.Header.Get("Content-Type"),
		Response: &struct {
			io.Reader
			io.Closer
		}{
			tee,
			resp.Body,
		},
		DoCache: func() {
			if err := l.CacheRepo.Set(host, name, dest.Bytes()); err != nil {
				log.Fatal(err)
			}
		},
	}, nil
}

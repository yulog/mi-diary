package logic

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
)

func (l *Logic) CacheLogic(ctx context.Context, profile, name string) (CacheOutput, error) {
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

	out, err := l.EmojiRepo.GetByName(ctx, profile, name)
	if err != nil {
		log.Fatal(err)
		return CacheOutput{}, err
	}
	if out.IsSymbol {
		return CacheOutput{}, fmt.Errorf("%s is symbol", name)
	}

	resp, err := http.Get(out.Image)
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

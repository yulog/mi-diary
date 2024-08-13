package color

import (
	"cmp"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"slices"

	"github.com/cenkalti/dominantcolor"
	"github.com/mattn/go-ciede2000"
	_ "golang.org/x/image/webp"
)

var (
	Red       = &color.RGBA{204, 0, 0, 255}
	Orange    = &color.RGBA{251, 148, 11, 255}
	Yellow    = &color.RGBA{255, 255, 0, 255}
	Green     = &color.RGBA{0, 204, 0, 255}
	Bluegreen = &color.RGBA{3, 192, 198, 255}
	Blue      = &color.RGBA{0, 0, 255, 255}
	Purple    = &color.RGBA{118, 44, 167, 255}
	Pink      = &color.RGBA{255, 152, 191, 255}
	White     = &color.RGBA{255, 255, 255, 255}
	Gray      = &color.RGBA{153, 153, 153, 255}
	Black     = &color.RGBA{0, 0, 0, 255}
	Brown     = &color.RGBA{136, 84, 24, 255}
)

var DefinedColors = []*color.RGBA{
	Red,
	Orange,
	Yellow,
	Green,
	Bluegreen,
	Blue,
	Purple,
	Pink,
	White,
	Gray,
	Black,
	Brown,
}

type colorDiff struct {
	color *color.RGBA
	diff  float64
}

// Color は支配色、色の分類を返す
func Color(url string) (string, string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return "", "", err
	}
	dc := dominantcolor.Find(img)
	s := []colorDiff{}
	for _, c := range DefinedColors {
		s = append(s, colorDiff{c, ciede2000.Diff(dc, c)})
	}
	slices.SortStableFunc(s, func(a, b colorDiff) int {
		return cmp.Compare(a.diff, b.diff)
	})
	if len(s) < 1 {
		return "", "", fmt.Errorf("illegal")
	}
	return dominantcolor.Hex(dc), dominantcolor.Hex(*s[0].color), nil
}

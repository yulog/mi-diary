package color

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"

	_ "golang.org/x/image/webp"

	"github.com/yulog/iroiro"
	"github.com/yulog/iroiro/classify"
)

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
	c, err := classify.Classify(img)
	if err != nil {
		return "", "", err
	}
	return iroiro.Hex(c.DominantColor), iroiro.Hex(c.ClassifiedColor), nil
}

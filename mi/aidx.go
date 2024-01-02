package mi

import (
	"strconv"
	"time"
)

const TIME2000 = 946684800000

// ParseAidx はaidxを日時に変換する(UTC)
//
// TODO: そのままだとタイムゾーンによって結果が変わるのでUTCにしたけど、それだと再度変換が必要になるので微妙かも
func ParseAidx(id string) (time.Time, error) {
	// fmt.Println(id[0:8])
	i, _ := strconv.ParseInt(id[0:8], 36, 64)
	// fmt.Println(i + TIME2000)
	t := time.UnixMilli(i + TIME2000).UTC()
	// fmt.Println(t)
	// jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	// fmt.Println(t.In(jst))
	return t, nil
}

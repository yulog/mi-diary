package mi

import (
	"strconv"
	"time"
)

const TIME2000 = 946684800000

// ParseAidx はaid/aidxを日時に変換する
func ParseAidx(id string) (time.Time, error) {
	// fmt.Println(id[0:8])
	i, err := strconv.ParseInt(id[0:8], 36, 64)
	if err != nil {
		return time.Time{}, err
	}
	// fmt.Println(i + TIME2000)
	t := time.UnixMilli(i + TIME2000)
	// fmt.Println(t)
	// jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	// fmt.Println(t.In(jst))
	return t, nil
}

func FormatTime(t time.Time, err error) string {
	return t.Format(time.DateTime)
}

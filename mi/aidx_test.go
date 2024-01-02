package mi

import (
	"reflect"
	"testing"
	"time"
)

func Test_parseAidx(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "",
			args:    args{id: "9nu5oj3jd1y902ml"},
			want:    time.Date(2023, 12, 29, 12, 9, 58, 447_000_000, time.UTC),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseAidx(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseAidx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseAidx() = %v, want %v", got, tt.want)
			}
		})
	}
}

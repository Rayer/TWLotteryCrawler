package TWLotteryCrawer

import (
	_ "embed"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type LotterySL638Suite struct {
	suite.Suite
}

//go:embed test_resources/sl638_history.html
var SL638HistoryPage string

func (l *LotterySL638Suite) SetupSuite() {
	httpmock.Activate()
	httpmock.RegisterResponder("GET", "https://www.taiwanlottery.com.tw/lotto/superlotto638/history.aspx", httpmock.NewStringResponder(200, SL638HistoryPage))
}

func (l *LotterySL638Suite) TearDownSuite() {
	httpmock.DeactivateAndReset()
}

func (l *LotterySL638Suite) TestParseResultSL638() {
	results, err := ParseSL638FromHistoryPage()
	if err != nil {
		l.Fail("Fail to parse SL638 history page!")
	}
	r := results[0]
	//Can't put time.Time, mock one
	t := time.Now()
	r.Date = t
	expectedResult := SuperLotto638Result{
		AZone:       []int{27, 8, 5, 19, 21, 37},
		AZoneSorted: []int{5, 8, 19, 21, 27, 37},
		BZone:       2,
		Serial:      "109000060",
		Date:        t,
		FirstPrice:  1562473476,
		SecondPrice: 3655333,
		RolloverFP:  0,
		RolloverSP:  0,
	}
	assert.Equal(l.T(), expectedResult, r)
}

func (l *LotterySL638Suite) TestSuperLotto638Result_RewardOf() {

	type args struct {
		aZone []int
		bZone int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "supposed be 普獎",
			args: args{
				aZone: []int{1, 2, 3, 4, 5, 6},
				bZone: 2,
			},
			want:    "普獎",
			wantErr: false,
		},
		{
			name: "supposed be 玖獎",
			args: args{
				aZone: []int{1, 2, 3, 21, 27, 37},
				bZone: 1,
			},
			want:    "玖獎",
			wantErr: false,
		},
		{
			name: "supposed be 捌獎",
			args: args{
				aZone: []int{5, 7, 19, 30, 31, 32},
				bZone: 2,
			},
			want:    "捌獎",
			wantErr: false,
		},
		{
			name: "supposed be 柒獎",
			args: args{
				aZone: []int{5, 8, 19, 20, 26, 36},
				bZone: 2,
			},
			want:    "柒獎",
			wantErr: false,
		},
		{
			name: "supposed be 陸獎",
			args: args{
				aZone: []int{5, 8, 19, 21, 26, 36},
				bZone: 3,
			},
			want:    "陸獎",
			wantErr: false,
		},
		{
			name: "supposed be 伍獎",
			args: args{
				aZone: []int{5, 8, 19, 21, 26, 36},
				bZone: 2,
			},
			want:    "伍獎",
			wantErr: false,
		},
		{
			name: "supposed be 肆獎",
			args: args{
				aZone: []int{5, 8, 18, 21, 27, 37},
				bZone: 3,
			},
			want:    "肆獎",
			wantErr: false,
		},
		{
			name: "supposed be 參獎",
			args: args{
				aZone: []int{5, 8, 18, 21, 27, 37},
				bZone: 2,
			},
			want:    "參獎",
			wantErr: false,
		},
		{
			name: "supposed be 貳獎",
			args: args{
				aZone: []int{5, 8, 19, 21, 27, 37},
				bZone: 3,
			},
			want:    "貳獎",
			wantErr: false,
		},
		{
			name: "supposed be 頭獎",
			args: args{
				aZone: []int{5, 8, 19, 21, 27, 37},
				bZone: 2,
			},
			want:    "頭獎",
			wantErr: false,
		},
		{
			name: "supposed be 沒中",
			args: args{
				aZone: []int{4, 7, 18, 20, 26, 36},
				bZone: 2,
			},
			want:    "沒中獎，再接再厲",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		l.Run(tt.name, func() {
			results, err := ParseSL638FromHistoryPage()
			if err != nil {
				l.Fail("Fail to parse SL638 Page")
			}
			result := results[0]
			reward, err := result.RewardOf(tt.args.aZone, tt.args.bZone)
			if (err != nil) != tt.wantErr {
				l.Failf("RewardOf() error = %s, wantErr %v", err.Error(), tt.wantErr)
				return
			}
			got := reward.Title
			if got != tt.want {
				l.Failf("RewardOf() got = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestLotterySL638Suite(t *testing.T) {
	suite.Run(t, new(LotterySL638Suite))
}

package TWLotteryCrawer

import (
	_ "embed"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

//go:embed test_resources/lotto649_history.html
var L649HistoryPage string

type LotteryL649Suite struct {
	suite.Suite
}

func (l *LotteryL649Suite) SetupSuite() {
	httpmock.Activate()
	httpmock.RegisterResponder("GET", "https://www.taiwanlottery.com.tw/lotto/Lotto649/history.aspx", httpmock.NewStringResponder(200, L649HistoryPage))
}

func (l *LotteryL649Suite) TearDownSuite() {
	httpmock.DeactivateAndReset()
}

func (l *LotteryL649Suite) TestParseResultL649() {
	results, err := ParseL649FromHistoryPage()
	if err != nil {
		l.Fail("Fail to parse L649 history page!")
	}
	r := results[0]
	//Can't put time.Time, mock one
	t := time.Now()
	r.Date = t

	expectedResult := L649Result{
		Regular:        []int{36, 1, 41, 17, 40, 5},
		RegularSorted:  []int{1, 5, 17, 36, 40, 41},
		Special:        38,
		Serial:         "110000077",
		Date:           t,
		FirstPrize:     451365562,
		SecondPrize:    2374773,
		ThirdPrize:     57470,
		ForthPrize:     14679,
		RolloverFP:     451365562,
		RolloverSP:     0,
		RolloverTP:     0,
		RolloverForthP: 0,
	}
	assert.Equal(l.T(), expectedResult, r)
}

func (l *LotteryL649Suite) TestSuperLotto638Result_RewardOf() {

	type args struct {
		numbers []int
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
				numbers: []int{1, 2, 3, 4, 5, 17},
			},
			want:    "普獎",
			wantErr: false,
		},
		{
			name: "supposed be 柒獎",
			args: args{
				numbers: []int{5, 8, 17, 20, 26, 38},
			},
			want:    "柒獎",
			wantErr: false,
		},
		{
			name: "supposed be 陸獎",
			args: args{
				numbers: []int{5, 8, 17, 20, 36, 38},
			},
			want:    "陸獎",
			wantErr: false,
		},
		{
			name: "supposed be 伍獎",
			args: args{
				numbers: []int{5, 8, 17, 26, 36, 41},
			},
			want:    "伍獎",
			wantErr: false,
		},
		{
			name: "supposed be 肆獎",
			args: args{
				numbers: []int{5, 17, 36, 38, 39, 41},
			},
			want:    "肆獎",
			wantErr: false,
		},
		{
			name: "supposed be 參獎",
			args: args{
				numbers: []int{1, 5, 16, 36, 40, 41},
			},
			want:    "參獎",
			wantErr: false,
		},
		{
			name: "supposed be 貳獎",
			args: args{
				numbers: []int{5, 17, 36, 38, 40, 41},
			},
			want:    "貳獎",
			wantErr: false,
		},
		{
			name: "supposed be 頭獎",
			args: args{
				numbers: []int{1, 5, 17, 36, 40, 41},
			},
			want:    "頭獎",
			wantErr: false,
		},
		{
			name: "supposed be 沒中",
			args: args{
				numbers: []int{4, 7, 18, 20, 26, 36},
			},
			want:    "沒中獎，再接再厲",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		l.Run(tt.name, func() {
			results, err := ParseL649FromHistoryPage()
			if err != nil {
				l.Fail("Fail to parse SL638 Page")
			}
			result := results[0]
			reward, err := result.RewardOf(tt.args.numbers)
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

func TestLotteryL649Suite(t *testing.T) {
	suite.Run(t, new(LotteryL649Suite))
}

package TWLotteryCrawer

import (
	"fmt"
	"github.com/jarcoal/httpmock"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"testing"
	"time"
)

type LotteryContextTestSuite struct {
	suite.Suite
	context *LotteryContext
	result  *LotteryResult
}

func (l *LotteryContextTestSuite) SetupSuite() {
	l.context = &LotteryContext{}
	httpmock.Activate()
	content, err := ioutil.ReadFile("./test_resources/lottery.html")
	if err != nil {
		fmt.Println(err.Error())
	}
	httpmock.RegisterResponder("GET", "https://www.taiwanlottery.com.tw/index_new.aspx", httpmock.NewStringResponder(200, string(content)))
	httpmock.RegisterResponder("GET", "http://210.71.254.181/index_new.htm", httpmock.NewStringResponder(200, string(content)))
	result, err := l.context.Fetch()
	if err != nil {
		l.Error(err, err.Error())
		l.FailNow("Error in parsing test html file!")
	}
	l.result = result
}

func (l *LotteryContextTestSuite) TearDownSuite() {
	httpmock.DeactivateAndReset()
}

func (l *LotteryContextTestSuite) TestParseResultSL638() {
	//Can't put time.Time, mock one
	t := time.Now()
	r := l.result.SuperLotto638Result
	r.Date = t
	expectedResult := SuperLotto638Result{
		AZone:       []int{35, 17, 9, 1, 18, 31},
		AZoneSorted: []int{1, 9, 17, 18, 31, 35},
		BZone:       2,
		Serial:      "109000059",
		Date:        t,
	}
	assert.Equal(l.T(), *r, expectedResult)
}

func (l *LotteryContextTestSuite) TestSuperLotto638Result_RewardOf() {

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
				aZone: []int{1, 2, 3, 4, 31, 35},
				bZone: 1,
			},
			want:    "玖獎",
			wantErr: false,
		},
		{
			name: "supposed be 捌獎",
			args: args{
				aZone: []int{1, 2, 3, 4, 31, 6},
				bZone: 2,
			},
			want:    "捌獎",
			wantErr: false,
		},
		{
			name: "supposed be 柒獎",
			args: args{
				aZone: []int{1, 2, 3, 4, 9, 17},
				bZone: 2,
			},
			want:    "柒獎",
			wantErr: false,
		},
		{
			name: "supposed be 陸獎",
			args: args{
				aZone: []int{1, 2, 3, 9, 17, 18},
				bZone: 3,
			},
			want:    "陸獎",
			wantErr: false,
		},
		{
			name: "supposed be 伍獎",
			args: args{
				aZone: []int{1, 2, 3, 9, 17, 18},
				bZone: 2,
			},
			want:    "伍獎",
			wantErr: false,
		},
		{
			name: "supposed be 肆獎",
			args: args{
				aZone: []int{1, 2, 17, 18, 31, 35},
				bZone: 3,
			},
			want:    "肆獎",
			wantErr: false,
		},
		{
			name: "supposed be 參獎",
			args: args{
				aZone: []int{1, 2, 17, 18, 31, 35},
				bZone: 2,
			},
			want:    "參獎",
			wantErr: false,
		},
		{
			name: "supposed be 貳獎",
			args: args{
				aZone: []int{1, 9, 17, 18, 31, 35},
				bZone: 3,
			},
			want:    "貳獎",
			wantErr: false,
		},
		{
			name: "supposed be 頭獎",
			args: args{
				aZone: []int{1, 9, 17, 18, 31, 35},
				bZone: 2,
			},
			want:    "頭獎",
			wantErr: false,
		},
		{
			name: "supposed be 沒中",
			args: args{
				aZone: []int{2, 4, 6, 8, 10, 12},
				bZone: 2,
			},
			want:    "沒中獎，再接再厲",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		l.Run(tt.name, func() {
			s := l.result.SuperLotto638Result
			result, err := s.RewardOf(tt.args.aZone, tt.args.bZone)
			if (err != nil) != tt.wantErr {
				l.Failf("RewardOf() error = %s, wantErr %v", err.Error(), tt.wantErr)
				return
			}
			got := result.Title
			if got != tt.want {
				l.Failf("RewardOf() got = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestLotteryContextSuite(t *testing.T) {
	suite.Run(t, new(LotteryContextTestSuite))
}

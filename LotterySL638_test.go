package TWLotteryCrawer

import (
	"fmt"
	"github.com/jarcoal/httpmock"
	"io/ioutil"
	"testing"
)

func TestLotteryContext_parseSL638FromHistoryPage(t *testing.T) {

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	content, err := ioutil.ReadFile("./test_resources/sl638_history.html")
	if err != nil {
		t.Error(err.Error())
	}
	httpmock.RegisterResponder("GET", "https://www.taiwanlottery.com.tw/lotto/superlotto638/history.aspx", httpmock.NewStringResponder(200, string(content)))

	l := &LotteryContext{}
	fmt.Println(l.ParseSL638FromHistoryPage())
}

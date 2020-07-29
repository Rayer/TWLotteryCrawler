package TWLotteryCrawer

import (
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
	httpmock.RegisterResponder("GET", "http://210.71.254.181/lotto/superlotto638/history.htm", httpmock.NewStringResponder(200, string(content)))

	l := &LotteryContext{}
	_, _ = l.parseSL638FromHistoryPage()
}
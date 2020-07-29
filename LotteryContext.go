package TWLotteryCrawer

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
)

type LotteryContext struct {
}

//All
type LotteryResult struct {
	SuperLotto638Result *SuperLotto638Result
}

func (l *LotteryContext) Fetch() (*LotteryResult, error) {
	url := "https://www.taiwanlottery.com.tw/index_new.aspx"
	//url := "http://210.71.254.181/index_new.htm"
	//Handle redirect
	//for ex :
	//<meta http-equiv="refresh" content="0; URL=http://210.71.254.181/index_new.htm" />
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	//Parse SuperLotto638
	superLotto638Result, err := l.parseSuperLotto638(doc)

	return &LotteryResult{
		SuperLotto638Result: superLotto638Result,
	}, err
}

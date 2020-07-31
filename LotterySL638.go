package TWLotteryCrawer

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//以下的Struct name都取自於台彩source code的命名...
//威力彩
type SuperLotto638Result struct {
	AZone       []int
	AZoneSorted []int
	BZone       int
	Serial      string
	Date        time.Time
	FirstPrice  int
	SecondPrice int
	RolloverFP  int
	RolloverSP  int
}

type SuperLottto638Reward struct {
	Reward      int
	Description string
	Title       string
}

func (s *SuperLotto638Result) RewardOf(aZone []int, bZone int) (*SuperLottto638Reward, error) {

	//Check aZone is unique and correct length
	aMap := make(map[int]int)
	for i, ia := range aZone {
		aMap[ia] = i
	}

	if len(aMap) != 6 {
		return nil, errors.New("In A Zone number count should be 6 and unique!")
	}

	aZoneCount := 0
	for _, ia := range aZone {
		for _, a := range s.AZone {
			if ia == a {
				aZoneCount++
			}
		}
	}

	var reward *SuperLottto638Reward
	if bZone == s.BZone {
		switch aZoneCount {
		case 1:
			reward = &SuperLottto638Reward{
				Reward:      100,
				Description: "第1區任1個+第2區",
				Title:       "普獎",
			}
		case 2:
			reward = &SuperLottto638Reward{
				Reward:      200,
				Description: "第1區任2個+第2區",
				Title:       "捌獎",
			}
		case 3:
			reward = &SuperLottto638Reward{
				Reward:      400,
				Description: "第1區任3個+第2區",
				Title:       "柒獎",
			}
		case 4:
			reward = &SuperLottto638Reward{
				Reward:      4000,
				Description: "第1區任4個+第2區",
				Title:       "伍獎",
			}
		case 5:
			reward = &SuperLottto638Reward{
				Reward:      150000,
				Description: "第1區任5個+第2區",
				Title:       "參獎",
			}
		case 6:
			reward = &SuperLottto638Reward{
				Reward:      s.FirstPrice,
				Description: "第1區6個+第2區",
				Title:       "頭獎",
			}
		}
	} else {
		switch aZoneCount {
		case 3:
			reward = &SuperLottto638Reward{
				Reward:      100,
				Description: "第1區任3個",
				Title:       "玖獎",
			}
		case 4:
			reward = &SuperLottto638Reward{
				Reward:      800,
				Description: "第1區任4個",
				Title:       "陸獎",
			}
		case 5:
			reward = &SuperLottto638Reward{
				Reward:      20000,
				Description: "第1區任5個",
				Title:       "肆獎",
			}
		case 6:
			reward = &SuperLottto638Reward{
				Reward:      s.SecondPrice,
				Description: "第1區任6個",
				Title:       "貳獎",
			}
		}
	}

	if reward == nil {
		reward = &SuperLottto638Reward{
			Reward:      0,
			Description: "...",
			Title:       "沒中獎，再接再厲",
		}
	}

	return reward, nil
}

func (l *LotteryContext) ParseSuperLotto638(doc *goquery.Document) (*SuperLotto638Result, error) {
	iconNode := doc.Find("div#contents_logo_02")
	if len(iconNode.Nodes) != 1 {
		return nil, errors.New(fmt.Sprintf("Wrong SuperLotto638 icon node : %d", len(iconNode.Nodes)))
	}
	parent := iconNode.Parent()
	aZoneBalls := make([]int, 6)

	var err error
	for i, a := range strings.Split(parent.Find("div.ball_tx").Text(), " ")[0:6] {
		aZoneBalls[i], err = strconv.Atoi(a)
		if err != nil {
			return nil, errors.Wrap(err, "Error while parsing AZone balls")
		}
	}
	aZoneBallsSorted := make([]int, 6)
	for i, a := range strings.Split(parent.Find("div.ball_tx").Text(), " ")[6:12] {
		aZoneBallsSorted[i], err = strconv.Atoi(a)
		if err != nil {
			return nil, errors.Wrap(err, "Error while parsing AZone Sorted balls")
		}
	}
	bZoneBall, err := strconv.Atoi(strings.Trim(parent.Find("div.ball_red").Text(), " "))

	if err != nil {
		return nil, errors.Wrap(err, "Error while parsing BZone balls")
	}
	//Parse date and serial
	//109/7/23&nbsp;第109000059期
	var serial string
	var date time.Time
	dateSerialString := strings.Trim(parent.Find("div.contents_mine_tx02").Find("span.font_black15").Text(), " ")
	r := regexp.MustCompile("(\\d+)\\/(\\d+)\\/(\\d+).第(\\d+)期")
	if find := r.FindStringSubmatch(dateSerialString); len(find) > 1 {
		year, _ := strconv.Atoi(find[1])
		year += 1911
		find[1] = strconv.Itoa(year)
		date, _ = time.Parse("2006 1 2", strings.Join(find[1:4], " "))
		serial = find[4]
	}

	sl638result := &SuperLotto638Result{
		AZone:       aZoneBalls,
		AZoneSorted: aZoneBallsSorted,
		BZone:       bZoneBall,
		Serial:      serial,
		Date:        date,
	}

	return sl638result, nil
}

func (l *LotteryContext) ParseSL638FromHistoryPage() ([]SuperLotto638Result, error) {
	//url := "http://210.71.254.181/lotto/superlotto638/history.htm"
	url := "https://www.taiwanlottery.com.tw/lotto/superlotto638/history.aspx"
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	ret := make([]SuperLotto638Result, 0)
	doc.Find("table#SuperLotto638Control_history1_dlQuery").Find("table.td_hm").Each(func(index int, selection *goquery.Selection) {
		td := selection.Find("table tr td")
		fmt.Printf("Index : %d\n", index)
		fieldTargets := strings.Fields(strings.Replace(td.Text(), "\n", "", -1))
		fmt.Println(fieldTargets)
		/**Per table, it will looks like :
		[期別 開獎日 兌獎截止(註5) 銷售金額 獎金總額 						(0-4) titles
		 109000060 109/07/27 109/10/27 1,999,389,700 3,633,111,827  (5-9) 期數，開獎日期，兌獎期限，本期銷售，獎金總額
		 獎號 第一區 第二區 開出順序 27 								(10-14) Title, A區1
		 08 05 19 21 37												(15-19) A區2-6
		 02 大小順序 05 08 19										(20-24) B區，title, 排序過的A區1-3
		 21 27 37 02 獎金分配										(25-29) 排序過的A區4-6，B區，Title
		 項目 頭獎 貳獎 參獎 肆獎										(30-34) Title
		 伍獎 陸獎 柒獎 捌獎 玖獎 										(35-59) Title
		 普獎 對中 獎號數 第一區6個 ＋第二區								(40-44) Title
		 第一區6個 第一區5個 ＋第二區 第一區5個 第一區4個 					(45-49) Title
		 ＋第二區 第一區4個 第一區3個 ＋第二區 第一區2個 					(50-54) Title
		 ＋第二區 第一區3個 第一區1個 ＋第二區 中獎 						(55-60) Title
		 注數 2 20 205 1,535 										(61-64) 中獎注數1-4
		 7,106 54,898 88,989 471,208 671,352 						(65-69) 中獎注數5-9
		 1,042,934 單注 獎金 1,562,473,476 3,655,333 				(70-74) 中獎注數10, Title, 獎金1-2
		 150,000 20,000 4,000 800 400								(75-79) 獎金3-7
		 200 100 100 累至次 期獎金									(80-84) 獎金8-10, title
		 0 0]														(85-  ) 累積至下期頭獎, 累積至下期貳獎
		*/

		aZone := make([]int, 0)
		aZoneSorted := make([]int, 0)
		for _, v := range fieldTargets[14:20] {
			d, err := strconv.Atoi(v)
			if err != nil {
				logrus.Warnf("Fail to parse aZone : %+v", v)
				return
			}
			aZone = append(aZone, d)
		}
		for _, v := range fieldTargets[22:28] {
			d, err := strconv.Atoi(v)
			if err != nil {
				logrus.Warnf("Fail to parse aZone : %+v", v)
				return
			}
			aZoneSorted = append(aZoneSorted, d)
		}

		bZone, err := strconv.Atoi(fieldTargets[28])

		if err != nil {
			logrus.Warnf("Fail to parse bZone : %+v", fieldTargets[28])
			return
		}

		firstPrice, err := strconv.Atoi(strings.Replace(fieldTargets[73], ",", "", -1))
		if err != nil {
			logrus.Warnf("Fail to parse firstPrice : %+v", fieldTargets[73])
			return
		}

		secondPrice, err := strconv.Atoi(strings.Replace(fieldTargets[74],",", "", -1))
		if err != nil {
			logrus.Warnf("Fail to parse secondPrice : %+v", fieldTargets[74])
			return
		}

		firstPriceRollover, err := strconv.Atoi(strings.Replace(fieldTargets[85],",", "", -1))
		if err != nil {
			logrus.Warnf("Fail to parse firstPriceRollover : %+v", fieldTargets[85])
			return
		}

		secondPriceRollover, err := strconv.Atoi(strings.Replace(fieldTargets[86],",", "", -1))
		if err != nil {
			logrus.Warnf("Fail to parse secondPriceRollover : %+v", fieldTargets[86])
			return
		}

		var resultDate time.Time
		r := regexp.MustCompile("(\\d+)\\/(\\d+)\\/(\\d+)")
		if find := r.FindStringSubmatch(fieldTargets[6]); len(find) > 1 {
			year, _ := strconv.Atoi(find[1])
			year += 1911
			find[1] = strconv.Itoa(year)
			resultDate, _ = time.Parse("2006 1 2", strings.Join(find[1:4], " "))
		}


		ret = append(ret, SuperLotto638Result{
			AZone:       aZone,
			AZoneSorted: aZoneSorted,
			BZone:       bZone,
			Serial:      fieldTargets[5],
			Date:        resultDate,
			FirstPrice:  firstPrice,
			SecondPrice: secondPrice,
			RolloverFP:  firstPriceRollover,
			RolloverSP:  secondPriceRollover,
		})

	})
	fmt.Printf("%+v", ret)
	return ret, nil
}

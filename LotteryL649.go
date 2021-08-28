package TWLotteryCrawer

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type L649Result struct {
	Regular        []int
	RegularSorted  []int
	Special        int
	Serial         string
	Date           time.Time
	FirstPrize     int
	SecondPrize    int
	ThirdPrize     int
	ForthPrize     int
	RolloverFP     int
	RolloverSP     int
	RolloverTP     int
	RolloverForthP int
}

type L649Reward struct {
	Reward      int
	Description string
	Title       string
}

func (l *L649Result) RewardOf(numbers []int) (*L649Reward, error) {
	haveSpecialNum := false
	numInRegular := 0

	if len(numbers) != 6 {
		return nil, errors.New("input number should exactly be 6")
	}

	for _, n := range numbers {
		if n == l.Special {
			haveSpecialNum = true
			continue
		}
		for _, r := range l.Regular {
			if n == r {
				numInRegular += 1
			}
		}
	}

	var ret *L649Reward
	if haveSpecialNum {
		switch numInRegular {
		case 2:
			ret = &L649Reward{
				Reward:      400,
				Description: "任兩個+特別號",
				Title:       "柒獎",
			}
		case 3:
			ret = &L649Reward{
				Reward:      1000,
				Description: "任三個+特別號",
				Title:       "陸獎",
			}
		case 4:
			ret = &L649Reward{
				Reward:      l.ForthPrize,
				Description: "任四個+特別號",
				Title:       "肆獎",
			}
		case 5:
			ret = &L649Reward{
				Reward:      l.SecondPrize,
				Description: "任五個+特別號",
				Title:       "貳獎",
			}
		}
	} else {
		switch numInRegular {
		case 3:
			ret = &L649Reward{
				Reward:      400,
				Description: "任三個",
				Title:       "普獎",
			}

		case 4:
			ret = &L649Reward{
				Reward:      2000,
				Description: "任四個",
				Title:       "伍獎",
			}

		case 5:
			ret = &L649Reward{
				Reward:      l.ThirdPrize,
				Description: "任五個",
				Title:       "參獎",
			}

		case 6:
			ret = &L649Reward{
				Reward:      l.FirstPrize,
				Description: "六個全中",
				Title:       "頭獎",
			}
		}
	}

	if ret == nil {
		ret = &L649Reward{
			Reward:      0,
			Description: "沒中",
			Title:       "沒中獎，再接再厲",
		}
	}

	return ret, nil
}

func ParseL649FromHistoryPage() ([]L649Result, error) {
	url := "https://www.taiwanlottery.com.tw/lotto/Lotto649/history.aspx"
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = res.Body.Close()
	}()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	ret := make([]L649Result, 0)
	doc.Find("table#Lotto649Control_history_dlQuery").Find("table.td_hm").Each(func(index int, selection *goquery.Selection) {
		td := selection.Find("table tr td")
		var fieldTargets []string
		td.Each(func(i int, selection *goquery.Selection) {
			s := selection.Text()
			s = strings.Replace(s, "\n", "", -1)
			s = strings.Trim(s, " ")
			fieldTargets = append(fieldTargets, s)
		})
		/**Per table, it will looks like :
		期別 開獎日 兌獎截止(註6) 銷售金額 獎金總額 						(0-4) titles
		110000077 110/08/24 110/11/24 265,428,850 540,088,355	  	(5-9) 期數，開獎日期，兌獎期限，本期銷售，獎金總額
		獎號 特別號 開出順序 36 01	 								(10-14) Title, 獎號1-2
		41 17 40 05 38												(15-19) 獎號3-6, 特別號
		大小順序 01 05 17 36										(20-24) Title，排序過的獎號1-4
		40 41 38 獎金分配 項目										(25-29) 排序過的獎號5-6，特別號，Title
		頭獎 貳獎 參獎 肆獎 伍獎										(30-34) Title
		陸獎 柒獎 普獎 對中獎號數 6個 									(35-39) Title
		任5個＋特別號 任5個 任4個＋特別號 任4個 任3個＋特別號				(40-44) Title
		任2個＋特別號 任3個 中獎注數 0 2			 					(45-49) Title, 頭獎注數, 二獎注數
		89 224 4,616 6,653 65,942				 					(50-54) 三四五六七獎注數
		83,606 單注獎金 0 2,374,773 57,470	 						(55-59) 普獎注數，title，頭二三獎金
		14,679 2,000 1,000 400 400									(60-64) 四五六七普獎金
		累至次期獎金 451,365,562 0 0 0		 						(65-69) title 累計次期頭二三四獎金
		*/

		prizeNum := make([]int, 0)
		prizeNumSorted := make([]int, 0)
		for _, v := range fieldTargets[13:19] {
			d, err := strconv.Atoi(v)
			if err != nil {
				logrus.Warnf("Fail to parse prizeNum : %+v", v)
				return
			}
			prizeNum = append(prizeNum, d)
		}
		for _, v := range fieldTargets[21:27] {
			d, err := strconv.Atoi(v)
			if err != nil {
				logrus.Warnf("Fail to parse prizeNumSorted : %+v", v)
				return
			}
			prizeNumSorted = append(prizeNumSorted, d)
		}

		specialNum, err := strconv.Atoi(fieldTargets[19])

		if err != nil {
			logrus.Warnf("Fail to parse specialNum : %+v", fieldTargets[19])
			return
		}

		firstPrize, err := strconv.Atoi(strings.Replace(fieldTargets[57], ",", "", -1))

		if err != nil {
			logrus.Warnf("Fail to parse firstPrize : %+v", fieldTargets[57])
			return
		}

		secondPrize, err := strconv.Atoi(strings.Replace(fieldTargets[58], ",", "", -1))
		if err != nil {
			logrus.Warnf("Fail to parse secondPrize : %+v", fieldTargets[58])
			return
		}

		thirdPrize, err := strconv.Atoi(strings.Replace(fieldTargets[59], ",", "", -1))
		if err != nil {
			logrus.Warnf("Fail to parse thirdPrize : %+v", fieldTargets[59])
			return
		}

		forthPrize, err := strconv.Atoi(strings.Replace(fieldTargets[60], ",", "", -1))
		if err != nil {
			logrus.Warnf("Fail to parse secondPrize : %+v", fieldTargets[60])
			return
		}

		firstPrizeRollover, err := strconv.Atoi(strings.Replace(fieldTargets[66], ",", "", -1))
		if err != nil {
			logrus.Warnf("Fail to parse firstPrizeRollover : %+v", fieldTargets[66])
			return
		}

		if firstPrize == 0 {
			firstPrize = firstPrizeRollover
		}

		secondPrizeRollover, err := strconv.Atoi(strings.Replace(fieldTargets[67], ",", "", -1))
		if err != nil {
			logrus.Warnf("Fail to parse secondPrizeRollover : %+v", fieldTargets[67])
			return
		}

		if secondPrize == 0 {
			secondPrize = secondPrizeRollover
		}

		thirdPrizeRollover, err := strconv.Atoi(strings.Replace(fieldTargets[68], ",", "", -1))
		if err != nil {
			logrus.Warnf("Fail to parse thirdPrizeRollover : %+v", fieldTargets[68])
			return
		}

		if thirdPrize == 0 {
			thirdPrize = thirdPrizeRollover
		}

		forthPrizeRollover, err := strconv.Atoi(strings.Replace(fieldTargets[69], ",", "", -1))
		if err != nil {
			logrus.Warnf("Fail to parse forthPrizeRollover : %+v", fieldTargets[69])
			return
		}

		if forthPrize == 0 {
			forthPrize = forthPrizeRollover
		}

		var resultDate time.Time
		r := regexp.MustCompile("(\\d+)\\/(\\d+)\\/(\\d+)")
		if find := r.FindStringSubmatch(fieldTargets[6]); len(find) > 1 {
			year, _ := strconv.Atoi(find[1])
			year += 1911
			find[1] = strconv.Itoa(year)
			resultDate, _ = time.Parse("2006 1 2", strings.Join(find[1:4], " "))
		}

		ret = append(ret, L649Result{
			Regular:        prizeNum,
			RegularSorted:  prizeNumSorted,
			Special:        specialNum,
			Serial:         fieldTargets[5],
			Date:           resultDate,
			FirstPrize:     firstPrize,
			SecondPrize:    secondPrize,
			ThirdPrize:     thirdPrize,
			ForthPrize:     forthPrize,
			RolloverFP:     firstPrizeRollover,
			RolloverSP:     secondPrizeRollover,
			RolloverTP:     thirdPrizeRollover,
			RolloverForthP: forthPrizeRollover,
		})

	})
	return ret, nil
}

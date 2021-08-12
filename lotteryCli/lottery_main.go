package main

import (
	"bufio"
	"fmt"
	Lottery "github.com/Rayer/TWLotteryCrawler"
	"os"
	"strconv"
	"strings"
)

func main() {
	l := Lottery.LotteryContext{}
	fmt.Println("Fetching data from server...")
	res, _ := l.ParseSL638FromHistoryPage()

	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("一共取得%d期資料:\n", len(res))
	for i, v := range res {
		fmt.Printf("[%d]\t%s期\t(%s)\n", i, v.Serial, v.Date.Format("2006-1-2"))
	}

	var sl Lottery.SuperLotto638Result
	fmt.Println("請選擇期數[0] : ")
	for {
		input, _ := reader.ReadString('\n')
		input = strings.Trim(input, "\n")
		if input == "" {
			input = "0"
		}
		idx, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Wrong input!")
			continue
		}
		if idx > len(res) {
			fmt.Printf("Out of bound!")
			continue
		}
		sl = res[idx]
		break
	}

	fmt.Println("Data fetched! ")
	fmt.Printf("本期樂透第%s期:\n開獎日期:%s\nA區:\t\t%v\n排序後A區:\t%v\nB區:%d\n\n", sl.Serial, sl.Date.Format("2006/1/2"), sl.AZone, sl.AZoneSorted, sl.BZone)

	for {
		fmt.Println("請輸入號碼，以空白隔開。前六個數字為Ａ區，最後一個數字為Ｂ區，共七組:")
		text, _ := reader.ReadString('\n')
		text = strings.Trim(text, "\n")
		text_slice := strings.Split(text, " ")
		if len(text_slice) == 0 {
			break
		}
		if len(text_slice) != 7 {
			fmt.Println("輸入錯誤，請重新輸入")
			continue
		}
		aZone := make([]int, 6)
		parseError := false
		for i, t := range text_slice[0:6] {
			parsed, err := strconv.Atoi(t)
			if err != nil {
				parseError = true
				break
			}
			aZone[i] = parsed
		}
		bZone, err := strconv.Atoi(text_slice[6])
		if err != nil || parseError {
			fmt.Println("輸入錯誤，請重新輸入")
			continue
		}
		reward, err := sl.RewardOf(aZone, bZone)
		if err != nil {
			fmt.Printf("Error occured : %s", err.Error())
			break
		}
		if reward.Reward != 0 {
			fmt.Printf("中獎了!(%s) 獎金:%d, 獎項: %s\n", reward.Description, reward.Reward, reward.Title)
		} else {
			fmt.Println("沒中獎....下一張!")
		}
	}

}

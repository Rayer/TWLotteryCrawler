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
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("[0] 威力彩")
	fmt.Println("[1] 大樂透")
	fmt.Print("請選擇服務 [0] : ")
	input, _ := reader.ReadString('\n')
	if input == "1" {
		execL649(reader)
	} else {
		execSL638(reader)
	}
}

func execSL638(reader *bufio.Reader) {
	fmt.Println("Fetching data from server...")
	res, _ := Lottery.ParseSL638FromHistoryPage()

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
	fmt.Printf("本期威力彩第%s期:\n開獎日期:%s\nA區:\t\t%v\n排序後A區:\t%v\nB區:%d\n\n", sl.Serial, sl.Date.Format("2006/1/2"), sl.AZone, sl.AZoneSorted, sl.BZone)

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

func execL649(reader *bufio.Reader) {
	fmt.Println("Fetching data from server...")
	res, _ := Lottery.ParseL649FromHistoryPage()

	fmt.Printf("一共取得%d期資料:\n", len(res))
	for i, v := range res {
		fmt.Printf("[%d]\t%s期\t(%s)\n", i, v.Serial, v.Date.Format("2006-1-2"))
	}

	var sl Lottery.L649Result
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
	fmt.Printf("本期樂透第%s期:\n開獎日期:%s\n開出號碼:\t%v\n排序後號碼:\t%v\n特別號:%d\n\n", sl.Serial, sl.Date.Format("2006/1/2"), sl.Regular, sl.RegularSorted, sl.Special)

	for {
		fmt.Println("請輸入號碼，以空白隔開，共六組 : ")
		text, _ := reader.ReadString('\n')
		text = strings.Trim(text, "\n")
		textSlice := strings.Split(text, " ")
		if len(textSlice) == 0 {
			break
		}
		if len(textSlice) != 6 {
			fmt.Println("輸入錯誤，請重新輸入")
			continue
		}
		numbers := make([]int, 6)
		parseError := false
		for i, t := range textSlice[0:6] {
			parsed, err := strconv.Atoi(t)
			if err != nil {
				parseError = true
				break
			}
			numbers[i] = parsed
		}
		if parseError {
			fmt.Println("輸入錯誤，請重新輸入")
			continue
		}
		reward, err := sl.RewardOf(numbers)
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

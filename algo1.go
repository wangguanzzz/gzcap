package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

var balance float32

const start float32 = 10000.0

func main() {
	files := os.Args[1:]
	acc := 0.0
	tradenum := 0

	for _, file := range files {
		result, num := dailyAlgo(file)
		acc += result
		tradenum += num
	}

	fmt.Println("total: ", acc, tradenum)
}

func dailyAlgo(filename string) (result float64, num int) {
	stage1 := 0.0
	stage2 := 0.0
	index := 0
	hold := 0.0
	balance := 0.0

	// 0 sell 1 buy
	direction := 0

	accumulate := 0.0

	tradetime := 0
	ok := true

	fs, err := os.Open(filename)
	checkError("cannot read file", err)
	defer fs.Close()
	r := csv.NewReader(fs)
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil && err != io.EOF {
			fmt.Println("reading error")
		}
		//date,hour,minute,begin,highest,lowest,price,volume,hold
		date := row[0]
		hour := row[1]
		minute := row[2]
		currentPrice, _ := strconv.ParseFloat(row[6], 64)
		//volume := strconv.ParseInt(row[7])

		if index != 0 && index != 1 {
			//algo part
			// buy
			if ok && (currentPrice > stage1) && (stage1 > stage2) {
				hold = currentPrice
				direction = 1
				tradetime += 1
				fmt.Println("buy ", "time: ", date, hour, minute, " price:", currentPrice, " stage1:", stage1, " stage2", stage2)
				ok = false
			}
			// sell
			if ok && currentPrice < stage1 && stage1 < stage2 {
				hold = currentPrice
				direction = 0
				tradetime += 1
				fmt.Println("sell ", "time: ", date, hour, minute, " price:", currentPrice, " stage1:", stage1, " stage2", stage2)
				ok = false
			}
			// calculate position
			if !ok {
				if direction == 1 {
					balance = float64(currentPrice) - hold
				} else {
					balance = hold - float64(currentPrice)
				}
				// sell strategy
				// lost 0.5%
				threshold := -(hold * 0.005)
				if balance < threshold/2 {
					accumulate += balance
					ok = true
					fmt.Println("lose settle down ", balance)
				}
				// win 0.5%
				if balance > -threshold/2 {
					accumulate += balance
					ok = true
					fmt.Println("win settle down", balance)
				}
				//end day close
				if hour == "14" && minute == "59" {
					accumulate += balance
					fmt.Println("end day settle down", balance)
				}

			}

			//----
			stage2 = stage1
			stage1 = currentPrice
		}
		if index == 0 {
			stage2 = currentPrice
		}
		if index == 1 {
			stage1 = currentPrice
		}

		index++
	}
	return accumulate, tradetime

}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

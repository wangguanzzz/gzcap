package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
)

func main() {
	files := os.Args[1:]
	acc := 0.0
	tradenum := 0

	days := 0
	dailyReturn := []float64{}

	// output result
	outputfile, err := os.Create("result/result.csv")
	checkError("cannot write file", err)
	defer outputfile.Close()

	writer := csv.NewWriter(outputfile)
	defer writer.Flush()

	for _, file := range files {
		result, num := dailyAlgo(file, writer)
		acc += result
		tradenum += num

		days++
		dailyReturn = append(dailyReturn, result)
	}

	averageReturn := acc / float64(days)

	fmt.Println("days: ", days, " total: ", acc, " trade times: ", tradenum)
	fmt.Println("average daily return: ", averageReturn)

	sharpe := 0.0
	for _, dr := range dailyReturn {
		sharpe += math.Pow(dr-averageReturn, 2.0)
	}
	sharpe = averageReturn / math.Sqrt(252*sharpe)

	fmt.Println("sharpe ratio:", sharpe)
}

func dailyAlgo(filename string, writer *csv.Writer) (float64, int) {
	lastprice := 0.0
	index := 0
	hold := 0.0
	balance := 0.0
	date := ""
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
		date = row[0]
		hour := row[1]
		minute := row[2]
		currentPrice, _ := strconv.ParseFloat(row[6], 64)
		//volume := strconv.ParseInt(row[7])

		if index != 0 {
			//algo part
			// buy
			if ok && ((currentPrice-lastprice)/lastprice > 0.005) {
				hold = currentPrice
				direction = 0
				tradetime += 1
				// fmt.Println("buy ", "time: ", date, hour, minute, " price:", currentPrice, " last:", lastprice)
				ok = false
			}
			// sell
			if ok && ((lastprice-currentPrice)/lastprice > 0.005) {
				hold = currentPrice
				direction = 1
				tradetime += 1
				// fmt.Println("sell ", "time: ", date, hour, minute, " price:", currentPrice, " last:", lastprice)
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

				threshold := -(hold * 0.005)
				// lost 0.5%
				if balance < threshold {
					accumulate += balance
					ok = true
					// fmt.Println("lose settle down ", balance)
				}
				// win 0.5%
				if balance > -threshold {
					accumulate += balance
					ok = true
					// fmt.Println("win settle down", balance)
				}
				//end day close
				if hour == "14" && minute == "59" {
					accumulate += balance
					// fmt.Println("end day settle down", balance)
				}

			}

			//----
			lastprice = currentPrice
		}
		if index == 0 {
			lastprice = currentPrice
		}

		index++
	}

	tempresult := (accumulate / lastprice) * 10.0
	writer.Write([]string{date, fmt.Sprintf("%f", tempresult)})
	return tempresult, tradetime

}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

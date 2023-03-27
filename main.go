package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/emersion/go-ical"
)

func isWeekend(t time.Time) bool {
	// 判断 t 所在的日期是否是周末
	weekday := t.Weekday()
	if weekday == time.Saturday || weekday == time.Sunday {
		return true
	}
	return false
}

func isHoliday(filenname string, date time.Time) bool {
	// 打开ICS文件
	file, err := os.Open(filenname)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	dec := ical.NewDecoder(file)
	for {
		cal, err := dec.Decode()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		for _, event := range cal.Events() {
			start, _ := event.DateTimeStart(time.Now().Location())
			end, _ := event.DateTimeEnd(time.Now().Location())
			// 如果在假期内, 则表示是Holiday
			if start.Before(date) && end.After(date) {
				return true
			}
		}
	}

	if isWeekend(date) {
		return true
	}
	return false
}

var (
	timeString *string
	debug      *bool
)

func init() {
	timeString = flag.String("time", time.Now().Format(time.DateOnly), "time")
	debug = flag.Bool("debug", false, "debug")
}
func main() {
	flag.Parse()
	date, err := time.Parse(time.DateOnly, *timeString)
	if err != nil {
		log.Fatal(err)
		return
	}
	if *debug {
		fmt.Printf("[DEBUG]: date: %v\n", date.Format(time.DateOnly))
	}
	if isHoliday("cal.ics", date) {
		fmt.Printf("holiday\n")
	} else {
		fmt.Printf("workday\n")
	}
	return
}

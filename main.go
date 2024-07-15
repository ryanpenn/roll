package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	// dropList := []*Drop{}
	// GetCsvUtilMgr().LoadCsv("data/Drop", dropList)

	var times int
	flag.IntVar(&times, "t", 0, "-t=10")
	flag.Parse()

	if times <= 0 {
		times = 1
	} else if times > 100_000_000 {
		times = 100_000_000
		fmt.Println("最大抽卡次数不超过", times)
	}

	begin := time.Now()
	Roll(times)
	fmt.Printf("共耗时%dms\n", time.Since(begin).Milliseconds())
}

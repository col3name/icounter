package main

import (
	"flag"
	"fmt"
	"github.com/col3name/ip-unique-addr/pkg/service"
	"runtime"
	"time"
)

func main() {
	fileName := flag.String("f", "ip_addresses_out", "input file1")
	countParallelTask := flag.Int("n", runtime.NumCPU(), "count parallel reader")
	flag.Parse()

	fmt.Println(*fileName)
	start := time.Now()
	filter := service.NewUniqueCounterHLL(*fileName, *countParallelTask)
	count, err := filter.Count()
	elapsed := time.Since(start).String()
	if err != nil {
		fmt.Println(elapsed, err)
		return
	}
	fmt.Println("Time taken", elapsed, "\nUnique", count)
}

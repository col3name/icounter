package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/clarkduvall/hyperloglog"
	"github.com/col3name/ip-unique-addr/pkg/hash"
	"os"
	"runtime"
	"sync"
	"time"
)

func main() {
	fileName := flag.String("f", "C:\\Users\\mikha\\go\\src\\github.com\\col3name\\ip-unique-addr\\ip_addresses_out", "input file1")
	countParallelTask := flag.Int("n", runtime.NumCPU(), "count parallel reader")
	flag.Parse()

	fmt.Println(*fileName)
	file, err := os.Open(*fileName)
	if err != nil {
		fmt.Println("cannot able to read the file1", err)
		return
	}

	stat, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(stat.Size())

	partSize := stat.Size() / int64(*countParallelTask)
	defer func() {
		_ = file.Close()
	}()
	h, _ := hyperloglog.New(16)
	start := time.Now()
	files := make([]*os.File, *countParallelTask)
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	for i := range files {
		file, err := os.Open(*fileName)
		if err != nil {
			fmt.Println(err)
			return
		}
		if i > 0 {
			_, err := file.Seek(partSize*int64(i), os.SEEK_CUR)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			hyperLogLog := readWriteFile(file, partSize)
			mu.Lock()
			err = h.Merge(hyperLogLog)
			mu.Unlock()
			if err != nil {
				fmt.Println(err)
				panic(err)
			}
		}()
		files[i] = file
	}
	wg.Wait()
	for _, file := range files {
		file.Close()
	}
	fmt.Println(time.Since(start).String(), h.Count())
}

func readWriteFile(input *os.File, count int64) *hyperloglog.HyperLogLog {
	r := bufio.NewReader(input)
	h, _ := hyperloglog.New(16)
	j := int64(0)
	for j < count-1 {
		bytes, err := r.ReadBytes('\n')
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		if len(bytes)/1024/1024 > 1 {
			fmt.Println(len(bytes))
		}
		j += int64(len(bytes))
		h.Add(hash.Hash32(string(bytes)))
	}
	fmt.Println(count, j)
	return h
}

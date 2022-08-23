package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/clarkduvall/hyperloglog"
	"hash"
	"hash/fnv"
	"io"
	"math"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	fileName := flag.String("f", "d:/ip_addresses/ip_addresses", "input file")
	flag.Parse()

	file, err := os.Open(*fileName)
	if err != nil {
		fmt.Println("cannot able to read the file", err)
		return
	}

	defer func() {
		_ = file.Close()
	}()

	fileStat, err := file.Stat()
	if err != nil {
		fmt.Println("Could not able to get the file stat")
		return
	}

	fileSize := fileStat.Size()
	offset := fileSize - 1
	lastLineSize := 0

	start := time.Now()
	fmt.Println("started at", start.String())

	for {
		b := make([]byte, 1)
		n, err := file.ReadAt(b, offset)
		if err != nil {
			fmt.Println("Error reading file ", err)
			break
		}
		char := string(b[0])
		if char == "\n" {
			break
		}
		offset--
		lastLineSize += n
	}

	lastLine := make([]byte, lastLineSize)
	_, err = file.ReadAt(lastLine, offset+1)

	if err != nil {
		fmt.Println("Could not able to read last line with offset", offset, "and lastline size", lastLineSize)
		return
	}

	err = process(file)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("\nTime taken - ", time.Since(start))
}

func process(f *os.File) error {
	h, _ := hyperloglog.New(16)

	linesPool := sync.Pool{New: func() interface{} {
		lines := make([]byte, 250*1024)
		return lines
	}}

	stringPool := sync.Pool{New: func() interface{} {
		lines := ""
		return lines
	}}

	r := bufio.NewReader(f)
	var wg sync.WaitGroup
	j := 0
	for {
		buf := linesPool.Get().([]byte)

		n, err := r.Read(buf)
		buf = buf[:n]

		if n == 0 {
			if err != nil {
				fmt.Println(err)
				break
			}
			if err == io.EOF {
				break
			}
			return err
		}

		nextUntilNewline, err := r.ReadBytes('\n')
		if err != io.EOF {
			buf = append(buf, nextUntilNewline...)
		}
		wg.Add(1)
		go func() {
			j += processChunk(buf, &linesPool, &stringPool, h)
			wg.Done()
		}()
	}

	fmt.Println("uniq", h.Count(), j)
	wg.Wait()
	return nil
}

func processChunk(chunk []byte, linesPool *sync.Pool, stringPool *sync.Pool, h *hyperloglog.HyperLogLog) int {
	logs := stringPool.Get().(string)
	logs = string(chunk)

	linesPool.Put(chunk)

	logsSlice := strings.Split(logs, "\n")
	stringPool.Put(logs)

	chunkSize := 300
	n := len(logsSlice)
	noOfThread := n / chunkSize

	if n%chunkSize != 0 {
		noOfThread++
	}
	count := 0
	var wg2 sync.WaitGroup
	for i := 0; i < (noOfThread); i++ {
		wg2.Add(1)
		go func(s int, e int) {
			defer wg2.Done()
			for j := s; j < e; j++ {
				text := logsSlice[j]
				if len(text) == 0 {
					continue
				}
				logSlice := strings.SplitN(text, ",", 2)
				line := logSlice[0]
				h.Add(hash32(line))
				count++
			}
		}(i*chunkSize, int(math.Min(float64((i+1)*chunkSize), float64(len(logsSlice)))))
	}

	wg2.Wait()
	logsSlice = nil
	return count
}

func hash32(s string) hash.Hash32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s))
	return h
}

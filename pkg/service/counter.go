package service

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/clarkduvall/hyperloglog"
	"github.com/col3name/ip-unique-addr/pkg/hash"
	"io"
	"os"
	"sync"
)

type CounterService interface {
	Count() (uint64, error)
}

var ErrorInvalidValue = errors.New("errorInvalidValue")

type uniqueCounterHyperLogLog struct {
	filePath          string
	countParallelTask int
}

func NewUniqueCounterHLL(filepath string, countParallelTask int) *uniqueCounterHyperLogLog {
	return &uniqueCounterHyperLogLog{filePath: filepath, countParallelTask: countParallelTask}
}

func (c *uniqueCounterHyperLogLog) CountInFile(index int, partSize int64) (*os.File, *hyperloglog.HyperLogLog, error) {
	if index < 0 {
		return nil, nil, ErrorInvalidValue
	}

	file, err := os.Open(c.filePath)
	if err != nil {
		return nil, nil, err
	}

	_, err = file.Seek(partSize*int64(index), io.SeekCurrent)
	if err != nil {
		return nil, nil, err
	}
	hll, err := c.readFile(file, partSize)

	return file, hll, err
}

func (c *uniqueCounterHyperLogLog) Count() (uint64, error) {
	partSize, err := c.getPartSize()
	if err != nil {
		return 0, err
	}
	h, _ := hyperloglog.New(16)
	files := make([]*os.File, c.countParallelTask)
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	for i := range files {
		c.handleFileAsync(&wg, partSize, &mu, h, files, i)
	}
	wg.Wait()
	c.closeAllFiles(files)

	count := h.Count()
	h.Clear()
	return count, nil
}

func (c *uniqueCounterHyperLogLog) closeAllFiles(files []*os.File) {
	for _, file := range files {
		file.Close()
	}
}

func (c *uniqueCounterHyperLogLog) getPartSize() (int64, error) {
	file, err := os.Open(c.filePath)
	if err != nil {
		return 0, err
	}
	stat, err := file.Stat()
	if err != nil {
		return 0, err
	}
	partSize := stat.Size() / int64(c.countParallelTask)
	defer func() {
		_ = file.Close()
	}()
	return partSize, nil
}

func (c *uniqueCounterHyperLogLog) handleFileAsync(wg *sync.WaitGroup, partSize int64, mu *sync.Mutex, h *hyperloglog.HyperLogLog, files []*os.File, i int) {
	wg.Add(1)
	go func(i int) {
		defer wg.Done()
		file, hll, err := c.CountInFile(i, partSize)
		if err != nil {
			fmt.Println(err)
			return
		}
		mu.Lock()
		defer mu.Unlock()
		err = h.Merge(hll)
		if err != nil {
			fmt.Println(err)
			return
		}
		files[i] = file
	}(i)
}

func (c *uniqueCounterHyperLogLog) readFile(input *os.File, count int64) (*hyperloglog.HyperLogLog, error) {
	r := bufio.NewReader(input)
	h, _ := hyperloglog.New(16)
	j := int64(0)
	for j < count-1 {
		bytes, err := r.ReadBytes('\n')
		if err != nil {
			return nil, err
		}
		j += int64(len(bytes))
		h.Add(hash.Hash32(string(bytes)))
	}
	return h, nil
}

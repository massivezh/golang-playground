package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

type ChunkInfo struct {
	fileOffset int64
	chunkRange string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func caluSegmentRange(size int64, segment_count int64) []string {
	fmt.Println("Total: ", size)
	// check length big than cnt
	segment_size := size / segment_count
	ranges := make([]string, segment_count)
	for i := int64(0); i < segment_count-1; i++ {
		ranges[i] = fmt.Sprintf("bytes=%d-%d", i*segment_size, (i+1)*segment_size-1)
	}
	ranges[segment_count-1] = fmt.Sprintf("bytes=%d-%d", segment_size*(segment_count-1), size-1)
	return ranges
}
func getUrlSize(url string) int64 {
	resp, err := http.Get(url)
	defer resp.Body.Close()

	check(err)

	return resp.ContentLength
}

func downloadFromUrl(url string, ch <-chan *ChunkInfo) {
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]
	fmt.Println("Downloading", url, "to", fileName)

	// TODO: check file existence first with io.IsExist
	output, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error while creating", fileName, "-", err)
		return
	}
	defer output.Close()

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)

	var range_str string
	range_str = fmt.Sprintf("bytes=%d-%d", start, end)
	req.Header.Add("Range", range_str)

	response, err := client.Do(req)

	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}

	fmt.Println(n, "bytes downloaded.")
}

// verify file seek opration
func WritePart(fileName string, ch <-chan int64, data []byte) {
	offset := <-ch
	fmt.Println(offset)
	f, err := os.OpenFile(fileName, os.O_RDWR, 0666)
	n1, err := f.Seek(offset, 0)
	check(err)

	nr := bytes.NewReader(data)
	n2, err := io.Copy(f, nr)
	check(err)
	_ = n1
	_ = n2

	defer f.Close()
}
func main() {
	url := "http://mirrors.sohu.com/centos/2/readme.txt"
	//url_big := "http://mirrors.sohu.com/centos/6.6/isos/x86_64/CentOS-6.6-x86_64-bin-DVD1.iso"

	size := getUrlSize(url_big)
	chunk_ranges := splitSegment(size, 5)
	ch := make(chan *ChunkInfo)

	var wg sync.WaitGroup

	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			DownloadPart(url, offchan)
		}()

	}
	go func() {
		for _, r_stru := range chunk_ranges {
			offchan <- &r_stru
		}
	}()
	wg.Wait()
}

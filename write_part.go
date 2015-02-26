// verify file seek opration
package main

import (
	//	"fmt"
	"os"
	"sync"
)

var off = []int64{0, 1, 2, 3, 4, 5}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func WritePart(f *os.File, ch <-chan int64, data []byte) {
	//	fmt.Println(offset)
	offset := <-ch
	m, err := f.WriteAt(data, offset)
	check(err)
	_ = m

	f.Sync()
}
func main() {

	//dt := []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99}

	offchan := make(chan int64)
	fileName := "/root/golang-playground/dat1"
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
	check(err)
	defer f.Close()

	var wg sync.WaitGroup

	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			WritePart(f, offchan, []byte{0x11})
		}()

	}
	go func() {
		for _, vOff := range off {
			offchan <- vOff
		}
	}()
	wg.Wait()
}

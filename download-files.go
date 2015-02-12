package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func downloadFromUrl(url string, start int, end int) {
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

func main() {
	url := "http://mirrors.sohu.com/centos/2/readme.txt"
	downloadFromUrl(url, 16, 31)
}

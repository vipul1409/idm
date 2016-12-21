package main

/*
	Build go build idm.go
	Usage Example : ./idm -url=http://9xmusiq.com/songs/Bollywood%20Songs/2016%20Hindi%20Mp3/A%20Flying%20Jatt%20\(2016\)/Beat%20Pe%20Booty%20%5bStarmusiq.xyz%5d.mp3 -parts=20 -output=beat-booty
*/
import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

func downloadPartFile(url string, startIndex int, endIndex int, ouptutFile string, wg *sync.WaitGroup) {
	defer wg.Done()
	client := &http.Client{}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Range", fmt.Sprintf("bytes=%d-%d", startIndex, endIndex))

	resp, _ := client.Do(req)

	defer resp.Body.Close()

	writeRespToFile(resp, ouptutFile)
}

func getTotalSize(url string) int {

	resp, _ := http.Head(url)
	maps := resp.Header

	defer resp.Body.Close()
	a, _ := strconv.Atoi(maps["Content-Length"][0])

	return a

}

func writeRespToFile(resp *http.Response, fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
}

func main() {

	var wg sync.WaitGroup

	url_ptr := flag.String("url", "https://www.bing.com", "URL to download.")
	parts_ptr := flag.Int("parts", 20, "Number of go routines to use for downloading.")
	output_ptr := flag.String("output", "output", "Output File Name.")

	flag.Parse()

	url := strings.Trim(*url_ptr, " ")
	parts := *parts_ptr
	output := strings.Trim(*output_ptr, " ")

	fmt.Println(url)
	fmt.Println(parts)
	fmt.Println(output)

	var extension string

	if strings.LastIndex(url, ".") != -1 {
		extension = url[strings.LastIndex(url, ".")+1 : len(url)]
	}
	size := getTotalSize(url)
	factor := size / parts
	if size%parts > 0 {
		parts = parts + 1
	}
	wg.Add(parts)
	os.Mkdir(fmt.Sprintf("%s-tmp", output), 0777)
	for i := 0; i < parts; i++ {
		go downloadPartFile(url, i*factor, ((i+1)*factor - 1), fmt.Sprintf("%s-tmp/part-%d", output, i), &wg)
	}

	wg.Wait()

	fmt.Println("Now merging parts")
	curr_dir, _ := os.Getwd()
	os.Chdir(fmt.Sprintf("curr_dir/%s-tmp", output))
	exec.Command("sh", "-c", fmt.Sprintf("cat `ls part-* | sort -t '-' -k 2n` > %s.%s", output, extension)).Output()
	os.Chdir(curr_dir)
	os.Remove(fmt.Sprintf("%s-tmp", output))

}

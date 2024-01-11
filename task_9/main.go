package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Args struct {
	O        []string
	maxDepth int

	addresses []string
}

func getArgs() (*Args, error) {
	O := flag.String("O", "", "new filename")
	maxDepth := flag.Int("depth", 1, "sets the maximum depth for recursively loading the entire site")

	flag.Parse()

	var files []string

	if *maxDepth < 1 {
		return nil, errors.New("depth should be positive")
	}

	if len(*O) > 0 {
		files = strings.Split(*O, " ")
	}

	args := &Args{
		O:        files,
		maxDepth: *maxDepth,
	}

	if len(flag.Args()) < 1 {
		return nil, errors.New("wrong address")
	}
	args.addresses = append(args.addresses, flag.Args()...)

	return args, nil
}

// getFileName - sets filename and suffix for existed one
func getFileName(filename, address, path string) string {
	var counter int

	if filename == "" {
		if strings.HasSuffix(address, "/") {
			filename = "index.html"
		} else {
			filename = filepath.Base(address)
			if !strings.Contains(filename, ".") {
				filename = fmt.Sprintf("%s.html", filename)
			}
		}
	}

	originalName := filename

	for {
		_, err := os.Stat(fmt.Sprintf("%s/%s", path, filename))

		if errors.Is(err, os.ErrNotExist) {
			break
		}

		counter++
		suffix := fmt.Sprintf(".%d", counter)

		newFilename := strings.Builder{}
		newFilename.Grow(len(originalName) + len(suffix))
		newFilename.WriteString(originalName)
		newFilename.WriteString(suffix)

		filename = newFilename.String()
	}

	return filename
}

// saveToFile - saving body to filename
func saveToFile(body []byte, filename, address string) error {
	parsed, err := url.Parse(address)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("%s%s", parsed.Host, parsed.Path)
	path = filepath.Dir(path)

	err = os.MkdirAll(path, 0644)
	if err != nil && os.IsNotExist(err) {
		return err
	}

	filename = getFileName(filename, address, path)

	err = ioutil.WriteFile(fmt.Sprintf("%s/%s", path, filename), body, 0644)
	if err != nil {
		return err
	}

	fmt.Printf("%s - '%s'saved\n\n", time.Now().Format("01/02/06 15:04:05"), filename)

	return nil
}

// getLink - return all links from page
func getLinks(address string) map[string]bool {
	links := make(map[string]bool)

	// parsing URL
	parsed, _ := url.Parse(address)
	host := parsed.Hostname()

	client := http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(address)
	if err != nil || resp == nil {
		return nil
	}
	defer resp.Body.Close()

	// Create a goquery document from the HTTP response
	document, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	// find all links 'a' and process them
	document.Find("a").Each(func(index int, element *goquery.Selection) {
		link, _ := element.Attr("href")
		parsed, err := url.Parse(link)
		if err != nil || parsed.Path == "" {
			return
		}

		// checking for same host name
		linkHost := parsed.Hostname()
		if linkHost != "" && linkHost != host {
			return
		}

		scheme := "https"
		if parsed.Scheme != "" {
			scheme = parsed.Scheme
		}

		newLink := fmt.Sprintf("%s://%s%s", scheme, host, parsed.Path)

		// save only unique
		if !uniqueLinks[newLink] {
			links[newLink] = true
			uniqueLinks[newLink] = true
		}
	})

	return links
}

func download(address, filename string, maxDepth int) error {
	if maxDepth < 1 {
		return nil
	}

	client := http.Client{}

	resp, err := client.Get(address)
	if err != nil || resp == nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("response status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = saveToFile(body, filename, address)
	if err != nil {
		return err
	}

	if maxDepth > 2 {
		links := getLinks(address)
		for link := range links {
			err := download(link, filename, maxDepth-1)
			if err != nil {
				continue
			}
		}
	}

	return nil
}

var uniqueLinks map[string]bool

func wget() error {
	if len(os.Args) < 2 {
		return errors.New("you need to specify a webaddress")
	}

	args, err := getArgs()
	if err != nil {
		return err
	}

	uniqueLinks = make(map[string]bool)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGSEGV)

	go func() {
		<-sigs
		os.Exit(1)
	}()

	for i, address := range args.addresses {
		var filename string

		if i < len(args.O) {
			filename = args.O[i]
		}

		err := download(address, filename, args.maxDepth)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	err := wget()
	if err != nil {
		log.Fatal(err)
	}
}

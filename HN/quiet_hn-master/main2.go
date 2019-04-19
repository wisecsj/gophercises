package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"./hn"
)

// items count rate
const RATE = 1.25

func main() {

	// fq
	os.Setenv("HTTP_PROXY", "http://127.0.0.1:12333")
	os.Setenv("HTTPS_PROXY", "http://127.0.0.1:12333")
	// parse flags
	var port, numStories int
	flag.IntVar(&port, "port", 8888, "the port to start the web server on")
	flag.IntVar(&numStories, "num_stories", 30, "the number of top stories to display")
	flag.Parse()

	tpl := template.Must(template.ParseFiles("./index.gohtml"))

	http.HandleFunc("/", handler(numStories, tpl))

	// Start the server
	go log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
	go log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", 8887), nil))
}
func getItem(id int, client hn.Client, wg *sync.WaitGroup, c chan item) {
	defer wg.Done()
	hnItem, err := client.GetItem(id)
	if err != nil {
		log.Printf("Get item with id %d failed", id)
		return
	}
	item := parseHNItem(hnItem)
	if isStoryLink(item) {
		c <- item
	}

}
func getItems(start, end int, ids []int, client hn.Client) []item {
	var stories []item
	var wg sync.WaitGroup
	chItem := make(chan item, end-start)

	for i := start; i < end; i++ {
		wg.Add(1)
		go getItem(ids[i], client, &wg, chItem)
	}

	// 用来确保所有story都已经录入到切片中
	over := make(chan int)
	go func() {
		for s := range chItem {
			stories = append(stories, s)
		}
		over <- 1
	}()
	wg.Wait()
	close(chItem)
	<-over

	return stories
}
func handler(numStories int, tpl *template.Template) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		var client hn.Client
		ids, err := client.TopItems()
		// 保存id顺序
		mp := make(map[int]int)
		for i, id := range ids {
			mp[id] = i
		}
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to load top stories:%s", err), http.StatusInternalServerError)
			return
		}

		itemCount := int(RATE * float32(numStories))
		left := 0
		right := itemCount
		stories := getItems(left, right, ids, client)
		// 当以RATE的倍率无法满足numStories需求时，循环每次请求一定数量story
		for {
			if len(stories) < numStories {
				left = right
				right = right + (itemCount - numStories)
				s := getItems(left, right, ids, client)
				stories = append(stories, s...)
			} else {
				break
			}
		}

		// 对stories排序，并取前numStories个item
		sort.Slice(stories, func(i, j int) bool {
			return mp[stories[i].ID] < mp[stories[j].ID]
		})

		data := templateData{
			Stories: stories[:numStories],
			Time:    time.Now().Sub(start),
		}
		err = tpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	})
}

func isStoryLink(item item) bool {
	return item.Type == "story" && item.URL != ""
}

func parseHNItem(hnItem hn.Item) item {
	ret := item{Item: hnItem}
	url, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(url.Hostname(), "www.")
	}
	return ret
}

// item is the same as the hn.Item, but adds the Host field
type item struct {
	hn.Item
	Host string
}

type templateData struct {
	Stories []item
	Time    time.Duration
}

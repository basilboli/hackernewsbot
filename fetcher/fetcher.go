/*Package fetcher fetches latest posts from hackernews and updating its cache
  use multiple workers to fetch posts in parallel
*/
package fetcher

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/basilboli/hackernewsbot/models"
)

func httpGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	byt := []byte(body)
	if err != nil {
		return nil, err
	}
	return byt, nil
}

func storyFetchWorker(id int, jobs <-chan string, results chan<- int) {
	for j := range jobs {

		byt, err := httpGet(j)
		if err != nil {
			log.Printf("Problem reading %s ", j)
		}
		var post models.Post
		// var dat map[string]interface{}
		if err := json.Unmarshal(byt, &post); err != nil {
			log.Printf("Problem decoding %s", j)
		}
		fmt.Println(post.Title)
		post.Update()
		results <- 1
	}
}

func fetchStories(URL string) []int64 {
	story := "https://hacker-news.firebaseio.com/v0/item/%d.json"

	byt, err := httpGet(URL)

	var stories []int64
	if err := json.Unmarshal(byt, &stories); err != nil {
		log.Printf("Problem decoding %s ", URL)
	}

	storiesNo := len(stories)
	fmt.Printf("Found %d stories \n", storiesNo)

	if storiesNo <= 0 {
		log.Fatal("No story can be fetched!")
	}

	jobs := make(chan string, storiesNo)
	results := make(chan int, storiesNo)

	for w := 1; w <= 20; w++ {
		go storyFetchWorker(w, jobs, results)
	}

	if err != nil {
		log.Fatal(err)
	}

	for _, num := range stories {
		itemURL := fmt.Sprintf(story, num)
		jobs <- itemURL
	}

	close(jobs)

	for a := 1; a <= storiesNo; a++ {
		<-results
	}
	fmt.Printf("Processed %d elements for %s", storiesNo, URL)
	return stories
}

func fetchNewStories() {
	// fetching new stories ids
	URL := "https://hacker-news.firebaseio.com/v0/newstories.json"
	fetchStories(URL)
}

// FetchTopStories fetches top stories from hacker news site
func FetchTopStories() {
	// fetching top stories
	for {
		fmt.Println("Fetching top stories")
		URL := "https://hacker-news.firebaseio.com/v0/topstories.json"
		ids := fetchStories(URL)
		err := models.UpdateTopStories(ids)
		if err != nil {
			log.Println(err)
		}
		time.Sleep(time.Minute * 10)
	}
}

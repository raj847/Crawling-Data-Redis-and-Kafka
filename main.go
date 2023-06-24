package main

// import (
// 	"fmt"
// 	"net/http"
// 	"strings"

// 	"github.com/PuerkitoBio/goquery"
// )

// func main() {
// 	resp, err := http.Get("https://nasional.tempo.co/read/1732991/pdip-bilang-sikap-megawati-soal-gagasan-tunda-pemilu-sudah-jelas")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer resp.Body.Close()

// 	doc, err := goquery.NewDocumentFromReader(resp.Body)
// 	if err != nil {
// 		panic(err)
// 	}

// 	//$ pada jquery = doc.find
// 	title := doc.Find("div.detail-title > h1.title").Text()
// 	image, _ := doc.Find(`div.foto-detail > figure > img[itemprop="image"]`).Attr("src")

// 	fmt.Println("title :", title)
// 	fmt.Println("image :", image)

// 	contents := []string{}
// 	doc.Find("div.detail-konten > p").Each(func(i int, s *goquery.Selection) {
// 		contents = append(contents, s.Text())
// 	})
// 	res := strings.Join(contents, "\n")

// 	fmt.Println(res)
// }

// next

// import (
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"strings"

// 	"github.com/PuerkitoBio/goquery"
// )

// type News struct {
// 	Title    string
// 	ImageURL string
// 	Content  string
// }

// type ScrapeOptions struct {
// 	Body            io.Reader
// 	TitleSelector   string
// 	ImageSelector   string
// 	ContentSelector string
// }

// func scrape(opts ScrapeOptions) (*News, error) {
// 	doc, err := goquery.NewDocumentFromReader(opts.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	title := doc.Find(opts.TitleSelector).Text()
// 	imageURL, _ := doc.Find(opts.ImageSelector).Attr("src")
// 	contents := []string{}
// 	doc.Find(opts.ContentSelector).Each(func(i int, s *goquery.Selection) {
// 		contents = append(contents, s.Text())
// 	})

// 	news := News{
// 		Title:    title,
// 		ImageURL: imageURL,
// 		Content:  strings.Join(contents, "\n"),
// 	}
// 	return &news, nil
// }

// func parseUrls(body io.Reader, pattern string) ([]string, error) {
// 	doc, err := goquery.NewDocumentFromReader(body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	urls := []string{}
// 	urlmap := map[string]bool{}
// 	doc.Find("a").Each(func(i int, s *goquery.Selection) {
// 		url, _ := s.Attr("href")
// 		if strings.Contains(url, pattern) {
// 			if _, ok := urlmap[url]; !ok {
// 				urls = append(urls, url)
// 				urlmap[url] = true
// 			}
// 		}
// 	})

// 	return urls, nil
// }

// func main() {
// 	resp, err := http.Get("https://nasional.tempo.co/read/1732903/hingga-april-2023-ada-11-kasus-kematian-karena-rabies-kemenkes-segera-ke-faskes-jika-digigit-anjing")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer resp.Body.Close()

// 	apapun, err := parseUrls(resp.Body, "tempo.co/read")
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(apapun)

// news, err := scrape(ScrapeOptions{
// 	Body:            resp.Body,
// 	TitleSelector:   "div.detail-title > h1.title",
// 	ImageSelector:   `div.foto-detail > figure > img[itemprop="image"]`,
// 	ContentSelector: "div.detail-konten > p",
// })

// if err != nil {
// 	panic(err)
// }

// fmt.Println(news.Title)

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/redis/go-redis/v9"
)

type News struct {
	Title    string
	ImageURL string
	Content  string
}

type ScrapeOptions struct {
	Body            io.Reader
	TitleSelector   string
	ImageSelector   string
	ContentSelector string
}

func parseNews(opts ScrapeOptions) (*News, error) {
	doc, err := goquery.NewDocumentFromReader(opts.Body)
	if err != nil {
		return nil, err
	}

	title := doc.Find(opts.TitleSelector).Text()
	imageURL, _ := doc.Find(opts.ImageSelector).Attr("src")
	contents := []string{}
	doc.Find(opts.ContentSelector).Each(func(i int, s *goquery.Selection) {
		contents = append(contents, s.Text())
	})

	news := News{
		Title:    title,
		ImageURL: imageURL,
		Content:  strings.Join(contents, "\n"),
	}
	return &news, nil
}

func isTempoURL(urlString string) bool {
	u, err := url.Parse(urlString)
	if err != nil {
		return false
	}

	domain := "tempo.co"
	// fmt.Println(u.Path)
	// fmt.Println(u.Hostname())
	if strings.Contains(u.Hostname(), domain) {
		if strings.HasPrefix(u.Path, "/read") {
			return true
		}
	}
	return false
}

func parseUrls(body io.Reader) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}

	urls := []string{}
	urlmap := map[string]bool{}
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		url, _ := s.Attr("href")
		if isTempoURL(url) {
			if _, ok := urlmap[url]; !ok {
				urls = append(urls, url)
				urlmap[url] = true
			}
		}
	})

	return urls, nil
}

func main() {
	var ctx = context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	resp, err := http.Get("https://tempo.co/")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	urls, err := parseUrls(resp.Body)
	if err != nil {
		panic(err)
	}

	for len(urls) > 0 {
		url := urls[0]
		// fmt.Println("checking on", url)
		v, _ := rdb.Exists(ctx, url).Result()

		if v == 1 {
			urls = urls[1:]
			continue
		}
		fmt.Println("Working on", url)
		resp, err := http.Get(url)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		urlsx, err := parseUrls(resp.Body)
		if err != nil {
			panic(err)
		}
		rdb.Set(ctx, url, 1, 5*time.Minute)

		urls = urls[1:]
		urls = append(urls, urlsx...)
		time.Sleep(800 * time.Millisecond)
	}

	// news, err := parseNews(ScrapeOptions{
	// 	Body:            resp.Body,
	// 	TitleSelector:   "div.detail-title > h1.title",
	// 	ImageSelector:   `div.foto-detail > figure > img[itemprop="image"]`,
	// 	ContentSelector: "div.detail-konten > p",
	// })

	// if err != nil {
	// 	panic(err)
	// }

}

package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

var wg sync.WaitGroup

func main() {

	now := time.Now()
	//Crawl("https://www.ibm.com/community/z/open-source-software/", 1)
	links, err := findLinks("https://www.ibm.com/community/z/open-source-software")
	if err != nil {
		fmt.Println("error occurred : ", err)
	}

	fmt.Println("total links found:", len(links))

	links = removeDuplicateValues(links)

	fmt.Println("total links found:", len(links))

	wg.Add(len(links))

	for _, link := range links {
		go validateLink(link)
	}

	wg.Wait()

	fmt.Println("time taken:", time.Since(now))
}

func findLinks(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}
	return visit(nil, doc), nil
}

// visit appends to links each link found in n, and returns the result.
func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c)
	}

	return links
}

func removeDuplicateValues(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			if !strings.HasPrefix(entry, "#") {
				list = append(list, entry)
			}
		}
	}

	return list
}

func validateLink(url string) {
	defer wg.Done()
	resp, err := http.Get(url)

	if err != nil {

		fmt.Println("Validating ", url, "Errored:  ", err)
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		panic(fmt.Errorf("parsing %s as HTML: %v", url, err))
	}
	fmt.Println("Validating ", url, "..Ok")
}

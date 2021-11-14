package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const IssuesURL = "https://api.github.com/search/issues"

type IssuesSearchResult struct {
	TotalCount int `json:"total_count"`
	Items      []*Issue
}

type Issue struct {
	Number    int
	HTMLURL   string `json:"html_url"`
	Title     string
	State     string
	User      *User
	CreatedAt time.Time
	Body      string
}

type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}

func SearchIssues(terms []string) (*IssuesSearchResult, error) {
	q := url.QueryEscape(strings.Join(terms, " "))
	resp, err := http.Get(IssuesURL + "?q=" + q)
	if err != nil {
		return nil, err
	}

	// We must close resp.Body on all execution paths.
	// (Chapter 5 presents 'defer', which makes this simpler.)
	if resp.StatusCode != http.StatusOK {
		err := resp.Body.Close()
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("search query failed: %s", resp.Status)
	}

	var result IssuesSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		err := resp.Body.Close()
		if err != nil {
			return nil, err
		}
		return nil, err
	}
	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func PrintIssuesGroupByTime(items []*Issue) {
	var PastMonth []*Issue
	var PastYear []*Issue
	var OverYear []*Issue

	for _, item := range items {
		if time.Since(item.CreatedAt).Hours() <= 24*30 {
			PastMonth = append(PastMonth, item)
		} else if time.Since(item.CreatedAt).Hours() <= 24*365 {
			PastYear = append(PastYear, item)
		} else {
			OverYear = append(OverYear, item)
		}
	}
	fmt.Println("Issues in past month:")
	printIssues(PastMonth)
	fmt.Println("Issues in past year:")
	printIssues(PastYear)
	fmt.Println("Issues over one year:")
	printIssues(OverYear)
}

func printIssues(items []*Issue) {
	for _, item := range items {
		fmt.Printf("%v #%-5d %9.9s %.55s\n",
			item.CreatedAt, item.Number, item.User.Login, item.Title)
	}
}

func main() {
	result, err := SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d issues:\n", result.TotalCount)
	for _, item := range result.Items {
		fmt.Printf("#%-5d %9.9s %.55s\n", item.Number, item.User.Login, item.Title)
	}
	PrintIssuesGroupByTime(result.Items)
}

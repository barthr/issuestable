package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/google/go-github/github"
	"github.com/olekukonko/tablewriter"
)

var (
	amount = flag.Int("amount", -1, "Set the maximum amount of issues to show")
	repo   = flag.String("repo", "", "Repository on http://www.github.com to list the issues from")

	client = github.NewClient(nil)
)

func main() {
	flag.Usage = func() {
		fmt.Println("Usage of Issues table")
		flag.PrintDefaults()
	}
	flag.Parse()

	githubRepoURL := *repo
	if githubRepoURL == "" {
		flag.Usage()
		os.Exit(1)
	}

	u, err := url.Parse(githubRepoURL)

	if err != nil {
		fmt.Printf("Invalid url! %s", err)
		os.Exit(1)
	}

	path := strings.Split(u.EscapedPath(), "/")

	if len(path) < 3 {
		fmt.Println("Invalid github URL")
		os.Exit(1)
	}

	opt := &github.IssueListByRepoOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	var issues []*github.Issue

	for len(issues) <= *amount || *amount < 0 {
		newIssues, resp, err := client.Issues.ListByRepo(context.Background(), path[1], path[2], opt)
		if err != nil {
			log.Fatalf("error fetching issues for %v/%v: %v", path[1], path[2], err)
		}

		issues = append(issues, newIssues...)
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}

	if len(issues) == 0 {
		fmt.Printf("No issues found for this repository")
		os.Exit(0)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"number", "created_at", "title"})

	for index, issue := range issues {
		if index == *amount {
			break
		}
		column := make([]string, 3)
		column[0] = strconv.Itoa(*issue.Number)
		column[1] = issue.CreatedAt.String()
		column[2] = *issue.Title
		table.Append(column)
	}
	table.SetRowLine(true)
	table.Render()
}

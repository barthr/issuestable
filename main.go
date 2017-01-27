package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/google/go-github/github"
	"github.com/olekukonko/tablewriter"

	"flag"

	"os"

	"net/url"
)

var (
	amount = flag.Int("amount", 50, "Set the maximum amount of issues to show")
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

	pageOptions := &github.IssueListByRepoOptions{}

	issues, page, err := client.Issues.ListByRepo(path[1], path[2], nil)

	if err != nil {
		fmt.Printf("Cannot find issues for repo %s/%s because %v", path[1], path[2], err)
		os.Exit(1)
	}

	for currentPage := 0; currentPage < page.LastPage && len(issues) != *amount; currentPage++ {
		pageAmount := github.ListOptions{Page: currentPage, PerPage: 100}
		pageOptions.ListOptions = pageAmount
		newIssues, _, err := client.Issues.ListByRepo(path[1], path[2], pageOptions)

		if err == nil {
			issues = append(issues, newIssues...)
		}
	}

	if len(issues) == 0 {
		fmt.Printf("No issues found for this repository")
		os.Exit(0)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"number", "created_at", "title"})

	for _, issue := range issues {
		row := make([]string, 3)
		row[0] = strconv.Itoa(*issue.Number)
		row[1] = issue.CreatedAt.String()
		row[2] = *issue.Title
		table.Append(row)
	}

	table.Render()
}
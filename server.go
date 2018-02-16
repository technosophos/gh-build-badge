package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
)

var (
	githubToken = ""
	lastStatus  = map[string]string{}
)

func main() {
	githubToken = os.Getenv("GITHUB_TOKEN")

	r := gin.New()
	r.Use(gin.Recovery())
	r.GET("/healthz", func(c *gin.Context) { c.String(http.StatusOK, "OK") })

	gh := r.Group("/v1/github")
	gh.Use(gin.Logger())
	gh.GET("/build/:owner/:project/badge.svg", badge)

	r.Run(":8181")
}

func badge(c *gin.Context) {
	owner := c.Param("owner")
	project := c.Param("project")
	branch := c.DefaultQuery("branch", "master")
	key := fmt.Sprintf("%s/%s/%s", owner, project, branch)

	status, err := ghStatus(owner, project, branch)
	if err != nil {
		if _, ok := err.(*github.RateLimitError); ok {
			status, ok = lastStatus[key]
			if !ok {
				status = "Unknown"
			}
		} else {
			fmt.Printf("failed ghStatus: %s", err)
			svg(c, fmt.Sprintf(other, "Unknown", "Unknown"))
			return
		}
	}

	// Trivial rate limit protection
	lastStatus[key] = status

	switch status {
	case "success":
		svg(c, pass)
	case "pending":
		svg(c, run)
	case "failure", "error":
		svg(c, fail)
	default:
		svg(c, fmt.Sprintf(other, status, status))
	}
}

func ghStatus(owner, project, branch string) (string, error) {
	ctx := context.TODO()
	client := github.NewClient(nil)
	status, _, err := client.Repositories.GetCombinedStatus(ctx, owner, project, branch, &github.ListOptions{})
	if err != nil {
		return "error", err
	}
	return status.GetState(), nil
}

func svg(c *gin.Context, data string) {
	c.Data(http.StatusOK, "image/svg+xml", []byte(data))
}

const (
	pass  = `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="86" height="20"><linearGradient id="b" x2="0" y2="100%"><stop offset="0" stop-color="#bbb" stop-opacity=".1"/><stop offset="1" stop-opacity=".1"/></linearGradient><clipPath id="a"><rect width="86" height="20" rx="3" fill="#fff"/></clipPath><g clip-path="url(#a)"><path fill="#555" d="M0 0h51v20H0z"/><path fill="#4c1" d="M51 0h35v20H51z"/><path fill="url(#b)" d="M0 0h86v20H0z"/></g><g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="110"><text x="265" y="150" fill="#010101" fill-opacity=".3" transform="scale(.1)" textLength="410">brigade</text><text x="265" y="140" transform="scale(.1)" textLength="410">brigade</text><text x="675" y="150" fill="#010101" fill-opacity=".3" transform="scale(.1)" textLength="250">pass</text><text x="675" y="140" transform="scale(.1)" textLength="250">pass</text></g> </svg>`
	fail  = `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="78" height="20"><linearGradient id="b" x2="0" y2="100%"><stop offset="0" stop-color="#bbb" stop-opacity=".1"/><stop offset="1" stop-opacity=".1"/></linearGradient><clipPath id="a"><rect width="78" height="20" rx="3" fill="#fff"/></clipPath><g clip-path="url(#a)"><path fill="#555" d="M0 0h51v20H0z"/><path fill="#e05d44" d="M51 0h27v20H51z"/><path fill="url(#b)" d="M0 0h78v20H0z"/></g><g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="110"><text x="265" y="150" fill="#010101" fill-opacity=".3" transform="scale(.1)" textLength="410">brigade</text><text x="265" y="140" transform="scale(.1)" textLength="410">brigade</text><text x="635" y="150" fill="#010101" fill-opacity=".3" transform="scale(.1)" textLength="170">fail</text><text x="635" y="140" transform="scale(.1)" textLength="170">fail</text></g> </svg>`
	run   = `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="104" height="20"><linearGradient id="b" x2="0" y2="100%"><stop offset="0" stop-color="#bbb" stop-opacity=".1"/><stop offset="1" stop-opacity=".1"/></linearGradient><clipPath id="a"><rect width="104" height="20" rx="3" fill="#fff"/></clipPath><g clip-path="url(#a)"><path fill="#555" d="M0 0h51v20H0z"/><path fill="#dfb317" d="M51 0h53v20H51z"/><path fill="url(#b)" d="M0 0h104v20H0z"/></g><g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="110"><text x="265" y="150" fill="#010101" fill-opacity=".3" transform="scale(.1)" textLength="410">brigade</text><text x="265" y="140" transform="scale(.1)" textLength="410">brigade</text><text x="765" y="150" fill="#010101" fill-opacity=".3" transform="scale(.1)" textLength="430">running</text><text x="765" y="140" transform="scale(.1)" textLength="430">running</text></g> </svg>`
	other = `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="104" height="20"><linearGradient id="b" x2="0" y2="100%%"><stop offset="0" stop-color="#bbb" stop-opacity=".1"/><stop offset="1" stop-opacity=".1"/></linearGradient><clipPath id="a"><rect width="104" height="20" rx="3" fill="#fff"/></clipPath><g clip-path="url(#a)"><path fill="#555" d="M0 0h51v20H0z"/><path fill="#dfb317" d="M51 0h53v20H51z"/><path fill="url(#b)" d="M0 0h104v20H0z"/></g><g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="110"><text x="265" y="150" fill="#010101" fill-opacity=".3" transform="scale(.1)" textLength="410">brigade</text><text x="265" y="140" transform="scale(.1)" textLength="410">brigade</text><text x="765" y="150" fill="#010101" fill-opacity=".3" transform="scale(.1)" textLength="430">%s</text><text x="765" y="140" transform="scale(.1)" textLength="430">%s</text></g> </svg>`
)

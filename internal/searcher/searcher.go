package searcher

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/ahmetartuc/dighub/internal/dorks"
)

func Run(org, token string) error {
	client := &http.Client{}

	for _, dork := range dorks.Dorks {
	retry:
		query := fmt.Sprintf("org:%s %s", org, dork)
		urlStr := fmt.Sprintf("https://api.github.com/search/code?q=%s", url.QueryEscape(query))

		req, _ := http.NewRequest("GET", urlStr, nil)
		req.Header.Set("Authorization", "token "+token)
		req.Header.Set("Accept", "application/vnd.github+json")

		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)

		if resp.StatusCode == 200 {
			var parsed map[string]interface{}
			json.Unmarshal(body, &parsed)
			items := parsed["items"].([]interface{})
			if len(items) > 0 {
				color.Green("[+] Dork matched: %s (%d results)", dork, len(items))
				for _, item := range items {
					it := item.(map[string]interface{})
					fmt.Println("    -", it["html_url"])
				}
			} else {
				color.Yellow("[-] No result for: %s", dork)
			}
		} else if resp.StatusCode == 403 && strings.Contains(string(body), "rate limit") {
			waitUntilReset(token)
			goto retry
		} else {
			color.Red("[!] GitHub API error: %d %s", resp.StatusCode, string(body))
		}

		time.Sleep(3 * time.Second)
	}
	return nil
}

func waitUntilReset(token string) {
	url := "https://api.github.com/rate_limit"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "token "+token)
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		time.Sleep(60 * time.Second)
		return
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&data)

	if resources, ok := data["resources"].(map[string]interface{}); ok {
		if search, ok := resources["search"].(map[string]interface{}); ok {
			if reset, ok := search["reset"].(float64); ok {
				resetTime := time.Unix(int64(reset), 0)
				waitTime := time.Until(resetTime) + 5*time.Second
				color.Cyan("[*] Search rate limit hit. Waiting %v until %s", waitTime.Round(time.Second), resetTime.Format("15:04:05"))
				time.Sleep(waitTime)
			}
		}
	}
}

package searcher

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/fatih/color"
	"github.com/ahmetartuc/dighub/internal/dorks"
)

func Run(org, token string) error {
	for _, dork := range dorks.Dorks {
		query := fmt.Sprintf("org:%s %s", org, dork)
		urlStr := fmt.Sprintf("https://api.github.com/search/code?q=%s", url.QueryEscape(query))

		req, _ := http.NewRequest("GET", urlStr, nil)
		req.Header.Set("Authorization", "token "+token)
		req.Header.Set("Accept", "application/vnd.github+json")

		client := &http.Client{}
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
		} else {
			color.Red("[!] GitHub API error: %d %s", resp.StatusCode, string(body))
		}

		time.Sleep(2 * time.Second)
	}
	return nil
}

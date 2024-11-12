package health

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/marshallku/statusy/config"
	"github.com/marshallku/statusy/utils"
)

func Check(cfg *config.Config) {
	var wg sync.WaitGroup
	for _, page := range cfg.Pages {
		wg.Add(1)
		go func(p config.Page) {
			defer wg.Done()
			checkPage(cfg, p)
		}(page)
	}
	wg.Wait()
}

func checkPage(cfg *config.Config, page config.Page) {
	client := &http.Client{
		Timeout: time.Duration(cfg.Timeout) * time.Millisecond,
	}

	var req *http.Request
	var err error

	if page.Request != nil {
		req, err = http.NewRequest(page.Request.Method, page.URL, strings.NewReader(page.Request.Body))
		if err != nil {
			utils.SendNotification(cfg, utils.NotificationParams{
				Description: "🚫 Failed to create request",
				Color:       "16007990",
				Fields: map[string]string{
					"URL": page.URL,
				},
			})
			return
		}
		for key, value := range page.Request.Headers {
			req.Header.Set(key, value)
		}
	} else {
		req, err = http.NewRequest("GET", page.URL, nil)
		if err != nil {
			utils.SendNotification(cfg, utils.NotificationParams{
				Description: "🚫 Failed to create request",
				Color:       "16007990",
				Fields: map[string]string{
					"URL": page.URL,
				},
			})
			return
		}
	}

	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		utils.SendNotification(cfg, utils.NotificationParams{
			Description: "🚫 Failed to connect to server",
			Color:       "16007990",
			Fields: map[string]string{
				"URL": page.URL,
			},
		})
		return
	}
	defer resp.Body.Close()

	duration := time.Since(start)
	body, _ := io.ReadAll(resp.Body)
	timeTaken := fmt.Sprintf("%.3f ms", float64(duration.Milliseconds()))

	if page.Speed > 0 && duration.Milliseconds() > int64(page.Speed) {
		utils.SendNotification(cfg, utils.NotificationParams{
			Description: "🐌 Server responded successfully, but it was too slow.",
			Color:       "16761095",
			Fields: map[string]string{
				"URL":         page.URL,
				"Status Code": fmt.Sprintf("%d", resp.StatusCode),
				"Time Taken":  timeTaken,
			},
		})
		return
	}

	expectedStatus := page.Status
	if expectedStatus == 0 {
		expectedStatus = 200
	}

	if expectedStatus != resp.StatusCode {
		utils.SendNotification(cfg, utils.NotificationParams{
			Description: fmt.Sprintf("🙅 Expected status is %d, but actual status is %d", page.Status, resp.StatusCode),
			Color:       "16007990",
			Fields: map[string]string{
				"URL":         page.URL,
				"Status Code": fmt.Sprintf("%d", resp.StatusCode),
				"Time Taken":  timeTaken,
			},
		})
		return
	}

	if page.TextToInclude != "" && !strings.Contains(string(body), page.TextToInclude) {
		utils.SendNotification(cfg, utils.NotificationParams{
			Description: fmt.Sprintf("😑 String `%s` not found in HTTP response", page.TextToInclude),
			Color:       "16007990",
			Fields: map[string]string{
				"URL":         page.URL,
				"Status Code": fmt.Sprintf("%d", resp.StatusCode),
				"Time Taken":  timeTaken,
			},
		})
		return
	}

	fmt.Printf("Succeeded: %s with status %d\n", page.URL, resp.StatusCode)
}

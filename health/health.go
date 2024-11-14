package health

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/marshallku/statusy/config"
	"github.com/marshallku/statusy/store"
	"github.com/marshallku/statusy/types"
	"github.com/marshallku/statusy/utils"
)

const (
	UP                         = "UP"
	DOWN                       = "DOWN"
	MicrosecondsInMilliSeconds = 1000
	MicrosecondsInSecond       = 1000000
)

func Check(cfg *config.Config, store *store.Store) {
	var wg sync.WaitGroup
	for _, page := range cfg.Pages {
		wg.Add(1)
		go func(p config.Page) {
			defer wg.Done()
			result := checkPage(cfg, p)
			if store != nil {
				store.UpdateResult(result)
			}
		}(page)
	}
	wg.Wait()
}

func checkPage(cfg *config.Config, page config.Page) types.CheckResult {
	client := &http.Client{
		Timeout: time.Duration(cfg.Timeout) * time.Millisecond,
	}

	var req *http.Request
	var err error
	checkedAt := time.Now()

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
			return types.CheckResult{
				URL:         page.URL,
				StatusCode:  0,
				TimeTaken:   "0",
				Status:      false,
				LastChecked: checkedAt,
			}
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
			return types.CheckResult{
				URL:         page.URL,
				StatusCode:  0,
				TimeTaken:   "0",
				Status:      false,
				LastChecked: checkedAt,
			}
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
		return types.CheckResult{
			URL:         page.URL,
			StatusCode:  0,
			TimeTaken:   "0",
			Status:      false,
			LastChecked: checkedAt,
		}
	}
	defer resp.Body.Close()

	duration := time.Since(start)
	body, _ := io.ReadAll(resp.Body)
	timeTakenInMicroseconds := duration.Microseconds()
	timeTaken := fmt.Sprintf("%.3f ms", float64(duration.Microseconds())/MicrosecondsInMilliSeconds)

	if timeTakenInMicroseconds > MicrosecondsInSecond {
		timeTaken = fmt.Sprintf("%.3f s", float64(timeTakenInMicroseconds)/MicrosecondsInSecond)
	}

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
		return types.CheckResult{
			URL:         page.URL,
			StatusCode:  resp.StatusCode,
			TimeTaken:   timeTaken,
			Status:      true,
			LastChecked: checkedAt,
		}
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
		return types.CheckResult{
			URL:         page.URL,
			StatusCode:  resp.StatusCode,
			TimeTaken:   timeTaken,
			Status:      false,
			LastChecked: checkedAt,
		}
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
		return types.CheckResult{
			URL:         page.URL,
			StatusCode:  resp.StatusCode,
			TimeTaken:   timeTaken,
			Status:      false,
			LastChecked: checkedAt,
		}
	}

	fmt.Printf("Succeeded: %s with status %d\n", page.URL, resp.StatusCode)
	return types.CheckResult{
		URL:         page.URL,
		StatusCode:  resp.StatusCode,
		TimeTaken:   timeTaken,
		Status:      true,
		LastChecked: checkedAt,
	}
}

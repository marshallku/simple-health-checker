// health/health_test.go
package health

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/marshallku/statusy/config"
	"github.com/stretchr/testify/assert"
)

func TestCheckPage(t *testing.T) {
	tests := []struct {
		name           string
		serverHandler  http.HandlerFunc
		cfg            *config.Config
		page           config.Page
		expectedStatus bool
	}{
		{
			name: "successful GET request",
			serverHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "GET", r.Method)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("success"))
			}),
			cfg: &config.Config{
				Timeout: 5000,
			},
			page: config.Page{
				URL: "", // Will be set to test server URL
			},
			expectedStatus: true,
		},
		{
			name: "custom status code check",
			serverHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusCreated)
			}),
			cfg: &config.Config{
				Timeout: 5000,
			},
			page: config.Page{
				URL:    "", // Will be set to test server URL
				Status: http.StatusCreated,
			},
			expectedStatus: true,
		},
		{
			name: "text inclusion check",
			serverHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("specific content"))
			}),
			cfg: &config.Config{
				Timeout: 5000,
			},
			page: config.Page{
				URL:           "", // Will be set to test server URL
				TextToInclude: "specific content",
			},
			expectedStatus: true,
		},
		{
			name: "custom headers check",
			serverHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "test-value", r.Header.Get("Test-Header"))
				w.WriteHeader(http.StatusOK)
			}),
			cfg: &config.Config{
				Timeout: 5000,
			},
			page: config.Page{
				URL: "", // Will be set to test server URL
				Request: &config.Request{
					Method: "GET",
					Headers: map[string]string{
						"Test-Header": "test-value",
					},
				},
			},
			expectedStatus: true,
		},
		{
			name: "speed threshold exceeded",
			serverHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(200 * time.Millisecond)
				w.WriteHeader(http.StatusOK)
			}),
			cfg: &config.Config{
				Timeout: 5000,
			},
			page: config.Page{
				URL:   "", // Will be set to test server URL
				Speed: 100,
			},
			expectedStatus: true, // Status is true but with slow response notification
		},
		{
			name: "incorrect status code",
			serverHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
			}),
			cfg: &config.Config{
				Timeout: 5000,
			},
			page: config.Page{
				URL:    "", // Will be set to test server URL
				Status: http.StatusOK,
			},
			expectedStatus: false,
		},
		{
			name: "text not found",
			serverHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("different content"))
			}),
			cfg: &config.Config{
				Timeout: 5000,
			},
			page: config.Page{
				URL:           "", // Will be set to test server URL
				TextToInclude: "expected content",
			},
			expectedStatus: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.serverHandler)
			defer server.Close()

			tt.page.URL = server.URL
			result := checkPage(tt.cfg, tt.page)

			assert.Equal(t, tt.expectedStatus, result.Status)
			assert.NotEmpty(t, result.TimeTaken)
			assert.NotZero(t, result.LastChecked)
			if tt.expectedStatus {
				assert.Greater(t, result.StatusCode, 0)
			}
		})
	}
}

func TestCheckPage_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := &config.Config{
		Timeout: 100, // 100ms timeout
	}
	page := config.Page{
		URL: server.URL,
	}

	result := checkPage(cfg, page)
	assert.False(t, result.Status)
	assert.Equal(t, 0, result.StatusCode)
	assert.Equal(t, "0", result.TimeTaken)
}

func TestCheckPage_InvalidURL(t *testing.T) {
	cfg := &config.Config{
		Timeout: 5000,
	}
	page := config.Page{
		URL: "invalid-url",
	}

	result := checkPage(cfg, page)
	assert.False(t, result.Status)
	assert.Equal(t, 0, result.StatusCode)
	assert.Equal(t, "0", result.TimeTaken)
}

func TestCheckPage_CustomRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := &config.Config{
		Timeout: 5000,
	}
	page := config.Page{
		URL: server.URL,
		Request: &config.Request{
			Method: "POST",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: `{"test": true}`,
		},
	}

	result := checkPage(cfg, page)
	assert.True(t, result.Status)
	assert.Equal(t, http.StatusOK, result.StatusCode)
}

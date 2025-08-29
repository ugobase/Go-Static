package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestEndpoints checks all routes in a table-driven test
func TestEndpoints(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		wantStatus int
		wantCT     string
	}{
		{"Home", "/home", http.StatusOK, "text/html"},
		{"About", "/about", http.StatusOK, "text/html"},
		{"Contact", "/contact", http.StatusOK, "text/html"},
		{"Health", "/health", http.StatusOK, ""},
	}

	// Setup mux using real handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/home", servePage("home"))
	mux.HandleFunc("/about", servePage("about"))
	mux.HandleFunc("/contact", servePage("contact"))
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("ok")); err != nil {
			t.Logf("write response error: %v", err)
		}
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			resp := w.Result()
			if resp.StatusCode != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, resp.StatusCode)
			}

			if tt.wantCT != "" {
				if ct := resp.Header.Get("Content-Type"); ct != tt.wantCT {
					t.Errorf("expected Content-Type %s, got %s", tt.wantCT, ct)
				}
			}
		})
	}
}


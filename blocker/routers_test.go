package blocker

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestEverythingBlocked(t *testing.T) {
	type TestCase struct {
		name   string
		wanted int
	}

	testCases := []TestCase{
		{
			name:   "everything is forbidden",
			wanted: http.StatusForbidden,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			b := NewBlocker("1994", true, []string{}, time.Now(), time.Now())

			host := "google.com"
			url := fmt.Sprintf("http://%s", host)

			r := httptest.NewRequest("GET", url, strings.NewReader(""))
			w := httptest.NewRecorder()

			http.HandlerFunc(b.blockAllHandler).ServeHTTP(w, r)

			if w.Code != tc.wanted {
				t.Fatalf("%s failed: wanted: %v, got: %v", tc.name, tc.wanted, w.Code)
			}
		})
	}
}

func TestAllowEverything(t *testing.T) {
	type TestCase struct {
		name   string
		wanted int
	}

	testCases := []TestCase{
		{
			name:   "everything allowed",
			wanted: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			b := NewBlocker("1994", false, []string{}, time.Time{}, time.Time{})

			host := "google.com"
			url := fmt.Sprintf("http://%s", host)

			r := httptest.NewRequest("GET", url, strings.NewReader(""))
			w := httptest.NewRecorder()

			http.HandlerFunc(b.routeAllRequestsHandler).ServeHTTP(w, r)

			if w.Code != tc.wanted {
				t.Fatalf("%s failed: wanted: %v, got: %v", tc.name, tc.wanted, w.Code)
			}
		})
	}
}

func TestIsHostInBlockList(t *testing.T) {
	type TestCase struct {
		name   string
		host   string
		wanted bool
	}

	testCases := []TestCase{
		{
			name:   "host is in block list",
			host:   "google.com",
			wanted: true,
		},
		{
			name:   "host is not in block list",
			host:   "collogne.de",
			wanted: false,
		},
		{
			name:   "empty host",
			host:   "",
			wanted: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			b := NewBlocker("1994", false, []string{}, time.Time{}, time.Time{})
			b.UrlList = append(b.UrlList, "google.com")
			isInList := b.isHostInBlockList(tc.host)

			if isInList != tc.wanted {
				t.Fatalf("%s failed: wanted: %v, got: %v", tc.name, tc.wanted, isInList)
			}
		})
	}
}

func TestRemoveHostFromBlockList(t *testing.T) {
	type TestCase struct {
		name         string
		hostToRemove string
		wanted       bool
	}

	testCases := []TestCase{
		{
			name:         "host is removed correctly",
			hostToRemove: "google.com",
			wanted:       false,
		},
		{
			name:         "host initial not in list",
			hostToRemove: "collogne.de",
			wanted:       false,
		},
		{
			name:         "empty host",
			hostToRemove: "",
			wanted:       false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			b := NewBlocker("1994", false, []string{}, time.Time{}, time.Time{})
			b.UrlList = append(b.UrlList, "google.com", "hildegard.de")
			b.removeHostFromBlockList(tc.hostToRemove)
			isInList := b.isHostInBlockList(tc.hostToRemove)
			if isInList != tc.wanted {
				t.Fatalf("%s failed: wanted: %v, got: %v", tc.name, tc.wanted, isInList)
			}
		})
	}
}

func TestTimeIsInWindow(t *testing.T) {
	type TestCase struct {
		name         string
		suppliedTime string
		wanted       bool
	}

	testCases := []TestCase{
		{
			name:         "time in window",
			suppliedTime: "10:00",
			wanted:       true,
		},
		{
			name:         "time not in window",
			suppliedTime: "08:00",
			wanted:       false,
		},
		{
			name:         "time is invalid",
			suppliedTime: "700",
			wanted:       false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			b := NewBlocker("1994", false, []string{}, time.Date(int(0000), time.January, int(1), int(9), int(0), int(0), int(0), time.Now().Local().Location()), time.Date(int(0000), time.January, int(1), int(17), int(0), int(0), int(0), time.Now().Local().Location()))
			ot, _ := parseTime(tc.suppliedTime)

			tiw := b.timeIsInWindow(ot)
			if tiw != tc.wanted {
				t.Fatalf("%s failed: wanted: %v, got: %v", tc.name, tc.wanted, tiw)
			}
		})
	}
}

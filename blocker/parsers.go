package blocker

import (
	"fmt"
	"strings"
	"time"
)

func parseTime(t string) (time.Time, error) {
	timeString := t + ":00"
	pt, err := time.Parse("15:04:05", timeString)
	if err != nil {
		return time.Time{}, err
	}
	return pt, nil
}

func parseList(urllist string) ([]string, error) {
	urlList := strings.Split(urllist, ",")
	if len(urlList) > 0 {
		return urlList, nil
	}
	return nil, fmt.Errorf("parsing url list failed")
}

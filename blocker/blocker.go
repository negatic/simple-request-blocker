package blocker

import "time"

type Blocker struct {
	Port              string
	BlockEveryRequest bool
	UrlList           []string
	StartTime         time.Time
	EndTime           time.Time
}

func NewBlocker(port string, blockEveryRequest bool, urlList []string, startTime time.Time, endTime time.Time) *Blocker {
	return &Blocker{
		Port:              port,
		BlockEveryRequest: blockEveryRequest,
		UrlList:           urlList,
		StartTime:         startTime,
		EndTime:           endTime,
	}
}

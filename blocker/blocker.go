package blocker

import (
	"flag"
	"time"
)

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

func Run() error {
	blockEveryRequest := flag.Bool("blockeveryrequest", false, "Set to true to block all requests")
	blockList := flag.String("blocklist", "", "Sites to block,seprate by a comma. Example: --blocklist google.com,github.com")
	port := flag.String("port", "1080", "the port to listen on, Example: --port 1080")
	startTime := flag.String("starttime", "10:00", "The time to start the blocking. Example: --starttime 10:00")
	endTime := flag.String("endtime", "12:00", "The time to end the blocking. Example: --endtime 12:00")
	flag.Parse()
}

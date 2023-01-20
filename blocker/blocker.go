package blocker

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

var (
	startBlockTime = "00:00"
	endBlockTime   = "17:00"
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

	list, err := parseList(*blockList)
	if err != nil {
		return err
	}

	if *startTime != "" {
		startBlockTime = *startTime
	}
	if *endTime != "" {
		endBlockTime = *endTime
	}

	st, err := parseTime(startBlockTime)
	if err != nil {
		return err
	}
	et, err := parseTime(endBlockTime)
	if err != nil {
		return err
	}

	blocker := NewBlocker(*port, *blockEveryRequest, list, st, et)
	router := blocker.CreateRouter()

	fmt.Printf("The proxy runs on port %s...", *port)
	fmt.Printf("Following sites will be blocked: %s", blocker.UrlList)

	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%s", blocker.Port), router))
	return nil

}

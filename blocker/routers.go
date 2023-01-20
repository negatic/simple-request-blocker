package blocker

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func (b *Blocker) CreateRouter() *mux.Router {
	r := mux.NewRouter()

	if b.BlockEveryRequest {
		r.HandleFunc("/", b.blockAllHandler)
		return r
	}

	r.HandleFunc("/admin/{command}/{host}", b.adminHostConfigurationHandler)
	r.HandleFunc("/", b.blockRequestsFromList)
	return r
}

func (b *Blocker) blockAllHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("the requested URL %s is blocked", r.URL.Path)
	http.Error(w, "FORBIDDEN REQUEST", http.StatusForbidden)
}

func (b *Blocker) timeIsInWindow(tnow time.Time) bool {
	currentTime := time.Date(int(0000), time.January, int(1), tnow.Hour(), tnow.Minute(), tnow.Second(), tnow.Nanosecond(), time.Now().Local().Location())
	if currentTime.After(b.StartTime) && currentTime.Before(b.EndTime) {
		return true
	}
	return false
}

func (b *Blocker) isHostInBlockList(host string) bool {
	for _, s := range b.UrlList {
		if s == host {
			return true
		}
	}
	return false
}

func (b *Blocker) blockRequestsFromList(w http.ResponseWriter, r *http.Request) {
	if b.timeIsInWindow(time.Now()) {
		fmt.Printf("Yes time is in window.")
		if b.isHostInBlockList(r.URL.Host) {
			fmt.Printf("Yes host is in blocklist.")
			b.blockAllHandler(w, r)
			return
		}
	}

	b.routeAllRequestsHandler(w, r)
}

func (b *Blocker) adminHostConfigurationHandler(w http.ResponseWriter, r *http.Request) {
	adminMethod := mux.Vars(r)["command"]
	switch adminMethod {
	case "":
		http.Error(w, "no valid admin command", http.StatusBadRequest)
		return
	case "block":
		host := mux.Vars(r)["host"]
		if host == "" {
			http.Error(w, "no valid host to block", http.StatusBadRequest)
			return
		}
		if !b.isHostInBlockList(host) {
			b.UrlList = append(b.UrlList, host)
		}
		w.Write([]byte(fmt.Sprintf("the host %s was successfully blocked", host)))
		return
	case "unblock":
		host := mux.Vars(r)["host"]
		if host == "" {
			http.Error(w, "no valid host to unblock", http.StatusBadRequest)
			return
		}
		if b.isHostInBlockList(host) {
			b.removeHostFromBlockList(host)
		}
		w.Write([]byte(fmt.Sprintf("the host %s was successfully unblocked", host)))
	}
}

func (b *Blocker) removeHostFromBlockList(host string) {
	for i, v := range b.UrlList {
		if v == host {
			b.UrlList = append(b.UrlList[:i], b.UrlList[i+1:]...)
			break
		}
	}
}

func (b *Blocker) routeAllRequestsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("the request with Host %s is routed", r.Host)
	res, err := http.Get(r.URL.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Write(body)
}

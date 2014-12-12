package main

import (
        "bytes"
        "fmt"
        "log"
        "net/http"
        "os"
        "time"

        "github.com/PuerkitoBio/goquery"
)

var result string

func handler(w http.ResponseWriter, r *http.Request) {
        log.Printf("Request from %s\n", r.RemoteAddr)

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(200)
        fmt.Fprint(w, result)
}

func crawl() {
        start := time.Now()
        var buffer bytes.Buffer

        doc, err := goquery.NewDocument("https://www.yahoo.com")
        if err != nil {
                log.Fatal(err)
        }

        buffer.WriteString("document.domain = 'naver.com';")
        doc.Find(".topics-wrapper").Each(func(i int, s1 *goquery.Selection) {
                buffer.WriteString("var obj =[")
                s1.Find(".d-b, .pl-l, .ell").Each(func(j int, s2 *goquery.Selection) {
                        buffer.WriteString("{K: \"")
                        buffer.WriteString(s2.Text())
                        buffer.WriteString("\"}")
                        if j != 9 {
                                buffer.WriteString(",")
                        }
                })
                buffer.WriteString("];")
                result = buffer.String()
                log.Printf("End crawling: %s elapsed\n", time.Since(start))
                return
        })
}

func startCrawler() {
        go func() {
                for {
                        select {
                        case <-time.After(time.Second * 60):
                                log.Print("Begin crawling")
                                go crawl()
                        }
                }
        }()
        crawl()
}

func main() {
        log.SetOutput(os.Stdout)
        log.Print("Starting server...")

        go startCrawler()

        for {
                if result != "" {
                        break
                }
                log.Print("Waiting for first crawl to finish...")
                time.Sleep(time.Second * 1)
        }

        http.HandleFunc("/", handler)

        err := http.ListenAndServe(":29100", nil)
        if err != nil {
                log.Fatal(err)
        }
}

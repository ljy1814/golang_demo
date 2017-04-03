package visit

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	visit_chan chan Visit = make(chan Visit, 1)
)

func main() {
	config := GetConfig()
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "error parameters num. Usage ./visit_analytics config.yaml")
		os.Exit(1)
	}
	config.parse(os.Args[1])

	InitDB()

	router := gin.Default()

	router.GET("/analytics.js", func(c *gin.Context) {
		go analyse(c)
		c.Header("Content-Type", "application/javascript")
		c.String(http.StatusOK, "console.log(\"https://github.com/nladuo/visit_analytics\")")
	})

	router.StaticFile("/test", "./www/test.html")
	router.StaticFile("/test2", "./www/test.html")

	MakeRouters(router)

	go func() {
		for {
			visit := <-visit_chan
			handleVisit(visit)
		}
	}()
	router.Run(config.RunAt)
}

func analyse(c *gin.Context) {
	referer := TrimUrl(c.Request.Referer())
	host_name := GetHostName(referer)
	if len(referer) == 0 || host_name == "" {
		return
	}
	title := GetTitle(referer)
	visit_chan <- Visit{
		ClientIp:  c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
		Referer:   referer,
		Title:     title,
		Host:      host_name,
	}
}

func handleVisit(visit Visit) {
	recordHost(visit)
	recordPage(visit)
	recordDailyRecord(visit)
	recordMonthlyRecord(visit)
}

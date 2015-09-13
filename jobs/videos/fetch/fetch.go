package fetch

import (
	"fmt"
	"github.com/desmondhume/furrytemple/carpenter"
	"github.com/desmondhume/furrytemple/crawler"
	"github.com/desmondhume/furrytemple/parser"
)

const (
	keywords = []string{"cat", "kitten", "cute kitten", "grumpy cat", "playing cat", "fighting cat", "meowing cats", "cat rescue"}

	channels = map[string]chan []byte{
		"youtube": make(chan []byte),
	}
)

func Run(output chan map[string]interface{}, jobsReports chan map[string]string) {
	for source, input := range channels {
		for _, keyword := range keywords {
			go func(src string, input chan []byte, kw string) {
				crawler.Crawl(kw, input)
				parser.Parse(src, input, output)
			}(source, input, keyword)
		}
	}

	for _ = range output {
		go func() {
			carpenter.Exec("populate", output, jobsReports)
		}()
		for report := range jobsReports {
			fmt.Println(report)
		}
	}
}

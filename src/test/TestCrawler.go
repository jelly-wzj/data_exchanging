package test

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

// Article 抓取blog数据
type Article struct {
	ID       int    `json:"id,omitempty"`
	Title    string `json:"title,omitempty"`
	URL      string `json:"url,omitempty"`
	Created  string `json:"created,omitempty"`
	Reads    string `json:"reads,omitempty"`
	Comments string `json:"comments,omitempty"`
	Feeds    string `json:"feeds,omitempty"`
}

// 数据持久化
func csvSave(fName string, data []Article) error {
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"ID", "Title", "URL", "Created", "Reads", "Comments", "Feeds"})
	for _, v := range data {
		writer.Write([]string{strconv.Itoa(v.ID), v.Title, v.URL, v.Created, v.Reads, v.Comments, v.Feeds})
	}
	return nil
}

func TestCrawler() {
	articles := make([]Article, 0, 200)
	// 1.准备收集器实例
	c := colly.NewCollector(
		// 开启本机debug
		// colly.Debugger(&debug.LogDebugger{}),
		colly.AllowedDomains("learnku.com"),
		// 防止页面重复下载
		// colly.CacheDir("./learnku_cache"),
	)

	// 2.分析页面数据
	c.OnHTML("div.blog-article-list > .event", func(e *colly.HTMLElement) {
		article := Article{
			Title: e.ChildText("div.content > div.summary"),
			URL:   e.ChildAttr("div.content a.title", "href"),
			Feeds: e.ChildText("div.item-meta > a:first-child"),
		}
		// 查找同一集合不同子项
		e.ForEach("div.content > div.meta > div.date>a", func(i int, el *colly.HTMLElement) {
			switch i {
			case 1:
				article.Created = el.Attr("data-tooltip")
			case 2:
				// 用空白切割字符串
				article.Reads = strings.Fields(el.Text)[1]
			case 3:
				article.Comments = strings.Fields(el.Text)[1]
			}
		})
		// 正则匹配替换,字符串转整型
		article.ID, _ = strconv.Atoi(regexp.MustCompile(`\d+`).FindAllString(article.URL, -1)[0])
		articles = append(articles, article)
	})

	// 下一页
	c.OnHTML("a[href].page-link", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	// 启动
	c.Visit("https://learnku.com/blog/pardon")

	// 输出
	csvSave("pardon.csv", articles)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(articles)

	// 显示收集器的打印信息
	log.Println(c)
}

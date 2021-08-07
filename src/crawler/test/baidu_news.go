package test

import (
	"data_exchanging/src/crawler"
	log "github.com/sirupsen/logrus"
)

func init() {
	crawler.Register(rule)
}

var (
	outputFields  = []string{"category", "title", "link"}
	outputFields2 = []string{"category", "category_link"}

	namespace1 = "baidu_news"
	namespace2 = "baidu_category"
)
var multNamespaceConf = map[string]*crawler.MultipleNamespaceConf{
	namespace1: {
		OutputFields:      outputFields,
		OutputConstraints: crawler.NewStringsConstraints(outputFields, 64, 128, 512),
	},
	namespace2: {
		OutputFields:      outputFields2,
		OutputConstraints: crawler.NewStringsConstraints(outputFields2, 64, 256),
	},
}

// 演示如何在一条规则里面，同时需要导出数据到两张表
var rule = &crawler.TaskRule{
	Name:                      "百度新闻规则",
	Description:               "抓取百度新闻各个分类的最新焦点新闻以及最新的新闻分类和链接",
	OutputToMultipleNamespace: true,
	MultipleNamespaceConf:     multNamespaceConf,
	Rule: &crawler.Rule{
		Head: func(ctx *crawler.Context) error {
			return ctx.VisitForNext("http://news.baidu.com")
		},
		Nodes: map[int]*crawler.Node{
			0: step1, // 第一步: 获取所有分类
			1: step2, // 第二步: 获取每个分类的新闻标题链接
		},
	},
}

var step1 = &crawler.Node{
	OnRequest: func(ctx *crawler.Context, req *crawler.Request) {
		log.Println("Visting", req.URL.String())
	},
	OnHTML: map[string]func(*crawler.Context, *crawler.HTMLElement) error{
		`#channel-all .menu-list a`: func(ctx *crawler.Context, el *crawler.HTMLElement) error { // 获取所有分类
			category := el.Text
			ctx.PutReqContextValue("category", category)
			link := el.Attr("href")

			if category != "首页" {
				err := ctx.Output(map[int]interface{}{
					0: category,
					1: ctx.AbsoluteURL(link),
				}, namespace2)
				if err != nil {
					return err
				}
			}

			return ctx.VisitForNextWithContext(link)
		},
	},
}

var step2 = &crawler.Node{
	OnRequest: func(ctx *crawler.Context, req *crawler.Request) {
		log.Println("Visting", req.URL.String())
	},
	OnHTML: map[string]func(*crawler.Context, *crawler.HTMLElement) error{
		`#col_focus a, .focal-news a, .auto-col-focus a, .l-common .fn-c a`: func(ctx *crawler.Context, el *crawler.HTMLElement) error {
			title := el.Text
			link := el.Attr("href")
			if title == "" || link == "javascript:void(0);" {
				return nil
			}

			category := ctx.GetReqContextValue("category")
			return ctx.Output(map[int]interface{}{
				0: category,
				1: title,
				2: link,
			}, namespace1)
		},
	},
}

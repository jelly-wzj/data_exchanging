package test

import (
	"data_exchanging/src/crawler"
	_ "data_exchanging/src/crawler/test"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strings"
	"time"
)

func TestRuleCrawler() {
	task := &Task{
		TaskRuleName:      "百度新闻规则",
		OptAllowedDomains: "news.baidu.com",
		OutputType:        "csv",
	}
	rule, _ := crawler.GetTaskRule(task.TaskRuleName)

	var optAllowedDomains []string
	if task.OptAllowedDomains != "" {
		optAllowedDomains = strings.Split(task.OptAllowedDomains, ",")
	}

	var urlFiltersReg []*regexp.Regexp
	if task.OptURLFilters != "" {
		urlFilters := strings.Split(task.OptURLFilters, ",")
		for _, v := range urlFilters {
			reg, err := regexp.Compile(v)
			if err != nil {
				log.Print(err)
			}
			urlFiltersReg = append(urlFiltersReg, reg)
		}
	}

	config := crawler.TaskConfig{
		CronSpec: task.CronSpec,
		Option: crawler.Option{
			UserAgent:              task.OptUserAgent,
			MaxDepth:               task.OptMaxDepth,
			AllowedDomains:         optAllowedDomains,
			URLFilters:             urlFiltersReg,
			AllowURLRevisit:        rule.AllowURLRevisit,
			MaxBodySize:            task.OptMaxBodySize,
			IgnoreRobotsTxt:        rule.IgnoreRobotsTxt,
			InsecureSkipVerify:     rule.InsecureSkipVerify,
			ParseHTTPErrorResponse: rule.ParseHTTPErrorResponse,
			DisableCookies:         rule.DisableCookies,
		},
		Limit: crawler.Limit{
			Enable:      task.LimitEnable,
			DomainGlob:  task.LimitDomainGlob,
			Delay:       time.Duration(task.LimitDelay) * time.Millisecond,
			RandomDelay: time.Duration(task.LimitRandomDelay) * time.Millisecond,
			Parallelism: task.LimitParallelism,
		},
		OutputConfig: crawler.OutputConfig{
			Type: task.OutputType,
		},
	}

	if task.OutputType == crawler.OutputTypeCSV {
		config.OutputConfig.CSVConf = crawler.CSVConf{
			CSVFilePath: "./csv_output",
		}
	}

	c := make(chan crawler.MTS)

	crawler.New(crawler.NewTask(1, *rule, config), c).Run()
}

type Task struct {
	ID               uint64             `json:"id,string" gorm:"column:id;type:bigint unsigned AUTO_INCREMENT;primary_key"`
	TaskName         string             `json:"task_name" gorm:"column:task_name;type:varchar(64);not null;unique_index:uk_task_name"`
	TaskRuleName     string             `json:"task_rule_name" gorm:"column:task_rule_name;type:varchar(64);not null"`
	TaskDesc         string             `json:"task_desc" gorm:"column:task_desc;type:varchar(512);not null;default:''"`
	Status           crawler.TaskStatus `json:"status" gorm:"column:status;type:tinyint;not null;default:'0'"`
	Counts           int                `json:"counts" gorm:"column:counts;type:int;not null;default:'0'"`
	CronSpec         string             `json:"cron_spec" gorm:"column:cron_spec;type:varchar(64);not null;default:''"`
	OutputType       string             `json:"output_type" gorm:"column:output_type;type:varchar(64);not null;"`
	OutputExportDBID uint64             `json:"output_exportdb_id" gorm:"column:output_exportdb_id;type:bigint;not null;default:'0'"`
	// 参数配置部分
	OptUserAgent      string `json:"opt_user_agent" gorm:"column:opt_user_agent;type:varchar(128);not null;default:''"`
	OptMaxDepth       int    `json:"opt_max_depth" gorm:"column:opt_max_depth;type:int;not null;default:'0'"`
	OptAllowedDomains string `json:"opt_allowed_domains" gorm:"column:opt_allowed_domains;type:varchar(512);not null;default:''"`
	OptURLFilters     string `json:"opt_url_filters" gorm:"column:opt_url_filters;type:varchar(512);not null;default:''"`
	OptMaxBodySize    int    `json:"opt_max_body_size" gorm:"column:opt_max_body_size;type:int;not null;default:'0'"`
	OptRequestTimeout int    `json:"opt_request_timeout" gorm:"column:opt_request_timeout;type:int;not null;default:'10'"`
	// auto migrate
	AutoMigrate bool `json:"auto_migrate" gorm:"column:auto_migrate;type:tinyint;not null;default:'0'"`

	// 频率限制
	LimitEnable       bool      `json:"limit_enable" gorm:"column:limit_enable;type:tinyint;not null;default:'0'"`
	LimitDomainRegexp string    `json:"limit_domain_regexp" gorm:"column:limit_domain_regexp;type:varchar(128);not null;default:''"`
	LimitDomainGlob   string    `json:"limit_domain_glob" gorm:"column:limit_domain_glob;type:varchar(128);not null;default:''"`
	LimitDelay        int       `json:"limit_delay" gorm:"column:limit_delay;type:int;not null;default:'0'"`
	LimitRandomDelay  int       `json:"limit_random_delay" gorm:"column:limit_random_delay;type:int;not null;default:'0'"`
	LimitParallelism  int       `json:"limit_parallelism" gorm:"column:limit_parallelism;type:int;not null;default:'0'"`
	ProxyURLs         string    `json:"proxy_urls" gorm:"column:proxy_urls;type:varchar(2048);not null;default:''"`
	CreatedAt         time.Time `json:"created_at" gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;index:idx_created_at"`
}

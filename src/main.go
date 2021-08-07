package main

import (
	"data_exchanging/src/crawler"
	_ "data_exchanging/src/crawler/test"
	"github.com/pkg/errors"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	//esapi := new(elasticsearch.EsApiV7)
	//esapi.Host = "http://192.168.31.200:9200"
	//health := esapi.ClusterHealth()
	//fmt.Println(health.Status)

	/*	{
		"settings": {
		"index": {
			"number_of_shards": 3,
				"number_of_replicas": 2
			}
		}
	}*/

	//tempIndexSettings := map[string]interface{}{}
	//tempIndexSettings["settings"] = map[string]interface{}{}
	//tempIndexSettings["settings"].(map[string]interface{})["index"] = map[string]interface{}{}
	//tempIndexSettings["settings"].(map[string]interface{})["index"].(map[string]interface{})["number_of_shards"] = 1
	//tempIndexSettings["settings"].(map[string]interface{})["index"].(map[string]interface{})["number_of_replicas"] = 0
	//esapi.CreateIndex("test_go_index", tempIndexSettings)

	//esapi.DeleteIndex("test_go_index")

	//crawler.TestCrawler()
	t := New()

	t.Run()
}

type GoCrawler struct {
	backend string
	mysql   crawler.MySQLConf
	//web     *web.Server
}

func (gs *GoCrawler) parseSettingsFromEnv() {
	for _, e := range os.Environ() {
		if !strings.HasPrefix(e, "GOCrawler_") {
			continue
		}
		pair := strings.SplitN(e[9:], "=", 2)
		if f, ok := envMap[pair[0]]; ok {
			f(gs, pair[1])
		} else {
			log.Println("Unknown env variable:", pair[0])
		}
	}
}

func (gs *GoCrawler) Run() error {
	gs.print()
	db, err := crawler.NewGormDB(gs.mysql)
	if err != nil {
		return errors.Wrap(err, "new gorm db failed")
	}
	//core.SetGormDB(db)

	return nil
}

func (gs *GoCrawler) print() {
	log.Printf("goCrawler backend conf:%+v\n", gs.mysql)
}

var envMap = map[string]func(*GoCrawler, string){
	"DB_HOST": func(gs *GoCrawler, val string) {
		gs.mysql.Host = val
	},
	"DB_PORT": func(gs *GoCrawler, val string) {
		port, err := strconv.Atoi(val)
		if err == nil {
			gs.mysql.Port = port
		}
	},
	"DB_USER": func(gs *GoCrawler, val string) {
		gs.mysql.User = val
	},
	"DB_PASSWORD": func(gs *GoCrawler, val string) {
		gs.mysql.Password = val
	},
	"DB_NAME": func(gs *GoCrawler, val string) {
		gs.mysql.DBName = val
	},
	//"WEB_IP": func(gs *GoCrawler, val string) {
	//	gs.web.IP = val
	//},
	//"WEB_PORT": func(gs *GoCrawler, val string) {
	//	port, err := strconv.Atoi(val)
	//	if err == nil {
	//		gs.web.Port = port
	//	}
	//},
}

func (gs *GoCrawler) init() {
	gs.backend = "mysql"
	gs.mysql.Host = "192.168.31.200"
	gs.mysql.Port = 3306
	gs.mysql.User = "root"
	gs.mysql.Password = "123456"
	gs.mysql.MaxIdleConns = 3
	gs.mysql.MaxOpenConns = 10
	//gs.web = &web.Server{Port: 8080}
}

func New(opts ...func(*GoCrawler)) *GoCrawler {
	gs := &GoCrawler{}
	gs.init()

	for _, f := range opts {
		f(gs)
	}

	gs.parseSettingsFromEnv()

	return gs
}

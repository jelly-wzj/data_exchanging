package main

import (
	"data_exchanging/src/elasticsearch"
	"fmt"
)

func main() {
	esapi := new(elasticsearch.EsApiV7)
	esapi.Host = "http://192.168.31.200:9200"
	health := esapi.ClusterHealth()
	fmt.Println(health.Status)
	//tempIndexSettings := map[string]interface{}{}
	//tempIndexSettings["settings"] = map[string]interface{}{}
	//tempIndexSettings["settings"].(map[string]interface{})["index"] = map[string]interface{}{}
	//esapi.CreateIndex("test_go_index", tempIndexSettings)

	//esapi.DeleteIndex("test_go_index")

}

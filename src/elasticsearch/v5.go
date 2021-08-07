package elasticsearch

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/cihub/seelog"
	"io"
	"io/ioutil"
	"strings"
)

type EsApiV5 struct {
	EsApiV0
}

func (s *EsApiV5) ClusterHealth() *ClusterHealth {
	return s.EsApiV0.ClusterHealth()
}

func (s *EsApiV5) Bulk(data *bytes.Buffer) {
	s.EsApiV0.Bulk(data)
}

func (s *EsApiV5) GetIndexSettings(indexNames string) (*Indexes, error) {
	return s.EsApiV0.GetIndexSettings(indexNames)
}

func (s *EsApiV5) UpdateIndexSettings(indexName string, settings map[string]interface{}) error {
	return s.EsApiV0.UpdateIndexSettings(indexName, settings)
}

func (s *EsApiV5) GetIndexMappings(copyAllIndexes bool, indexNames string) (string, int, *Indexes, error) {
	return s.EsApiV0.GetIndexMappings(copyAllIndexes, indexNames)
}

func (s *EsApiV5) UpdateIndexMapping(indexName string, settings map[string]interface{}) error {
	return s.EsApiV0.UpdateIndexMapping(indexName, settings)
}

func (s *EsApiV5) DeleteIndex(name string) (err error) {
	return s.EsApiV0.DeleteIndex(name)
}

func (s *EsApiV5) CreateIndex(name string, settings map[string]interface{}) (err error) {
	return s.EsApiV0.CreateIndex(name, settings)
}

func (s *EsApiV5) Refresh(name string) (err error) {
	return s.EsApiV0.Refresh(name)
}

func (s *EsApiV5) NewScroll(indexNames string, scrollTime string, docBufferCount int, query string, slicedId, maxSlicedCount int, fields string) (scroll interface{}, err error) {
	url := fmt.Sprintf("%s/%s/_search?scroll=%s&size=%d", s.Host, indexNames, scrollTime, docBufferCount)

	jsonBody := ""
	if len(query) > 0 || maxSlicedCount > 0 || len(fields) > 0 {
		queryBody := map[string]interface{}{}

		if len(fields) > 0 {
			if !strings.Contains(fields, ",") {
				queryBody["_source"] = fields
			} else {
				queryBody["_source"] = strings.Split(fields, ",")
			}
		}

		if len(query) > 0 {
			queryBody["query"] = map[string]interface{}{}
			queryBody["query"].(map[string]interface{})["query_string"] = map[string]interface{}{}
			queryBody["query"].(map[string]interface{})["query_string"].(map[string]interface{})["query"] = query
		}

		if maxSlicedCount > 1 {
			log.Tracef("sliced scroll, %d of %d", slicedId, maxSlicedCount)
			queryBody["slice"] = map[string]interface{}{}
			queryBody["slice"].(map[string]interface{})["id"] = slicedId
			queryBody["slice"].(map[string]interface{})["max"] = maxSlicedCount
		}

		jsonArray, err := json.Marshal(queryBody)
		if err != nil {
			log.Error(err)

		} else {
			jsonBody = string(jsonArray)
		}
	}

	resp, body, errs := Post(url, s.Auth, jsonBody, s.HttpProxy)

	if resp != nil && resp.Body != nil {
		io.Copy(ioutil.Discard, resp.Body)
		defer resp.Body.Close()
	}

	if errs != nil {
		log.Error(errs)
		return nil, errs[0]
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(body)
	}

	log.Trace("new scroll,", body)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	scroll = &Scroll{}
	err = DecodeJson(body, scroll)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return scroll, err
}

func (s *EsApiV5) NextScroll(scrollTime string, scrollId string) (interface{}, error) {
	id := bytes.NewBufferString(scrollId)

	url := fmt.Sprintf("%s/_search/scroll?scroll=%s&scroll_id=%s", s.Host, scrollTime, id)
	resp, body, errs := Get(url, s.Auth, s.HttpProxy)

	if resp != nil && resp.Body != nil {
		io.Copy(ioutil.Discard, resp.Body)
		defer resp.Body.Close()
	}

	if errs != nil {
		log.Error(errs)
		return nil, errs[0]
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(body)
	}

	// decode elasticsearch scroll response
	scroll := &Scroll{}
	err := DecodeJson(body, &scroll)
	if err != nil {
		log.Error(body)
		log.Error(err)
		return nil, err
	}

	return scroll, nil
}

package elasticsearch

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/parnurzeal/gorequest"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func Get(url string, auth *Auth, proxy string) (*http.Response, string, []error) {
	request := gorequest.New()

	tr := &http.Transport{
		DisableKeepAlives: true,
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
	}
	request.Transport = tr

	if auth != nil {
		request.SetBasicAuth(auth.User, auth.Pass)
	}

	request.Header["Content-Type"] = []string{"application/json"}
	//request.Header.Set("Content-Type", "application/json")

	if len(proxy) > 0 {
		request.Proxy(proxy)
	}

	resp, body, errs := request.Get(url).End()
	return resp, body, errs

}

func Post(url string, auth *Auth, body string, proxy string) (*http.Response, string, []error) {
	request := gorequest.New()
	tr := &http.Transport{
		DisableKeepAlives: true,
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
	}
	request.Transport = tr

	if auth != nil {
		request.SetBasicAuth(auth.User, auth.Pass)
	}

	//request.Header.Set("Content-Type", "application/json")
	request.Header["Content-Type"] = []string{"application/json"}

	if len(proxy) > 0 {
		request.Proxy(proxy)
	}

	request.Post(url)

	if len(body) > 0 {
		request.Send(body)
	}

	return request.End()
}

func newDeleteRequest(client *http.Client, method, urlStr string) (*http.Request, error) {
	if method == "" {
		// We document that "" means "GET" for Request.Method, and people have
		// relied on that from NewRequest, so keep that working.
		// We still enforce validMethod for non-empty methods.
		method = "GET"
	}
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	req := &http.Request{
		Method:     method,
		URL:        u,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		Host:       u.Host,
	}
	return req, nil
}

func Request(method string, r string, auth *Auth, body *bytes.Buffer, proxy string) (string, error) {

	//TODO use global client
	var client *http.Client
	client = &http.Client{}
	if len(proxy) > 0 {
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			log.Error(err)
		} else {
			transport := &http.Transport{
				Proxy:             http.ProxyURL(proxyURL),
				DisableKeepAlives: true,
			}
			client = &http.Client{Transport: transport}
		}
	}

	tr := &http.Transport{
		DisableKeepAlives: true,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client.Transport = tr

	var err error
	var reqest *http.Request
	if body != nil {
		reqest, err = http.NewRequest(method, r, body)
	} else {
		reqest, err = newDeleteRequest(client, method, r)
	}

	if err != nil {
		panic(err)
	}

	if auth != nil {
		reqest.SetBasicAuth(auth.User, auth.Pass)
	}

	reqest.Header.Set("Content-Type", "application/json")

	resp, errs := client.Do(reqest)
	if errs != nil {
		log.Error(errs)
		return "", errs
	}

	if resp != nil && resp.Body != nil {
		//io.Copy(ioutil.Discard, resp.Body)
		defer resp.Body.Close()
	}

	if resp.StatusCode != 200 {
		b, _ := ioutil.ReadAll(resp.Body)
		return "", errors.New("server error: " + string(b))
	}

	respBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Error(err)
		return string(respBody), err
	}

	log.Trace(r, string(respBody))

	if err != nil {
		return string(respBody), err
	}
	io.Copy(ioutil.Discard, resp.Body)
	defer resp.Body.Close()
	return string(respBody), nil
}

func DecodeJson(jsonStream string, o interface{}) error {
	decoder := json.NewDecoder(strings.NewReader(jsonStream))
	// UseNumber causes the Decoder to unmarshal a number into an interface{} as a Number instead of as a float64.
	decoder.UseNumber()

	if err := decoder.Decode(o); err != nil {
		fmt.Println("error:", err)
		return err
	}
	return nil
}

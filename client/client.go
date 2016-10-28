package client

import (
	"encoding/json"
	"fmt"
	"github.com/sergeyignatov/simpleipam/common"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ApiClient struct {
	URL  string
	Http http.Client
}

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, 500*time.Millisecond)
}
func NewClient(url string) (*ApiClient, error) {
	if url == "" {
		return nil, fmt.Errorf("need an url")
	}
	transport := http.Transport{
		Dial: dialTimeout,
		ResponseHeaderTimeout: time.Second,
	}
	hc := http.Client{
		Transport: &transport,
	}
	c := &ApiClient{URL: url, Http: hc}
	return c, nil
}

func (c *ApiClient) url(elem ...string) string {
	path := strings.Join(elem, "/")
	uri := c.URL + "/" + path
	return strings.TrimSuffix(uri, "/")
}

func (c *ApiClient) GetIP(args interface{}) (*common.ApiResponse, error) {
	var v common.ApiResponse
	body, err := c.post("getip", args)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &v)
	if err != nil {
		return nil, err
	}
	return &v, nil
}
func (c *ApiClient) ReleaseIP(args interface{}) (*common.ApiResponseInt, error) {
	var v common.ApiResponseInt
	body, err := c.post("release", args)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &v)
	if err != nil {
		return nil, err
	}
	return &v, nil
}
func (c *ApiClient) post(base string, args interface{}) ([]byte, error) {
	uri := c.url("api", common.ApiVersion, base)
	val := url.Values{}
	fmt.Printf("%T, %+v\n", args, args)
	rr := make(map[string]string)

	switch args := args.(type) {
	case map[string]interface{}:
		for k, v := range args {
			rr[k] = v.(string)
		}
	case map[string]string:
		rr = args
	case common.Response:
		rr["subnet"] = args.Subnet
		rr["mac"] = args.Mac
		rr["ip"] = args.Ip
	default:
		return nil, fmt.Errorf("Unsupported format")
	}
	for k, v := range rr {
		val.Add(k, v)
	}
	resp, err := c.Http.PostForm(uri, val)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(string(body))
	}
	//var v common.ApiResponse2
	fmt.Println(string(body))
	/*err = json.Unmarshal(body, &v)
	if err != nil {
		return nil, err
	}*/
	return body, nil

}

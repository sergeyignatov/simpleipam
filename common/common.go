package common

import (
	"bytes"

	"fmt"

	mrand "math/rand"
	"time"
)

type HardwareAddr string
type IPAddr string

type Client struct {
	Hostname   string `yaml:"hostname"`
	Ip         string `yaml:"ip"`
	Mac        string `yaml:"mac"`
	CreateTime int64  `yaml:"unixtime"`
}

var ApiVersion = "1.0"

type Response struct {
	Ip       string `json:"ip"`
	Gateway  string `json:"gateway"`
	Hostname string `json:"hostname"`
	Mac      string `json:"mac"`
	Subnet   string `json:"subnet"`
}
type ApiResponseInt struct {
	Status string      `json:"status"`
	Resp   interface{} `json:"data"`
}
type ApiResponse struct {
	Status string   `json:"status"`
	Resp   Response `json:"data"`
}

func NewApiResponse(r interface{}) *ApiResponseInt {
	if t, ok := r.(error); ok {
		return &ApiResponseInt{Status: "error", Resp: t.Error()}
	}
	return &ApiResponseInt{Status: "ok", Resp: r}
}

func Generatemac() string {
	mrand.Seed(time.Now().UnixNano())
	buf := bytes.Buffer{}
	mask := "00:xx:xx:xx:xx:xx"
	for _, c := range mask {
		if c == 'x' {
			s := mrand.Intn(16)

			buf.WriteString(fmt.Sprintf("%x", s))
		} else {
			buf.WriteString(string(c))
		}
	}
	return buf.String()
}

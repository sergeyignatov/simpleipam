package client

import (
	"fmt"
	//"github.com/sergeyignatov/simpleipam/common"
	"testing"
)

/*func TestErr(t *testing.T) {
	c, err := NewClient("http://127.0.0.1:4567")
	if err != nil {
		t.Error(err)
	}
	args := make(map[string]string)
	args["mac"] = "02:10:08:0c:c8:9e"
	args["fqdn"] = "test23"
	//args["subnet"] = "10.30.32.0/20"
	_, err = c.GetIP(args)
	if err == nil {
		t.Error(err)
	}
}
*/
func TestOk(t *testing.T) {
	c, err := NewClient("http://127.0.0.1:4567")
	if err != nil {
		t.Error(err)
	}
	args := make(map[string]string)
	args["mac"] = "02:10:08:0c:c8:9e"
	args["fqdn"] = "test23"
	args["subnet"] = "10.30.32.0/20"
	d, err := c.GetIP(args)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", d)
	fmt.Printf("%+v\n", d.Resp.Ip)

	/*_, err = c.ReleaseIP(d.Resp)
	if err != nil {
		fmt.Println(err)
	}*/
}

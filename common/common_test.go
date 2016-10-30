package common

import (
	"fmt"
	"testing"
	"time"
)

/*func TestMacC(t *testing.T) {
	tm := time.Now()
	fmt.Println(deviceNextInterfaceHWAddr())
	fmt.Println(time.Now().Sub(tm).Nanoseconds())
}*/

func TestMacM(t *testing.T) {
	tm := time.Now()
	i := 0
	for i = 0; i < 100; i++ {
		fmt.Println(Generatemac())
	}
	fmt.Printf("release %s %15s %s\n", "00:16:3e:48:ff:b1", "10.30.33.32", "test-51477640295.dev-alcfd.com")
	fmt.Printf("commit  %s %15s %s\n", "00:16:3e:48:ff:b1", "255.255.255.255", "test-51477640295.dev-alcfd.com")

	fmt.Println(time.Now().Sub(tm).Nanoseconds())
}

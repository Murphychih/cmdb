package connectivity_test

import (
	"testing"

	"github.com/Murphychih/cmdb/provider/txyun/connectivity"
)

func TestTencentCloudClient(t *testing.T) {
	conn := connectivity.C()
	if err := conn.Check(); err != nil {
		t.Fatal(err)
	}
	t.Log(conn.AccountID())
}

func init() {
	//  初始化client
	err := connectivity.LoadClientFromEnv()
	if err != nil {
		panic(err)
	}
}

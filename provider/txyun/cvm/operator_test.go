package cvm_test

import (
	"context"
	"testing"

	"github.com/Murphychih/cmdb/apps/host"
	"github.com/Murphychih/cmdb/provider/txyun/connectivity"
	"github.com/Murphychih/cmdb/provider/txyun/cvm"
	tx_cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"go.uber.org/zap"
)

var (
	op *cvm.CVMOperator
)

func TestQuery(t *testing.T) {
	var (
		offset int64 = 0
		limit  int64 = 20
	)

	req := tx_cvm.NewDescribeInstancesRequest()
	req.Offset = &offset
	req.Limit = &limit
	set, err := op.Query(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(set)
}

func TestPaggerQuery(t *testing.T) {
	p := cvm.NewPagger(5, op)
	for p.Next() {
		set := host.NewHostSet()
		if err := p.Scan(context.Background(), set); err != nil {
			panic(err)
		}
		t.Log("page query result: ", set)
	}
}

func TestPagerV2Query(t *testing.T) {
	p := cvm.NewPagerV2(op)

	for p.Next() {
		set := host.NewHostSet()
		if err := p.Scan(context.Background(), set); err != nil {
			panic(err)
		}
		t.Log("page query result: ", set)
	}
}

func init() {
	//  初始化client
	err := connectivity.LoadClientFromEnv()
	if err != nil {
		panic(err)
	}

	// 初始化全局logger
	zap.NewDevelopment()
	op = cvm.NewCVMOperator(connectivity.C().CvmClient())
}

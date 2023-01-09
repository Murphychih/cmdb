package impl_test

import (
	"context"
	"fmt"
	"github.com/Murphychih/cmdb/apps/host"
	"github.com/Murphychih/cmdb/apps/host/impl"
	"github.com/Murphychih/cmdb/conf"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	// 定义对象是满足该接口的实例
	service host.Service
)

func TestCreate(t *testing.T) {
	should := assert.New(t)
	ins := host.NewHost()
	ins.Id = "ins-05"
	ins.Name = "test"
	ins.Region = "cn-hangzhou"
	ins.Type = "sm1"
	ins.CPU = 1
	ins.Memory = 2048
	ins, err := service.CreateHost(context.Background(), ins)
	if should.NoError(err) {
		fmt.Println(ins)
	}
}

func TestQuery(t *testing.T) {
	should := assert.New(t)

	req := host.NewQueryHostRequest()
	req.Keywords = "test"
	set, err := service.QueryHost(context.Background(), req)
	if should.NoError(err) {
		for i := range set.Items {
			fmt.Println(set.Items[i].Id)
		}
	}
}

func TestDescribe(t *testing.T) {
	should := assert.New(t)

	req := host.NewDescribeHostRequestWithId("ins-05")
	ins, err := service.DescribeHost(context.Background(), req)
	if should.NoError(err) {
		fmt.Println(ins.Id)
	}
}

func TestUpdate(t *testing.T) {
	should := assert.New(t)

	req := host.NewPutUpdateHostRequest("ins-05")
	req.Name = "更新测试02"
	req.Region = "rg 02"
	req.Type = "small"
	req.CPU = 1
	req.Memory = 2048
	req.Description = "测试更新"
	ins, err := service.UpdateHost(context.Background(), req)
	if should.NoError(err) {
		fmt.Println(ins.Id)
	}
}

func TestDelete(t *testing.T) {
	should := assert.New(t)

	req := host.NewDeleteHostRequest("ins-01")
	ins, err := service.DeleteHost(context.Background(), req)
	if should.NoError(err) {
		fmt.Println(*ins)
	}

}

func init() {
	// 测试用例的配置文件
	err := conf.LoadConfigFromToml("../../../etc/demo.toml")
	if err != nil {
		panic(err)
	}

	// host service 的具体实现
	service = impl.NewHostSerivice()
}

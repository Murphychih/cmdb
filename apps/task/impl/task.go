package impl

import (
	"context"
	"fmt"
	"time"

	"github.com/Murphychih/cmdb/apps/resource"
	"github.com/Murphychih/cmdb/apps/secret"
	"github.com/Murphychih/cmdb/apps/task"
	"github.com/Murphychih/cmdb/conf"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// 创建任务的业务逻辑
func (i *impl) CreateTask(ctx context.Context, req *task.CreateTaskRequst) (*task.Task, error) {
	// 创建Task实例
	t, err := task.CreateTask(req)
	if err != nil {
		return nil, err
	}

	// 1. 查询secret
	s, err := i.secret.DescribeSecret(ctx, secret.NewDescribeSecretRequest(req.SecretId))
	if err != nil {
		return nil, err
	}
	t.SecretDescription = s.Data.Description

	// 并解密api sercret
	if err := s.Data.DecryptAPISecret(conf.LoadGloabal().App.EncryptKey); err != nil {
		return nil, err
	}

	// 需要把Task 标记为Running, 修改一下Task对象的状态
	t.Run()

	var taskCancel context.CancelFunc
	switch req.Type {
	case task.Type_RESOURCE_SYNC:
		// 根据secret所属的厂商, 初始化对应厂商的operator
		switch s.Data.Vendor {
		case resource.Vendor_TENCENT:
			// 操作那种资源:
			switch req.ResourceType {
			case resource.Type_HOST:
				// 直接使用goroutine 把最耗时的逻辑
				// ctx 是不是传递Http 的ctx
				taskExecCtx, cancel := context.WithTimeout(
					context.Background(),
					time.Duration(req.Timeout)*time.Second,
				)
				taskCancel = cancel

				go i.syncHost(taskExecCtx, newSyncHostRequest(s, t))
			case resource.Type_RDS:
			case resource.Type_BILL:
			}
		case resource.Vendor_ALIYUN:

		case resource.Vendor_HUAWEI:

		case resource.Vendor_AMAZON:

		case resource.Vendor_VSPHERE:
		default:
			return nil, fmt.Errorf("unknow resource type: %s", s.Data.Vendor)
		}

		// 2. 利用secret的信息, 初始化一个operater
		// 使用operator进行资源的操作, 比如同步

		// 调用host service 把数据入库
	case task.Type_RESOURCE_RELEASE:
	default:
		return nil, fmt.Errorf("unknow task type: %s", req.Type)
	}

	// 需要保存到数据库

	if err := i.insert(ctx, t); err != nil {
		if taskCancel != nil {
			taskCancel()
		}
		return nil, err
	}

	return t, nil
}

func (i *impl) QueryTask(context.Context, *task.QueryTaskRequest) (*task.TaskSet, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryBook not implemented")
}

func (i *impl) DescribeTask(context.Context, *task.DescribeTaskRequest) (*task.Task, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DescribeBook not implemented")
}

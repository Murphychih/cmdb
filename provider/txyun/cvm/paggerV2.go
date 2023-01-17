package cvm

import (
	"context"

	"github.com/Murphychih/cmdb/apps/host"
	"github.com/Murphychih/cmdb/common/pagination"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	tx_cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"go.uber.org/zap"
)

type paggerV2 struct {
	req *cvm.DescribeInstancesRequest
	log *zap.Logger
	op  *CVMOperator
	*pagination.BasePagerV2
}

func NewPagerV2(op *CVMOperator) *paggerV2 {
	return &paggerV2{
		op:          op,
		log:         zap.L().Named("CVM"),
		BasePagerV2: pagination.NewBasePagerV2(),
		req:         tx_cvm.NewDescribeInstancesRequest(),
	}
}

// 修改Req 执行真正的下一页的offset
func (p *paggerV2) nextReq() *cvm.DescribeInstancesRequest {
	os := p.Offset()
	ps := p.PageSize()
	p.req.Offset = &os
	p.req.Limit = &ps
	return p.req
}

func (p *paggerV2) Scan(ctx context.Context, set *host.HostSet) error {
	p.log.Sugar().Debugf("query page: %d", p.PageNumber())
	hs, err := p.op.Query(ctx, p.nextReq())
	if err != nil {
		return err
	}

	// 把查询出来的数据赋值给set
	for i := range hs.Items {
		set.Add(hs.Items[i])
	}

	// 可以根据当前一页是满页来决定是否有下一页
	p.CheckHasNext(hs.Length())
	return nil
}

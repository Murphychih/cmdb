package pagination

import (
	"context"

	"github.com/infraboard/mcube/flowcontrol/tokenbucket"
)

// 抽象通用Pager

type Set interface {
	// 往Set里面添加元素, 任何类型都可以
	Add(any)
	// 当前的集合里面有多个元素
	Length() int64
}

type PagerV2 interface {
	Next() bool
	Scan(context.Context, Set) error
	Offset() int64
	SetPageSize(ps int64)
	SetRate(r float64)
	PageSize() int64
	PageNumber() int64
}

func NewBasePagerV2() *BasePagerV2 {
	return &BasePagerV2{
		hasNext:    true,
		tb:         tokenbucket.NewBucketWithRate(1, 1),
		pageNumber: 1,
		pageSize:   20,
	}
}

// 面向组合, 用他来实现一个模板, 除了Scan的其他方法都实现

type BasePagerV2 struct {
	// 令牌桶
	hasNext bool
	tb      *tokenbucket.Bucket

	// 控制分页的核心参数
	pageNumber int64
	pageSize   int64
}

func (p *BasePagerV2) Next() bool {
	// 等待分配令牌
	p.tb.Wait(1)

	return p.hasNext
}

func (p *BasePagerV2) Offset() int64 {
	return (p.pageNumber - 1) * p.pageSize
}

func (p *BasePagerV2) SetPageSize(ps int64) {
	p.pageSize = ps
}

func (p *BasePagerV2) PageSize() int64 {
	return p.pageSize
}

func (p *BasePagerV2) PageNumber() int64 {
	return p.pageNumber
}

func (p *BasePagerV2) SetRate(r float64) {
	p.tb.SetRate(r)
}

func (p *BasePagerV2) CheckHasNext(current int64) {
	// 可以根据当前一页是满页来决定是否有下一页
	if current < p.pageSize {
		p.hasNext = false
	} else {
		// 直接调整指针到下一页
		p.pageNumber++
	}
}

package impl

import (
	"context"

	"github.com/Murphychih/cmdb/apps/host"
	"github.com/infraboard/mcube/sqlbuilder"

)

// 网络服务业务处理层（controller层）


// 录入主机
func (i *HostService) CreateHost(ctx context.Context, ins *host.Host) (*host.Host, error){
	// 首先校验数据合法性
	if err := ins.Validate(); err != nil {
		return nil, err
	}
	
	// 默认值填充
	ins.InjectDefault()

	// 由dao.go模块负责将对象入库
	if err := i.save(ctx, ins); err != nil {
		return nil, err
	}

	return ins, nil
}

// 查询主机列表
func (i *HostService) QueryHost(ctx context.Context, req *host.QueryHostRequest) (*host.HostSet, error){
	 b := sqlbuilder.NewBuilder(QueryHostSQL)
	if req.Keywords != ""{
		b.Where("r.`name`LIKE ? OR r.description LIKE ? OR r.private_ip LIKE ? OR r.public_ip LIKE ?",
			"%"+req.Keywords+"%",
			"%"+req.Keywords+"%",
			req.Keywords+"%",
			req.Keywords+"%",
		)
	}

	b.Limit(req.OffSet(), req.GetPageSize())
	querySQL, args := b.Build()
	i.l.Sugar().Debugf("query sql: %s, args: %v", querySQL, args)

	stmt, err :=i.db.PrepareContext(ctx, querySQL)
	if err !=nil{
		return nil, err
	}
	defer stmt.Close()

	set := host.NewHostSet()
	rows, err := stmt.QueryContext(ctx, args...)
	for rows.Next(){
		// 每扫描一行,就需要读取出来
		// h.cpu, h.memory, h.gpu_spec, h.gpu_amount, h.os_type, h.os_name, h.serial_number
		ins := host.NewHost()
		if err := rows.Scan(
			&ins.Id, &ins.Vendor, &ins.Region, &ins.CreateAt, &ins.ExpireAt,
			&ins.Type, &ins.Name, &ins.Description, &ins.Status, &ins.UpdateAt, &ins.SyncAt,
			&ins.Account, &ins.PublicIP, &ins.PrivateIP,
			&ins.CPU, &ins.Memory, &ins.GPUSpec, &ins.GPUAmount, &ins.OSType, &ins.OSName, &ins.SerialNumber,
		); err != nil {
			return nil, err
		}
		set.Add(ins)
	}

	// total统计
	countSQL, args := b.BuildCount()
	i.l.Sugar().Debugf("count sql: %s, args: %v", countSQL, args)
	countStmt, err := i.db.PrepareContext(ctx, countSQL)
	if err != nil {
		return nil, err
	}
	defer countStmt.Close()
	// 返回一行
	if err := countStmt.QueryRowContext(ctx, args...).Scan(&set.Total); err != nil {
		return nil, err
	}

	return set, nil

}
// 查询主机详情
func (i *HostService) DescribeHost(context.Context, *host.DescribeHostRequest) (*host.Host, error)
// 主机更新
func (i *HostService) UpdateHost(context.Context, *host.UpdateHostRequest) (*host.Host, error)
// 主机删除 如前端需要 打印当前删除主机的Ip或者其他信息
func (i *HostService) DeleteHost(context.Context, *host.DeleteHostRequest) (*host.Host, error)

package impl

import (
	"context"
	"fmt"

	"github.com/Murphychih/cmdb/apps/host"
	"go.uber.org/zap"
)

// 该文件负责完成对象和数据库之间的转换
// 把Host对象保存到数据内 维持数据的一致性

func (i *HostService) save(ctx context.Context, ins *host.Host) error {
	var err error

	// 数据需同时入库两张表
	// 此处使用事务维护数据一致性
	tx, err := i.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("start transaction <%s> failed", err)
	}

	// 通过Defer处理事务提交方式
	// 1. 无报错，则Commit 事务
	// 2. 有报错, 则Rollback 事务
	defer func() {
		if err != nil {
			if err = tx.Rollback(); err != nil {
				i.l.Error("rollback transaction failed",
					zap.String("error", err.Error()))
			}
		} else {
			if err = tx.Commit(); err != nil {
				i.l.Error("commit transaction failed",
					zap.String("error", err.Error()))
			}
		}
	}()

	// 插入Resource表数据
	rstmt, err := tx.PrepareContext(ctx, InsertResourceSQL)
	if err != nil {
		return err
	}
	defer rstmt.Close()

	_, err = rstmt.ExecContext(ctx,
		ins.Id, ins.Vendor, ins.Region, ins.CreateAt, ins.ExpireAt, ins.Type,
		ins.Name, ins.Description, ins.Status, ins.UpdateAt, ins.SyncAt, ins.Account, ins.PublicIP,
		ins.PrivateIP,
	)

	if err != nil {
		return err
	}

	// 插入Describe数据
	dstmt, err := tx.PrepareContext(ctx, InsertDescribeSQL)
	if err != nil {
		return err
	}
	defer dstmt.Close()

	_, err = dstmt.ExecContext(ctx,
		ins.Id, ins.CPU, ins.Memory, ins.GPUAmount, ins.GPUSpec,
		ins.OSType, ins.OSName, ins.SerialNumber,
	)

	if err != nil {
		return err
	}

	return nil
}

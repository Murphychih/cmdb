package impl

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Murphychih/cmdb/apps/secret"
	"github.com/Murphychih/cmdb/common/response/exception"
	"github.com/Murphychih/cmdb/conf"
	"github.com/infraboard/mcube/sqlbuilder"
)

func (s *impl) CreateSecret(ctx context.Context, req *secret.CreateSecretRequest) (
	*secret.Secret, error) {
	ins, err := secret.NewSecret(req)
	if err != nil {
		return nil, fmt.Errorf("vaildate secret creation failed: %v", err)
	}

	stmt, err := s.db.PrepareContext(ctx, insertSecretSQL)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// 入库之前需要加密
	// 判断文本是否已经加密 并加密
	if err := ins.Data.EncryptAPISecret(conf.LoadGloabal().App.EncryptKey); err != nil {
		s.log.Sugar().Warnf("encrypt api key error, %s", err)
	}

	_, err = stmt.ExecContext(ctx,
		ins.Id, ins.CreateAt, ins.Data.Description, ins.Data.Vendor, ins.Data.Address,
		ins.Data.AllowRegionString(), ins.Data.CrendentialType, ins.Data.ApiKey, ins.Data.ApiSecret,
		ins.Data.RequestRate,
	)
	if err != nil {
		return nil, err
	}

	return ins, nil
}

func (s *impl) QuerySecret(ctx context.Context, req *secret.QuerySecretRequest) (
	*secret.SecretSet, error) {
	query := sqlbuilder.NewQuery(querySecretSQL)

	if req.Keywords != "" {
		query.Where("description LIKE ? OR api_key = ?",
			"%"+req.Keywords+"%",
			req.Keywords,
		)
	}

	querySQL, args := query.Order("create_at").Desc().Limit(req.Page.ComputeOffset(), uint(req.Page.PageSize)).BuildQuery()
	s.log.Sugar().Debugf("sql: %s, args: %v", querySQL, args)

	queryStmt, err := s.db.PrepareContext(ctx, querySQL)
	if err != nil {
		return nil, exception.NewInternalServerError("prepare query secret error, %s", err.Error())
	}
	defer queryStmt.Close()

	rows, err := queryStmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, exception.NewInternalServerError(err.Error())
	}
	defer rows.Close()

	set := secret.NewSecretSet()
	allowRegions := ""
	for rows.Next() {
		ins := secret.NewDefaultSecret()
		err := rows.Scan(
			&ins.Id, &ins.CreateAt, &ins.Data.Description, &ins.Data.Vendor, &ins.Data.Address,
			&allowRegions, &ins.Data.CrendentialType, &ins.Data.ApiKey, &ins.Data.ApiSecret,
			&ins.Data.RequestRate,
		)
		if err != nil {
			return nil, exception.NewInternalServerError("query secret error, %s", err.Error())
		}
		ins.Data.LoadAllowRegionFromString(allowRegions)
		ins.Data.Desense()
		set.Add(ins)
	}

	// 获取total SELECT COUNT(*) FROMT t Where ....
	countSQL, args := query.BuildCount()
	countStmt, err := s.db.PrepareContext(ctx, countSQL)
	if err != nil {
		return nil, exception.NewInternalServerError(err.Error())
	}

	defer countStmt.Close()
	err = countStmt.QueryRowContext(ctx, args...).Scan(&set.Total)
	if err != nil {
		return nil, exception.NewInternalServerError(err.Error())
	}

	return set, nil
}
func (s *impl) DescribeSecret(ctx context.Context, req *secret.DescribeSecretRequest) (
	*secret.Secret, error) {
	query := sqlbuilder.NewQuery(querySecretSQL)
	querySQL, args := query.Where("id = ?", req.Id).BuildQuery()
	s.log.Sugar().Debugf("sql: %s", querySQL)

	queryStmt, err := s.db.PrepareContext(ctx, querySQL)
	if err != nil {
		return nil, exception.NewInternalServerError("prepare query secret error, %s", err.Error())
	}
	defer queryStmt.Close()

	ins := secret.NewDefaultSecret()
	allowRegions := ""
	err = queryStmt.QueryRowContext(ctx, args...).Scan(
		&ins.Id, &ins.CreateAt, &ins.Data.Description, &ins.Data.Vendor, &ins.Data.Address,
		&allowRegions, &ins.Data.CrendentialType, &ins.Data.ApiKey, &ins.Data.ApiSecret,
		&ins.Data.RequestRate,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.NewNotFound("%#v not found", req)
		}
		return nil, exception.NewInternalServerError("describe secret error, %s", err.Error())
	}

	ins.Data.LoadAllowRegionFromString(allowRegions)
	return ins, nil
}
func (s *impl) DeleteSecret(ctx context.Context, req *secret.DeleteSecretRequest) (
	*secret.Secret, error) {
	ins, err := s.DescribeSecret(ctx, secret.NewDescribeSecretRequest(req.Id))
	if err != nil {
		return nil, err
	}

	if err := s.deleteSecret(ctx, ins); err != nil {
		return nil, err
	}

	return ins, nil
}

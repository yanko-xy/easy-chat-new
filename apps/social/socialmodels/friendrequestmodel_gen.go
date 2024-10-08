// Code generated by goctl. DO NOT EDIT.

package socialmodels

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	friendRequestFieldNames          = builder.RawFieldNames(&FriendRequest{})
	friendRequestRows                = strings.Join(friendRequestFieldNames, ",")
	friendRequestRowsExpectAutoSet   = strings.Join(stringx.Remove(friendRequestFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	friendRequestRowsWithPlaceHolder = strings.Join(stringx.Remove(friendRequestFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"

	cacheFriendRequestIdPrefix = "cache:friendRequest:id:"
)

type (
	friendRequestModel interface {
		Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error ) error
		Insert(ctx context.Context, data *FriendRequest) (sql.Result, error)
		FindOne(ctx context.Context, id uint64) (*FriendRequest, error)
		FindByReqUidAndUserId(ctx context.Context, reqUid, userId string) (*FriendRequest, error)
		ListByUserId(ctx context.Context, userId string) ([]*FriendRequest, error)
		ListByNoHandle(ctx context.Context, userId string) ([]*FriendRequest, error)
		Update(ctx context.Context, data *FriendRequest) error
		UpdateWithTrans(ctx context.Context, session sqlx.Session, data *FriendRequest) error
		Delete(ctx context.Context, id uint64) error
	}

	defaultFriendRequestModel struct {
		sqlc.CachedConn
		table string
	}

	FriendRequest struct {
		Id           uint64         `db:"id"`
		UserId       string         `db:"user_id"`
		ReqUid       string         `db:"req_uid"`
		ReqMsg       sql.NullString `db:"req_msg"`
		ReqTime      time.Time      `db:"req_time"`
		HandleResult sql.NullInt64  `db:"handle_result"`
		HandleMsg    sql.NullString `db:"handle_msg"`
		HandledAt    sql.NullTime   `db:"handled_at"`
		CreatedAt    time.Time      `db:"created_at"`
		UpdatedAt    time.Time      `db:"updated_at"`
		DeletedAt    sql.NullTime   `db:"deleted_at"`
	}
)

func newFriendRequestModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) *defaultFriendRequestModel {
	return &defaultFriendRequestModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "`friend_request`",
	}
}

func (m *defaultFriendRequestModel) Trans(ctx context.Context, fn func(ctx context.Context,
	session sqlx.Session) error) error {
	return m.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}

func (m *defaultFriendRequestModel) Delete(ctx context.Context, id uint64) error {
	friendRequestIdKey := fmt.Sprintf("%s%v", cacheFriendRequestIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, friendRequestIdKey)
	return err
}

func (m *defaultFriendRequestModel) FindOne(ctx context.Context, id uint64) (*FriendRequest, error) {
	friendRequestIdKey := fmt.Sprintf("%s%v", cacheFriendRequestIdPrefix, id)
	var resp FriendRequest
	err := m.QueryRowCtx(ctx, &resp, friendRequestIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", friendRequestRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultFriendRequestModel) FindByReqUidAndUserId(ctx context.Context, reqUid,
	userId string) (*FriendRequest, error) {
	query := fmt.Sprintf("select %s from %s where `req_uid` = ? and `user_id` = ? ",
		friendRequestRows, m.table)

	var resp FriendRequest
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, reqUid, userId)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultFriendRequestModel) ListByUserId(ctx context.Context, userId string) ([]*FriendRequest, error)  {
	var resp []*FriendRequest
	query := fmt.Sprintf("select %s from %s where `user_id` = ? ", friendRequestRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, userId)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultFriendRequestModel) ListByNoHandle(ctx context.Context, userId string) ([]*FriendRequest, error)  {
	var resp []*FriendRequest
	query := fmt.Sprintf("select %s from %s where `user_id` = ? and `handle_result` = 1", friendRequestRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, userId)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultFriendRequestModel) Insert(ctx context.Context, data *FriendRequest) (sql.Result, error) {
	friendRequestIdKey := fmt.Sprintf("%s%v", cacheFriendRequestIdPrefix, data.Id)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?)", m.table, friendRequestRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.UserId, data.ReqUid, data.ReqMsg, data.ReqTime, data.HandleResult, data.HandleMsg, data.HandledAt, data.DeletedAt)
	}, friendRequestIdKey)
	return ret, err
}

func (m *defaultFriendRequestModel) Update(ctx context.Context, data *FriendRequest) error {
	friendRequestIdKey := fmt.Sprintf("%s%v", cacheFriendRequestIdPrefix, data.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, friendRequestRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, data.UserId, data.ReqUid, data.ReqMsg, data.ReqTime, data.HandleResult, data.HandleMsg, data.HandledAt, data.DeletedAt, data.Id)
	}, friendRequestIdKey)
	return err
}

func (m *defaultFriendRequestModel) UpdateWithTrans(ctx context.Context, session sqlx.Session,
	data *FriendRequest) error  {
	friendRequestIdKey := fmt.Sprintf("%s%v", cacheFriendRequestIdPrefix, data.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, friendRequestRowsWithPlaceHolder)
		return session.ExecCtx(ctx, query, data.UserId, data.ReqUid, data.ReqMsg, data.ReqTime,
			data.HandleResult, data.HandleMsg, data.HandledAt, data.DeletedAt, data.Id)
	}, friendRequestIdKey)
	return err
}

func (m *defaultFriendRequestModel) formatPrimary(primary any) string {
	return fmt.Sprintf("%s%v", cacheFriendRequestIdPrefix, primary)
}

func (m *defaultFriendRequestModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary any) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", friendRequestRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultFriendRequestModel) tableName() string {
	return m.table
}

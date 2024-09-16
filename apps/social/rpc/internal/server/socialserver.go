// Code generated by goctl. DO NOT EDIT.
// Source: social.proto

package server

import (
	"context"

	"github.com/yanko-xy/easy-chat/apps/social/rpc/internal/logic"
	"github.com/yanko-xy/easy-chat/apps/social/rpc/internal/svc"
	"github.com/yanko-xy/easy-chat/apps/social/rpc/social"
)

type SocialServer struct {
	svcCtx *svc.ServiceContext
	social.UnimplementedSocialServer
}

func NewSocialServer(svcCtx *svc.ServiceContext) *SocialServer {
	return &SocialServer{
		svcCtx: svcCtx,
	}
}

// 好友业务：请求好友、通过或拒绝申请、好友列表
func (s *SocialServer) FriendPutIn(ctx context.Context, in *social.FriendPutInReq) (*social.FriendPutInResp, error) {
	l := logic.NewFriendPutInLogic(ctx, s.svcCtx)
	return l.FriendPutIn(in)
}

func (s *SocialServer) FriendPutInHandle(ctx context.Context, in *social.FriendPutInHandleReq) (*social.FriendPutInHandleResp, error) {
	l := logic.NewFriendPutInHandleLogic(ctx, s.svcCtx)
	return l.FriendPutInHandle(in)
}

func (s *SocialServer) FriendPutInList(ctx context.Context, in *social.FriendPutInListReq) (*social.FriendPutInListResp, error) {
	l := logic.NewFriendPutInListLogic(ctx, s.svcCtx)
	return l.FriendPutInList(in)
}

func (s *SocialServer) FriendList(ctx context.Context, in *social.FriendListReq) (*social.FriendListResp, error) {
	l := logic.NewFriendListLogic(ctx, s.svcCtx)
	return l.FriendList(in)
}

// 群业务：创建群，修改群，群公告，申请群，用户群列表，群成员，申请群，群退出..
func (s *SocialServer) GroupCreate(ctx context.Context, in *social.GroupCreateReq) (*social.GroupCreateResp, error) {
	l := logic.NewGroupCreateLogic(ctx, s.svcCtx)
	return l.GroupCreate(in)
}

func (s *SocialServer) GroupPutIn(ctx context.Context, in *social.GroupPutInReq) (*social.GroupPutInResp, error) {
	l := logic.NewGroupPutInLogic(ctx, s.svcCtx)
	return l.GroupPutIn(in)
}

func (s *SocialServer) GroupPutInList(ctx context.Context, in *social.GroupPutInListReq) (*social.GroupPutInListResp, error) {
	l := logic.NewGroupPutInListLogic(ctx, s.svcCtx)
	return l.GroupPutInList(in)
}

func (s *SocialServer) GroupPutInHandle(ctx context.Context, in *social.GroupPutInHandleReq) (*social.GroupPutInHandleResp, error) {
	l := logic.NewGroupPutInHandleLogic(ctx, s.svcCtx)
	return l.GroupPutInHandle(in)
}

func (s *SocialServer) GroupList(ctx context.Context, in *social.GroupListReq) (*social.GroupListResp, error) {
	l := logic.NewGroupListLogic(ctx, s.svcCtx)
	return l.GroupList(in)
}

func (s *SocialServer) GroupUsers(ctx context.Context, in *social.GroupUsersReq) (*social.GroupUsersResp, error) {
	l := logic.NewGroupUsersLogic(ctx, s.svcCtx)
	return l.GroupUsers(in)
}

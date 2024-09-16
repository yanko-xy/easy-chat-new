package logic

import (
	"context"
	"github.com/yanko-xy/easy-chat/apps/user/rpc/user"
	"reflect"
	"testing"
)

func TestGetUserInfoLogic_GetUserInfo(t *testing.T) {

	type args struct {
		in *user.GetUserInfoReq
	}
	tests := []struct {
		name    string
		args    args
		want    *user.GetUserInfoResp
		wantErr bool
	}{
		{
			"get userInfo successfully",
			args{
				in: &user.GetUserInfoReq{
					Id: "0x0000001000000001",
				},
			},
			&user.GetUserInfoResp{
				User: &user.UserEntity{
					Id:       "0x0000001000000001",
					Avatar:   "test.png",
					Nickname: "yanko",
					Phone:    "18758004746",
					Sex:      0,
				},
			},
			false,
		},
		{
			"user not found",
			args{
				in: &user.GetUserInfoReq{
					Id: "0x0000000000000000",
				},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewGetUserInfoLogic(context.Background(), svcCtx)
			got, err := l.GetUserInfo(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserInfoLogic.GetUserInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserInfoLogic.GetUserInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

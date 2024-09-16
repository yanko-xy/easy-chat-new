package logic

import (
	"context"
	"github.com/yanko-xy/easy-chat/apps/user/rpc/user"
	"testing"
)

func TestLoginLogic_Login(t *testing.T) {
	type args struct {
		in *user.LoginReq
	}
	tests := []struct {
		name      string
		args      args
		wantPrint bool
		wantErr   bool
	}{
		{
			"login successfully",
			args{
				in: &user.LoginReq{
					Phone:    "18758004746",
					Password: "123456",
				},
			},
			true,
			false,
		},
		{
			"login err with phone not registered",
			args{
				in: &user.LoginReq{
					Phone:    "18758001111",
					Password: "123456",
				},
			},
			false,
			true,
		},
		{
			"login err with password error",
			args{
				in: &user.LoginReq{
					Phone:    "18758004746",
					Password: "1234567",
				},
			},
			false,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLoginLogic(context.Background(), svcCtx)
			got, err := l.Login(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoginLogic.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantPrint {
				t.Log(tt.name, got)
			}
		})
	}
}

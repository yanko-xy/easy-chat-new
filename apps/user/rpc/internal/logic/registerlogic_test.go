package logic

import (
	"context"
	"github.com/yanko-xy/easy-chat/apps/user/rpc/user"
	"testing"
)

func TestRegisterLogic_Register(t *testing.T) {
	type args struct {
		in *user.RegisterReq
	}
	tests := []struct {
		name      string
		args      args
		wantPrint bool
		wantErr   bool
	}{
		{
			"register successfully",
			args{
				in: &user.RegisterReq{
					Phone:    "18758000002",
					Nickname: "testname3",
					Password: "123456",
					Avatar:   "test.png",
					Sex:      0,
				},
			},
			true,
			false,
		},
		{
			"phone is register",
			args{
				in: &user.RegisterReq{
					Phone:    "18758004743",
					Nickname: "yanko",
					Password: "123456",
					Avatar:   "test.png",
					Sex:      0,
				},
			},
			false,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewRegisterLogic(context.Background(), svcCtx)
			got, err := l.Register(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterLogic.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantPrint {
				t.Log(tt.name, got)
			}
		})
	}
}

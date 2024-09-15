package logic

import (
	"context"
	"reflect"
	"rpc/user"
	"testing"
)

func TestFindUserLogic_FindUser(t *testing.T) {

	type args struct {
		in *user.FindUserReq
	}
	tests := []struct {
		name    string
		args    args
		want    *user.FindUserResp
		wantErr bool
	}{
		{
			"find user by phone successfully",
			args{
				in: &user.FindUserReq{
					Phone: "18758000000",
				},
			},
			&user.FindUserResp{
				Users: []*user.UserEntity{
					{
						Id:       "0x000000e000000001",
						Avatar:   "test.png",
						Nickname: "testname1",
						Phone:    "18758000000",
						Sex:      0,
					},
				},
			},
			false,
		},
		{
			"find user by name successfully",
			args{
				in: &user.FindUserReq{
					Name: "testname1",
				},
			},
			&user.FindUserResp{
				Users: []*user.UserEntity{
					{
						Id:       "0x000000e000000001",
						Avatar:   "test.png",
						Nickname: "testname1",
						Phone:    "18758000000",
						Sex:      0,
					},
				},
			},
			false,
		},
		{
			"find user by ids successfully",
			args{
				in: &user.FindUserReq{
					Ids: []string{"0x000000e000000001", "0x000000f000000001"},
				},
			},
			&user.FindUserResp{
				Users: []*user.UserEntity{
					{
						Id:       "0x000000e000000001",
						Avatar:   "test.png",
						Nickname: "testname1",
						Phone:    "18758000000",
						Sex:      0,
					},
					{
						Id:       "0x000000f000000001",
						Avatar:   "test.png",
						Nickname: "testname2",
						Phone:    "18758000001",
						Sex:      0,
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewFindUserLogic(context.Background(), svcCtx)
			got, err := l.FindUser(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindUserLogic.FindUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindUserLogic.FindUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

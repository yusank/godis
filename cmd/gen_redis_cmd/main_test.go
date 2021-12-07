package main

import (
	"testing"
)

func Test_regexCmdFuncStr(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name      string
		args      args
		wantFName string
		wantErr   bool
	}{
		{
			name:      "1",
			args:      args{str: "func sAdd(c *Command) (*protocol.Response, error) {"},
			wantErr:   false,
			wantFName: "sAdd",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFName := regexCmdFuncStr(tt.args.str)
			if gotFName != tt.wantFName {
				t.Errorf("regexCmdFuncStr() gotFName = %v, want %v", gotFName, tt.wantFName)
			}
		})
	}
}

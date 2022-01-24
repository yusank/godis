package protocol

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeFromReader(t *testing.T) {
	tests := []struct {
		name    string
		args    Reader
		wantRec *Receive
		wantErr bool
	}{
		{
			name:    "simple string",
			args:    bytes.NewBuffer([]byte("+ok\r\n")),
			wantErr: false,
			wantRec: &Receive{
				Elements: []string{"ok"},
				OrgStr:   "+ok\r\n",
			},
		},
		{
			name:    "error",
			args:    bytes.NewBuffer([]byte("-ERR unknown\r\n")),
			wantErr: false,
			wantRec: &Receive{
				Elements: []string{"ERR unknown"},
				OrgStr:   "-ERR unknown\r\n",
			},
		},
		{
			name:    "bulk string",
			args:    bytes.NewBuffer([]byte("$5\r\nhello\r\n")),
			wantErr: false,
			wantRec: &Receive{
				Elements: []string{"hello"},
				OrgStr:   "$5\r\nhello\r\n",
			},
		},
		{
			name:    "bulk string array",
			args:    bytes.NewBuffer([]byte("*3\r\n$5\r\nhello\r\n$0\r\n\r\n$5\r\nworld\r\n")),
			wantErr: false,
			wantRec: &Receive{Elements: []string{"hello", "", "world"},
				OrgStr: "*3\r\n$5\r\nhello\r\n$0\r\n\r\n$5\r\nworld\r\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRec, err := DecodeFromReader(tt.args)
			if !assert.NoError(t, err) {
				return
			}
			if !assert.Equal(t, tt.wantRec, gotRec) {
				return
			}
		})
	}
}

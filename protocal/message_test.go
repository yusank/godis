package protocal

import (
	"reflect"
	"testing"
)

func Test_readBulkString(t *testing.T) {
	type args struct {
		startAt int
		data    []byte
	}
	tests := []struct {
		name           string
		args           args
		wantEle        *Element
		wantNewStartAt int
		wantErr        bool
	}{
		// TODO: Add test cases.
		{
			name: "test_right_bulk",
			args: args{
				startAt: 1,
				data:    []byte("$12\r\nabcdabcdabcd\r\n"),
			},
			wantEle: &Element{
				Type:            ElementTypeString,
				DescriptionType: BulkStringsPrefix,
				Value:           "abcdabcdabcd",
			},
			wantNewStartAt: 19,
			wantErr:        false,
		},
		{
			name: "test_nil_bulk",
			args: args{
				startAt: 1,
				data:    []byte("$-1\r\n"),
			},
			wantEle: &Element{
				Type:            ElementTypeNil,
				DescriptionType: BulkStringsPrefix,
			},
			wantNewStartAt: 5,
			wantErr:        false,
		},
		{
			name: "test_empty_bulk",
			args: args{
				startAt: 1,
				data:    []byte("$0\r\n\r\n"),
			},
			wantEle: &Element{
				Type:            ElementTypeString,
				DescriptionType: BulkStringsPrefix,
				Value:           "",
			},
			wantNewStartAt: 6,
			wantErr:        false,
		},
		{
			name: "test_invalid_bulk_1",
			args: args{
				startAt: 1,
				data:    []byte("$-2\r\n"),
			},
			wantErr: true,
		},
		{
			name: "test_invalid_bulk_2",
			args: args{
				startAt: 1,
				data:    []byte("$12\r\na\r\n"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEle, gotNewStartAt, err := readBulkString(tt.args.startAt, tt.args.data)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("readBulkString() error = %v, wantErr %v", err, tt.wantErr)
				}

				return
			}

			if !reflect.DeepEqual(gotEle, tt.wantEle) {
				t.Errorf("readBulkString() gotEle = %v, want %v", gotEle, tt.wantEle)
			}
			if gotNewStartAt != tt.wantNewStartAt {
				t.Errorf("readBulkString() gotNewStartAt = %v, want %v", gotNewStartAt, tt.wantNewStartAt)
			}
		})
	}
}

func TestMessage_Decode(t *testing.T) {
	type fields struct {
		OriginalData []byte
		Elements     []*Element
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "normal_test",
			fields: fields{
				OriginalData: []byte("*2\r\n*3\r\n:1\r\n:2\r\n:3\r\n*2\r\n+Foo\r\n-Bar\r\n"),
				Elements:     make([]*Element, 0),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Message{
				OriginalData: tt.fields.OriginalData,
				Elements:     tt.fields.Elements,
			}
			if err := m.Decode(); (err != nil) != tt.wantErr {
				t.Errorf("Message.Decode() error = %v, wantErr %v", err, tt.wantErr)
			}

			for _, v := range m.Elements {
				t.Log(v.String())
			}
		})
	}
}

func TestMessage_Encode(t *testing.T) {
	type fields struct {
		OriginalData []byte
		Elements     []*Element
	}
	tests := []struct {
		name       string
		fields     fields
		wantErr    bool
		wantResult []byte
	}{
		// TODO: Add test cases.
		{
			name: "normal_test",
			fields: fields{
				Elements: []*Element{
					{
						Type:            ElementTypeArray,
						Len:             4,
						Value:           "4",
						DescriptionType: ArrarysPrefix,
					},
					{
						Type:            ElementTypeString,
						Value:           "hello",
						DescriptionType: SimpleStringsPrefix,
					},
					{
						Type:            ElementTypeInt,
						Value:           "12",
						DescriptionType: IntegersPrefix,
					},
					{
						Type:            ElementTypeString,
						Len:             4,
						Value:           "abcd",
						DescriptionType: BulkStringsPrefix,
					},
					{
						Type:            ElementTypeNil,
						Value:           "-1",
						DescriptionType: BulkStringsPrefix,
					},
				},
			},
			/*
				*4\r\n
				+hello\r\n
				:12\r\n
				$4\r\n
				abcd\r\n
				$-1\r\n
			*/
			wantErr:    false,
			wantResult: []byte("*4\r\n+hello\r\n:12\r\n$4\r\nabcd\r\n$-1\r\n"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Message{
				OriginalData: tt.fields.OriginalData,
				Elements:     tt.fields.Elements,
			}
			if err := m.Encode(); (err != nil) != tt.wantErr {
				t.Errorf("Message.Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(m.OriginalData, tt.wantResult) {
				t.Errorf("Want:%v, got:%v", tt.wantErr, m.OriginalData)
			}
		})
	}
}

func Test_validArray(t *testing.T) {
	type args struct {
		elements []*Element
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "normal_test",
			args: args{
				elements: []*Element{
					{
						Type:            ElementTypeArray,
						Len:             3,
						Value:           "3",
						DescriptionType: ArrarysPrefix,
					},
					{
						Type:            ElementTypeArray,
						Len:             2,
						Value:           "2",
						DescriptionType: ArrarysPrefix,
					},
					{
						Type:            ElementTypeString,
						Value:           "hello",
						DescriptionType: SimpleStringsPrefix,
					},
					{
						Type:            ElementTypeInt,
						Value:           "12",
						DescriptionType: IntegersPrefix,
					},
					{
						Type:            ElementTypeArray,
						Len:             2,
						Value:           "2",
						DescriptionType: ArrarysPrefix,
					},
					{
						Type:            ElementTypeString,
						Len:             4,
						Value:           "abcd",
						DescriptionType: BulkStringsPrefix,
					},

					{
						Type:            ElementTypeString,
						Len:             4,
						Value:           "cdab",
						DescriptionType: BulkStringsPrefix,
					},
					{
						Type:            ElementTypeString,
						Len:             4,
						Value:           "eeee",
						DescriptionType: BulkStringsPrefix,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validArray(tt.args.elements); (err != nil) != tt.wantErr {
				t.Errorf("validArray() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

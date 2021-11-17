package protocol

import (
	"reflect"
	"testing"
)

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
						Type:        ElementTypeArray,
						Len:         4,
						Value:       "4",
						Description: DescriptionArray,
					},
					{
						Type:        ElementTypeString,
						Value:       "hello",
						Description: DescriptionSimpleStrings,
					},
					{
						Type:        ElementTypeInt,
						Value:       "12",
						Description: DescriptionIntegers,
					},
					{
						Type:        ElementTypeString,
						Len:         4,
						Value:       "abcd",
						Description: DescriptionBulkStrings,
					},
					{
						Type:        ElementTypeNil,
						Value:       "-1",
						Description: DescriptionBulkStrings,
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
				originalData: tt.fields.OriginalData,
				Elements:     tt.fields.Elements,
			}
			if err := m.Encode(); (err != nil) != tt.wantErr {
				t.Errorf("Message.Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(m.originalData, tt.wantResult) {
				t.Errorf("Want:%v, got:%v", tt.wantErr, m.originalData)
			}
		})
	}
}
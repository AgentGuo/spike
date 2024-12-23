/*
@author: panfengguo
@since: 2024/11/17
@desc: desc
*/
package storage

import (
	"testing"
)

func TestMysqlClient_GetFuncMetaDataByFunctionName(t *testing.T) {
	type args struct {
		functionName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test1", args{"test"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMysqlClient()
			got, err := m.GetFuncMetaDataByFunctionName(tt.args.functionName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFuncMetaDataByFunctionName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(got)
		})
	}
}

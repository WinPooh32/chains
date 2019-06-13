package merkle

import (
	"reflect"
	"testing"
)

func Test_makeAuditNode(t *testing.T) {
	type args struct {
		hash   string
		branch Direction
	}
	tests := []struct {
		name string
		args args
		want *AuditNode
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makeAuditNode(tt.args.hash, tt.args.branch); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("makeAuditNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

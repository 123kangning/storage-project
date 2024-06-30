package dal

import "testing"

func TestAdd(t *testing.T) {
	type args struct {
		name string
		hash string
		size int64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestAdd",
			args: args{
				name: "test.txt",
				hash: "hash",
				size: 100,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Add(tt.args.name, tt.args.hash, tt.args.size)
			if err != nil {
				t.Errorf("Add() error = %v", err)
			}
		})
	}
}
func TestDel(t *testing.T) {
	type args struct {
		hash string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestDel",
			args: args{
				hash: "hash",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Del(tt.args.hash)
		})
	}
}

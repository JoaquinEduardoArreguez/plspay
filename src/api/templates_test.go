package main

import (
	"testing"
	"time"
)

func Test_humanDate(t *testing.T) {
	t.Parallel()

	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "UTC",
			args: args{t: time.Date(2024, 1, 9, 0, 0, 0, 0, time.UTC)},
			want: "09 Jan 2024",
		},
		{
			name: "Empty",
			args: args{t: time.Time{}},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := humanDate(tt.args.t); got != tt.want {
				t.Errorf("humanDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

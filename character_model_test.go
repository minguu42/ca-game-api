package ca_game_api

import (
	"testing"
)

func TestCalculateLevel(t *testing.T) {
	type args struct {
		experience int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "experience = 100",
			args: args{experience: 100},
			want: 1,
		},
		{
			name: "experience = 400",
			args: args{experience: 400},
			want: 2,
		},
		{
			name: "experience = 899",
			args: args{experience: 899},
			want: 2,
		},
		{
			name: "experience = 999999",
			args: args{experience: 999999},
			want: 99,
		},
		{
			name: "experience = 1000000",
			args: args{experience: 1000000},
			want: 100,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := calculateLevel(test.args.experience); got != test.want {
				t.Errorf("level should be %v, but %v", test.want, got)
			}
		})
	}
}

func TestCalculatePower(t *testing.T) {
	type args struct {
		experience int
		basePower int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "experience = 100, basePower = 200",
			args: args{
				experience: 100,
				basePower: 200,
			},
			want: 200,
		},
		{
			name: "experience = 400, basePower = 300",
			args: args{
				experience: 400,
				basePower: 300,
			},
			want: 600,
		},
		{
			name: "experience = 899, basePower = 400",
			args: args{
				experience: 899,
				basePower: 400,
			},
			want: 800,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := calculatePower(test.args.experience, test.args.basePower); got != test.want {
				t.Errorf("power should be %v, but %v", test.want, got)
			}
		})
	}
}

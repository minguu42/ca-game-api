package ca_game_api

import (
	"testing"
)

func TestDecideRarity(t *testing.T) {
	tests := []string{"test1", "test2", "test3", "test4", "test5"}
	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			got := decideRarity(10, 30, 60)
			if got != 3 && got != 4 && got != 5 {
				t.Errorf("rarity should be between 3 and 5, but %v", got)
			}
		})
	}
}

func TestCalculateCharacterExperience(t *testing.T) {
	type args struct {
		level int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "level = 1",
			args: args{level: 1},
			want: 100,
		},
		{
			name: "level = 2",
			args: args{level: 2},
			want: 400,
		},
		{
			name: "level = 9",
			args: args{level: 9},
			want: 8100,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := calculateCharacterExperience(test.args.level); got != test.want {
				t.Errorf("experience should be %v, but %v", test.want, got)
			}
		})
	}
}

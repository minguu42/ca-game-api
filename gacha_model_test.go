package ca_game_api

import (
	"testing"
)

func TestDecideRarity(t *testing.T) {
	tests := []string{"test1", "test2", "test3", "test4", "test5"}
	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			got := decideRarity()
			if got != 3 && got != 4 && got != 5 {
				t.Errorf("rarity should be between 3 and 5, but %v", got)
			}
		})
	}
}

func TestDecideCharacterId(t *testing.T) {
	tests := []string{"test1", "test2", "test3", "test4", "test5"}
	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			got, err := decideCharacterId()
			if err != nil {
				t.Errorf("decideCharacterId failed: %v", err)
			}
			if 30000001 > got && got >= 60000000 {
				t.Errorf("rarity should be between 30000001 and 60000000, but %v", got)
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

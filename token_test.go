package ca_game_api

import "testing"

func TestGenerateRandomString(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "length of generated string is 1",
			args: args{n: 1},
			want: 1,
		},
		{
			name: "length of generated string is 32",
			args: args{n: 32},
			want: 32,
		},
		{
			name: "length of generated string is 64",
			args: args{n: 64},
			want: 64,
		},
		{
			name: "length of generated string is 128",
			args: args{n: 128},
			want: 128,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s, err := generateRandomString(test.args.n)
			if err != nil {
				t.Errorf("generateRandomString failed: %v", err)
			}
			if len(s) != test.want {
				t.Errorf("generated string should be %v, but %v", test.want, len(s))
			}
		})
	}
}

func TestHash(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "token length 1",
			args: args{token: "a"},
			want: "ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb",
		},
		{
			name: "token length 64",
			args: args{token: "ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb"},
			want: "da3811154d59c4267077ddd8bb768fa9b06399c486e1fc00485116b57c9872f5",
		},
		{
			name: "token length 128",
			args: args{token: "da3811154d59c4267077ddd8bb768fa9b06399c486e1fc00485116b57c9872f5da3811154d59c4267077ddd8bb768fa9b06399c486e1fc00485116b57c9872f5"},
			want: "047727487a77ff783255a0642dcf450d592b3c4a3f6610aecaafb1bfcda834ac",
		},
		{
			name: "token length 256",
			args: args{token: "047727487a77ff783255a0642dcf450d592b3c4a3f6610aecaafb1bfcda834ac3f2178a31ebf7bd6b7d3529b425b970f6ecaedcdfe021734bc57d889970f73747797385f7220e54c5bfefdd7706cd8b6832ead022dc6408ffeabadeb3c20a08cf26330b4cf60ecca802f4ad0ba2f12d138dcb3cd4c333a8b7fed0a8bfae20e34"},
			want: "9770b4c658eedffb4a8e46e115de31ce0844133749beb8b0e797d33758bc128d",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := hash(test.args.token); got != test.want {
				t.Errorf("hash should be %v, but %v", test.want, got)
			}
		})
	}
}

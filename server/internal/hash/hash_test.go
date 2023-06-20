package hash

import (
	"testing"
)

func Test_generateHash(t *testing.T) {
	type args struct {
		token string
		sol   uint64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test1",
			args: args{
				token: "token1",
				sol:   1,
			},
			want: "254c714ad1de6ac2a58b995585803f03da548f8716b02db1911b214223532d26",
		},
		{
			name: "Test2",
			args: args{
				token: "token2",
				sol:   2,
			},
			want: "64f499e1852d205fc53d2f339e350f7178a4c8cdc1ff3dd9eca123f437080bb8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateHash(tt.args.token, tt.args.sol); got != tt.want {
				t.Errorf("generateHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkHash(t *testing.T) {
	type args struct {
		hash string
		diff uint8
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test1",
			args: args{
				hash: "0000f7f5a65fb1247754e45c28aad6a4ec42d8a538f78b76d3c2181c0ea0b3b7",
				diff: 4,
			},
			want: true,
		},
		{
			name: "Test2",
			args: args{
				hash: "0000f7f5a65fb1247754e45c28aad6a4ec42d8a538f78b76d3c2181c0ea0b3b7",
				diff: 5,
			},
			want: false,
		},
		{
			name: "Test3",
			args: args{
				hash: "0f7f5a65fb1247754e45c28aad6a4ec42d8a538f78b76d3c2181c0ea0b3b7",
				diff: 1,
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckHash(tt.args.hash, tt.args.diff); got != tt.want {
				t.Errorf("checkHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

package service

import (
	"strings"
	"testing"
)

func Test_newToken(t *testing.T) {
	type args struct {
		ip string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "New token creation",
			args: args{
				ip: "192.0.2.1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newToken(tt.args.ip)
			if !strings.HasPrefix(string(got), tt.args.ip+separatorToken) {
				t.Errorf("newToken() = %v, want prefix %v", got, tt.args.ip+separatorToken)
			}
		})
	}
}

func TestToken_ip(t *testing.T) {
	tests := []struct {
		name string
		tr   Token
		want string
	}{
		{
			name: "Token IP extraction",
			tr:   Token("192.0.2.1-1234567890-abcdef"),
			want: "192.0.2.1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.ip(); got != tt.want {
				t.Errorf("Token.ip() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateRandomPart(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Generate Random Part",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateRandomPart()
			if len(got) != 32 { // MD5 hash is always 32 characters
				t.Errorf("generateRandomPart() = %v, want length 32", got)
			}
		})
	}
}

package utils

import (
	"testing"
)

func TestGetDotEnvVar(t *testing.T) {

	tests := []struct {
		name string
		key  string
		want string
	}{
		{name: "db URI not empty string", key: "DATABASE_URI", want: ""},
		{name: "db DB_PWD not empty string", key: "DB_PWD", want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// t.Errorf("got: %s, want: %s", GetDotEnvVar(tt.key), tt.want)
			if got := GetDotEnvVar(tt.key); got == tt.want {
				t.Errorf("GetDotEnvVar() = %v, want %v", got, tt.want)
			}
		})
	}
}
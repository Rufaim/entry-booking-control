package main

import (
	"bytes"
	"testing"
)

func Test_intToBytes(t *testing.T) {
	tests := []struct {
		name string
		in   int
		want []byte
	}{
		{"Zero", 0, []byte("0")},
		{"One", 1, []byte("1")},
		{"FortyTwo", 42, []byte("42")},
		{"NegativeEleven", -11, []byte("-11")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := intToBytes(tt.in); !bytes.Equal(got, tt.want) {
				t.Errorf("intToBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bytesToInt(t *testing.T) {
	tests := []struct {
		name string
		in   []byte
		want int
		err  bool
	}{
		{"Zero", []byte("0"), 0, false},
		{"One", []byte("1"), 1, false},
		{"FortyTwo", []byte("42"), 42, false},
		{"NegativeEleven", []byte("-11"), -11, false},
		{"WrongInput", []byte("qwerty"), 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := bytesToInt(tt.in)
			if (err != nil) != tt.err {
				t.Errorf("bytesToInt() error = %v, wantErr %v", err, tt.err)
				return
			}
			if got != tt.want {
				t.Errorf("bytesToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

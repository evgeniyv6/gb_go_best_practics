package parser

import "testing"

func TestToken_IsEmpty(t *testing.T) {
	token := Token{}
	if !token.IsEmpty() {
		t.Error("token.IsEmpty() - false, want - true for empty token")
	}
}

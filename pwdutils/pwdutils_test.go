package pwdutils

import (
	"testing"
)

func TestPasswordHashAndVerify(t *testing.T) {
	password := "mySecret123!"

	hash, err := PasswordHash(password)
	if err != nil {
		t.Fatalf("PasswordHash error: %v", err)
	}
	if hash == "" {
		t.Fatal("PasswordHash returned empty string")
	}

	// 正确密码应验证通过
	if !PasswordVerify(hash, password) {
		t.Error("PasswordVerify failed for correct password")
	}

	// 错误密码应验证失败
	if PasswordVerify(hash, "wrongPassword") {
		t.Error("PasswordVerify passed for wrong password")
	}
}

func TestPasswordHash_UniqueHashes(t *testing.T) {
	password := "samePassword"
	hash1, err1 := PasswordHash(password)
	hash2, err2 := PasswordHash(password)
	if err1 != nil || err2 != nil {
		t.Fatalf("PasswordHash error: %v, %v", err1, err2)
	}
	if hash1 == hash2 {
		t.Error("PasswordHash should generate different hashes for the same password (due to salt)")
	}
}

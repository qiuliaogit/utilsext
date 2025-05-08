package pwdutils

import "golang.org/x/crypto/bcrypt"

// php里面用的password_hash()和password_verify()函数，golang里面也有相应的库，比如golang.org/x/crypto/bcrypt
// 这里我们使用golang.org/x/crypto/bcrypt库来实现密码哈希和验证
// 密码哈希：将密码加密成不可逆的字符串，一般用于存储密码
// 密码验证：将用户输入的密码与数据库中存储的密码哈希进行比较，如果相同，则验证通过，否则验证失败
// 下面的两个函数，分别实现了密码哈希和验证的功能

// 生成密码哈希
func PasswordHash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// 验证密码哈希
func PasswordVerify(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

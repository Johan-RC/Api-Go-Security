package bcrypt

import "golang.org/x/crypto/bcrypt"

// HashPassword genera el hash bcrypt de una contraseña en texto plano.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// ComparePassword verifica que la contraseña coincida con su hash.
func ComparePassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

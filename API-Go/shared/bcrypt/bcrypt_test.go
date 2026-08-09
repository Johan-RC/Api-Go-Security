package bcrypt

import "testing"

func TestHashAndCompare(t *testing.T) {
	hash, err := HashPassword("Admin#2026")
	if err != nil {
		t.Fatalf("HashPassword error: %v", err)
	}
	if hash == "" {
		t.Fatal("hash vacío")
	}
	if hash == "Admin#2026" {
		t.Fatal("el hash no debe ser el texto plano")
	}
	if !ComparePassword("Admin#2026", hash) {
		t.Error("ComparePassword debería ser true con la contraseña correcta")
	}
	if ComparePassword("wrong-pass", hash) {
		t.Error("ComparePassword debería ser false con contraseña incorrecta")
	}
}

func TestCompareWithKnownHash(t *testing.T) {
	// Hash conocido generado con DefaultCost ($2a$10$).
	const knownHash = "$2a$10$B5QpFbxXTLuASMtERZxZBuGh4zOKqNHL/nUB60/5KkAErtbV4LsE6"
	if !ComparePassword("Admin#2026", knownHash) {
		t.Error("el hash conocido debería verificar la contraseña seguida")
	}
}
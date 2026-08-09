package validation

import (
	"regexp"
	"strings"
	"unicode"
)

// emailRegex basado en el RFC 5322 para la parte local, con dominio al menos
// de dos niveles y TLD alfabético de 2 a 63 caracteres.
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*\.[a-zA-Z]{2,63}$`)

// reservedLocalParts son partes locales que suelen ser sistemas o listas de
// distribución y no bandejas personales.
var reservedLocalParts = map[string]bool{
	"abuse": true, "admin": true, "bounce": true, "contact": true,
	"devnull": true, "info": true, "list": true, "mailer-daemon": true,
	"no-reply": true, "noreply": true, "postmaster": true, "root": true,
	"sales": true, "security": true, "support": true, "team": true, "webmaster": true,
}

// ValidateEmail valida el formato del correo y descarta partes locales
// reservadas (no bandejas personales).
func ValidateEmail(email string) bool {
	email = strings.TrimSpace(email)
	if len(email) > 254 || !strings.Contains(email, "@") {
		return false
	}

	if !emailRegex.MatchString(email) {
		return false
	}

	local := strings.Split(email, "@")[0]
	return !reservedLocalParts[strings.ToLower(local)]
}

// ValidatePassword valida la política de contraseña: mínimo 8 caracteres,
// al menos una mayúscula, una minúscula, un dígito y un símbolo.
func ValidatePassword(password string) bool {
	password = strings.TrimSpace(password)
	if len(password) < 8 {
		return false
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r) || unicode.IsSpace(r):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasDigit && hasSpecial
}

// ValidateName valida que el texto solo tenga letras, espacios, tildes,
// ñ/Ñ y apóstrofes (nombres propios), con un límite de longitud.
func ValidateName(name string, maxLen int) bool {
	name = strings.TrimSpace(name)
	if name == "" || len(name) > maxLen {
		return false
	}

	for _, r := range name {
		if !unicode.IsLetter(r) && !unicode.IsSpace(r) && r != '\'' && r != '-' {
			return false
		}
	}
	return true
}
package apperror

import "errors"

// Error representa un error de dominio con código y estado HTTP.
type Error struct {
	Code    string
	Message string
	Status  int
	Err     error
}

// Error implementa la interfaz error.
func (e *Error) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

// Unwrap permite usar errors.Is/As sobre el error raíz.
func (e *Error) Unwrap() error {
	return e.Err
}

// New crea un Error de dominio.
func New(code, message string, status int) *Error {
	return &Error{Code: code, Message: message, Status: status}
}

// Wrap envuelve un error raíz dentro de un Error de dominio.
func Wrap(code, message string, status int, err error) *Error {
	return &Error{Code: code, Message: message, Status: status, Err: err}
}

// AsError extrae el *Error de dominio de una cadena de errores.
func AsError(err error) (*Error, bool) {
	var e *Error
	if errors.As(err, &e) {
		return e, true
	}
	return nil, false
}

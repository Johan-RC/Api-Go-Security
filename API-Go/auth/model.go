package auth

// El módulo auth no define modelos propios: reutiliza los de otros módulos:
//
//   - user.User
//   - role.Role, role.UserRole
//   - session.RefreshToken, session.PasswordResetRequest
//   - audit.AuditLogin

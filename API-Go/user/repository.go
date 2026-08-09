package user

import (
	"time"

	"gorm.io/gorm"
)

// Repository encapsula el acceso a datos de usuarios (solo GORM).
type Repository struct {
	db *gorm.DB
}

// NewRepository crea el repositorio de usuarios.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create inserta un nuevo usuario.
func (r *Repository) Create(u *User) error {
	// Select incluye explícitamente campos con zero-value (is_active, etc.)
	// para que la BD no aplique los defaults sobre ellos.
	return r.db.Select(
		"id", "email", "password_hash", "first_name", "last_name",
		"actor_type", "actor_id", "is_active", "failed_attempts",
	).Create(u).Error
}

// Update actualiza los campos non-zero del usuario por su ID.
func (r *Repository) Update(u *User) error {
	return r.db.Save(u).Error
}

// Delete elimina un usuario por su ID.
func (r *Repository) Delete(id string) error {
	return r.db.Delete(&User{}, "id = ?", id).Error
}

// FindByID obtiene un usuario por su ID.
func (r *Repository) FindByID(id string) (*User, error) {
	var u User
	if err := r.db.Where("id = ?", id).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// FindByEmail obtiene un usuario por su email.
func (r *Repository) FindByEmail(email string) (*User, error) {
	var u User
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// List obtiene usuarios paginados.
func (r *Repository) List(offset, limit int) ([]User, error) {
	var users []User
	if err := r.db.
		Order("last_name ASC, first_name ASC").
		Offset(offset).Limit(limit).
		Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// Count devuelve el total de usuarios.
func (r *Repository) Count() (int64, error) {
	var count int64
	if err := r.db.Model(&User{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// ExistsByEmail indica si ya existe un usuario con ese email.
func (r *Repository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.db.Model(&User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// UpdatePassword actualiza el hash de contraseña de un usuario.
func (r *Repository) UpdatePassword(id, passwordHash string) error {
	return r.db.Model(&User{}).Where("id = ?", id).
		Update("password_hash", passwordHash).Error
}

// UpdateLastLogin registra la fecha del último inicio de sesión.
func (r *Repository) UpdateLastLogin(id string, lastLoginAt time.Time) error {
	return r.db.Model(&User{}).Where("id = ?", id).
		Update("last_login_at", lastLoginAt).Error
}

// IncrementFailedAttempts suma 1 al contador de intentos fallidos.
func (r *Repository) IncrementFailedAttempts(id string) error {
	return r.db.Model(&User{}).Where("id = ?", id).
		UpdateColumn("failed_attempts", gorm.Expr("failed_attempts + 1")).Error
}

// ResetFailedAttempts pone el contador de intentos fallidos en 0.
func (r *Repository) ResetFailedAttempts(id string) error {
	return r.db.Model(&User{}).Where("id = ?", id).
		Update("failed_attempts", 0).Error
}

// SetLocked bloquea a un usuario hasta la fecha indicada (NULL desbloquea).
func (r *Repository) SetLocked(id string, lockedUntil *time.Time) error {
	return r.db.Model(&User{}).Where("id = ?", id).
		Update("locked_until", lockedUntil).Error
}

// SetActive activa o desactiva un usuario.
func (r *Repository) SetActive(id string, active bool) error {
	return r.db.Model(&User{}).Where("id = ?", id).
		Update("is_active", active).Error
}
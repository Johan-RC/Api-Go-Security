package feature

// CreateModuleRequest es la entrada para crear un módulo.
type CreateModuleRequest struct {
	Code         string `json:"code" binding:"required,max=30"`
	Name         string `json:"name" binding:"required,max=100"`
	Description  *string `json:"description"`
	DisplayOrder int16   `json:"display_order" binding:"omitempty"`
	IconKey      *string `json:"icon_key"`
}

// CreateFeatureRequest es la entrada para crear una funcionalidad.
type CreateFeatureRequest struct {
	ModuleID    string `json:"module_id" binding:"required"`
	Code        string `json:"code" binding:"required,max=60"`
	Name        string `json:"name" binding:"required,max=120"`
	Description *string `json:"description"`
	ActionLevel string `json:"action_level" binding:"required,oneof=READ WRITE DELETE PUBLISH APPROVE"`
}

// UpdateModuleRequest permite actualizar los campos editables de un módulo.
type UpdateModuleRequest struct {
	Name         *string `json:"name"`
	Description  *string `json:"description"`
	DisplayOrder *int16  `json:"display_order"`
	IconKey      *string `json:"icon_key"`
	IsActive     *bool   `json:"is_active"`
}

// UpdateFeatureRequest permite actualizar los campos editables de una funcionalidad.
type UpdateFeatureRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	ActionLevel *string `json:"action_level" binding:"omitempty,oneof=READ WRITE DELETE PUBLISH APPROVE"`
	IsActive    *bool   `json:"is_active"`
}

// ModuleResponse es la salida de un módulo.
type ModuleResponse struct {
	ID           string  `json:"id"`
	Code         string  `json:"code"`
	Name         string  `json:"name"`
	DisplayOrder int16   `json:"display_order"`
	IsActive     bool    `json:"is_active"`
}

// FeatureResponse es la salida de una funcionalidad.
type FeatureResponse struct {
	ID          string `json:"id"`
	ModuleID    string `json:"module_id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	ActionLevel string `json:"action_level"`
	IsActive    bool   `json:"is_active"`
}

// ListModulesResponse es la salida paginada de módulos.
type ListModulesResponse struct {
	Items []ModuleResponse `json:"items"`
	Total int64            `json:"total"`
	Page  int              `json:"page"`
}

// ListFeaturesResponse es la salida paginada de funcionalidades.
type ListFeaturesResponse struct {
	Items []FeatureResponse `json:"items"`
	Total int64             `json:"total"`
	Page  int               `json:"page"`
}

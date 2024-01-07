package model

// FieldError kullanıcıya döndürülecek hata formatını tanımlar
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

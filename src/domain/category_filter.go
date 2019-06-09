package domain

type CategoryFilter struct {
	Name   string `form:"name"`
	Hidden bool   `form:"hidden"`
	Limit  int64  `form:"limit"`
}

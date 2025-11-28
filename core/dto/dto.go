package dto

type IPatchDTO interface {
	ToPatchMap() map[string]interface{}

	IsEmpty() bool
}

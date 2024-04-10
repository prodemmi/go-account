package models

type App struct {
	BaseModel

	Name string `json:"name"`
	Url  string `json:"url"`
}

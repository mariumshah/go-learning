package domain

type Author struct {
	ID       int    `json:"id"`
	Name     string `json:"author_name"`
	Password string `json:"-"`
}

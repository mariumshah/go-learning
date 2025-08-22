package domain

type Book struct {
	ID            int
	Title         string
	Author        int //author id
	AuthorDetails *Author
}

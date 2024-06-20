package repository

type Authorization interface {
}

type Post interface{}

type Subscription interface {
}

type Repository struct {
	Authorization
	Post
	Subscription
}

func NewRepository() *Repository {
	return &Repository{}
}

package main

import "context"

type Repository interface {
	Select(ctx context.Context, id int64) (*Person, error)
	SelectAll(ctx context.Context) ([]*Person, error)
	Insert(ctx context.Context, p *Person) error
	Update(ctx context.Context, p *Person) error
}

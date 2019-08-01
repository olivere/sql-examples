package main

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

var _ Repository = (*MySQLRepository)(nil)

type MySQLRepository struct {
	db *sql.DB
}

func NewMySQLRepository(urn string) (*MySQLRepository, error) {
	db, err := sql.Open("mysql", urn)
	if err != nil {
		return nil, errors.Wrap(err, "unable to open DB")
	}

	if _, err := db.Query(`SELECT 1 FROM people`); err != nil {
		// Create "people" table
		_, err = db.Exec(`CREATE TABLE people (
			id BIGINT NOT NULL AUTO_INCREMENT,
			name VARCHAR(100) NOT NULL,
			created DATETIME(6) NOT NULL,
			PRIMARY KEY (id)
			) ENGINE=InnoDB;
			`)
		if err != nil {
			return nil, errors.Wrap(err, "unable to create people table")
		}
		// Insert Alice as first person in that table
		alice := &Person{
			ID:      1,
			Name:    "Alice",
			Created: time.Now().UTC(),
		}
		_, err = db.Exec(`INSERT INTO people (id,name,created) VALUES (?,?,?)`, alice.ID, alice.Name, alice.Created)
		if err != nil {
			return nil, errors.Wrap(err, "unable to insert Alice into people table")
		}
	}

	return &MySQLRepository{
		db: db,
	}, nil
}

func (r *MySQLRepository) Select(ctx context.Context, id int64) (*Person, error) {
	panic("not implemented yet")
}

func (r *MySQLRepository) SelectAll(ctx context.Context) ([]*Person, error) {
	var people []*Person
	rows, err := r.db.QueryContext(ctx, `SELECT * FROM people ORDER BY id`)
	if err != nil {
		return nil, errors.Wrap(err, "unable to query rows")
	}
	defer rows.Close()
	for rows.Next() {
		p := new(Person)
		if err := rows.Scan(&p.ID, &p.Name, &p.Created); err != nil {
			return nil, errors.Wrap(err, "unable to scan row")
		}
		select {
		default:
			people = append(people, p)
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "unable to iterate over resultset")
	}
	return people, nil
}

func (r *MySQLRepository) Insert(ctx context.Context, p *Person) error {
	if p.Created.IsZero() {
		p.Created = time.Now().UTC()
	}
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO people (name,created) VALUES (?,?)`,
		p.Name,
		p.Created,
	)
	if err != nil {
		return errors.Wrap(err, "unable to insert person")
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "unable to insert person")
	}
	if rowsAffected != 1 {
		return errors.New("unable to insert person")
	}
	p.ID, err = res.LastInsertId()
	if err != nil {
		return errors.Wrapf(err, "unable to insert person")
	}
	return nil
}

func (r *MySQLRepository) Update(ctx context.Context, p *Person) error {
	panic("not implemented yet")
}

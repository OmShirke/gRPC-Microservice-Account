package account //This ensures its types and functions can be used throughout the account module.

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
)

type Repo interface {
	Close()
	PutAccount(ctx context.Context, a Account) error
	GetAccountByID(ctx context.Context, id string) (*Account, error)
	ListAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error)
}

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(url string) (Repo, error) {
	db, err := sql.Open("postgres", url) //string matches the driver you’ve registered in your application
	if err != nil {
		return nil, err
	}

	err = db.Ping() //check if connection in active
	if err != nil {
		return nil, err
	}

	return &postgresRepo{db}, nil
}

func (r *postgresRepo) Close() {
	r.db.Close()
}

func (r *postgresRepo) Ping() { //checks if connection is active
	r.db.Ping()
}

func (r *postgresRepo) PutAccount(ctx context.Context, a Account) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO accounts(id, name) VALUES($1, $2)", a.ID, a.Name) //Executes given SQL query in context of provided ctx.
	return err
}

func (r *postgresRepo) GetAccountByID(ctx context.Context, id string) (*Account, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, name FROM account WHERE id=$1", id) // Executes a query that is expected to return exactly one row (or none)
	a := &Account{}
	if err := row.Scan(&a.ID, &a.Name); err != nil { //Copies the values from the database result into the fields of the Account struct
		return nil, err
	}
	return a, nil
}

func (r *postgresRepo) ListAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {
	rows, err := r.db.QueryContext( // Executes a query expected to return multiple rows
		ctx,
		"SELECT id, name FROM accounts ORDER BY id DESC OFFSET $1 LIMIT $2",
		skip,
		take,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := []Account{}

	for rows.Next() {
		a := &Account{}
		if err = rows.Scan(&a.ID, &a.Name); err == nil {
			accounts = append(accounts, *a)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return accounts, nil
}

package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Gonzapepe/bank-api/helper"
	"github.com/Gonzapepe/bank-api/types"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)


type Storage interface {
	CreateAccount(*types.Account) error
	DeleteAccount(int) error
	UpdateAccount(*types.Account) error
	GetAccounts() ([]*types.Account, error)
	GetAccountByID(int) (*types.Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	password := os.Getenv("DB_PASSWORD")
	sslmode := os.Getenv("DB_SSLMODE")

	connStr := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=%s", user, dbname, password, sslmode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	return s.createAccountTable()

}

func (s *PostgresStore) createAccountTable() error {
	query := `create table if not exists account (
		id serial primary key not null,
		first_name varchar(50),
		last_name varchar(50),
		gender int,
		dni bigint,
		cuit bigint,
		balance bigint,
		created_at timestamp default current_timestamp,
		updated_at timestamp default current_timestamp
	)`

	_, err := s.db.Exec(query)

	return err
}

func (s *PostgresStore) CreateAccount(account *types.Account) error {
	query := `
	insert into account 
	(first_name, last_name, gender, dni, cuit, balance)
	values ($1, $2, $3, $4, $5, $6)
	`
	_, err := s.db.Exec(query, account.FirstName, account.LastName, account.Gender, account.Dni, account.Cuit, account.Balance)

	if err != nil {
		return err
	}


	return nil
}

func (s *PostgresStore) UpdateAccount(account *types.Account) error {
	query := `
	update account 
	set first_name = $1, last_name = $2, gender = $3, dni = $4, cuit = $5, balance = $6, updated_at = $7
	where id = $8
	`
	cuit, err := helper.Cuit(account.Dni, account.Gender)

	if err != nil {
		return err
	}
	
	_, err = s.db.Exec(query, account.FirstName, account.LastName, account.Gender, account.Dni, cuit, account.Balance, account.UpdatedAt, account.ID)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	query := `delete from account where id = $1`

	_, err := s.db.Exec(query, id)

	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) GetAccountByID(id int) (*types.Account, error) {
	query := `select * from account where id = $1`

	account := &types.Account{}

	err := s.db.QueryRow(query, id).Scan(&account)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("No matching found")
		}
	}

	return account, nil
}

func (s *PostgresStore) GetAccounts() ([]*types.Account, error) {

	rows, err := s.db.Query("select * from account")

	if err != nil {
		return nil, err
	}

	accounts := []*types.Account{}

	for rows.Next() {
		account := &types.Account{}
		err := rows.Scan(
			&account.ID,
			&account.FirstName,
			&account.LastName,
			&account.Gender,
			&account.Dni,
			&account.Cuit,
			&account.Balance,
			&account.CreatedAt,
			&account.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}
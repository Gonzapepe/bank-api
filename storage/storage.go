package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Gonzapepe/bank-api/types"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)


type Storage interface {
	CreateAccount(*types.Account) error
	DeleteAccount(int) error
	UpdateAccount(*types.Account) error
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
		id serial primary key not null auto_increment,
		first_name varchar(50),
		last_name varchar(50),
		gender int,
		dni bigint,
		cbu bigint,
		balance bigint,
		created_at timestamp default current_timestamp,
		updated_at timestamp default current_timestamp
	)`

	_, err := s.db.Exec(query)

	return err
}

func (s *PostgresStore) CreateAccount(*types.Account) error {
	return nil
}

func (s *PostgresStore) UpdateAccount(*types.Account) error {
	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	return nil
}

func (s *PostgresStore) GetAccountByID(id int) (*types.Account, error) {
	return nil, nil
}
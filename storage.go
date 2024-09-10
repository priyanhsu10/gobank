package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"log"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccont(int) error
	UpdateAccount(*Account) error
	GetAccountById(int) (*Account, error)
	GetAllAccounts() ([]*Account, error)
}

type PostgressStore struct {
	db *sql.DB
}

func NewPostgressStore() (*PostgressStore, error) {
	connStr := "user=fedpet dbname=tododb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgressStore{db: db}, nil
}

func (s *PostgressStore) init() error {
	return s.createAccountTable()
}

func (s *PostgressStore) createAccountTable() error {

	query := `create table  if not exists account(
	 id serial primary key,
	 first_name varchar(50),
	 last_name varchar(50),
	 number int,
	 balance serial,
	 created_at timestamp
	)`

	_, err := s.db.Exec(query)
	return err

}

func (s *PostgressStore) CreateAccount(account *Account) error {

	query := `insert into account(first_name,last_name,number,balance,created_at)
	 values($1,$2,$3,$4,$5)`
	var id int
	r, err := s.db.Query(query, account.FirstName, account.LastName, account.Number, account.Balance, account.CreatedAt)
	r.Scan(&id)
	fmt.Println("create id ", id)
	fmt.Printf("%+v", r)
	return err
}
func (s *PostgressStore) DeleteAccont(id int) error {

	query := `delete from account where id=$1`
	var deletedid int
	result, err := s.db.Query(query, id)
	result.Scan(&deletedid)
	fmt.Println("delete id ", deletedid)
	return err
}
func (s *PostgressStore) UpdateAccount(account *Account) error {
	query := `update account set first_name=$1,last_name=$2,number=$3,balance=$4 where ID=$5`
	_, err := s.db.Query(query, account.FirstName, account.LastName, account.Number, account.Balance, account.ID)

	if err != nil {
		return err
	}
	return nil
}
func (s *PostgressStore) GetAccountById(id int) (*Account, error) {

	query := `select * from  account where id=$1`
	r, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	r.Next()
	account, err := bind(r)
	if err != nil {
		return nil, fmt.Errorf("account not found")
	}

	return account, nil
}
func bind(result *sql.Rows) (*Account, error) {
	account := &Account{}
	err := result.Scan(&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt)

	if err != nil {
		return nil, err
	}
	return account, nil
}
func (s *PostgressStore) GetAllAccounts() ([]*Account, error) {

	query := `select * from account`
	result, err := s.db.Query(query)

	if err != nil {
		return nil, err
	}

	accounts := []*Account{}
	for result.Next() {

		account, err := bind(result)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)

	}
	return accounts, nil
}

package database

import (
	"database/sql"

	"github.com.br/devfullcycle/fc-ms-balances/internal/entity"
)

type BalanceDB struct {
	DB *sql.DB
}

func NewBalanceDB(db *sql.DB) *BalanceDB {
	return &BalanceDB{
		DB: db,
	}
}

func (b *BalanceDB) FindByAccountID(accountID string) (*entity.Balance, error) {
	var balance entity.Balance

	stmt, err := b.DB.Prepare("SELECT * FROM balances WHERE account_id = ? ORDER BY created_at DESC LIMIT 1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(accountID)
	err = row.Scan(
		&balance.ID,
		&balance.AccountID,
		&balance.Balance,
		&balance.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &balance, nil
}

func (b *BalanceDB) FindByID(id string) (*entity.Balance, error) {
	var balance entity.Balance

	stmt, err := b.DB.Prepare("SELECT * FROM balances WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(id)
	err = row.Scan(
		&balance.ID,
		&balance.AccountID,
		&balance.Balance,
		&balance.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &balance, nil
}

func (b *BalanceDB) Save(balance *entity.Balance) error {
	stmt, err := b.DB.Prepare("INSERT INTO balances (id, account_id, balance, created_at) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(balance.ID, balance.AccountID, balance.Balance, balance.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (b *BalanceDB) Update(balance *entity.Balance) error {
	stmt, err := b.DB.Prepare("UPDATE balances SET balance = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(balance.Balance, balance.ID)
	if err != nil {
		return err
	}
	return nil
}

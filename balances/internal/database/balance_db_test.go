package database

import (
	"database/sql"
	"testing"

	"github.com.br/devfullcycle/fc-ms-balances/internal/entity"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
)

type BalanceDBTestSuite struct {
	suite.Suite
	db        *sql.DB
	balanceDB *BalanceDB
}

func (s *BalanceDBTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.Nil(err)
	s.db = db
	db.Exec("Create TABLE balances (id varchar(255), account_id varchar(255), balance int, created_at timestamp)")
	s.balanceDB = NewBalanceDB(db)
}

func (s *BalanceDBTestSuite) TearDownSuite() {
	defer s.db.Close()
	s.db.Exec("DROP TABLE balances")
}

func TestBalanceDBTestSuite(t *testing.T) {
	suite.Run(t, new(BalanceDBTestSuite))
}

func (s *BalanceDBTestSuite) TestSave() {
	balance := entity.NewBalance("1", 100.00)
	err := s.balanceDB.Save(balance)
	s.Nil(err)
}

func (s *BalanceDBTestSuite) TestFindByID() {
	balance := entity.NewBalance("1", 100.00)
	s.balanceDB.Save(balance)

	balanceDB, err := s.balanceDB.FindByID(balance.ID)
	s.Nil(err)
	s.Equal(balance.ID, balanceDB.ID)
	s.Equal(balance.AccountID, balanceDB.AccountID)
	s.Equal(balance.Balance, balanceDB.Balance)
}

func (s *BalanceDBTestSuite) TestUpdate() {
	balance := entity.NewBalance("1", 100.00)
	s.balanceDB.Save(balance)

	balance.Balance = 200.00
	err := s.balanceDB.Update(balance)
	s.Nil(err)

	balanceDB, err := s.balanceDB.FindByID(balance.ID)
	s.Nil(err)
	s.Equal(balance.ID, balanceDB.ID)
	s.Equal(balance.AccountID, balanceDB.AccountID)
	s.Equal(balance.Balance, balanceDB.Balance)
}

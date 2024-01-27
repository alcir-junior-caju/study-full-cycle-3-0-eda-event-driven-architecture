package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com.br/devfullcycle/fc-ms-wallet/internal/database"
	"github.com.br/devfullcycle/fc-ms-wallet/internal/event"
	"github.com.br/devfullcycle/fc-ms-wallet/internal/event/handler"
	createaccount "github.com.br/devfullcycle/fc-ms-wallet/internal/usecase/create_account"
	"github.com.br/devfullcycle/fc-ms-wallet/internal/usecase/create_client"
	"github.com.br/devfullcycle/fc-ms-wallet/internal/usecase/create_transaction"
	"github.com.br/devfullcycle/fc-ms-wallet/internal/web"
	"github.com.br/devfullcycle/fc-ms-wallet/internal/web/webserver"
	"github.com.br/devfullcycle/fc-ms-wallet/pkg/events"
	"github.com.br/devfullcycle/fc-ms-wallet/pkg/kafka"
	"github.com.br/devfullcycle/fc-ms-wallet/pkg/uow"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", "root", "root", "wallet-mysql", "3306", "wallet"))
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS clients (id varchar(255), name varchar(255), email varchar(255), created_at date)")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS accounts (id varchar(255), client_id varchar(255), balance integer, created_at date)")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS transactions (id varchar(255), account_id_from varchar(255), account_id_to varchar(255), amount integer, created_at date)")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("INSERT INTO clients (id, name, email, created_at) VALUES ('d6a3b836-bbd0-11ee-8a24-0242ac120003', 'John', 'johndoe@email.com', '2020-01-01'), ('d6a3bac8-bbd0-11ee-8a24-0242ac120003', 'Jane', 'janedoe@email.com', '2020-01-02'), ('d6a3bb72-bbd0-11ee-8a24-0242ac120003', 'Alice', 'alice@email.com', '2020-01-03'), ('d6a3bbc1-bbd0-11ee-8a24-0242ac120003', 'Bob', 'bob@email.com', '2020-01-04'), ('d6a3bc01-bbd0-11ee-8a24-0242ac120003', 'Charlie', 'charlie@email.com', '2020-01-05'), ('d6a3bc3e-bbd0-11ee-8a24-0242ac120003', 'David', 'david@email.com', '2020-01-06'), ('d6a3bc81-bbd0-11ee-8a24-0242ac120003', 'Eva', 'eva@email.com', '2020-01-07'), ('d6a3bcb9-bbd0-11ee-8a24-0242ac120003', 'Frank', 'frank@email.com', '2020-01-08'), ('d6a3bd0a-bbd0-11ee-8a24-0242ac120003', 'Grace', 'grace@email.com', '2020-01-09'), ('d6a3bd5c-bbd0-11ee-8a24-0242ac120003', 'Harry', 'harry@email.com', '2020-01-10')")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("INSERT INTO accounts (id, client_id, balance, created_at) VALUES ('d6a4053e-bbd0-11ee-8a24-0242ac120003', (SELECT id FROM clients WHERE name = 'John'), 1000, '2020-01-01'), ('d6a40ce9-bbd0-11ee-8a24-0242ac120003', (SELECT id FROM clients WHERE name = 'Jane'), 1000, '2020-01-02'), ('d6a40f3a-bbd0-11ee-8a24-0242ac120003', (SELECT id FROM clients WHERE name = 'Alice'), 1000, '2020-01-03'), ('d6a41083-bbd0-11ee-8a24-0242ac120003', (SELECT id FROM clients WHERE name = 'Bob'), 1000, '2020-01-04'), ('d6a411ce-bbd0-11ee-8a24-0242ac120003', (SELECT id FROM clients WHERE name = 'Charlie'), 1000, '2020-01-05'), ('d6a412fe-bbd0-11ee-8a24-0242ac120003', (SELECT id FROM clients WHERE name = 'David'), 1000, '2020-01-06'), ('d6a41403-bbd0-11ee-8a24-0242ac120003', (SELECT id FROM clients WHERE name = 'Eva'), 800, '2020-01-07'), ('d6a41504-bbd0-11ee-8a24-0242ac120003', (SELECT id FROM clients WHERE name = 'Frank'), 1200, '2020-01-08'), ('d6a41604-bbd0-11ee-8a24-0242ac120003', (SELECT id FROM clients WHERE name = 'Grace'), 1000, '2020-01-09'), ('d6a416ee-bbd0-11ee-8a24-0242ac120003', (SELECT id FROM clients WHERE name = 'Harry'), 1000, '2020-01-10')")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("INSERT INTO transactions (id, account_id_from, account_id_to, amount, created_at) VALUES (UUID(), (SELECT id FROM accounts WHERE client_id = (SELECT id FROM clients WHERE name = 'John')), (SELECT id FROM accounts WHERE client_id = (SELECT id FROM clients WHERE name = 'Jane')), 100, '2020-01-01'), (UUID(), (SELECT id FROM accounts WHERE client_id = (SELECT id FROM clients WHERE name = 'Alice')), (SELECT id FROM accounts WHERE client_id = (SELECT id FROM clients WHERE name = 'Bob')), 100, '2020-01-02'), (UUID(), (SELECT id FROM accounts WHERE client_id = (SELECT id FROM clients WHERE name = 'Charlie')), (SELECT id FROM accounts WHERE client_id = (SELECT id FROM clients WHERE name = 'David')), 100, '2020-01-03'), (UUID(), (SELECT id FROM accounts WHERE client_id = (SELECT id FROM clients WHERE name = 'Eva')), (SELECT id FROM accounts WHERE client_id = (SELECT id FROM clients WHERE name = 'Frank')), 100, '2020-01-04'), (UUID(), (SELECT id FROM accounts WHERE client_id = (SELECT id FROM clients WHERE name = 'Grace')), (SELECT id FROM accounts WHERE client_id = (SELECT id FROM clients WHERE name = 'Harry')), 100, '2020-01-05')")
	if err != nil {
		panic(err)
	}

	configMap := ckafka.ConfigMap{
		"bootstrap.servers": "kafka:29092",
		"group.id":          "wallet",
	}
	kafkaProducer := kafka.NewKafkaProducer(&configMap)

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("TransactionCreated", handler.NewTransactionCreatedKafkaHandler(kafkaProducer))
	eventDispatcher.Register("BalanceUpdated", handler.NewUpdateBalanceKafkaHandler(kafkaProducer))
	transactionCreatedEvent := event.NewTransactionCreated()
	balanceUpdatedEvent := event.NewBalanceUpdated()

	clientDb := database.NewClientDB(db)
	accountDb := database.NewAccountDB(db)

	ctx := context.Background()
	uow := uow.NewUow(ctx, db)

	uow.Register("AccountDB", func(tx *sql.Tx) interface{} {
		return database.NewAccountDB(db)
	})

	uow.Register("TransactionDB", func(tx *sql.Tx) interface{} {
		return database.NewTransactionDB(db)
	})
	createTransactionUseCase := create_transaction.NewCreateTransactionUseCase(uow, eventDispatcher, transactionCreatedEvent, balanceUpdatedEvent)
	createClientUseCase := create_client.NewCreateClientUseCase(clientDb)
	createAccountUseCase := createaccount.NewCreateAccountUseCase(accountDb, clientDb)

	webserver := webserver.NewWebServer(":8080")

	clientHandler := web.NewWebClientHandler(*createClientUseCase)
	accountHandler := web.NewWebAccountHandler(*createAccountUseCase)
	transactionHandler := web.NewWebTransactionHandler(*createTransactionUseCase)

	webserver.AddHandler("/clients", clientHandler.CreateClient)
	webserver.AddHandler("/accounts", accountHandler.CreateAccount)
	webserver.AddHandler("/transactions", transactionHandler.CreateTransaction)

	fmt.Println("Server is running at port 8080")
	webserver.Start()
}

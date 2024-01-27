package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com.br/devfullcycle/fc-ms-balances/internal/database"
	"github.com.br/devfullcycle/fc-ms-balances/internal/usecase/create_balance"
	"github.com.br/devfullcycle/fc-ms-balances/internal/usecase/find_account"
	"github.com.br/devfullcycle/fc-ms-balances/internal/web"
	"github.com.br/devfullcycle/fc-ms-balances/internal/web/webserver"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/go-sql-driver/mysql"
)

type BalanceEvent struct {
	Name    string      `json:"Name"`
	Payload BalanceData `json:"Payload"`
}

type BalanceData struct {
	AccountIDFrom      string  `json:"account_id_from"`
	AccountIDTo        string  `json:"account_id_to"`
	BalanceAccountFrom float64 `json:"balance_account_id_from"`
	BalanceAccountTo   float64 `json:"balance_account_id_to"`
}

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", "root", "root", "balances-mysql", "3306", "balances"))
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS balances (id varchar(255), account_id varchar(255), balance integer, created_at timestamp)")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("INSERT INTO balances (id, account_id, balance, created_at) VALUES ('b8bbd006-0318-48ba-b7a4-30dc5920f4c2', 'd6a41403-bbd0-11ee-8a24-0242ac120003', 850, '2024-01-27 17:13:48'), ('ab37cbe9-b212-446b-a7ca-3dc8117281d3', 'd6a41504-bbd0-11ee-8a24-0242ac120003', 1150, '2024-01-27 17:13:48'), ('d10d6ec7-92ed-4387-a689-422fe1559a30', 'd6a41403-bbd0-11ee-8a24-0242ac120003', 800, '2024-01-27 17:13:55'), ('5af44c65-fc2f-4cfc-9b8b-5a908bda18b1', 'd6a41504-bbd0-11ee-8a24-0242ac120003', 1200, '2024-01-27 17:13:55')")
	if err != nil {
		panic(err)
	}

	balanceDb := database.NewBalanceDB(db)

	createBalanceUseCase := create_balance.NewCreateBalanceUseCase(balanceDb)
	findAccountUseCase := find_account.NewFindAccountUseCase(balanceDb)

	go func() {
		webserver := webserver.NewWebServer(":3003")
		accountHandler := web.NewWebBalanceHandler(*findAccountUseCase)
		webserver.AddHandler("/accounts/{account_id}", accountHandler.FindAccount)
		fmt.Println("Server running at port 3003")
		webserver.Start()
	}()

	configMap := &kafka.ConfigMap{
		"bootstrap.servers": "kafka:29092",
		"client.id":         "balances",
		"group.id":          "balances",
		"auto.offset.reset": "earliest",
	}
	c, err := kafka.NewConsumer(configMap)
	if err != nil {
		fmt.Println("error consumer", err.Error())
	}
	topics := []string{"balances"}
	c.SubscribeTopics(topics, nil)

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			fmt.Println(string(msg.Value), msg.TopicPartition)

			var event BalanceEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				fmt.Println("Error JSON:", err)
				continue
			}

			input := create_balance.CreateBalanceInputDTO{
				AccountID: event.Payload.AccountIDFrom,
				Balance:   event.Payload.BalanceAccountFrom,
			}
			createBalanceUseCase.Execute(input)

			input = create_balance.CreateBalanceInputDTO{
				AccountID: event.Payload.AccountIDTo,
				Balance:   event.Payload.BalanceAccountTo,
			}
			createBalanceUseCase.Execute(input)
		}
		c.CommitMessage(msg)
	}
}

package main

import (
	"Homework-1/internal/service"
	"Homework-1/internal/storage"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	createOrderID := createCmd.Int("id", 0, "id принимаемого заказа")
	createCustomerID := createCmd.Int("customer_id", 0, "id получателя заказа")
	createExpireDateStr := createCmd.String("expire_date", "", "срок хранения заказа")

	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	deleteOrderID := deleteCmd.Int("id", 0, "id удаляемого заказа")

	finishCmd := flag.NewFlagSet("finish", flag.ExitOnError)
	finishOrderIDsStr := finishCmd.String("ids", "", "список id выдаваемых заказов в формате [1,2,3]")

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	listUserID := listCmd.Int("customer_id", 0, "id клиента")
	listLimit := listCmd.Int("limit", 0, "ограничение на количество последних выводимых заказов (необязательно)")
	listOnlyNotFinished := listCmd.Bool("only_not_finished", false, "выводить только те заказы, что находятся в ПВЗ (необязательно)")

	returnCmd := flag.NewFlagSet("return", flag.ExitOnError)
	returnUserID := returnCmd.Int("customer_id", 0, "id клиента, возвращающего заказ")
	returnOrderID := returnCmd.Int("id", 0, "id возвращаемого заказа")

	returnsCmd := flag.NewFlagSet("returns", flag.ExitOnError)
	returnsResultsPerPage := returnsCmd.Int("results_per_page", 5, "количество возвратов, отображаемых на одной странице (необязательно)")

	if len(os.Args) == 1 {
		log.Println("необходимо указать хотя бы одну команду")
		return
	}

	arguments := os.Args[1:]
	command := arguments[0]

	stor, err := storage.New()
	if err != nil {
		log.Println("не удалось подключиться к хранилищу")
		return
	}
	serv := service.New(&stor)
	switch command {

	case "create":
		err = createCmd.Parse(os.Args[2:])
		if err != nil {
			log.Println(err)
			return
		}
		err = serv.Create(*createOrderID, *createCustomerID, *createExpireDateStr)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println("заказ успешно принят и добавлен в базу")

	case "delete":
		err = deleteCmd.Parse(os.Args[2:])
		if err != nil {
			log.Println(err)
			return
		}
		err = serv.Delete(*deleteOrderID)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println("заказ был успешно возвращён курьеру")

	case "list":
		err = listCmd.Parse(os.Args[2:])
		list, err := serv.List(*listUserID, *listLimit, *listOnlyNotFinished)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println("Список заказов для клиента с id=" + strconv.Itoa(*listUserID) + ":")
		for _, order := range list {
			s := fmt.Sprintf("id заказа: %d, ", order.ID)
			if !order.IsFinished {
				s += fmt.Sprintf("на складе, срок хранения до: %d-%d-%d", order.ExpireDate.Year(), order.ExpireDate.Month(), order.ExpireDate.Day())
			} else {
				s += fmt.Sprintf("выдан, дата выдачи: %d-%d-%d", order.DateFinished.Year(), order.DateFinished.Month(), order.DateFinished.Day())
			}
			if order.IsReturnedByClient {
				s += ", клиент оформил возврат"
			}
			fmt.Println(s)
		}

	case "return":
		err = returnCmd.Parse(os.Args[2:])
		if err != nil {
			log.Println(err)
			return
		}
		err = serv.Return(*returnOrderID, *returnUserID)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println("возврат от клиента был успешно принят")

	case "returns":
		err = returnsCmd.Parse(os.Args[2:])
		if err != nil {
			log.Println(err)
			return
		}
		paginatedReturns, err := serv.Returns(*returnsResultsPerPage)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println("Список возвратов:")
		fmt.Printf(paginatedReturns)

	case "finish":
		err = finishCmd.Parse(os.Args[2:])
		if err != nil {
			log.Println(err)
			return
		}
		err = serv.Finish(*finishOrderIDsStr)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println("заказы успешно выданы")

	case "help":
		fmt.Println("Команда 'create' используется при приёме заказа от курьера и добавляет его в базу данных")
		createCmd.SetOutput(os.Stdout)
		createCmd.PrintDefaults()
		fmt.Println("")

		fmt.Println("Команда 'delete' используется при возвращении заказа курьеру и удаляет его из базы данных")
		deleteCmd.SetOutput(os.Stdout)
		deleteCmd.PrintDefaults()
		fmt.Println("")

		fmt.Println("Команда 'finish' используется при выдаче заказа клиенту и помечает заказ как завершённый")
		finishCmd.SetOutput(os.Stdout)
		finishCmd.PrintDefaults()
		fmt.Println("")

		fmt.Println("Команда 'list' используется для получения списка заказов определённого клиента")
		listCmd.SetOutput(os.Stdout)
		listCmd.PrintDefaults()
		fmt.Println("")

		fmt.Println("Команда 'return' используется при принятии возврата от клиента и помечает заказ как возвращённый")
		returnCmd.SetOutput(os.Stdout)
		returnCmd.PrintDefaults()
		fmt.Println("")

		fmt.Println("Команда 'returns' пагинированно выводит список всех возвратов принятых ПВЗ")
		returnsCmd.SetOutput(os.Stdout)
		returnsCmd.PrintDefaults()
		fmt.Println("")

	default:
		fmt.Println("неизвестная команда")
	}

}

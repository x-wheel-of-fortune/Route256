package cli

import (
	"Homework-1/internal/service"
	"flag"
	"fmt"
	"os"
	"strconv"
)

type CLI struct {
	service *service.Service
}

func New(s *service.Service) *CLI {
	return &CLI{service: s}
}

func (c *CLI) HandleCreate(args []string) error {
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	createOrderID := createCmd.Int("id", 0, "id принимаемого заказа")
	createCustomerID := createCmd.Int("customer_id", 0, "id получателя заказа")
	createExpireDateStr := createCmd.String("expire_date", "", "срок хранения заказа")
	err := createCmd.Parse(args)
	if err != nil {
		return err
	}
	err = c.service.Create(*createOrderID, *createCustomerID, *createExpireDateStr)
	if err != nil {
		return err
	}
	fmt.Println("заказ успешно принят и добавлен в базу")

	return nil
}

func (c *CLI) HandleDelete(args []string) error {
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	deleteOrderID := deleteCmd.Int("id", 0, "id удаляемого заказа")

	err := deleteCmd.Parse(args)
	if err != nil {
		return err
	}
	err = c.service.Delete(*deleteOrderID)
	if err != nil {
		return err
	}
	fmt.Println("заказ был успешно возвращён курьеру")
	return nil
}

func (c *CLI) HandleFinish(args []string) error {
	finishCmd := flag.NewFlagSet("finish", flag.ExitOnError)
	finishOrderIDsStr := finishCmd.String("ids", "", "список id выдаваемых заказов в формате [1,2,3]")

	err := finishCmd.Parse(args)
	if err != nil {
		return err
	}
	err = c.service.Finish(*finishOrderIDsStr)
	if err != nil {
		return err
	}
	fmt.Println("заказы успешно выданы")
	return nil
}

func (c *CLI) HandleList(args []string) error {
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	listUserID := listCmd.Int("customer_id", 0, "id клиента")
	listLimit := listCmd.Int("limit", 0, "ограничение на количество последних выводимых заказов (необязательно)")
	listOnlyNotFinished := listCmd.Bool("only_not_finished", false, "выводить только те заказы, что находятся в ПВЗ (необязательно)")

	err := listCmd.Parse(args)
	list, err := c.service.List(*listUserID, *listLimit, *listOnlyNotFinished)
	if err != nil {
		return err
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
	return nil
}

func (c *CLI) HandleReturn(args []string) error {
	returnCmd := flag.NewFlagSet("return", flag.ExitOnError)
	returnUserID := returnCmd.Int("customer_id", 0, "id клиента, возвращающего заказ")
	returnOrderID := returnCmd.Int("id", 0, "id возвращаемого заказа")

	err := returnCmd.Parse(args)
	if err != nil {
		return err
	}
	err = c.service.Return(*returnOrderID, *returnUserID)
	if err != nil {
		return err
	}
	fmt.Println("возврат от клиента был успешно принят")
	return nil
}

func (c *CLI) HandleReturns(args []string) error {
	returnsCmd := flag.NewFlagSet("returns", flag.ExitOnError)
	returnsResultsPerPage := returnsCmd.Int("results_per_page", 5, "количество возвратов, отображаемых на одной странице (необязательно)")

	err := returnsCmd.Parse(args)
	if err != nil {
		return err
	}
	paginatedReturns, err := c.service.Returns(*returnsResultsPerPage)
	if err != nil {
		return err
	}
	fmt.Println("Список возвратов:")
	fmt.Printf(paginatedReturns)
	return nil
}

func (c *CLI) HandleHelp() error {
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	createCmd.Int("id", 0, "id принимаемого заказа")
	createCmd.Int("customer_id", 0, "id получателя заказа")
	createCmd.String("expire_date", "", "срок хранения заказа")

	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	deleteCmd.Int("id", 0, "id удаляемого заказа")

	finishCmd := flag.NewFlagSet("finish", flag.ExitOnError)
	finishCmd.String("ids", "", "список id выдаваемых заказов в формате [1,2,3]")

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	listCmd.Int("customer_id", 0, "id клиента")
	listCmd.Int("limit", 0, "ограничение на количество последних выводимых заказов (необязательно)")
	listCmd.Bool("only_not_finished", false, "выводить только те заказы, что находятся в ПВЗ (необязательно)")

	returnCmd := flag.NewFlagSet("return", flag.ExitOnError)
	returnCmd.Int("customer_id", 0, "id клиента, возвращающего заказ")
	returnCmd.Int("id", 0, "id возвращаемого заказа")

	returnsCmd := flag.NewFlagSet("returns", flag.ExitOnError)
	returnsCmd.Int("results_per_page", 5, "количество возвратов, отображаемых на одной странице (необязательно)")

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
	return nil
}

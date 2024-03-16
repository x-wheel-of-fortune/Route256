package main

import (
	"Homework-1/internal/cli"
	"Homework-1/internal/service"
	"Homework-1/internal/storage"
	"fmt"
	"log"
	"os"
)

func main() {

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
	cli := cli.New(&serv)

	switch command {

	case "create":
		err = cli.HandleCreate(os.Args[2:])
	case "delete":
		err = cli.HandleDelete(os.Args[2:])
	case "list":
		err = cli.HandleList(os.Args[2:])
	case "return":
		err = cli.HandleReturn(os.Args[2:])
	case "returns":
		err = cli.HandleReturns(os.Args[2:])
	case "finish":
		err = cli.HandleFinish(os.Args[2:])
	case "help":
		err = cli.HandleHelp()
	case "interactive":
		err = cli.HandleInteractive()

	default:
		fmt.Println("неизвестная команда")
	}
	if err != nil {
		log.Println(err)
	}

}

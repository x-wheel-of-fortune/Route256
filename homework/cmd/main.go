package main

import (
	"fmt"
	"log"
	"os"

	"homework/internal/cli"
	"homework/internal/model"
	"homework/internal/service"
	"homework/internal/service/packaging"
	"homework/internal/storage"
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

	packagingVariants := map[model.PackagingType]packaging.PackagingVariant{
		packaging.BagPackagingType:  packaging.BagPackaging{},
		packaging.BoxPackagingType:  packaging.BoxPackaging{},
		packaging.FilmPackagingType: packaging.FilmPackaging{},
	}

	serv := service.New(&stor, packagingVariants)
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

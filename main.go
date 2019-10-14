package main

import (
	"bufio"
	"os"
	"strings"

	"github.com/4tyTwo/parking/client"
	"github.com/4tyTwo/parking/utils"
	"github.com/fatih/color"
)

func getGeneralHelp() string {
	return `Available commands:
	help  - print this help message
	exit  - exit the client
	free  - get number of free places on the parking
	place - place car to the parking lot
	take {code} - take car from the parking lot 
	`
}

func getCommandHelp(command string) string {
	switch command {
	case "free":
		return "get the number of free spaces"
	case "place":
		return "place inserted car in the parking lot, returns the code, which must be provided to take the cat back"
	case "take":
		return "syntax: take {code}. Take car, specified by code"
	default:
		return "unknown command"
	}
}

func validateCode(code string) bool {
	if len(code) != 5 && strings.ToUpper(code) != code {
		return false
	}
	return true
}

func main() {
	if len(os.Args) < 2 {
		color.Red("Unspecified host")
		color.White("Usage: parking-client {hostname}")
		os.Exit(1)
	}
	host := os.Args[1]
	reader := bufio.NewReader(os.Stdin)
	color.Green("Parking system client started")
	color.Green("Type help to see list of available commands")
	color.Green("Type help {command} to see info about command")
	color.Green("Enter command: ")
	for true {
		text, err := reader.ReadString('\n')
		utils.CheckErr(err)
		text = strings.TrimRight(text, "\n")
		tokens := strings.Split(text, " ")
		if tokens[0] == "exit" {
			os.Exit(0)
		}
		if tokens[0] == "help" {
			if len(tokens) == 1 {
				color.Yellow(getGeneralHelp())
			} else {
				color.Yellow(getCommandHelp(tokens[1]))
			}
		}
		if tokens[0] == "free" {
			resp, err := client.GetFreePlaces(host)
			if err != nil {
				color.Red(err.Error())
			} else {
				if resp.FreePlaces == 0 {
					color.Cyan("No free places left")
				} else {
					color.Cyan("There are %v free places", resp.FreePlaces)
				}
			}
		}
		if tokens[0] == "place" {
			resp, err := client.PlaceCar(host)
			if err != nil {
				color.Red(err.Error())
			} else {
				color.Cyan("Car places, your code is: %v", resp.PlaceCode)
			}
		}
		if tokens[0] == "take" {
			if len(tokens) == 1 {
				color.Red("No code provided, try again")
			} else {
				if validateCode(tokens[1]) {
					err = client.GetCar(host, tokens[1])
					if err != nil {
						color.Red(err.Error())
					} else {
						color.Cyan("Car is getting taken")
					}
				} else {
					color.Red("Code fromat mismatch")
				}
			}
		}
	}

	// fmt.Println(text)
	// resp, err := client.GetFreePlaces("http://127.0.0.1:4242")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("Free Places: %v\n", resp.FreePlaces)
}

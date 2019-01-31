package main

import (
	"flag"
	"fmt"
	"github.com/docker-compose/executor"
	"github.com/docker-compose/interfaces"
	"os"
	"time"
)

// Global variables
var executorClient interface{}
var err error

var usageStr = `
Usage: ./main service operation		Start app with requested {service} and {operation}
       -h							Show this message
`

//var executorClient interface{}

const (
	START           = "start"
	STOP            = "stop"
	RESTART         = "restart"
	ApplicationName = "docker-compose-executor"
	AppOpenMsg      = "This is the docker-compose-executor application!"
)

// usage will print out the flag options for the app.
func HelpCallback() {
	msg := fmt.Sprintf(usageStr, os.Args[0])
	fmt.Printf("%s\n", msg)
	os.Exit(0)
}

func main() {

	start := time.Now()

	flag.Usage = HelpCallback
	flag.Parse()

	fmt.Println(fmt.Sprintf("Starting the %s application...", ApplicationName))
	fmt.Println(AppOpenMsg)

	// Time it took to start service
	fmt.Println("Application started in: " + time.Since(start).String())

	executorClient, err = newExecutorClient()

	var service = ""
	var operation = ""

	if len(os.Args) > 2 {
		service = os.Args[1]
		operation = os.Args[2]

		switch operation {
		case START:
			if starter, ok := executorClient.(interfaces.ServiceStarter); ok {
				err := starter.Start(service)
				if err != nil {
					fmt.Println("error starting service: ", service)
				}
			} else {
				fmt.Println("success in starting service: ", service)
			}
			break

		case STOP:
			if stopper, ok := executorClient.(interfaces.ServiceStopper); ok {
				err := stopper.Stop(service)
				if err != nil {
					fmt.Println("error stopping service: ", service)
				}
			} else {
				fmt.Println("success in stopping service: ", service)
			}
			break

		case RESTART:
			if restarter, ok := executorClient.(interfaces.ServiceRestarter); ok {
				err := restarter.Restart(service)
				if err != nil {
					fmt.Println("error restarting service: ", service)
				}
			} else {
				fmt.Println("success in restarting service: ", service)
			}
			break

		default:
			fmt.Println("unknown operation was requested: ", operation)
			break

		}
	}
}

func newExecutorClient() (interface{}, error) {

	return &executor.ExecuteDockerCompose{}, nil
}

/*type ExecuteApp struct {
}
*/

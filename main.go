package docker_compose

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Global variables
var executorClient interface{}
var err error

var usageStr = `
Usage: ./main service operation		Start app with requested {service} and {operation}
       -h							Show this message
`

const (
	START           = "start"
	STOP            = "stop"
	RESTART         = "restart"
	ApplicationName = "docker-compose-executor"
	AppOpenMsg      = "This is the docker-compose-executor application!"
)

// TODO: Externalize these by putting in a properties file.
const (
	ConfigSeedServiceKey            = "edgex-config-seed"
	CoreCommandServiceKey           = "edgex-core-command"
	CoreDataServiceKey              = "edgex-core-data"
	CoreMetaDataServiceKey          = "edgex-core-metadata"
	ExportClientServiceKey          = "edgex-export-client"
	ExportDistroServiceKey          = "edgex-export-distro"
	SupportLoggingServiceKey        = "edgex-support-logging"
	SupportNotificationsServiceKey  = "edgex-support-notifications"
	SystemManagementAgentServiceKey = "edgex-sys-mgmt-agent"
	SupportSchedulerServiceKey      = "edgex-support-scheduler"
)

// TODO: Externalize these by putting in a properties file.
var services = map[string]string{
	SupportNotificationsServiceKey: "Notifications",
	CoreCommandServiceKey:          "Command",
	CoreDataServiceKey:             "CoreData",
	CoreMetaDataServiceKey:         "Metadata",
	ExportClientServiceKey:         "Export",
	ExportDistroServiceKey:         "Distro",
	SupportLoggingServiceKey:       "Logging",
	SupportSchedulerServiceKey:     "Scheduler",
}

type ServiceStarter interface {
	Start(service string) error
}

type ServiceStopper interface {
	Stop(service string, ) error
}

type ServiceRestarter interface {
	Restart(service string) error
}

func newExecutorClient() (interface{}, error) {

	return &ExecuteDockerCompose{}, nil
}

type ExecuteDockerCompose struct {
}

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
			if starter, ok := executorClient.(ServiceStarter); ok {
				err := starter.Start(service)
				if err != nil {
					fmt.Println("error starting service: ", service)
				}
			} else {
				fmt.Println("success in starting service: ", service)
			}
			break

		case STOP:
			if stopper, ok := executorClient.(ServiceStopper); ok {
				err := stopper.Stop(service)
				if err != nil {
					fmt.Println("error stopping service: ", service)
				}
			} else {
				fmt.Println("success in stopping service: ", service)
			}
			break

		case RESTART:
			if restarter, ok := executorClient.(ServiceRestarter); ok {
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

func (ec *ExecuteDockerCompose) Start(service string) error {
	err := ExecuteDockerCommands(service, "start")
	return err
}

func (ec *ExecuteDockerCompose) Stop(service string) error {
	err := ExecuteDockerCommands(service, "stop")
	return err
}

func (ec *ExecuteDockerCompose) Restart(service string) error {
	err := ExecuteDockerCommands(service, "restart")
	return err
}

func findDockerContainerStatus(service string, status string) bool {

	var (
		cmdOut []byte
		err    error
	)
	cmdName := "docker"
	cmdArgs := []string{"ps"}
	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		fmt.Println("error running the docker command", "error message", err.Error())
		os.Exit(1)
	}

	dockerOutput := string(cmdOut)

	// Find whether the container to start has started.
	for _, line := range strings.Split(strings.TrimSuffix(dockerOutput, "\n"), "\n") {
		if strings.Contains(line, service) {

			if status == "Up" {
				if strings.Contains(line, "Up") {
					fmt.Println("container started", "service name", service, "details", line)
					return true
				} else {
					fmt.Println("container NOT started", "service name", service)
					return false
				}
			} else if status == "Exited" {
				if strings.Contains(line, "Exited") {
					fmt.Println("container stopped", "service name", service, "details", line)
					return true
				} else {
					fmt.Println("container NOT stopped", "service name", service)
					return false
				}
			}
		}
	}
	return false
}

func ExecuteDockerCommands(service string, operation string) error {
	_, knownService := services[service]

	if knownService {
		runDockerCommands(service, services[service], operation)

		return nil
	} else {
		newError := fmt.Errorf("unknown service: %v", service)
		fmt.Println(newError.Error())

		return newError
	}
}

func runDockerCommands(service string, dockerService string, operation string) {

	var (
		err    error
		cmdDir string
	)

	cmdName := "docker"

	cmdArgs := []string{operation, dockerService}
	cmd := exec.Command(cmdName, cmdArgs...)
	cmd.Dir = cmdDir

	/*
	A call to exec.CombinedOutput will return a nil or 1 (via err), along with the standard
	output, back to the SMA.
	A return of string 0 indicates that the execution completed its task and exited “normally” or without issue.
	A return of string 1 indicates that the execution did not complete “normally” and the caller should check
	the standard error for more information. The Executor should always return some information string indicting
	why the non-normal return when 1 is returned on the standard out.
	*/
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("docker command failed", "error message", err.Error())
		fmt.Println("associated output", "error message", string(out))
	} else if operation == "start" {
		if !findDockerContainerStatus(service, "Up") {
			fmt.Println("docker start operation failed", "service name", service)
		}
	} else if operation == "stop" {
		if !findDockerContainerStatus(service, "Exited") {
			fmt.Println("docker stop operation failed", "service name", service)
		}
	} else if operation == "restart" {
		if !findDockerContainerStatus(service, "Up") {
			fmt.Println("docker restart operation failed", "service name", service)
		}
	}
}

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types"

	"github.com/docker/docker/client"
	"github.com/thebsdbox/diver/pkg/ucp"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

type containerStatistics struct {
	MemoryStats memoryStatistics `json:"memory_stats"`
	Name        string           `json:"name"`
	ID          string           `json:"id"`
}

type memoryStatistics struct {
	Usage     int `json:"usage"`
	MaxUsages int `json:"max_usage"`
	Stats     struct {
	} `json:"stats"`
	Limit int `json:"limit"`
}

func main() {

	cmd := &cobra.Command{
		Use:   "diver",
		Short: "This tool uses the native APIs to \"dive\" into Docker EE",
	}

	client := ucp.Client{}

	cmd.Flags().StringVar(&client.Username, "username", os.Getenv("DIVER_USERNAME"), "Username that has permissions to authenticate to Docker EE")
	cmd.Flags().StringVar(&client.Password, "password", os.Getenv("DIVER_PASSWORD"), "Password allowing a user to authenticate to Docker EE")
	cmd.Flags().StringVar(&client.UCPURL, "url", os.Getenv("DIVER_URL"), "URL for Docker EE, e.g. https://10.0.0.1")
	ignoreCert := strings.ToLower(os.Getenv("DIVER_INSECURE")) == "true"

	cmd.Flags().BoolVar(&client.IgnoreCert, "ignorecert", ignoreCert, "Ignore x509 certificate")

	var logLevel = 5
	cmd.Flags().IntVar(&logLevel, "logLevel", 4, "Set the logging level [0=panic, 3=warning, 5=debug]")

	cmd.Execute()
	log.SetLevel(log.Level(logLevel))

	err := client.Connect()
	if err != nil {
		log.Errorf("%v", err)
	} else {
		// err = client.ListNetworks()
		// if err != nil {
		// 	log.Errorf("%v\n", err)
		// }
		// err = client.ListContainerJSON()
		// if err != nil {
		// 	log.Errorf("%v\n", err)
		// }
		// err = client.GetClientBundle()
		// if err != nil {
		// 	log.Errorf("%v\n", err)
		// }

		user := ucp.NewUser("dan finneran", "dan", "password", true, true, false)
		err = client.AddAccount(user)
		if err != nil {
			log.Errorf("%v\n", err)
		}
	}
}

func dockerSocketConnect() {
	cli, err := client.NewClient("unix:///var/run/docker.sock", "", nil, nil)
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for index := range containers {
		go func(currentContainer *types.Container) {
			statsPayload, _ := cli.ContainerStats(context.Background(), currentContainer.ID[:10], false)
			buf := new(bytes.Buffer)
			buf.ReadFrom(statsPayload.Body)
			statistics := containerStatistics{}
			err := json.Unmarshal(buf.Bytes(), &statistics)
			if err != nil {
				fmt.Printf("%v\n", err)
				os.Exit(-1)
			}
			//A lot of casting is required to ensure that the division of integers yeilds a float value
			percentage := ((float64(statistics.MemoryStats.Usage) / float64(statistics.MemoryStats.Limit)) * 100)
			if percentage > 90 {
				fmt.Printf("\033[35m%0.2f%%\033[m  %s %s\n", percentage, currentContainer.ID[:10], currentContainer.Image)
			} else if percentage > 75 {
				fmt.Printf("\033[35m%0.2f%%\033[m  %s %s\n", percentage, currentContainer.ID[:10], currentContainer.Image)
			} else {
				fmt.Printf("\033[32m%0.2f%%\033[m  %s %s\n", percentage, currentContainer.ID[:10], currentContainer.Image)
			}
		}(&containers[index])
		time.Sleep(time.Minute)
		//fmt.Printf("%d\t %d \t  %0.2f%% \n", statistics.MemoryStats.Limit, statistics.MemoryStats.Usage, percentage)
	}
}

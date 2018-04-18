package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/docker/docker/api/types"

	"github.com/docker/docker/client"
	"github.com/thebsdbox/diver/cmd"
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
	cmd.Execute()
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

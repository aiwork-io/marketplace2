package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"aiwork.io/aiworkclient/helpers"
	"aiwork.io/aiworkclient/internal"
	"github.com/robfig/cron/v3"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
			os.Exit(1)
		}
	}()

	provider := internal.NewConfigProvider(helpers.GetRootDir(), ".", "./secrets")
	configs := internal.NewConfigs(provider)

	job := newJob(configs)
	// start first round
	job()

	cronjob := cron.New()
	cronjob.AddFunc(configs.JobSchedule, job)
	go cronjob.Start()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}

func newJob(configs *internal.Configs) func() {
	tc := internal.NewTaskClient(configs)
	ec := internal.NewEngineClient(configs)

	return func() {
		task, err := tc.Get()
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("[%s] starting", task.Id)

		ok, err := helpers.Lock(task.Id)
		if err != nil {
			panic(err)
		}

		assetIds := []string{}
		// processing asset
		for _, asset := range task.Assets {
			if asset.Results != "" {
				log.Printf("[%s] completed %s", task.Id, asset.Id)
				assetIds = append(assetIds, asset.Id)
				continue
			}
			if !ok {
				log.Printf("[%s] locked - ignore %s", task.Id, asset.Id)
				continue
			}

			if err := ec.Push(&asset); err != nil {
				log.Println(err)
			}

			log.Printf("[%s] processing %s", task.Id, asset.Id)
		}

		completed := len(assetIds) == len(task.Assets)
		// completed -> no more task
		if completed {
			completedTask, err := tc.Complete(task, assetIds)
			if err != nil {
				log.Println(err)
			}
			helpers.Unlock(completedTask.Id)
			log.Printf("[%s] completed with txn %s", completedTask.Id, completedTask.RewardTxn)
		}
	}
}

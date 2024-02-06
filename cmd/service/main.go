package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"khranity/app/cloud"
	"khranity/app/log"
	"khranity/app/lore"
	"khranity/app/shell"
)

func main() {
	log.Start()

	if err := run(); err != nil {
		log.Get().Logger.Fatal(err.Error())
	}
}

func run() error {
	ctx := context.Background()
	lore := lore.Get()

	err := cloud.TestClouds(ctx, lore.Clouds)
	if err != nil {
		return err
	}

	s, err := shell.New(ctx, log.Get())
	if err != nil {
		return err
	}

	err = startJobs(ctx, s)
	if err != nil {
		return err
	}

	done := make(chan struct{})
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return nil
		case <-interrupt:
			log.Warn("interrupt")

			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return nil
		}
	}
}

func startJobs(ctx context.Context, s *shell.OS) error {
	jobs := lore.GetJobs()
	for i, job := range jobs {
		err := s.Shell.Start(ctx, &jobs[i])
		if err != nil {
			log.Get().Warn(fmt.Sprintf("schedule for job %00d skip", i), log.String("err", err.Error()))
			continue
		}
		log.Get().Info(fmt.Sprintf("schedule for job %00d operational", i), log.Object("job", &job))
	}

	return nil
}
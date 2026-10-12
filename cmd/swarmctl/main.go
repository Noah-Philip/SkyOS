package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"

	"github.com/noah-philip/skyos/internal"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: swarmctl start --config <path>")
	}

	switch os.Args[1] {
	case "start":
		startSwarm(os.Args[2:])
	default:
		log.Fatalf("unknown command: %s", os.Args[1])
	}
}

func startSwarm(args []string) {
	flags := flag.NewFlagSet("start", flag.ExitOnError)

	configPath := flags.String(
		"config",
		"",
		"path to swarm configuration file",
	)

	dronedPath := flags.String(
		"droned",
		"bin/droned.exe",
		"path to droned executable",
	)
	flags.Parse(args)

	if *configPath == "" {
		log.Fatal("--config is required")
	}

	absoluteConfigPath, err := filepath.Abs(*configPath)

	if err != nil {
		log.Fatal(err)
	}

	absoluteDronedPath, err := filepath.Abs(*dronedPath)
	if err != nil {
		log.Fatal(err)
	}

	config, err := internal.LoadSwarmConfig(absoluteConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	var processes []*exec.Cmd

	for _, node := range config.Nodes {
		command := exec.Command(
			absoluteDronedPath,
			"--config", absoluteConfigPath,
			"--id", node.ID,
		)

		command.Stdout = os.Stdout
		command.Stderr = os.Stderr

		if err := command.Start(); err != nil {
			stopProcesses(processes)
			log.Fatalf("start %s: %v", node.ID, err)
		}

		processes = append(processes, command)

		log.Printf(
			"started %s (PID %d)",
			node.ID,
			command.Process.Pid,
		)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	log.Print("swarm is running; press Ctrl+C to stop all simulated drones")
	<-ctx.Done()

	log.Print("stopping drones")
	stopProcesses(processes)
}

func stopProcesses(processes []*exec.Cmd) {
	for _, command := range processes {
		if command.Process == nil {
			continue
		}

		_ = command.Process.Kill()
		_ = command.Wait()
	}
}

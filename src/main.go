package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"shape-harvester/src/helpers/browser"
	"shape-harvester/src/helpers/utils"
)

func main() {
	browserType := "Opera"
	logger := utils.NewColorizedLogger(true)

	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		logger.Error("Stop Signal Received, Stopping Shape Harvester...")
		cancel()
	}()

	logger.Verbose(fmt.Sprintf("Starting Shape Harvester Using %s...", browserType))

	if err := browser.OpenBrowser(ctx, strings.ToLower(browserType), logger); err != nil {
		logger.Error(fmt.Sprintf("Shape Harvester Has Received An Error: %v", err))
	}
}

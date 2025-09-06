package browser

import (
	"context"
	"fmt"
	"sync"
	"time"

	browserutils "shape-harvester/src/helpers/browser/utils"
	"shape-harvester/src/helpers/utils"

	"github.com/chromedp/chromedp"
)

func OpenBrowser(ctx context.Context, browserType string, logger *utils.ColorizedLogger) error {
	browserPath, err := browserutils.FindBrowser(browserType)
	if err != nil {
		return fmt.Errorf("failed to find %s browser: %w", browserType, err)
	}

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(browserPath),
		chromedp.Flag("headless", false),
		chromedp.Flag("no-first-run", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-default-apps", true),
		chromedp.Flag("disable-background-networking", true),
		chromedp.Flag("disable-gpu", false),
		chromedp.Flag("window-size", "1024,768"),
		chromedp.Flag("disable-features", "TranslateUI"),
		chromedp.Flag("no-sandbox", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	browserCtx, cancelChrome := chromedp.NewContext(allocCtx)
	defer cancelChrome()

	var wg sync.WaitGroup
	wg.Add(1)

	interceptorDone := make(chan error, 1)

	go func() {
		defer wg.Done()
		logger.Info("Starting Network Interceptor...")
		if err := NetworkInterceptor(browserCtx, logger); err != nil {
			interceptorDone <- fmt.Errorf("failed to setup network interceptor: %w", err)
			return
		}

		logger.Silly("Network Interceptor Has Successfully Initialized...")
		interceptorDone <- nil
	}()

	select {
	case err := <-interceptorDone:
		if err != nil {
			return err
		}
	case <-time.After(120 * time.Second):
		return fmt.Errorf("network interceptor setup timed out")
	}

	time.Sleep(1 * time.Second)

	if err := chromedp.Run(browserCtx,
		chromedp.EmulateViewport(1024, 768),
		chromedp.Navigate("https://www.target.com/"),
		chromedp.WaitVisible("#headerPrimary", chromedp.ByID),
		chromedp.Sleep(8*time.Second),
		chromedp.Navigate("https://www.target.com/p/2024-nba-mosaic-collectible-trading-cards/-/A-94864074"),
		chromedp.WaitVisible("#addToCartButtonOrTextIdFor94864074", chromedp.ByID),
	); err != nil {
		return fmt.Errorf("failed to navigate and wait for selectors: %w", err)
	}

	if err := HarvestHeaders(browserCtx); err != nil {
		return fmt.Errorf("clicking loop failed: %w", err)
	}

	return nil
}

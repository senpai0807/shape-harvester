package browser

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
)

func HarvestHeaders(ctx context.Context) error {
	for {
		if err := chromedp.Run(ctx,
			chromedp.Click("#addToCartButtonOrTextIdFor94864074", chromedp.ByID),
			chromedp.WaitVisible(`button[aria-label="close"]`, chromedp.ByQuery),
			chromedp.Sleep(7500*time.Millisecond),
			chromedp.Click(`button[aria-label="close"]`, chromedp.ByQuery),
			chromedp.Sleep(7500*time.Millisecond),
		); err != nil {
			return fmt.Errorf("clicking loop iteration failed: %w", err)
		}
	}
}

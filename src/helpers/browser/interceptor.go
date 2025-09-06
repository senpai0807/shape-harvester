package browser

import (
	"context"
	"fmt"
	"strings"

	"shape-harvester/src/helpers/requests"
	"shape-harvester/src/helpers/utils"

	"github.com/chromedp/cdproto/fetch"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func NetworkInterceptor(ctx context.Context, logger *utils.ColorizedLogger) error {
	if err := chromedp.Run(ctx, fetch.Enable()); err != nil {
		return fmt.Errorf("failed to enable fetch domain: %w", err)
	}

	if err := chromedp.Run(ctx, network.Enable()); err != nil {
		return fmt.Errorf("failed to enable network domain: %w", err)
	}

	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *fetch.EventRequestPaused:
			go RequestPaused(ctx, ev, logger)
		}
	})

	return nil
}

func RequestPaused(ctx context.Context, ev *fetch.EventRequestPaused, logger *utils.ColorizedLogger) {
	if strings.Contains(ev.Request.URL, "web_checkouts/v1/cart_items") && ev.Request.Method == "POST" {
		logger.Http(fmt.Sprintf("Successfully Intercepted ATC Request: %s - %s", ev.Request.URL, ev.Request.Method))

		go ParseHeaders(ctx, ev.Request.URL, ev.Request.Method, ev.Request.Headers, logger)

		err := chromedp.Run(ctx, fetch.FailRequest(ev.RequestID, network.ErrorReasonFailed))
		if err != nil {
			logger.Error(fmt.Sprintf("Failed To Block Intercepted Request: %v", err))
		}
		return
	}

	err := chromedp.Run(ctx, fetch.ContinueRequest(ev.RequestID))
	if err != nil {
		logger.Error(fmt.Sprintf("Failed To Continue Request %s: %v", ev.Request.URL, err))
	}
}

func ParseHeaders(ctx context.Context, url, method string, headers network.Headers, logger *utils.ColorizedLogger) {
	headerMap := make(map[string]interface{})

	for key, value := range headers {
		headerMap[key] = value
	}

	fullData := map[string]interface{}{
		"url":     url,
		"method":  method,
		"headers": headerMap,
	}

	if err := requests.SendHeaders(ctx, fullData, logger); err != nil {
		logger.Error(fmt.Sprintf("Failed To Send Headers To Server: %v", err)
	}
}

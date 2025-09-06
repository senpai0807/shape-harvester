package requests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"shape-harvester/src/helpers/utils"
)

var (
	Client = &http.Client{
		Timeout: 0,
		Transport: &http.Transport{
			DisableKeepAlives:     false,
			MaxIdleConns:          100,
			IdleConnTimeout:       0,
			TLSHandshakeTimeout:   0,
			ExpectContinueTimeout: 0,
		},
	}
)

func SendHeaders(ctx context.Context, data map[string]any, logger *utils.ColorizedLogger) error {
	logger.Info("Sending Shape Headers To Backend Server...")

	if ctx.Err() != nil {
		return fmt.Errorf("request aborted: context canceled or deadline exceeded")
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:8721/shape", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/json")

	resp, err := Client.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return fmt.Errorf("request failed: context canceled or deadline exceeded")
		}
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			return
		}
	}()

	if resp.ContentLength == 0 {
		return nil
	}

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		logger.Silly("Backend Server Has Successfully Received Shape Headers...")
		return nil
	} else {
		return fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}
}

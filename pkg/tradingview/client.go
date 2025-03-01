package tradingview

import (
	"context"
	"fmt"
	"sync/atomic"

	tradingviewlb "github.com/VictorVictini/tradingview-lib"
	"github.com/bytedance/sonic"
)

type Client interface {
	AddSymbol(symbol string) error
	RemoveSymbol(symbol string) error
}

type client struct {
	ctx       context.Context
	connected *atomic.Bool

	priceCh chan Chart
	errorCh chan error

	symbols map[string]int

	tradingAPI *tradingviewlb.API
}

func NewClient(ctx context.Context, priceCh chan Chart, errCh chan error) *client {
	return &client{
		ctx:     ctx,
		priceCh: priceCh,
		errorCh: errCh,

		symbols:   make(map[string]int),
		connected: &atomic.Bool{},
	}
}

func (c *client) init() error {
	api := &tradingviewlb.API{}

	err := api.OpenConnection(nil)
	if err != nil {
		return fmt.Errorf("connect to tradingview client: %w", err)
	}

	go func() {
		for {
			select {
			case <-c.ctx.Done():
				return
			case data := <-api.Channels.Read:
				var chart Chart

				jsonData, err := sonic.ConfigFastest.Marshal(data)
				if err != nil {
					c.errorCh <- fmt.Errorf("marshal map data: %w", err)

					return
				}

				err = sonic.ConfigFastest.Unmarshal(jsonData, &chart)
				if err != nil {
					c.errorCh <- fmt.Errorf("unmarshal raw chart: %w", err)

					return
				}

				if !chart.Validate() {
					continue
				}

				c.priceCh <- chart
			case err := <-api.Channels.Error:
				c.errorCh <- err
			}
		}
	}()

	c.connected.Store(true)

	c.tradingAPI = api

	return nil
}

func (c *client) AddSymbol(symbol string) error {
	if !c.connected.Load() {
		if err := c.init(); err != nil {
			return fmt.Errorf("init: %w", err)
		}
	}

	if _, ok := c.symbols[symbol]; ok {
		c.symbols[symbol]++

		return nil
	}

	err := c.tradingAPI.AddRealtimeSymbols([]string{symbol})
	if err != nil {
		return fmt.Errorf("add realtime symbols: %w", err)
	}

	c.symbols[symbol]++

	return nil
}

func (c *client) RemoveSymbol(symbol string) error {
	if !c.connected.Load() {
		return nil
	}

	num := c.symbols[symbol]

	if num == 0 {
		return nil
	}

	num -= 1

	if num > minSubscriptionCount {
		c.symbols[symbol] = num

		return nil
	}

	err := c.tradingAPI.RemoveRealtimeSymbols([]string{symbol})
	if err != nil {
		return fmt.Errorf("remove realtime symbols: %w", err)
	}

	c.symbols[symbol] = num

	return nil
}

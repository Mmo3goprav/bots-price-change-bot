package main

import (
	"context"
	"log"

	"github.com/Mmo3goprav/bots-price-change-bot/internal/bot"
	"github.com/Mmo3goprav/bots-price-change-bot/pkg/tradingview"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conf := ReadConfig()

	errCh := make(chan error, 1)
	priceCh := make(chan tradingview.Chart)

	cl := tradingview.NewClient(ctx, priceCh, errCh)

	tgBot := bot.NewBot(conf.Token, cl, priceCh)

	go func() {
		err := tgBot.Run(ctx)
		if err != nil {
			log.Fatal(err)
		}

		errCh <- err
	}()

	select {
	case <-ctx.Done():
		return
	case err := <-errCh:
		log.Fatal(err)
	}
}

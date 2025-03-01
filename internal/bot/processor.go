package bot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Mmo3goprav/bots-price-change-bot/internal/subscription"
	"github.com/Mmo3goprav/bots-price-change-bot/pkg/tradingview"
)

type Bot struct {
	token string
	api   *tgbotapi.BotAPI

	chartsCh <-chan tradingview.Chart

	tradingviewClient tradingview.Client
	subscriptions     map[int64]*subscription.Subscription
}

func NewBot(token string, cl tradingview.Client, chartsCh <-chan tradingview.Chart) *Bot {
	return &Bot{
		token:             token,
		chartsCh:          chartsCh,
		tradingviewClient: cl,
		subscriptions:     make(map[int64]*subscription.Subscription),
	}
}

func (b *Bot) Run(ctx context.Context) error {
	err := b.init()
	if err != nil {
		return fmt.Errorf("init: %w", err)
	}

	updatesConfig := tgbotapi.NewUpdate(0)
	updatesConfig.Timeout = 60

	ch := b.api.GetUpdatesChan(updatesConfig)

	errCh := make(chan error)

	go func() {
		err := b.processCharts(ctx)
		if err != nil {
			errCh <- fmt.Errorf("process charts: %w", err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errCh:
			return err
		case update, ok := <-ch:
			if update.Message == nil {
				continue
			}

			if !ok {
				return fmt.Errorf("channel closed unexpectedly")
			}

			err = b.processCommands(update)
			if err != nil {
				return fmt.Errorf("process commands: %w", err)
			}

			err = b.processMessage(update)
			if err != nil {
				return fmt.Errorf("process message: %w", err)
			}
		}
	}
}

func (b *Bot) processCharts(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case chart, ok := <-b.chartsCh:
			if !ok {
				return nil
			}

			for _, sub := range b.subscriptions {
				if sub.GetState() != subscription.StateOnUpdates || !sub.SatisfyChart(chart.Symbol) {
					continue
				}

				text := fmt.Sprintf(
					chartMsg,
					chart.Symbol,
					chart.CurrentPrice,
					chart.PriceChange,
					chart.PriceChangePercentage,
				)

				err := b.sendMessage(newTextMsg(sub.GetID(), text))
				if err != nil {
					return fmt.Errorf("send message: %w", err)
				}
			}
		}
	}
}

func (b *Bot) init() error {
	bot, err := tgbotapi.NewBotAPI(b.token)
	if err != nil {
		return fmt.Errorf("create bot instance: %w", err)
	}

	b.api = bot

	return nil
}

func (b *Bot) sendMessage(msg tgbotapi.MessageConfig) error {
	_, err := b.api.Send(msg)
	if err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func (b *Bot) processCommands(update tgbotapi.Update) error {
	if update.Message.Command() == "" {
		return nil
	}

	switch update.Message.Command() {
	case commandStart:
		err := b.sendMessage(newChartKeyboard(update.Message.From.ID, initMsg))
		if err != nil {
			return fmt.Errorf("send message: %w", err)
		}
	}

	return nil
}

func (b *Bot) processMessage(update tgbotapi.Update) error {
	if update.Message.Text == "" || update.Message.Command() != "" {
		return nil
	}

	sub, ok := b.subscriptions[update.Message.From.ID]
	if !ok {
		sub = subscription.NewSubscription(update.Message.From.ID)
	}

	var msg tgbotapi.MessageConfig

	switch update.Message.Text {
	case callbackAddChart:
		sub.SetState(subscription.StateAwaitAddChart)

		b.subscriptions[update.Message.From.ID] = sub

		msg = newTextMsg(update.Message.From.ID, addChartMsg)
	case callbackRemoveChart:
		sub.SetState(subscription.StateAwaitRemoveChart)

		b.subscriptions[update.Message.From.ID] = sub

		msg = newTextMsg(update.Message.From.ID, removeChartMsg)
	default:
		awaitMsg, err := b.processAwaitMessage(update)
		if err != nil {
			return fmt.Errorf("process await message: %w", err)
		}

		msg = awaitMsg
	}

	err := b.sendMessage(msg)
	if err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func (b *Bot) processAwaitMessage(update tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	sub, ok := b.subscriptions[update.Message.From.ID]
	if !ok {
		sub = subscription.NewSubscription(update.Message.From.ID)
	}

	state := sub.GetState()

	var msg tgbotapi.MessageConfig

	switch state {
	case subscription.StateAwaitAddChart:
		sub.SetState(subscription.StateOnUpdates)

		if sub.SatisfyChart(update.Message.Text) {
			b.subscriptions[update.Message.From.ID] = sub

			return newTextMsg(update.Message.From.ID, fmt.Sprintf(chartAlreadyAdded, update.Message.Text)), nil
		}

		sub.AddChart(update.Message.Text)

		err := b.tradingviewClient.AddSymbol(update.Message.Text)
		if err != nil {
			return msg, fmt.Errorf("add symbol: %w", err)
		}

		msg = newTextMsg(update.Message.From.ID, fmt.Sprintf(chartAdded, update.Message.Text))

		b.subscriptions[update.Message.From.ID] = sub
	case subscription.StateAwaitRemoveChart:
		sub.SetState(subscription.StateOnUpdates)

		if !sub.SatisfyChart(update.Message.Text) {
			b.subscriptions[update.Message.From.ID] = sub

			return newTextMsg(update.Message.From.ID, fmt.Sprintf(chartAlreadyRemoved, update.Message.Text)), nil
		}

		sub.RemoveChart(update.Message.Text)

		err := b.tradingviewClient.RemoveSymbol(update.Message.Text)
		if err != nil {
			return msg, fmt.Errorf("remove symbol: %w", err)
		}

		msg = newTextMsg(update.Message.From.ID, fmt.Sprintf(chartRemoved, update.Message.Text))

		b.subscriptions[update.Message.From.ID] = sub
	default:
		msg = newTextMsg(update.Message.From.ID, helpMsg)
	}

	return msg, nil
}

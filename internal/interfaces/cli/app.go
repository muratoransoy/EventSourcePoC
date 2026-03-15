package cli

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strings"

	"eventsourcepoc/internal/application"
	"eventsourcepoc/internal/domain/cart"

	"github.com/google/uuid"
)

type App struct {
	logger *log.Logger
	store  application.EventStore
	carts  *application.CartService
	orders *application.OrderService
}

func NewApp(logger *log.Logger, store application.EventStore, carts *application.CartService, orders *application.OrderService) *App {
	return &App{
		logger: logger,
		store:  store,
		carts:  carts,
		orders: orders,
	}
}

func (a *App) Run(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing command")
	}

	switch args[0] {
	case "create-cart":
		return a.runCreateCart(ctx, args[1:])
	case "add-item":
		return a.runAddItem(ctx, args[1:])
	case "remove-item":
		return a.runRemoveItem(ctx, args[1:])
	case "create-order":
		return a.runCreateOrder(ctx, args[1:])
	case "read-stream":
		return a.runReadStream(ctx, args[1:])
	case "subscribe":
		return a.runSubscribe(ctx, args[1:])
	case "sample-flow":
		return a.runSampleFlow(ctx, args[1:])
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func (a *App) runCreateCart(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("create-cart", flag.ContinueOnError)
	cartIDValue := fs.String("cart-id", uuid.NewString(), "cart aggregate id")
	userIDValue := fs.String("user-id", uuid.NewString(), "customer id")
	currencyValue := fs.String("currency", "TRY", "cart currency: TRY, USD, EUR")
	if err := fs.Parse(args); err != nil {
		return err
	}

	cartID, err := uuid.Parse(*cartIDValue)
	if err != nil {
		return fmt.Errorf("parse cart-id: %w", err)
	}
	userID, err := uuid.Parse(*userIDValue)
	if err != nil {
		return fmt.Errorf("parse user-id: %w", err)
	}
	currency, err := cart.ParseCurrency(*currencyValue)
	if err != nil {
		return err
	}

	if err := a.carts.CreateCart(ctx, cartID, userID, currency); err != nil {
		return err
	}

	a.logger.Printf("cart created: cartId=%s stream=%s", cartID, application.CartStreamName(cartID))
	return nil
}

func (a *App) runAddItem(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("add-item", flag.ContinueOnError)
	cartIDValue := fs.String("cart-id", "", "cart aggregate id")
	productIDValue := fs.String("product-id", uuid.NewString(), "product id")
	productName := fs.String("name", "", "product name")
	quantity := fs.Int("quantity", 1, "quantity to add")
	unitPrice := fs.Float64("price", 0, "unit price")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if strings.TrimSpace(*productName) == "" {
		return fmt.Errorf("name is required")
	}

	cartID, err := uuid.Parse(*cartIDValue)
	if err != nil {
		return fmt.Errorf("parse cart-id: %w", err)
	}
	productID, err := uuid.Parse(*productIDValue)
	if err != nil {
		return fmt.Errorf("parse product-id: %w", err)
	}

	if err := a.carts.AddItem(ctx, cartID, productID, *productName, *quantity, *unitPrice); err != nil {
		return err
	}

	a.logger.Printf("item added: cartId=%s productId=%s quantity=%d", cartID, productID, *quantity)
	return nil
}

func (a *App) runRemoveItem(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("remove-item", flag.ContinueOnError)
	cartIDValue := fs.String("cart-id", "", "cart aggregate id")
	productIDValue := fs.String("product-id", "", "product id")
	productName := fs.String("name", "", "product name for event readability")
	quantity := fs.Int("quantity", 1, "quantity to remove")
	if err := fs.Parse(args); err != nil {
		return err
	}

	cartID, err := uuid.Parse(*cartIDValue)
	if err != nil {
		return fmt.Errorf("parse cart-id: %w", err)
	}
	productID, err := uuid.Parse(*productIDValue)
	if err != nil {
		return fmt.Errorf("parse product-id: %w", err)
	}

	if err := a.carts.RemoveItem(ctx, cartID, productID, *productName, *quantity); err != nil {
		return err
	}

	a.logger.Printf("item removed: cartId=%s productId=%s quantity=%d", cartID, productID, *quantity)
	return nil
}

func (a *App) runCreateOrder(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("create-order", flag.ContinueOnError)
	cartIDValue := fs.String("cart-id", "", "cart aggregate id")
	orderIDValue := fs.String("order-id", uuid.NewString(), "order aggregate id")
	if err := fs.Parse(args); err != nil {
		return err
	}

	cartID, err := uuid.Parse(*cartIDValue)
	if err != nil {
		return fmt.Errorf("parse cart-id: %w", err)
	}
	orderID, err := uuid.Parse(*orderIDValue)
	if err != nil {
		return fmt.Errorf("parse order-id: %w", err)
	}

	event, err := a.orders.CreateOrderFromCart(ctx, orderID, cartID)
	if err != nil {
		return err
	}

	a.logger.Printf("order created: orderId=%s cartId=%s total=%.2f stream=%s", event.OrderID, event.CartID, event.TotalAmount, application.OrderStreamName(event.OrderID))
	return nil
}

func (a *App) runReadStream(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("read-stream", flag.ContinueOnError)
	streamID := fs.String("stream", "", "stream id, for example cart-<uuid>")
	count := fs.Uint64("count", 100, "maximum number of events to read")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if strings.TrimSpace(*streamID) == "" {
		return fmt.Errorf("stream is required")
	}

	records, err := a.store.ReadStream(ctx, *streamID, *count)
	if err != nil {
		return err
	}

	for _, record := range records {
		if err := a.printRecord(record); err != nil {
			return err
		}
	}

	if len(records) == 0 {
		a.logger.Printf("stream %s is empty", *streamID)
	}
	return nil
}

func (a *App) runSubscribe(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("subscribe", flag.ContinueOnError)
	streamID := fs.String("stream", "", "stream id, for example cart-<uuid>")
	liveOnly := fs.Bool("live-only", false, "start from the end of the stream")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if strings.TrimSpace(*streamID) == "" {
		return fmt.Errorf("stream is required")
	}

	a.logger.Printf("subscribing to stream %s", *streamID)
	return a.store.Subscribe(ctx, *streamID, *liveOnly, a.printRecord)
}

func (a *App) runSampleFlow(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("sample-flow", flag.ContinueOnError)
	currencyValue := fs.String("currency", "TRY", "cart currency")
	if err := fs.Parse(args); err != nil {
		return err
	}

	currency, err := cart.ParseCurrency(*currencyValue)
	if err != nil {
		return err
	}

	cartID := uuid.New()
	userID := uuid.New()
	orderID := uuid.New()
	macbookID := uuid.New()
	pistachioID := uuid.New()

	if err := a.carts.CreateCart(ctx, cartID, userID, currency); err != nil {
		return err
	}
	if err := a.carts.AddItem(ctx, cartID, macbookID, "MacBook Pro", 1, 190000); err != nil {
		return err
	}
	if err := a.carts.AddItem(ctx, cartID, pistachioID, "Antep Pistachio", 3, 799.99); err != nil {
		return err
	}
	if err := a.carts.RemoveItem(ctx, cartID, pistachioID, "Antep Pistachio", 1); err != nil {
		return err
	}
	if _, err := a.orders.CreateOrderFromCart(ctx, orderID, cartID); err != nil {
		return err
	}

	a.logger.Printf("sample flow completed")
	a.logger.Printf("cart stream:  %s", application.CartStreamName(cartID))
	a.logger.Printf("order stream: %s", application.OrderStreamName(orderID))
	return nil
}

func (a *App) printRecord(record application.EventRecord) error {
	payload, err := application.DecodePayload(record)
	if err != nil {
		return fmt.Errorf("decode event %s: %w", record.EventType, err)
	}

	formatted, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return fmt.Errorf("format event %s: %w", record.EventType, err)
	}

	a.logger.Printf("event #%d %s", record.EventNumber, record.EventType)
	fmt.Println(string(formatted))
	return nil
}

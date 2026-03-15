package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"eventsourcepoc/internal/application"
	"eventsourcepoc/internal/infrastructure/config"
	"eventsourcepoc/internal/infrastructure/kurrent"
	"eventsourcepoc/internal/interfaces/cli"
)

func main() {
	ctx := context.Background()

	rootFlags := flag.NewFlagSet("eventsource-poc", flag.ExitOnError)
	connectionString := rootFlags.String("connection", config.Load().ConnectionString, "KurrentDB connection string")
	rootFlags.Usage = func() {
		fmt.Fprintf(rootFlags.Output(), "Usage: %s [global options] <command> [command options]\n\n", os.Args[0])
		fmt.Fprintln(rootFlags.Output(), "Global options:")
		rootFlags.PrintDefaults()
		fmt.Fprintln(rootFlags.Output(), "\nCommands:")
		fmt.Fprintln(rootFlags.Output(), "  create-cart    Append a CartCreated event")
		fmt.Fprintln(rootFlags.Output(), "  add-item       Append an ItemAddedToCart event")
		fmt.Fprintln(rootFlags.Output(), "  remove-item    Append an ItemRemovedFromCart event")
		fmt.Fprintln(rootFlags.Output(), "  create-order   Rebuild cart state and append an OrderCreated event")
		fmt.Fprintln(rootFlags.Output(), "  read-stream    Read and decode a stream")
		fmt.Fprintln(rootFlags.Output(), "  subscribe      Subscribe to a stream")
		fmt.Fprintln(rootFlags.Output(), "  sample-flow    Run a complete cart/order walkthrough")
	}

	if err := rootFlags.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
	if len(rootFlags.Args()) == 0 {
		rootFlags.Usage()
		return
	}

	store, err := kurrent.NewStore(*connectionString)
	if err != nil {
		log.Fatalf("connect to KurrentDB: %v", err)
	}
	defer store.Close()

	logger := log.New(os.Stdout, "", log.LstdFlags)
	cartService := application.NewCartService(store)
	orderService := application.NewOrderService(store, cartService)
	app := cli.NewApp(logger, store, cartService, orderService)

	if err := app.Run(ctx, rootFlags.Args()); err != nil {
		logger.Fatalf("command failed: %v", err)
	}
}

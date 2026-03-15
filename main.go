package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/kurrent-io/KurrentDB-Client-Go/kurrentdb"
)

var (
	kurrentConnection = "kurrentdb://admin:changeit@localhost:2113?tls=false"
	db                *kurrentdb.Client
	settings          *kurrentdb.Configuration
)

func main() {
	guidFlag := flag.String("guid", "", "GUID for the cart")
	listenFlag := flag.Bool("listen", false, "Listen to cart events")

	flag.Parse()

	createKurrentConnection()
	defer db.Close()
	log.Printf("Connected to KurrentDB at %s", kurrentConnection)

	if *guidFlag == "" {
		log.Fatalf("GUID Must be defined")
	}

	strUUID := *guidFlag

	inputID, err := uuid.Parse(strUUID)
	if err != nil {
		log.Fatalf("Invalid GUID: %v", err)
	}

	streamID = fmt.Sprintf("shoppingCart-%s", inputID.String)

	if *listenFlag {
		SubscribeToCartEvents(strUUID)
		return
	}
	//CreateOrders()
	//CreateProducts("Cart-427933cf-d8b7-4381-8924-a703ccf00eb6")
	//CreateProducts("Cart-f11282d7-05f8-4282-b14e-a8387abf5632")
	//CreateCarts()
	SetStreamMockCustomerId(strUUID)
	StreamMockCart()
}

func createKurrentConnection() {
	var err error
	settings, err = kurrentdb.ParseConnectionString(kurrentConnection)
	if err != nil {
		panic(err)
	}

	db, err = kurrentdb.NewClient(settings)
	if err != nil {
		panic(err)
	}
}

package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/cretz/tor-dht-poc/go-i2p-dht-poc/i2pdht"
	"github.com/cretz/tor-dht-poc/go-i2p-dht-poc/i2pdht/ipfs"
	"github.com/eyedeekay/sam3"
)

var impl i2pdht.Impl = ipfs.Impl

var command = flag.String("c", "help", "provide, find, rawid")

func main() {
	flag.Parse()
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	switch *command {
	case "provide":
		return provide(flag.Args())
	case "find":
		return find(flag.Args())
	case "rawid":
		return rawid(flag.Args())
	case "help":
		return help(flag.Args())
	default:
		return fmt.Errorf("Invalid command '%v'", *command)
	}
}

func provide(args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("No args accepted for 'provide' currently")
	}
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	// Fire up tor
	samI2P, err := sam3.NewSAM("127.0.0.1:7656")
	if err != nil {
		return fmt.Errorf("Failed starting SAM connection: %v", err)
	}
	defer samI2P.Close()

	// Make multiple DHTs, passing the known set to the other ones for connecting
	log.Printf("Creating %v peers", *participatingPeerCount)
	dhts := make([]i2pdht.DHT, *participatingPeerCount)
	prevPeers := []*i2pdht.PeerInfo{}
	for i := 0; i < len(dhts); i++ {
		// Start DHT
		conf := &i2pdht.DHTConf{
			I2P:            samI2P,
			Verbose:        *debug,
			BootstrapPeers: make([]*i2pdht.PeerInfo, len(prevPeers)),
		}
		copy(conf.BootstrapPeers, prevPeers)
		dht, err := impl.NewDHT(ctx, conf)
		if err != nil {
			return fmt.Errorf("Failed starting DHT: %v", err)
		}
		defer dht.Close()
		dhts[i] = dht
		prevPeers = append(prevPeers, dht.PeerInfo())
		log.Printf("Created peer #%v: %v\n", i+1, dht.PeerInfo())
	}

	// Have a couple provide our key
	log.Printf("Providing key on the first one (%v)\n", dhts[0].PeerInfo())
	if err = dhts[0].Provide(ctx, []byte(*dataID)); err != nil {
		return fmt.Errorf("Failed providing on first: %v", err)
	}
	log.Printf("Providing key on the last one (%v)\n", dhts[len(dhts)-1].PeerInfo())
	if err = dhts[len(dhts)-1].Provide(ctx, []byte(*dataID)); err != nil {
		return fmt.Errorf("Failed providing on last: %v", err)
	}

	// Wait for key press...
	log.Printf("Press enter to quit...\n")
	_, err = fmt.Scanln()
	return err
}

func find(args []string) error {
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	// Get all the peers from the args
	var err error
	dhtConf := &i2pdht.DHTConf{
		ClientOnly:     true,
		Verbose:        *debug,
		BootstrapPeers: make([]*i2pdht.PeerInfo, len(args)),
	}
	for i := 0; i < len(args); i++ {
		if dhtConf.BootstrapPeers[i], err = i2pdht.NewPeerInfo(args[i]); err != nil {
			return fmt.Errorf("Failed parsing arg #%v: %v", i+1, err)
		}
	}

	// Fire up i2p
	if dhtConf.I2P, err = sam3.NewSAM("127.0.0.1:7656"); err != nil {
		return fmt.Errorf("Failed connecting to SAM: %v", err)
	}
	defer dhtConf.I2P.Close()

	// Make a client-only DHT
	log.Printf("Creating DHT and connecting to peers\n")
	dht, err := impl.NewDHT(ctx, dhtConf)
	if err != nil {
		return fmt.Errorf("Failed creating DHT: %v", err)
	}

	// Now find who is providing the id
	providers, err := dht.FindProviders(ctx, []byte(*dataID), 2)
	if err != nil {
		return fmt.Errorf("Failed finding providers: %v", err)
	}
	for _, provider := range providers {
		log.Printf("Found data ID on %v\n", provider)
	}
	return nil
}

func rawid(args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("No args accepted for 'rawid' currently")
	}
	str, err := impl.RawStringDataID([]byte(*dataID))
	if err != nil {
		return err
	}
	fmt.Printf("Raw string ID: %v\n", str)
	return nil
}

func help(args []string) error {
	x := "help"
	if len(flag.Args()) > 0 {
		x = flag.Args()[0]
	}
	switch x {
	case "provide":
		fmt.Println("provide")
	case "find":
		fmt.Println("find:")
	case "rawid":
		fmt.Println("rawid:")
	default:
		fmt.Println("# Help")
		fmt.Println("")
		fmt.Println("  -c: DHT command to invoke")
		fmt.Println("  -d: debug mode")
		fmt.Println("  -p: create p participating peers -n: tunnel name to use")
		fmt.Println("")
		fmt.Println("## Command")
		fmt.Println("")
		fmt.Println("  provide:")
		fmt.Println("  find:")
		fmt.Println("  rawid:")
	}
	return nil
}

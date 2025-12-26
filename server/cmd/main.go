package main

import (
	"log"
	"log/slog"
	"os"
)

func main() {
	cfg := config{
		addr: ":8080",
		db:   dbConfig{},
	}
	api := application{
		config: cfg,
	}
	// logger
	// htto kay andr middleware lga diye jo info get kr rhy hn
	// so hm aik global logger bna rhy hn jo zada sahi sy hmy btay kay us request pr jo bhi  log ho rha hy us ka pattern kia ho
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	if err := api.run(api.mount()); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

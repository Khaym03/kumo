package main

import (
	"fmt"
	"log"

	"github.com/Khaym03/kumo/internal/adapters/config"
	"github.com/Khaym03/kumo/internal/adapters/storage"
	badger "github.com/dgraph-io/badger/v4"
)

func main() {
	conf := config.LoadKumoConfig()

	db, err := storage.NewBadgerDB(conf.StorageDir, conf.AllowBadgerLogger)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false // Don't fetch values, just keys
		it := txn.NewIterator(opts)
		defer it.Close()

		fmt.Println("Keys in the database:")
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			key := item.Key()
			fmt.Printf("Key: %s\n", key)
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

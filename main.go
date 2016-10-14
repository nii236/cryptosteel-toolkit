package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/nii236/go-bip39"
)

var passphrase string
var password string
var truncated bool
var childIndex int

// Example: pizz knoc taxi bris quar tuna much mang okay twel edge brok occu base corn
func init() {
	flag.StringVar(&passphrase, "p", "", "Passphrase")
	flag.BoolVar(&truncated, "t", false, "Truncated passphrase")
	flag.StringVar(&password, "w", "", "Password")
	flag.IntVar(&childIndex, "c", 0, "Child key")
}

func main() {
	flag.Parse()
	passphrase := strings.Split(passphrase, " ")
	if password != "" {
		fmt.Println("Password:", password)
	}
	mnemonic := []string{}
	if truncated {
		for _, truncatedWord := range passphrase {
			untruncatedWord := bip39.TruncatedWordMap[truncatedWord]
			if bip39.TruncatedWordMap[truncatedWord] == "" {
				log.Fatalln("No matching word found for:", truncatedWord)
			}
			mnemonic = append(mnemonic, untruncatedWord)
		}
	} else {
		for _, untruncatedWord := range passphrase {
			mnemonic = append(mnemonic, untruncatedWord)
		}
	}

	seed := bip39.NewSeed(strings.Join(mnemonic, " "), password)
	fmt.Println("mnemonic:", strings.Join(mnemonic, " "))
	fmt.Println("seed:", hex.EncodeToString(seed))

	hdkey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("BIP32 Root Key:", hdkey.String())
	child, _ := hdkey.Child(uint32(childIndex))
	scriptpubkey, err := child.Address(&chaincfg.MainNetParams)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("scriptpubkey:", scriptpubkey)
}

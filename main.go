package main

import (
	"bytes"
	"crypto/sha256"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/btcsuite/btcutil/base58"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	bip39 "github.com/nii236/go-bip39"
)

var mnemonic string
var password string
var truncated bool
var generate bool

// Example: pizz knoc taxi bris quar tuna much mang okay twel edge brok occu base corn
func init() {
	flag.StringVar(&mnemonic, "m", "", "Mnemonic")
	flag.BoolVar(&truncated, "t", false, "Truncated mnemonic")
	flag.StringVar(&password, "p", "", "Password")
	flag.BoolVar(&generate, "gen", false, "Generate a mnemonic")
}

func generateMnemonic() {
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		panic(err)
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		panic(err)
	}
	fmt.Println(mnemonic)
}

func prepareSeed(mnemonic, password string, truncated bool) *hdkeychain.ExtendedKey {
	seed := bip39.NewSeed(mnemonic, password)
	fmt.Println("mnemonic:", mnemonic)
	hdkey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		log.Fatalln(err)
	}
	child, _ := hdkey.Child(0x80000000 + 44)
	child, _ = child.Child(0x80000000)
	child, _ = child.Child(0x80000000)
	child, _ = child.Child(0)
	child, _ = child.Child(0)

	return child
}

func mnemonicToScriptPubKey(child *hdkeychain.ExtendedKey) {
	scriptpubkey, err := child.Address(&chaincfg.MainNetParams)
	fmt.Println("Extended Key:", child.String())
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("scriptPubKey:", scriptpubkey)
}

func mnemonicToWIF(child *hdkeychain.ExtendedKey) {
	privKey, err := child.ECPrivKey()

	if err != nil {
		log.Fatalln(err)
	}
	var finalBytes bytes.Buffer
	finalBytes.WriteByte(chaincfg.MainNetParams.PrivateKeyID)
	finalBytes.Write(privKey.Serialize())
	finalBytes.WriteByte(0x01)

	round1 := sha256.New()
	round1.Write(finalBytes.Bytes())
	result1 := round1.Sum(nil)

	round2 := sha256.New()
	round2.Write(result1)
	result2 := round2.Sum(nil)

	finalBytes.Write(result2[0:4])
	wif := base58.Encode(finalBytes.Bytes())
	fmt.Println("WIF:", wif)
}

func main() {
	flag.Parse()

	if generate {
		generateMnemonic()
		return
	}

	mnemonicArray := strings.Split(mnemonic, " ")
	mnemonicFinalArray := []string{}
	var mnemonicFinal string
	if truncated {
		fmt.Println("Untruncating mnemonic...")
		for _, truncatedWord := range mnemonicArray {
			fmt.Println("Untruncating word:", truncatedWord)
			untruncatedWord := bip39.TruncatedWordMap[truncatedWord]
			if bip39.TruncatedWordMap[truncatedWord] == "" {
				log.Fatalln("No matching word found for:", truncatedWord)
			}
			mnemonicFinalArray = append(mnemonicFinalArray, untruncatedWord)
		}
		mnemonicFinal = strings.Join(mnemonicFinalArray, " ")
	} else {
		fmt.Println("Using full mnemonic...")
		mnemonicFinal = mnemonic
	}

	child := prepareSeed(mnemonicFinal, password, truncated)
	mnemonicToScriptPubKey(child)
	mnemonicToWIF(child)
}

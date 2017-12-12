package main

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/btcsuite/btcutil/base58"
	"gopkg.in/urfave/cli.v2"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	bip39 "github.com/nii236/go-bip39"
)

func prepareSeed(mnemonic string, flags *RecoverFlags) (*hdkeychain.ExtendedKey, error) {
	seed := bip39.NewSeed(flags.Mnemonic, flags.Password)
	fmt.Println("mnemonic:", mnemonic)
	hdkey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}

	chainOffset := uint32(0)
	if flags.WalletTestnet {
		chainOffset = 1
	}

	child, _ := hdkey.Child(0x80000000 + 44)
	child, _ = child.Child(0x80000000 + chainOffset)
	child, _ = child.Child(uint32(0x80000000 + flags.WalletAccountIndex))
	child, _ = child.Child(uint32(0 + flags.WalletChainIndex))
	child, _ = child.Child(uint32(flags.WalletAddressIndex))

	return child, nil
}

func mnemonicToScriptPubKey(child *hdkeychain.ExtendedKey) error {
	scriptpubkey, err := child.Address(&chaincfg.MainNetParams)
	fmt.Println("Extended Key:", child.String())
	if err != nil {
		return err
	}
	fmt.Println("scriptPubKey:", scriptpubkey)
	return nil
}

func mnemonicToWIF(child *hdkeychain.ExtendedKey) error {
	privKey, err := child.ECPrivKey()

	if err != nil {
		return err
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
	return nil
}

// RecoverFlags are the settings needed to recover the keys from the mnemonic
type RecoverFlags struct {
	Mnemonic           string
	TruncatedMnemonic  bool
	Password           string
	WalletTestnet      bool
	WalletAccountIndex int
	WalletAddressIndex int
	WalletChainIndex   int
}

func main() {
	app := &cli.App{
		Name:  "Cryptosteel Toolkit",
		Usage: "Create and recover WIFs and scriptpubkeys from Cryptosteel mnemonics",

		Commands: []*cli.Command{
			{
				Name:   "generate",
				Usage:  "Generate a mnemonic",
				Action: generate,
			},
			{
				Name:   "recover",
				Usage:  "Recover WIF and scriptpubkeys from a mneomnic",
				Action: recover,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "mnemonic",
						Value: "",
						Usage: "Full mnemonic",
					},
					&cli.BoolFlag{
						Name:  "truncated",
						Value: false,
						Usage: "Truncated mnemonic (4 letters per word)",
					},
					&cli.StringFlag{
						Name:  "password",
						Value: "",
						Usage: "Mnemonic password",
					},
					&cli.BoolFlag{
						Name:  "wallet-testnet",
						Value: false,
						Usage: "Use testnet",
					},
					&cli.IntFlag{
						Name:  "wallet-account-index",
						Value: 0,
						Usage: "The wallet account index",
					},
					&cli.IntFlag{
						Name:  "wallet-address-index",
						Value: 0,
						Usage: "The wallet address index",
					},
					&cli.IntFlag{
						Name:  "wallet-chain-index",
						Value: 0,
						Usage: "The wallet chain index",
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}

func generate(c *cli.Context) error {
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return err
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return err
	}
	fmt.Println(mnemonic)
	return nil
}

func recover(c *cli.Context) error {
	flags := &RecoverFlags{
		Mnemonic:           c.String("mnemonic"),
		TruncatedMnemonic:  c.Bool("truncated"),
		Password:           c.String("password"),
		WalletTestnet:      c.Bool("wallet-testnet"),
		WalletAccountIndex: c.Int("wallet-account-index"),
		WalletAddressIndex: c.Int("wallet-address-index"),
		WalletChainIndex:   c.Int("wallet-chain-index"),
	}

	NetworkCoinIndex := 0

	if flags.WalletTestnet {
		NetworkCoinIndex = 1
	}

	fmt.Printf("Derivation Path: m / 44' / %d' / %d' / %d / %d\n", NetworkCoinIndex, flags.WalletAccountIndex, flags.WalletChainIndex, flags.WalletAddressIndex)

	mnemonicArray := strings.Split(flags.Mnemonic, " ")
	mnemonicFinalArray := []string{}
	var mnemonicFinal string
	if flags.TruncatedMnemonic {
		fmt.Println("Untruncating mnemonic...")
		for _, truncatedWord := range mnemonicArray {
			fmt.Println("Untruncating word:", truncatedWord)
			untruncatedWord := bip39.TruncatedWordMap[truncatedWord]
			if bip39.TruncatedWordMap[truncatedWord] == "" {
				return errors.New("No matching word found for: " + truncatedWord)
			}
			mnemonicFinalArray = append(mnemonicFinalArray, untruncatedWord)
		}
		mnemonicFinal = strings.Join(mnemonicFinalArray, " ")
	} else {
		fmt.Println("Using full mnemonic...")
		mnemonicFinal = flags.Mnemonic
	}

	child, err := prepareSeed(mnemonicFinal, flags)
	if err != nil {
		return err
	}
	err = mnemonicToScriptPubKey(child)
	if err != nil {
		return err
	}
	err = mnemonicToWIF(child)
	if err != nil {
		return err
	}
	return nil
}

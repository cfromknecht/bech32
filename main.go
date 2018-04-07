package main

import (
	"fmt"
	"os"

	"github.com/roasbeef/btcutil/bech32"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "bech32"
	app.Version = "0.0.0"
	app.Usage = "bech32 encode/decode"

	app.Commands = []cli.Command{
		encodeCommand,
		decodeCommand,
	}

	if err := app.Run(os.Args); err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "[bech32] %v\n", err)
	os.Exit(1)
}

var encodeCommand = cli.Command{
	Name:      "encode",
	Usage:     "bech32 encode payload with a human-readable prefix",
	ArgsUsage: "hrp payload",
	Description: `
	Converts the given payload into base 32, then encodes the base 32 bytes
	using bech32. The specified human-readable prefix will be prepended to
	the encoding.
	`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name: "hrp",
			Usage: "the human-readable prefix to be prepended to " +
				"the bech32 encoding",
		},
		cli.StringFlag{
			Name:  "payload",
			Usage: "the string payload to be bech32 encoded",
		},
	},
	Action: encode,
}

func encode(ctx *cli.Context) error {
	if ctx.NArg() == 0 && ctx.NumFlags() == 0 {
		cli.ShowCommandHelp(ctx, "encode")
		return nil
	}

	var (
		hrp     string
		payload []byte
	)

	args := ctx.Args()

	switch {
	case ctx.IsSet("hrp"):
		hrp = ctx.String("hrp")
	case args.Present():
		hrp = args.First()
		args = args.Tail()
	default:
		return fmt.Errorf("Human-readable prefix argument missing")
	}

	switch {
	case ctx.IsSet("payload"):
		payload = []byte(ctx.String("payload"))
	case args.Present():
		payload = []byte(args.First())
	default:
		return fmt.Errorf("Payload argument missing")
	}

	// Transform the requested payload into base 32.
	payload32, err := bech32.ConvertBits(payload, 8, 5, true)
	if err != nil {
		return fmt.Errorf("unable to convert bits: %v", err)
	}

	// Encode the base 32 payload.
	encoded, err := bech32.Encode(hrp, payload32)
	if err != nil {
		return fmt.Errorf("unable to encode: %v\n", err)
	}

	fmt.Println(encoded)

	return nil
}

var decodeCommand = cli.Command{
	Name:      "decode",
	Usage:     "bech32 decode human-readable prefix and payload",
	ArgsUsage: "encoding",
	Description: `
	Decodes the given bech32 encoding, and converts it back to base 256.
	Both the human-readable prefix and payload are returned.
	`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "encoding",
			Usage: "the bech32 encoding to decode",
		},
	},
	Action: decode,
}

func decode(ctx *cli.Context) error {
	if ctx.NArg() == 0 && ctx.NumFlags() == 0 {
		cli.ShowCommandHelp(ctx, "decode")
		return nil
	}

	var encoding string

	args := ctx.Args()

	switch {
	case ctx.IsSet("encoding"):
		encoding = ctx.String("encoding")
	case args.Present():
		encoding = args.First()
	default:
		return fmt.Errorf("Encoding argument missing")
	}

	// Decode the base 32 payload.
	decodedHrp, decodedPayload32, err := bech32.Decode(encoding)
	if err != nil {
		return fmt.Errorf("unable to decode: %v\n", err)
	}

	// Transform the decoded payload back into base 256.
	decodedPayload, err := bech32.ConvertBits(decodedPayload32, 5, 8, false)
	if err != nil {
		return fmt.Errorf("unable to convert bits: %v", err)
	}

	fmt.Printf("%s %s\n", decodedHrp, string(decodedPayload))

	return nil
}

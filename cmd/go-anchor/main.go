package main

import (
	"flag"
	"fmt"
	"os"

	idlcmd "go-solana-anchor/cmd/idl"
)

func main() {
	idlCmd := flag.NewFlagSet("idl", flag.ExitOnError)
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "idl":
		idlCmd.Parse(os.Args[2:])
		if idlCmd.NArg() < 1 {
			fmt.Fprintln(os.Stderr, "Usage: go-anchor idl <fetch|validate|convert|gen> [args...]")
			os.Exit(1)
		}
		runIDL(idlCmd.Args())
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "Usage: go-anchor <command> [args...]")
	fmt.Fprintln(os.Stderr, "Commands:")
	fmt.Fprintln(os.Stderr, "  idl fetch <program_id> -o idl.json    Fetch IDL from chain")
	fmt.Fprintln(os.Stderr, "  idl validate idl.json                 Validate IDL file")
	fmt.Fprintln(os.Stderr, "  idl convert legacy.json -o v30.json   Convert legacy IDL to v0.30")
	fmt.Fprintln(os.Stderr, "  idl gen -i idl.json -o pkg/ -p pkg    Generate Go client from IDL")
}

func runIDL(args []string) {
	sub := args[0]
	rest := args[1:]

	switch sub {
	case "fetch":
		out := "idl.json"
		var programID string
		for i := 0; i < len(rest); i++ {
			if rest[i] == "-o" && i+1 < len(rest) {
				out = rest[i+1]
				i++
			} else if len(rest[i]) > 3 && rest[i][:3] == "-o=" {
				out = rest[i][3:]
			} else if rest[i] != "" && programID == "" {
				programID = rest[i]
			}
		}
		if programID == "" {
			fmt.Fprintln(os.Stderr, "Usage: go-anchor idl fetch <program_id> -o idl.json")
			os.Exit(1)
		}
		rpcEndpoint := os.Getenv("RPC_URL")
		if rpcEndpoint == "" {
			rpcEndpoint = "https://api.mainnet-beta.solana.com"
		}
		if err := idlcmd.FetchIDL(rpcEndpoint, programID, out); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("IDL written to %s\n", out)
	case "validate":
		if len(rest) < 1 {
			fmt.Fprintln(os.Stderr, "Usage: go-anchor idl validate idl.json")
			os.Exit(1)
		}
		if err := idlcmd.ValidateIDL(rest[0]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("IDL is valid")
	case "convert":
		out := "idl.json"
		var input string
		for i := 0; i < len(rest); i++ {
			if rest[i] == "-o" && i+1 < len(rest) {
				out = rest[i+1]
				i++
			} else if len(rest[i]) > 3 && rest[i][:3] == "-o=" {
				out = rest[i][3:]
			} else if rest[i] != "" && input == "" {
				input = rest[i]
			}
		}
		if input == "" {
			fmt.Fprintln(os.Stderr, "Usage: go-anchor idl convert legacy.json -o v30.json")
			os.Exit(1)
		}
		if err := idlcmd.ConvertIDL(input, out); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Converted IDL written to %s\n", out)
	case "gen":
		input := ""
		out := "."
		pkg := ""
		for i := 0; i < len(rest); i++ {
			if rest[i] == "-i" && i+1 < len(rest) {
				input = rest[i+1]
				i++
			} else if rest[i] == "-o" && i+1 < len(rest) {
				out = rest[i+1]
				i++
			} else if rest[i] == "-p" && i+1 < len(rest) {
				pkg = rest[i+1]
				i++
			}
		}
		if input == "" {
			fmt.Fprintln(os.Stderr, "Usage: go-anchor idl gen -i idl.json -o <dir> [-p package]")
			os.Exit(1)
		}
		if err := idlcmd.Generate(input, idlcmd.GenOpts{Package: pkg, Output: out}); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Generated Go client in %s\n", out)
	default:
		fmt.Fprintf(os.Stderr, "Unknown idl subcommand: %s\n", sub)
		os.Exit(1)
	}
}

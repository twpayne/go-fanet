package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"reflect"
	"regexp"

	"github.com/twpayne/go-fanet"
)

var fanetRx = regexp.MustCompile(`#[A-Z]{3}.*`)

type commandJSON struct {
	Type    string
	Command fanet.Command
}

type responseJSON struct {
	Type          string
	Response      fanet.Response
	FNFPayload    any   `json:",omitempty"`
	FNFPayloadErr error `json:",omitempty"`
}

func run() error {
	flag.Parse()

	encoder := json.NewEncoder(os.Stdout)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		m := fanetRx.FindStringSubmatch(scanner.Text())
		if m == nil {
			continue
		}
		if command, err := fanet.ParseCommandString(m[0] + "\n"); err == nil {
			if err := encoder.Encode(commandJSON{
				Type:    reflect.TypeOf(command).String(),
				Command: command,
			}); err != nil {
				return err
			}
		}
		if response, err := fanet.ParseResponseString(m[0] + "\n"); err == nil {
			var fnfPayload any
			var fnfPayloadErr error
			if fnfResponse, ok := response.(*fanet.FNFResponse); ok {
				fnfPayload, fnfPayloadErr = fnfResponse.ParsePayload()
			}
			if err := encoder.Encode(responseJSON{
				Type:          reflect.TypeOf(response).String(),
				Response:      response,
				FNFPayload:    fnfPayload,
				FNFPayloadErr: fnfPayloadErr,
			}); err != nil {
				return err
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

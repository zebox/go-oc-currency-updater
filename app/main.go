// This is CLI util implement simple external OpenCart 3.x.x currency updater (OCCU) call OpenCart webhook which trigger update currency
// It can be scheduling with cron job or other schedulers
package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
)

// Options contain main parameters for the util
var opts struct {
	BaseOpenCartUrl string `short:"b" long:"base-url" env:"OCCU_BASE_URL" required:"true" description:"Base URL for OpenCart site (https://example.org)"`
	Login           string `short:"u" long:"username" env:"OCCU_LOGIN" required:"true" description:"Username for access to OpenCart admin panel"`
	Password        string `short:"p" long:"password" env:"OCCU_PASSWORD" required:"true" description:"Password for access to OpenCart admin panel"`
}

func main() {
	os.Exit(run())
}

func run() int {

	p := flags.NewParser(&opts, flags.PassDoubleDash|flags.HelpFlag)
	if _, err := p.Parse(); err != nil {
		if err.(*flags.Error).Type != flags.ErrHelp {
			fmt.Printf("%v\n", err)
			return 1
		}
		p.WriteHelp(os.Stderr)
		return 2
	}

	creds, err := getUserToken()
	if err != nil {
		fmt.Println(err)
		return 1
	}

	if err = doUpdateCurrency(creds); err != nil {
		fmt.Println(err)
		return 1
	}

	return 0
}

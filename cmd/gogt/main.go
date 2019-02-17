package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/umaumax/gogt"
)

var (
	lang string
)

func init() {
	flag.StringVar(&lang, "l", "en", "to lang e.g. en, ja")
}

func main() {
	flag.Parse()

	client := gogt.Client{}
	client.CacheVq()

	for _, v := range flag.Args() {
		ret, err := client.Translate(lang, v)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println(ret)
	}
}

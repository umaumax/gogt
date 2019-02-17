package main

import (
	"fmt"
	"log"
	"sync"

	fuzzyfinder "github.com/umaumax/go-fuzzyfinder"
	"github.com/umaumax/gogt"
)

func main() {
	client := gogt.Client{}
	go func() {
		client.CacheVq()
	}()

	_, err := fuzzyfinder.FindMulti(
		[]string{},
		func(i int) string {
			return ""
		},
		fuzzyfinder.WithPreviewWindow(func(i, width, height int, query string) string {
			if query == "" {
				return "No query"
			}
			wg := &sync.WaitGroup{}
			_ = wg
			langs := []string{"en", "ja"}
			rets := make([]string, len(langs))
			for i, lang := range langs {
				wg.Add(1)
				go func(i int, lang string) {
					ret, err := client.Translate(lang, query)
					if err != nil {
						rets[i] = fmt.Sprintln(err)
						return
					}
					rets[i] = ret
					wg.Done()
				}(i, lang)
			}
			wg.Wait()
			var ret string
			for i, lang := range langs {
				ret += fmt.Sprintf("[%s]\n", lang)
				ret += fmt.Sprintf("%s\n", rets[i])
				ret += "\n"
				ret += "\n"
			}
			return ret
		}))
	if err != nil {
		log.Fatal(err)
	}
}

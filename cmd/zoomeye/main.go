package main

import (
	"fmt"

	"github.com/zan8in/gologger"
	"github.com/zan8in/zoomeye/pkg/runner"
)

func main() {
	options, err := runner.ParseOptions()
	if err != nil {
		gologger.Fatal().Msgf("Parse options error: %v", err)
	}

	r, err := runner.New(options)
	if err != nil {
		gologger.Fatal().Msg(err.Error())
	}

	// result, err := r.Run()
	// if err != nil {
	// 	gologger.Fatal().Msg(err.Error())
	// }

	// fmt.Println("total: ", result[0].Total)
	// for _, r := range result {
	// 	for _, v := range r.Results {
	// 		var ip, service string
	// 		var port float64
	// 		ip = v["ip"].(string)
	// 		portinfo := v["portinfo"].(map[string]any)
	// 		if portinfo != nil {
	// 			port = portinfo["port"].(float64)
	// 			service = portinfo["service"].(string)
	// 		}
	// 		fmt.Printf("%s://%s:%d\n", service, ip, int(port))
	// 	}
	// }

	resultChan, err := r.RunChan()
	if err != nil {
		gologger.Fatal().Msg(err.Error())
	}

	var results []string
	var currentTotal, total int
	for r := range resultChan {
		if total == 0 {
			total = r.Total
		}
		for _, v := range r.Results {
			var ip, port, service string
			ip = v["ip"].(string)
			portinfo := v["portinfo"].(map[string]any)
			if portinfo != nil {
				switch portinfo["port"].(type) {
				case float64:
					port = fmt.Sprintf("%f", portinfo["port"].(float64))
				case string:
					port = portinfo["port"].(string)
				default:
					port = fmt.Sprintf("%v", portinfo["port"])
				}
				service = portinfo["service"].(string)
			}

			url := ""
			if len(port) != 0 {
				port = fmt.Sprintf(":%s", port)
			}
			if service == "http" || service == "https" {
				url = fmt.Sprintf("%s://%s%s", service, ip, port)
			} else {
				url = fmt.Sprintf("%s%s", ip, port)
			}
			results = append(results, url)
			currentTotal++
			if currentTotal == options.Count {
				break
			}
		}
		fmt.Printf("\rZoomEye Searching... Total: %d, Count: %d, Current: %d\n", total, options.Count, currentTotal)
	}

	fmt.Println("")

	// if len(results) > 0 {
	// 	for _, url := range results {
	// 		fmt.Println(url)
	// 	}
	// }
}

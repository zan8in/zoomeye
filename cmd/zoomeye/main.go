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

	var urlSlice []string
	var currentTotal, total int
	for r := range resultChan {
		currentTotal += len(r.Results)
		if total == 0 {
			total = r.Total
		}
		for _, v := range r.Results {
			var ip, service string
			var port float64
			ip = v["ip"].(string)
			portinfo := v["portinfo"].(map[string]any)
			if portinfo != nil {
				port = portinfo["port"].(float64)
				service = portinfo["service"].(string)
			}

			url := ""
			strPort := ""
			if int(port) != 0 {
				strPort = fmt.Sprintf(":%d", int(port))
			}
			if service == "http" || service == "https" {
				url = fmt.Sprintf("%s://%s%s", service, ip, strPort)
			} else {
				url = fmt.Sprintf("%s%s", ip, strPort)
			}
			urlSlice = append(urlSlice, url)
		}
		fmt.Printf("\rZoomEye Searching... Total: %d, Current: %d\n", total, currentTotal)
	}

	fmt.Println("-----------------------------")

	if len(urlSlice) > 0 {
		for _, url := range urlSlice {
			fmt.Println(url)
		}
	}
}

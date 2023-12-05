package zoomeye

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/zan8in/zoomeye/pkg/retryhttp"
)

type Zoomeye struct {
	Options *Options
	ApiURL  string
}

func New(options *Options) (z *Zoomeye, err error) {
	z = &Zoomeye{Options: options}
	return z, err
}

func (z *Zoomeye) Get(apiKey, page string) (result *Result, err error) {

	zoomeyeUrl := fmt.Sprintf(URL, url.QueryEscape("after:'"+GetLastMonthDate(-10)+"'+"+z.Options.Search), page)

	body, statusCode, err := retryhttp.GetWithApiKey(zoomeyeUrl, apiKey)
	if err != nil {
		return nil, err
	}
	fmt.Println(zoomeyeUrl, statusCode)

	if err = json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result, err
}

func GetLastMonthDate(n int) string {
	nowTime := time.Now()
	getTime := nowTime.AddDate(0, n, 0)
	return getTime.Format("2006-01-02")
}

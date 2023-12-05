package runner

import (
	"errors"
	"strconv"

	randutil "github.com/zan8in/pins/rand"
	"github.com/zan8in/zoomeye/pkg/config"
	"github.com/zan8in/zoomeye/pkg/result"
	"github.com/zan8in/zoomeye/pkg/retryhttp"
	"github.com/zan8in/zoomeye/pkg/zoomeye"
)

type Runner struct {
	Options      *Options
	Config       *config.Config
	Zoomeye      *zoomeye.Zoomeye
	Result       *result.Result
	CurrentTotal int // 表示批量导出数据进度中，当时已获取的数据总数
}

func New(options *Options) (runner *Runner, err error) {
	runner = &Runner{
		Options:      options,
		Result:       result.NewResult(),
		CurrentTotal: 0,
	}

	zoomeyeOpt := &zoomeye.Options{
		Search: options.Search,
	}

	if runner.Zoomeye, err = zoomeye.New(zoomeyeOpt); err != nil {
		return nil, err
	}

	if config, err := config.NewConfig(); err != nil {
		return nil, err
	} else {
		runner.Config = config
	}

	retryhttp.Init(&retryhttp.Options{
		Timeout: DefaultTimeout,
		Retries: DefaultRetries,
	})

	return runner, err
}

func (runner *Runner) Run() (zoomeyeResults []*zoomeye.Result, err error) {
	key := config.GetApiKey()
	if len(key) == 0 {
		return nil, errors.New("api key is expired")
	}

	var currentPage, lenResults, totalResults int
	currentPage = 1
	for {
		result, err := runner.Query(key, runner.Options.Search, strconv.Itoa(currentPage))
		if err != nil || result == nil {
			break
		}
		zoomeyeResults = append(zoomeyeResults, result)

		currentPage++
		lenResults += len(result.Results)
		totalResults = result.Total

		if lenResults > totalResults || lenResults == 0 {
			break
		}

		randutil.RandSleep(zoomeye.DefaultSleepTime)
	}

	return zoomeyeResults, err
}

func (runner *Runner) RunChan() (chan zoomeye.Result, error) {
	key := config.GetApiKey()
	if len(key) == 0 {
		return nil, errors.New("api key is expired")
	}

	results := make(chan zoomeye.Result)

	go func() {
		defer close(results)

		var currentPage, lenResults, totalResults int
		currentPage = 1
		for {
			result, err := runner.Query(key, runner.Options.Search, strconv.Itoa(currentPage))
			if err != nil || result == nil {
				break
			}

			results <- *result

			currentPage++
			lenResults += len(result.Results)
			if totalResults == 0 {
				totalResults = result.Total
			}

			if runner.Options.Count > 0 && lenResults >= runner.Options.Count {
				break
			}

			if lenResults >= totalResults || lenResults == 0 {
				break
			}

			randutil.RandSleep(zoomeye.DefaultSleepTime)
		}
	}()

	return results, nil
}

func (runner *Runner) Query(key, search, page string) (result *zoomeye.Result, err error) {

	result, err = runner.Zoomeye.Get(key, page)
	if err != nil {
		return nil, err
	}

	return result, err
}

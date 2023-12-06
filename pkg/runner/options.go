package runner

import (
	"errors"

	"github.com/zan8in/goflags"
	"github.com/zan8in/zoomeye/pkg/config"
)

type Options struct {
	Search string
	Page   int
	Count  int
	ApiKey string
	After  string
}

func ParseOptions() (*Options, error) {
	options := &Options{}

	flagSet := goflags.NewFlagSet()
	flagSet.SetDescription(`ZoomEye`)

	flagSet.CreateGroup("input", "Input",
		flagSet.StringVarP(&options.Search, "search", "s", "", "query conditions"),
		flagSet.StringVar(&options.ApiKey, "api", "", "api key"),
		flagSet.IntVar(&options.Count, "count", DefaultCount, "query count"),
	)

	_ = flagSet.Parse()

	if err := options.validateOptions(); err != nil {
		return nil, err
	}

	return options, nil
}

func (options *Options) validateOptions() (err error) {

	if len(options.Search) == 0 {
		return errors.New("search query is empty")
	}

	if len(options.ApiKey) > 0 {
		config.ValidKeys = append(config.ValidKeys, options.ApiKey)
	}

	return err
}

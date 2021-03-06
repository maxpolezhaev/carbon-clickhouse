package uploader

import (
	"github.com/lomik/carbon-clickhouse/helper/config"
	"time"
)

type Config struct {
	Type      string           `toml:"type"`  // points, series, points-reverse, series-reverse
	TableName string           `toml:"table"` // keep empty for same as key
	Timeout   *config.Duration `toml:"timeout"`
	Date      string           `toml:"date"` // for tree table
	Threads   int              `toml:"threads"`
	URL       string           `toml:"url"`
	CacheTTL  *config.Duration `toml:"cache-ttl"`
	TreeDate  time.Time        `toml:"-"`
}

func (cfg *Config) Parse() error {
	var err error

	if cfg.Date != "" {
		cfg.TreeDate, err = time.ParseInLocation("2006-01-02", cfg.Date, time.Local)
		if err != nil {
			return err
		}
	}

	return nil
}

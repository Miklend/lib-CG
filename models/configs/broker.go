package configs

import (
	"time"
)

type Broker struct {
	BrockerType  string        `yaml:"brocker_type"`
	Brokers      []string      `yaml:"brokers"`
	GroupID      string        `yaml:"group_id"`
	StartOffset  int64         `yaml:"start_offset"`
	BatchSize    int           `yaml:"batch_size"`
	BatchTimeout time.Duration `yaml:"batch_timeout"`
	Async        bool          `yaml:"async"`
}

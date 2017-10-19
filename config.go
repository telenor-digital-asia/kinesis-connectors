package connector

import (
	"os"
	"time"

	redis "gopkg.in/redis.v5"

	log "github.com/sirupsen/logrus"
)

const (
	defaultBufferSize = 500
	defaultRedisAddr  = "127.0.0.1:6379"
)

// Config vars for the application
type Config struct {
	// AppName is the application name and checkpoint namespace.
	AppName string

	// StreamName is the Kinesis stream.
	StreamName string

	// StreamRegion is the Kinesis stream.
	StreamRegion string

	// FlushInterval is a regular interval for flushing the buffer. Defaults to 1s.
	FlushInterval time.Duration

	// BufferSize determines the batch request size. Must not exceed 500. Defaults to 500.
	BufferSize int

	// Logger is the logger used. Defaults to log.Log.
	Logger *log.Logger

	// Checkpoint for tracking progress of consumer.
	Checkpoint Checkpoint
}

// defaults for configuration.
func (c *Config) setDefaults() {
	if c.Logger == nil {
		c.Logger = log.New()
	}

	c.Logger.WithFields(log.Fields{
		"package": "kinesis-connectors",
	})

	if c.AppName == "" {
		c.Logger.WithField("type", "config").Error("AppName required")
		os.Exit(1)
	}

	if c.StreamName == "" {
		c.Logger.WithField("type", "config").Error("StreamName required")
		os.Exit(1)
	}

	if c.StreamRegion == "" {
		c.Logger.WithField("type", "config").Error("StreamRegion required")
		os.Exit(1)
	}

	c.Logger.WithFields(log.Fields{
		"app":    c.AppName,
		"stream": c.StreamName,
		"region": c.StreamRegion,
	})

	if c.BufferSize == 0 {
		c.BufferSize = defaultBufferSize
	}

	if c.FlushInterval == 0 {
		c.FlushInterval = time.Second
	}

	if c.Checkpoint == nil {
		client, err := defaultRedisClient()
		if err != nil {
			c.Logger.WithError(err).Error("Redis connection failed")
			os.Exit(1)
		}
		c.Checkpoint = &RedisCheckpoint{
			AppName:    c.AppName,
			StreamName: c.StreamName,
			Client:     client,
		}
	}
}

func defaultRedisClient() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: defaultRedisAddr,
	})
	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}

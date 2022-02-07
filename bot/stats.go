package bot

import (
	"context"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/termora/berry/common/log"
)

// StatsClient is an InfluxDB client
type StatsClient struct {
	Client api.WriteAPI

	queriesMu sync.Mutex
	queries   uint32

	cmdsMu sync.Mutex
	cmds   uint32

	apiCallsMu sync.Mutex
	apiCalls   uint32

	guildCount func() int
}

// IncQuery increments the query count by one
func (c *StatsClient) IncQuery() {
	if c == nil {
		return
	}

	c.queriesMu.Lock()
	c.queries++
	c.queriesMu.Unlock()
}

// IncCommand increments the command count by one
func (c *StatsClient) IncCommand() {
	if c == nil {
		return
	}

	c.cmdsMu.Lock()
	c.cmds++
	c.cmdsMu.Unlock()
}

// IncAPICall increments the API call count by one
func (c *StatsClient) IncAPICall() {
	if c == nil {
		return
	}

	c.apiCallsMu.Lock()
	c.apiCalls++
	c.apiCallsMu.Unlock()
}

func (c *StatsClient) submit() {
	if c == nil {
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	defer stop()

	ticker := time.NewTicker(time.Minute)

	for {
		select {
		case <-ticker.C:
			// submit metrics
			go c.submitInner()
		case <-ctx.Done():
			// break if we're shutting down
			ticker.Stop()
			c.Client.Flush()
			return
		}
	}
}

func (c *StatsClient) submitInner() {
	if c == nil {
		return
	}

	log.Info("Submitting metrics to InfluxDB")

	c.queriesMu.Lock()
	queries := c.queries
	c.queries = 0
	c.queriesMu.Unlock()

	c.cmdsMu.Lock()
	cmds := c.cmds
	c.cmds = 0
	c.cmdsMu.Unlock()

	c.apiCallsMu.Lock()
	apiCalls := c.apiCalls
	c.apiCalls = 0
	c.apiCallsMu.Unlock()

	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)

	data := map[string]interface{}{
		"api_calls":   apiCalls,
		"queries":     queries,
		"commands":    cmds,
		"alloc":       stats.Alloc,
		"sys":         stats.Sys,
		"total_alloc": stats.TotalAlloc,
		"goroutines":  runtime.NumGoroutine(),
	}

	if c.guildCount != nil {
		data["guilds"] = c.guildCount()
	}

	p := influxdb2.NewPoint("statistics", nil, data, time.Now())
	c.Client.WritePoint(p)
}

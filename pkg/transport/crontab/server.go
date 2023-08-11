package crontab

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"

	"github.com/go-eagle/eagle/pkg/config"
)

// Config crontab config
type Config struct {
	Timezone string
	Tasks    []Task
}

// Task crontab task
type Task struct {
	Name     string
	Schedule string
}

// Server crontab server
type Server struct {
	conf *Config

	// cron schedule
	schedule *cron.Cron
	jobs     map[string]cron.Job

	logger cron.Logger
	stop   chan struct{}
}

// NewServer new a crontab server
func NewServer(jobs map[string]cron.Job, logger cron.Logger) *Server {
	// load config
	cfg, err := loadConf()
	if err != nil {
		panic(err)
	}

	if len(cfg.Tasks) == 0 {
		panic("crontab config is empty")
	}
	if len(jobs) == 0 || jobs == nil {
		panic("crontab jobs is empty")
	}
	if len(cfg.Tasks) != len(jobs) {
		panic("crontab tasks and jobs not match")
	}

	loc, err := time.LoadLocation(cfg.Timezone)
	if err != nil {
		panic(err)
	}

	// new server
	return &Server{
		conf: cfg,
		schedule: cron.New(
			cron.WithLocation(loc),
			cron.WithLogger(logger),
		),
		jobs: jobs,
		stop: make(chan struct{}, 1),
	}
}

// Start the crontab server
func (s *Server) Start(ctx context.Context) error {
	for _, task := range s.conf.Tasks {
		task := task
		// get job
		job, ok := s.jobs[task.Name]
		if !ok {
			return fmt.Errorf("[crontab] job not found: %s", task.Name)
		}
		_, err := s.schedule.AddJob(task.Schedule, job)
		if err != nil {
			return errors.Wrapf(err, "[crontab] add job [%s] error", task.Name)
		}
	}

	s.schedule.Start()

	select {
	case <-s.stop:
		s.schedule.Stop()
	}

	return nil
}

// Stop the crontab server
func (s *Server) Stop(ctx context.Context) error {
	log.Printf("[crontab] server stopping...")
	s.stop <- struct{}{}
	return nil
}

// loadConf load config
func loadConf() (ret *Config, err error) {
	v, err := config.LoadWithType("crontab", "yaml")
	if err != nil {
		return nil, err
	}

	c := Config{}
	err = v.Unmarshal(&c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

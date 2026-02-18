package asynq

import (
	"github.com/hibiken/asynq"
)

type Client struct {
	client *asynq.Client
}

type Server struct {
	server *asynq.Server
}

func NewClient(redisAddr string) *Client {
	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr: redisAddr,
	})

	return &Client{
		client: client,
	}
}

func NewServer(redisAddr string, concurrency int) *Server {
	server := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr: redisAddr,
		},
		asynq.Config{
			Concurrency: concurrency,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	return &Server{
		server: server,
	}
}

func (c *Client) Enqueue(task *asynq.Task, opts ...asynq.Option) error {
	_, err := c.client.Enqueue(task, opts...)
	return err
}

func (c *Client) Close() error {
	return c.client.Close()
}

func (s *Server) Start(mux *asynq.ServeMux) error {
	return s.server.Start(mux)
}

func (s *Server) Stop() {
	s.server.Stop()
}

func (s *Server) Shutdown() {
	s.server.Shutdown()
}

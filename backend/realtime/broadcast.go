package realtime

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"backend/config"

	"github.com/redis/go-redis/v9"
)

type Broadcaster interface {
	BroadcastEvent(projectID uint, eventType string, payload interface{})
	Subscribe(projectID uint) chan string
	Unsubscribe(projectID uint, ch chan string)
}

type Service struct {
	redisClient *redis.Client
	ctx         context.Context
	subscribers map[uint]map[chan string]bool
	subMu       sync.RWMutex
}

func New(cfg *config.Config) *Service {
	client := redis.NewClient(&redis.Options{
		Addr: cfg.RedisURL,
	})

	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Redis connection established")

	s := &Service{
		redisClient: client,
		ctx:         ctx,
		subscribers: make(map[uint]map[chan string]bool),
	}

	// Start a global subscriber for this instance to listen to Redis Pub/Sub
	go s.listenToRedis()

	return s
}

type Event struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func (s *Service) BroadcastEvent(projectID uint, eventType string, payload interface{}) {
	event := Event{
		Type:    eventType,
		Payload: payload,
	}
	data, _ := json.Marshal(event)

	channel := fmt.Sprintf("project:%d:events", projectID)
	err := s.redisClient.Publish(s.ctx, channel, data).Err()
	if err != nil {
		log.Printf("Failed to publish to Redis: %v", err)
	}
}

func (s *Service) Subscribe(projectID uint) chan string {
	s.subMu.Lock()
	defer s.subMu.Unlock()

	if s.subscribers[projectID] == nil {
		s.subscribers[projectID] = make(map[chan string]bool)
	}

	ch := make(chan string, 10)
	s.subscribers[projectID][ch] = true
	return ch
}

func (s *Service) Unsubscribe(projectID uint, ch chan string) {
	s.subMu.Lock()
	defer s.subMu.Unlock()

	if s.subscribers[projectID] != nil {
		delete(s.subscribers[projectID], ch)
		close(ch)
	}
}

func (s *Service) listenToRedis() {
	// Pattern subscribe to all project channels
	pubsub := s.redisClient.PSubscribe(s.ctx, "project:*:events")
	defer pubsub.Close()

	ch := pubsub.Channel()
	for msg := range ch {
		var projectID uint
		fmt.Sscanf(msg.Channel, "project:%d:events", &projectID)

		s.subMu.RLock()
		if subs, ok := s.subscribers[projectID]; ok {
			for c := range subs {
				select {
				case c <- msg.Payload:
				default:
					// Skip slow consumers
				}
			}
		}
		s.subMu.RUnlock()
	}
}

func (s *Service) Shutdown() {
	if s.redisClient != nil {
		s.redisClient.Close()
	}
}

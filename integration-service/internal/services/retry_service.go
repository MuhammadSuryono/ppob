package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type CompensationJob struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	JobID         string         `gorm:"uniqueIndex;size:100" json:"job_id"`
	TransactionID string         `gorm:"index;size:100" json:"transaction_id"`
	Action        string         `gorm:"size:50;not null" json:"action"`
	Payload       string         `gorm:"type:text" json:"payload"`
	Status        string         `gorm:"size:20;default:pending" json:"status"`
	RetryCount    int            `gorm:"default:0" json:"retry_count"`
	MaxRetries    int            `gorm:"default:3" json:"max_retries"`
	ErrorMessage  string         `gorm:"type:text" json:"error_message"`
	NextRetryAt   *time.Time     `json:"next_retry_at"`
	CompletedAt   *time.Time     `json:"completed_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

type RetryPolicy struct {
	MaxRetries      int
	InitialInterval time.Duration
	MaxInterval     time.Duration
	Multiplier      float64
}

var defaultRetryPolicy = RetryPolicy{
	MaxRetries:      3,
	InitialInterval: 1 * time.Second,
	MaxInterval:     30 * time.Second,
	Multiplier:      2.0,
}

type CompensationService struct {
	db     *gorm.DB
	redis  *redis.Client
	policy RetryPolicy
}

func NewCompensationService(db *gorm.DB, redis *redis.Client) *CompensationService {
	return &CompensationService{
		db:     db,
		redis:  redis,
		policy: defaultRetryPolicy,
	}
}

func (s *CompensationService) CreateJob(ctx context.Context, transactionID string, action string, payload map[string]interface{}) (*CompensationJob, error) {
	payloadJSON, _ := json.Marshal(payload)

	job := &CompensationJob{
		JobID:         uuid.New().String(),
		TransactionID: transactionID,
		Action:        action,
		Payload:       string(payloadJSON),
		Status:        "pending",
		RetryCount:    0,
		MaxRetries:    s.policy.MaxRetries,
	}

	if err := s.db.Create(job).Error; err != nil {
		return nil, err
	}

	return job, nil
}

func (s *CompensationService) ScheduleRetry(ctx context.Context, jobID string) error {
	var job CompensationJob
	if err := s.db.Where("job_id = ?", jobID).First(&job).Error; err != nil {
		return err
	}

	if job.RetryCount >= job.MaxRetries {
		job.Status = "failed"
		job.ErrorMessage = "Max retries exceeded"
		now := time.Now()
		job.CompletedAt = &now
		return s.db.Save(&job).Error
	}

	delay := s.calculateBackoff(job.RetryCount)
	nextRetry := time.Now().Add(delay)

	job.RetryCount++
	job.Status = "scheduled"
	job.NextRetryAt = &nextRetry

	if err := s.db.Save(&job).Error; err != nil {
		return err
	}

	jobKey := fmt.Sprintf("compensation:%s", jobID)
	s.redis.Set(ctx, jobKey, job.JobID, delay)

	return nil
}

func (s *CompensationService) calculateBackoff(retryCount int) time.Duration {
	delay := s.policy.InitialInterval
	for i := 0; i < retryCount; i++ {
		delay = time.Duration(float64(delay) * s.policy.Multiplier)
		if delay > s.policy.MaxInterval {
			delay = s.policy.MaxInterval
		}
	}
	return delay
}

func (s *CompensationService) MarkSuccess(ctx context.Context, jobID string) error {
	var job CompensationJob
	if err := s.db.Where("job_id = ?", jobID).First(&job).Error; err != nil {
		return err
	}

	job.Status = "completed"
	now := time.Now()
	job.CompletedAt = &now

	return s.db.Save(&job).Error
}

func (s *CompensationService) MarkFailed(ctx context.Context, jobID string, errorMessage string) error {
	var job CompensationJob
	if err := s.db.Where("job_id = ?", jobID).First(&job).Error; err != nil {
		return err
	}

	job.Status = "failed"
	job.ErrorMessage = errorMessage
	now := time.Now()
	job.CompletedAt = &now

	return s.db.Save(&job).Error
}

func (s *CompensationService) GetPendingJobs(ctx context.Context) ([]CompensationJob, error) {
	var jobs []CompensationJob
	err := s.db.Where("status = ? AND (next_retry_at IS NULL OR next_retry_at <= ?)", "pending", time.Now()).Find(&jobs).Error
	return jobs, err
}

func (s *CompensationService) GetScheduledJobs(ctx context.Context) ([]CompensationJob, error) {
	var jobs []CompensationJob
	err := s.db.Where("status = ? AND next_retry_at <= ?", "scheduled", time.Now()).Find(&jobs).Error
	return jobs, err
}

func (s *CompensationService) GetDeadLetterJobs(ctx context.Context) ([]CompensationJob, error) {
	var jobs []CompensationJob
	err := s.db.Where("status = ? AND retry_count >= max_retries", "failed").Find(&jobs).Error
	return jobs, err
}

func (s *CompensationService) ReprocessJob(ctx context.Context, jobID string, processor func(context.Context, *CompensationJob) error) error {
	var job CompensationJob
	if err := s.db.Where("job_id = ?", jobID).First(&job).Error; err != nil {
		return err
	}

	if err := processor(ctx, &job); err != nil {
		return s.ScheduleRetry(ctx, jobID)
	}

	return s.MarkSuccess(ctx, jobID)
}

type DeadLetterQueue struct {
	redis *redis.Client
}

func NewDeadLetterQueue(redis *redis.Client) *DeadLetterQueue {
	return &DeadLetterQueue{redis: redis}
}

func (q *DeadLetterQueue) Push(ctx context.Context, jobID string, error error) error {
	key := "dead_letter_queue"
	entry := fmt.Sprintf("%s:%d:%s", jobID, time.Now().Unix(), error.Error())

	return q.redis.LPush(ctx, key, entry).Err()
}

func (q *DeadLetterQueue) Pop(ctx context.Context) (string, error) {
	key := "dead_letter_queue"
	return q.redis.RPop(ctx, key).Result()
}

func (q *DeadLetterQueue) Size(ctx context.Context) (int64, error) {
	key := "dead_letter_queue"
	return q.redis.LLen(ctx, key).Result()
}

func (q *DeadLetterQueue) Clear(ctx context.Context) error {
	key := "dead_letter_queue"
	return q.redis.Del(ctx, key).Err()
}
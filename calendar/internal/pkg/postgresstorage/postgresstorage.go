package postgresstorage

import (
	"context"
	"time"

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"

	"github.com/google/uuid"

	"github.com/AndreyAndreevich/otus_go/calendar/internal/domain"
	"go.uber.org/zap"
)

// PostgresStorage is storage which used sqlx
type PostgresStorage struct {
	logger *zap.Logger
	db     *sqlx.DB
}

type dbEvent struct {
	ID          string    `db:"id"`
	Heading     string    `db:"heading"`
	StartDate   time.Time `db:"start_date"`
	StartTime   time.Time `db:"start_time"`
	EndDate     time.Time `db:"end_date"`
	EndTime     time.Time `db:"end_time"`
	Description string    `db:"descr"`
	Owner       string    `db:"owner"`
}

// New created new PostgresStorage
func New(logger *zap.Logger, dsn string, maxOpenConn, maxIdleConn int) (*PostgresStorage, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		logger.Error("connect to db error", zap.Error(err))
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConn)
	db.SetMaxIdleConns(maxIdleConn)

	return &PostgresStorage{
		logger: logger,
		db:     db,
	}, nil
}

// HealthCheck is ping to db
func (s *PostgresStorage) HealthCheck() error {
	return s.db.Ping()
}

// Insert into events
func (s *PostgresStorage) Insert(ctx context.Context, event domain.Event) error {
	query := `INSERT INTO events (id, heading, start_date, start_time, end_date, end_time, descr, owner)
				VALUES (:id, :heading, :start_date, :start_time, :end_date, :end_time, :descr, :owner)`
	_, err := s.db.NamedExecContext(ctx, query, map[string]interface{}{
		"id":         event.ID.String(),
		"heading":    event.Heading,
		"start_date": event.DateTime,
		"start_time": event.DateTime,
		"end_date":   event.DateTime.Add(event.Duration),
		"end_time":   event.DateTime.Add(event.Duration),
		"descr":      event.Description,
		"owner":      event.Owner,
	})

	return err
}

// Remove from events
func (s *PostgresStorage) Remove(ctx context.Context, id domain.EventID) error {
	query := `DELETE FROM events WHERE id = :id`
	_, err := s.db.NamedExecContext(ctx, query, map[string]interface{}{
		"id": id.String(),
	})

	return err
}

// Update event
func (s *PostgresStorage) Update(ctx context.Context, event domain.Event) error {
	query := `UPDATE events SET heading = :heading, start_date = :start_date, start_time = :start_time, 
				end_date = :end_date, end_time = :end_time, descr = :descr, owner = :owner
				WHERE id = :id`

	_, err := s.db.NamedExecContext(ctx, query, map[string]interface{}{
		"id":         event.ID.String(),
		"heading":    event.Heading,
		"start_date": event.DateTime,
		"start_time": event.DateTime,
		"end_date":   event.DateTime.Add(event.Duration),
		"end_time":   event.DateTime.Add(event.Duration),
		"descr":      event.Description,
		"owner":      event.Owner,
	})

	return err
}

// Listing all events
func (s *PostgresStorage) Listing(ctx context.Context) ([]domain.Event, error) {
	query := `SELECT * FROM events`
	rows, err := s.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}

	return s.parse(rows)
}

// GetEventsInTime get events from time to time + duration
func (s *PostgresStorage) GetEventsInTime(ctx context.Context,
	time time.Time,
	duration time.Duration) ([]domain.Event, error) {

	query := `SELECT * FROM events WHERE (start_time + start_date) >= :start_time AND (end_time + end_date) <= :end_time`

	rows, err := s.db.NamedQueryContext(ctx, query, map[string]interface{}{
		"start_time": time,
		"end_time":   time.Add(duration),
	})

	if err != nil {
		return nil, err
	}

	return s.parse(rows)
}

func (s *PostgresStorage) parse(rows *sqlx.Rows) ([]domain.Event, error) {
	defer rows.Close()

	var events []domain.Event
	for rows.Next() {
		event := &dbEvent{}

		if err := rows.StructScan(event); err != nil {
			return nil, err
		}

		uuidID, err := uuid.Parse(event.ID)
		if err != nil {
			return nil, err
		}

		dateTime := event.StartTime.AddDate(
			event.StartDate.Year(),
			int(event.StartDate.Month()-1),
			event.StartDate.Day()-1)

		duration := event.EndTime.AddDate(
			event.EndDate.Year(),
			int(event.EndDate.Month()-1),
			event.EndDate.Day()-1).Sub(dateTime)

		events = append(events, domain.Event{
			ID:          uuidID,
			Heading:     event.Heading,
			DateTime:    dateTime,
			Duration:    duration,
			Description: event.Description,
			Owner:       event.Owner,
		})
	}

	return events, nil
}

func (s *PostgresStorage) Close() error {
	return s.db.Close()
}

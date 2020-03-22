package sqlxstorage

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"

	"github.com/AndreyAndreevich/otus_go/calendar/internal/domain"
	"github.com/gobuffalo/packr"
	"github.com/jmoiron/sqlx"
	migrate "github.com/rubenv/sql-migrate"
	"go.uber.org/zap"
)

// SQLXStorage is storage which used sqlx
type SQLXStorage struct {
	logger *zap.Logger
	db     *sqlx.DB
}

// New created new SQLXStorage
func New(logger *zap.Logger, db *sqlx.DB) *SQLXStorage {
	return &SQLXStorage{
		logger: logger,
		db:     db,
	}
}

// Migrate db
func (s *SQLXStorage) Migrate(dialect string) error {
	migrate.SetTable("_calendar_migrations")
	migrations := &migrate.PackrMigrationSource{
		Box: packr.NewBox("./migrations"),
	}
	s.logger.Debug("Storage migrations: start")
	n, err := migrate.Exec(s.db.DB, dialect, migrations, migrate.Up)
	if err != nil {
		return err
	}

	rows, err := migrate.GetMigrationRecords(s.db.DB, dialect)
	if err != nil {
		return err
	}
	cnt := len(rows)
	last := ""
	if cnt > 0 {
		last = rows[cnt-1].Id
	}

	s.logger.Info("Storage migrations: migrated", zap.Int("count", n), zap.String("current", last))
	return nil
}

// Insert into events
func (s *SQLXStorage) Insert(ctx context.Context, event domain.Event) error {
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
func (s *SQLXStorage) Remove(ctx context.Context, id domain.EventID) error {
	query := `DELETE FROM events WHERE id = :id`
	_, err := s.db.NamedExecContext(ctx, query, map[string]interface{}{
		"id": id.String(),
	})

	return err
}

// Update event
func (s *SQLXStorage) Update(ctx context.Context, event domain.Event) error {
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
func (s *SQLXStorage) Listing(ctx context.Context) ([]domain.Event, error) {
	query := `SELECT * FROM events`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	return s.parse(rows)
}

// GetEventsInTime get events from time to time + duration
func (s *SQLXStorage) GetEventsInTime(ctx context.Context,
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

	return s.parse(rows.Rows)
}

func (s *SQLXStorage) parse(rows *sql.Rows) ([]domain.Event, error) {
	defer rows.Close()

	var events []domain.Event
	for rows.Next() {
		var id string
		var heading string
		var startDate time.Time
		var startTime time.Time
		var endDate time.Time
		var endTime time.Time
		var descr string
		var owner string
		if err := rows.Scan(&id, &heading, &startDate, &startTime, &endDate, &endTime, &descr, &owner); err != nil {
			return nil, err
		}

		uuidID, err := uuid.Parse(id)
		if err != nil {
			return nil, err
		}

		dateTime := startTime.AddDate(startDate.Year(), int(startDate.Month()-1), startDate.Day()-1)
		duration := endTime.AddDate(endDate.Year(), int(endDate.Month()-1), endDate.Day()-1).Sub(dateTime)

		events = append(events, domain.Event{
			ID:          uuidID,
			Heading:     heading,
			DateTime:    dateTime,
			Duration:    duration,
			Description: descr,
			Owner:       owner,
		})
	}

	return events, nil
}

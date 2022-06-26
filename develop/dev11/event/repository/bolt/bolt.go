package bolt

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"

	"go.etcd.io/bbolt"

	"calendar/event"
)

type boltEventRepository struct {
	db *bbolt.DB
}

func NewBoltEventRepository(db *bbolt.DB) event.EventRepository {
	return &boltEventRepository{
		db: db,
	}
}

func NewBoltDB(path string) (*bbolt.DB, error) {
	db, err := bbolt.Open(path, 0666, nil)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (b *boltEventRepository) Create(user_id uint64, e event.Event) (event.Event, error) {
	var result event.Event
	err := b.db.Update(func(tx *bbolt.Tx) error {
		user, err := tx.CreateBucketIfNotExists(itob(user_id))
		if err != nil {
			return err
		}
		eBkt, err := user.CreateBucketIfNotExists([]byte("events"))
		if err != nil {
			return err
		}
		eventID, err := eBkt.NextSequence()
		if err != nil {
			return err
		}
		e.ID = eventID

		if buf, err := json.Marshal(e); err != nil {
			return err
		} else if err := eBkt.Put(itob(e.ID), buf); err != nil {
			return err
		}
		result = e

		return nil
	})

	if err != nil {
		return event.Event{}, err
	}

	return result, nil
}

func (b *boltEventRepository) Update(user_id uint64, e event.Event) error {
	return b.db.Update(func(tx *bbolt.Tx) error {
		user := tx.Bucket(itob(user_id))
		if user == nil {
			return fmt.Errorf("%w: user %d does not exist", event.ErrNotFound, user_id)
		}

		eBkt := user.Bucket([]byte("events"))
		if eBkt == nil {
			return fmt.Errorf("%w: user %d has no events", event.ErrNotFound, user_id)
		}

		v := eBkt.Get(itob(e.ID))
		if v == nil {
			return fmt.Errorf("%w: user %d has no %d event", event.ErrNotFound, user_id, e.ID)
		}

		buf, err := json.Marshal(e)
		if err != nil {
			return fmt.Errorf("%w: %s", event.ErrInternalServerError, err.Error())
		}

		return eBkt.Put(itob(e.ID), buf)
	})
}

func (b *boltEventRepository) Delete(user_id uint64, event_id uint64) error {
	return b.db.Update(func(tx *bbolt.Tx) error {
		user := tx.Bucket(itob(user_id))
		if user == nil {
			return fmt.Errorf("%w: user %d does not exist", event.ErrNotFound, user_id)
		}

		eBkt := user.Bucket([]byte("events"))
		if eBkt == nil {
			return fmt.Errorf("%w: user %d has no events", event.ErrNotFound, user_id)
		}

		v := eBkt.Get(itob(event_id))
		if v == nil {
			return fmt.Errorf("%w: user %d has no %d event", event.ErrNotFound, user_id, event_id)
		}

		return eBkt.Delete(itob(event_id))
	})
}

func (b *boltEventRepository) GetForDay(user_id uint64, day time.Time) ([]event.Event, error) {
	events := make([]event.Event, 0)
	err := b.db.View(func(tx *bbolt.Tx) error {
		user := tx.Bucket(itob(user_id))
		if user == nil {
			return fmt.Errorf("%w: user %d does not exist", event.ErrNotFound, user_id)
		}

		eBkt := user.Bucket([]byte("events"))
		if eBkt == nil {
			return fmt.Errorf("%w: user %d has no events", event.ErrNotFound, user_id)
		}
		c := eBkt.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var ev event.Event
			err := json.Unmarshal(v, &ev)
			if err != nil {
				return err
			}

			if ev.Date.Year() == day.Year() && ev.Date.Month() == day.Month() && ev.Date.Day() == day.Day() {
				events = append(events, ev)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return events, nil
}

func (b *boltEventRepository) GetForWeek(user_id uint64, week time.Time) ([]event.Event, error) {
	events := make([]event.Event, 0)
	err := b.db.View(func(tx *bbolt.Tx) error {
		user := tx.Bucket(itob(user_id))
		if user == nil {
			return fmt.Errorf("%w: user %d does not exist", event.ErrNotFound, user_id)
		}

		eBkt := user.Bucket([]byte("events"))
		if eBkt == nil {
			return fmt.Errorf("%w: user %d has no events", event.ErrNotFound, user_id)
		}
		c := eBkt.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var ev event.Event
			err := json.Unmarshal(v, &ev)
			if err != nil {
				return err
			}
			ew1, ey1 := ev.Date.ISOWeek()
			ew2, ey2 := week.ISOWeek()
			if ew1 == ew2 && ey1 == ey2 {
				events = append(events, ev)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return events, nil
}

func (b *boltEventRepository) GetForMonth(user_id uint64, month time.Time) ([]event.Event, error) {
	events := make([]event.Event, 0)
	err := b.db.View(func(tx *bbolt.Tx) error {
		user := tx.Bucket(itob(user_id))
		if user == nil {
			return fmt.Errorf("%w: user %d does not exist", event.ErrNotFound, user_id)
		}

		eBkt := user.Bucket([]byte("events"))
		if eBkt == nil {
			return fmt.Errorf("%w: user %d has no events", event.ErrNotFound, user_id)
		}
		c := eBkt.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var ev event.Event
			err := json.Unmarshal(v, &ev)
			if err != nil {
				return err
			}
			em1 := ev.Date.Month()
			ey1 := ev.Date.Year()
			em2 := month.Month()
			ey2 := month.Year()
			if em1 == em2 && ey1 == ey2 {
				events = append(events, ev)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return events, nil
}

// itob returns an 8-byte big endian representation of v.
func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

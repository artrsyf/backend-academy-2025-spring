# Пакетная обработка

Пакетная (batch) обработка — это подход к обработке данных, при котором данные обрабатываются не по одному, а группами (батчами). Это повышает производительность и снижает нагрузку на систему.


```go
ctx := context.Background()

for {
    rows, err := db.Query(ctx, `
        SELECT id, payload
        FROM tasks
        WHERE status = 'pending'
        ORDER BY id
        LIMIT $1`, batchSize)
    if err != nil {
        log.Fatalf("failed to query tasks: %v", err)
    }
    defer rows.Close()

	hasRows := false

    for rows.Next() {
        var id int
        var payload string
        err := rows.Scan(&id, &payload)
        if err != nil {
            log.Printf("failed to scan row: %v", err)
            continue
        }
        go processTask(id, payload)
    }

	if !hasRows {
		time.Sleep(10*time.Second)
	}

    if rows.Err() != nil {
        log.Printf("error during rows iteration: %v", rows.Err())
    }
}
```

## Параллельная обработка

SQL-конструкция `SELECT... FOR UPDATE SKIP LOCKED` позволяет забирать записи из таблицы, блокируя их для других транзакций, но пропуская уже заблокированные (SKIP LOCKED). Это важно при параллельной пакетной обработке задач, чтобы избежать гонок и дублирующей работы.


```go
ctx := context.Background()

for {
    rows, err := db.Query(ctx, `
        SELECT id, payload
        FROM tasks
        WHERE status = 'pending'
        ORDER BY id
        FOR UPDATE SKIP LOCKED
        LIMIT $1`, batchSize)
    if err != nil {
        log.Fatalf("failed to query tasks: %v", err)
    }
    defer rows.Close()

    for rows.Next() {
        var id int
        var payload string
        err := rows.Scan(&id, &payload)
        if err != nil {
            log.Printf("failed to scan row: %v", err)
            continue
        }
        processTask(id, payload)
    }

    if rows.Err() != nil {
        log.Printf("error during rows iteration: %v", rows.Err())
    }
}
```


Пример:

```
- Запись 0 <- w1
- Запись 1 <- w1
- Запись 2 <- w1
- Запись 3 <- w2
- Запись 4 <- w2
- Запись 5 <- w2
```


## gocron

```go
package main

import (
	"context"
	"log"
	"math/rand/v2"
	"time"

	"github.com/go-co-op/gocron"
)

func Repeat(fn func(context.Context) bool) func(context.Context) {
	return func(ctx context.Context) {
		repeat := true

		for repeat {
			select {
			case <-ctx.Done():
				return
			default:
				repeat = fn(ctx)
			}
		}
	}
}

func main() {
	ctx := context.Background()

	s := gocron.NewScheduler(time.UTC)

	s.Every(20*time.Second).Do(Repeat(func(ctx context.Context) bool {
		log.Println("Processor 1")
		return false
	}), ctx)

	s.Every(10*time.Second).Do(Repeat(func(ctx context.Context) bool {
		log.Println("Processor 2")
		if rand.N(2) == 0 {
			return false
		}

		return true
	}), ctx)

	s.Cron("0 0 * * *").Do(func() {
		log.Println("Daily cleanup at midnight")
	})

	s.StartBlocking()
}
```

## Outbox

Outbox — шаблон интеграции, при котором события для внешних систем (например, Kafka, RabbitMQ) сначала сохраняются в локальную таблицу БД, а затем асинхронно отправляются, гарантируя согласованность между БД и шиной сообщений.

```sql
CREATE TABLE outbox (
    id SERIAL PRIMARY KEY,
    aggregate_type TEXT,
    aggregate_id TEXT,
    event_type TEXT,
    payload JSONB,
    created_at TIMESTAMP DEFAULT now(),
    sent BOOLEAN DEFAULT false
);
```

```go
func saveEventToOutbox(tx *sql.Tx, event OutboxEvent) error {
	_, err := tx.Exec(`
		INSERT INTO outbox (aggregate_type, aggregate_id, event_type, payload)
		VALUES ($1, $2, $3, $4)
	`, event.Type, event.ID, event.Event, event.Payload)
	return err
}
```

```go
func sendOutboxEvents() {
	rows, _ := db.Query(`
		SELECT id, payload FROM outbox
		WHERE sent = false
		ORDER BY id
		FOR UPDATE SKIP LOCKED
		LIMIT 100
	`)
	for rows.Next() {
		var id int
		var payload []byte
		rows.Scan(&id, &payload)

		// Отправка события во внешний сервис
		if sendToKafka(payload) {
			db.Exec("UPDATE outbox SET sent = true WHERE id = $1", id)
		}
	}
}
```


## Удаление из БД батчами

```go
	for {
		res, err := pool.Exec(ctx, "DELETE FROM tasks WHERE status = 'completed' LIMIT $1", batchSize)
		if err != nil {
			log.Printf("Failed to delete batch: %v", err)
			break
		}
		rowsAffected += int(res.RowsAffected())
		if res.RowsAffected() < batchSize {
			break
		}
	}
```
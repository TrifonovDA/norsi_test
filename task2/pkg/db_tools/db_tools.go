package db_tools

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"simple_test/task2/config"
	"simple_test/task2/internal/handlers"
)

func NewConnection(ctx context.Context) *pgxpool.Pool {

	url := fmt.Sprintf("postgres://%v:%v@%v:%v/%v", config.BdCred.Username,
		config.BdCred.Password, config.BdCred.Host, config.BdCred.Port, config.BdCred.Database)

	dbpool, err := pgxpool.New(ctx, url)
	if err != nil {
		log.Printf("Unable to first connect to database_tools: %v\n", err)
		return nil
	}
	fmt.Println("Database connection is good!")
	return dbpool
}

const query_up_cache = "select id, title from notes;"

func Up_cache(ctx context.Context, cache *handlers.Simple_cache, pool *pgxpool.Pool) (err error) {
	rows, err := pool.Query(ctx, query_up_cache)
	defer rows.Close()
	if err != nil {
		log.Printf("There's error: %v", err)
		if pgErr, ok := err.(*pgconn.PgError); ok { //обработка ошибок бд
			newErr := fmt.Sprintf("SQL Error: %s, Detail: %s, Code: %s, SQLState: %%", pgErr.Message, pgErr.Detail, pgErr.Code, pgErr.SQLState())
			log.Println(newErr)
		}
		return err
	}
	for rows.Next() {
		var note handlers.Note
		err = rows.Scan(&note.ID, &note.Title)
		if err != nil {
			return err
		}
		cache.Insert(note)
	}
	return nil
}

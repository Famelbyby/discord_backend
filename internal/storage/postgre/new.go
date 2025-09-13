package postgre

import (
	"fmt"
	"idea-store-auth/configs/postgres"
	"os"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

const emptyValue = -1

func New() (*Storage, error) {
	pool, err := postgres.LoadPgxPool()
	//conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return &Storage{db: pool}, nil
}

func ParseIdsString(str string) ([]int64, error) {
	if len(str) == 0 {
		return []int64{}, nil
	}
	slice := strings.Split(str, " ")
	var ids []int64

	for _, i := range slice {
		val, err := strconv.ParseInt(i, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing id \"%v\"", i)
		}
		ids = append(ids, val)
	}
	return ids, nil
}

package repo

import (
	"book/internal/adapters/storage"
	"book/internal/config"
	"book/internal/logger"
	"context"
	"fmt"
	"testing"
)

func BenchmarkBookRepository_Book(b *testing.B) {
	config.MustLoad("C:\\Users\\Sanfo\\GolandProjects\\Arch\\.env")
	log := logger.Load()
	cfg := config.NewDbConfig()
	d := storage.NewDB(context.Background(), logger.Load(), cfg)
	r := NewRepository(log, d, nil)
	fmt.Println(r)

	_, err := r.Book(context.Background(), "56776e97-b0e0-11ef-82a7-74563c3486c4")
	if err != nil {
		b.Fatal(err)
	}
}

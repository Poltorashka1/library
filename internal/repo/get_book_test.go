package repo

import (
	"testing"
)

func BenchmarkBookRepository_Book(b *testing.B) {
	b.ResetTimer()

	// Запуск бенчмарка
	for i := 0; i < b.N; i++ {
		_, err := s.r.Book(s.ctx, "56776e97-b0e0-11ef-82a7-74563c3486c4")
		if err != nil {
			b.Fatal(err)
		}
	}
}

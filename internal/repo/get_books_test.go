package repo

import (
	"book/internal/entities"
	"testing"
)

func BenchmarkBookRepository_Books(b *testing.B) {
	b.ResetTimer()
	payload := &entities.BooksFilter{
		Start:      0,
		Stop:       10,
		BooksCount: 45,
	}
	// Запуск бенчмарка
	for i := 0; i < b.N; i++ {
		_, err := s.r.Books(s.ctx, payload)
		if err != nil {
			b.Fatal(err)
		}
	}
}

//func BenchmarkBookRepository_BooksV2(b *testing.B) {
//	b.ResetTimer()
//	payload := &entities.BooksFilter{
//		Start:      0,
//		Stop:       10,
//		BooksCount: 45,
//	}
//	// Запуск бенчмарка
//	for i := 0; i < b.N; i++ {
//		_, err := s.r.BooksV2(s.ctx, payload)
//		if err != nil {
//			b.Fatal(err)
//		}
//	}
//}

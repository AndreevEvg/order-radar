package fetcher 

import(
	"context"
	"order-radar/internal/models"
)

// Fetcher — интерфейс, который должен реализовать каждый адаптер площадки.
// Это позволяет легко добавлять новые источники без изменения основной логики.
type Fetcher interface {
	// Name возвращает человекочитаемое имя источника (для логов).
	Name() string 

	// Fetch получает список заказов с площадки.
	// Контекст позволяет отменять запросы (например, при graceful shutdown).
	Fetch(ctx context.Context) ([]models.Order, error)
}
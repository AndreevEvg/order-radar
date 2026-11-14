package fetcher

import (
	"context"
	"fmt"
	"time"

	"github.com/rtexty/gokwork/pkg/kwork"
	"order-radar/internal/models"
)

// KworkFetcher реализует Fetcher для kwork.ru через библиотеку GoKwork.
// GoKwork — асинхронная обёртка над закрытым API Kwork, поддерживающая
// аутентификацию, получение проектов и WebSocket для real-time событий[reference:6].
type KworkFetcher struct {
	client *kwork.Client
}

// NewKworkFetcher создаёт клиент Kwork с указанными учётными данными.
// proxyURL опционален; если задан, все запросы пойдут через прокси.
func NewKworkFetcher(login, password, phoneLast, proxyURL string) (*KworkFetcher, error) {
	cfg := kwork.Config {
		Login: login,
		Password: password,
		PhoneLast: phoneLast,
	}

	if proxyURL != "" {
		cfg.ProxyURL = proxyURL
	}

	client, err := kwork.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("kwork client: %w", err)
	}

	return &KworkFetcher{client: client}, nil
}

func (f *KworkFetcher) Name() string { return "kwork" }

// Fetch получает проекты (заказы) с Kwork.
// В GoKwork для этого есть метод получения проектов; ниже — типовой вызов.
// Точное имя метода уточните в документации библиотеки на актуальной версии.
func (f *KworkFetcher) Fetch(ctx context.Context) ([]models.Order, error) {
	// Пример: projects, err := f.client.GetProjects(ctx, ...)
	// Поскольку API Kwork закрытый и может меняться,
	// здесь приведён псевдокод с типовой сигнатурой.
	//
	// В реальной реализации используйте метод из GoKwork,
	// например GetProjects или GetExchangeProjects.
	projects, err := f.client.GetProjects(ctx, kwork.ProjectsParams{}) // <-- уточнить сигнатуру
	if err != nil {
		return nil, fmt.Errorf("kwork get projects: %w", err)
	}

	var orders []models.Order 
	for _, p := range projects {
		orders = append(orders, models.Order{
			Source: "kwork",
			ExternalID: fmt.Sprintf("%d", p.ID),
			Title: p.Title,
			Description: p.Description,
			Budget: string(p.Price),
			URL:          fmt.Sprintf("https://kwork.ru/project/%d", p.ID),
			PublishedAt:  time.Now(),
			FetchedAt:    time.Now(),
		})
	}
	
	return orders, nil
}
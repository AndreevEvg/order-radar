package fetcher 

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"order-radar/internal/models"
)

// FreelanceRuFetcher реализует Fetcher для freelance.ru через HTML-парсинг.
// API у площадки нет, поэтому используем goquery — jQuery-подобную библиотеку,
// стандарт де-факто для скрапинга на Go[reference:9].
//
// ⚠️ ВНИМАНИЕ: селекторы ниже основаны на типовой вёрстке Freelance.ru.
// При изменении вёрстки их нужно будет обновить.
type FreelanceRuFetcher struct {
	searchURL string 
}

// NewFreelanceRuFetcher создаёт адаптер.
// searchURL — страница поиска/категории с заказами.
func NewFreelanceRuFetcher(searchURL string) *FreelanceRuFetcher {
	if searchURL == "" {
		// Пример URL категории "Разработка" — уточните актуальный URL.
		searchURL = "https://freelance.ru/project/search?cat=11"
	}

	return &FreelanceRuFetcher{ searchURL: searchURL }
}

func (f *FreelanceRuFetcher) Name() string {
	return "freelanceru"
}

func (f *FreelanceRuFetcher) Fetch(ctx context.Context) ([]models.Order, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, f.searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("freelanceru request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; OrderRadar/1.0)")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("freelanceru fetch: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("freelanceru: status %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("freelanceru parse html: %w", err)
	}

	var orders []models.Order

	// ⚠️ Селекторы-ЗАГЛУШКИ. Замените на реальные, inspect элементов на freelance.ru.
	// Пример: карточки заказов могут быть в .project-item или .task-card.
	doc.Find(".project-item, .task-card").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find(".title, .project-title").First().Text())
		link, _ := s.Find("a").First().Attr("href")
		desc := strings.TrimSpace(s.Find(".description, .project-desc").First().Text())
		budget := strings.TrimSpace(s.Find(".budget, .price").First().Text())

		if title == "" || link == "" {
			return // пропускаем пустые карточки
		}

		// Делаем ссылку абсолютной
		if strings.HasPrefix(link, "/") {
			link = "https://freelance.ru" + link
		}

		// ID извлекаем из URL (последний сегмент) или используем сам URL.
		parts := strings.Split(strings.TrimRight(link, "/"), "/")
		externalID := parts[len(parts)-1]
		if externalID == "" {
			externalID = link 
		}

		orders = append(orders, models.Order{
			Source: "freelanceru",
			ExternalID: externalID,
			Title: title,
			Description: desc,
			Budget: budget,
			URL: link,
			PublishedAt: time.Now(), // точную дату из карточки извлечь сложнее
		})
	})

	return orders, nil 
}
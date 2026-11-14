package fetcher 

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"time"
	"order-radar/internal/models"
)

// FLRUFetcher реализует Fetcher для fl.ru через открытую RSS-ленту.
// Официального API у FL.ru нет, но RSS-лента стабильна и не требует авторизации[reference:8].
type FLRUFetcher struct {
	rssURL string 
}

// NewFLRUFetcher создаёт адаптер с указанным URL RSS-ленты.
// Можно фильтровать по категориям/подкатегориям, добавляя параметры к URL.
func NewFLRUFetcher(rssURL string) *FLRUFetcher {
	if rssURL == "" {
		rssURL = "https://www.fl.ru/rss/all.xml"
	}

	return &FLRUFetcher{ rssURL: rssURL }
}

func (f *FLRUFetcher) Name() string {
	return "flru"
}

// rssFeed — минимальная структура для разбора RSS.
// FL.ru отдаёт стандартный RSS 2.0.
type rssFeed struct {
	Channel struct {
		Items []struct {
			Title string `xml:"title"`
			Link string `xml:"link"`
			Description string `xml:"description"`
			PubDate string `xml:"pubDate"`
			GUID string `xml:"guid"`
		} `xml:"item"`
	} `xml:"channel"`
}

func (f *FLRUFetcher) Fetch(ctx context.Context) ([]models.Order, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, f.rssURL, nil)
	if err != nil {
		return nil, fmt.Errorf("flru request: %w", err)
	}

	// Устанавливаем User-Agent, чтобы выглядеть как обычный клиент.
	req.Header.Set("User-Agent", "OrderRadar/1.0 (+https://example.com)")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("flru fetch: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("flru: unexpected status %d", resp.StatusCode)
	}

	var feed rssFeed 
	if err := xml.NewDecoder(resp.Body).Decode(&feed); err != nil {
		return nil, fmt.Errorf("flru decode rss: %w", err)
	}

	var orders []models.Order 
	for _, item := range feed.Channel.Items {
		// Парсим дату публикации. Формат RFC1123Z обычно используется в RSS.
		pub, _ := time.Parse(time.RFC1123Z, item.PubDate)
		if pub.IsZero() {
			pub, _ = time.Parse(time.RFC1123, item.PubDate)
		}

		// GUID может отсутствовать; в качестве fallback используем Link.
		externalID := item.GUID 
		if externalID == "" {
			externalID = item.Link 
		}

		orders = append(orders, models.Order{
			Source: "flru",
			ExternalID: externalID,
			Title: item.Title,
			Description: item.Description,
			Budget: "", // RSS FL.ru обычно не содержит бюджета
			URL: item.Link,
			PublishedAt: pub,
		})
	}

	return orders, nil 
}
package models 

import "time"

// Order — единая нормализованная структура заказа,
// в которую каждый адаптер приводит данные со своей площадки.
type Order struct {
	// Source — идентификатор площадки: "kwork", "flru", "freelanceru".
	Source string `json:"source"`

	// ExternalID — уникальный ID заказа на площадке.
	// Используется для дедупликации.
	ExternalID string `json:"external_id"`

	Title string `json:"title"`
	Description string `json:"description"`
	Budget string `json:"budget"` // текстовое представление, т.к. у площадок разные форматы
	URL string `json:"url"`
	PublishedAt time.Time `json:"published_at"`

	// FetchedAt — когда заказ был получен ботом (для сортировки в UI).
	FetchedAt time.Time `json:"fetched_at"`
}
package entity

import "time"

type CreatePresignedPostUrlParams struct {
	Filename    string        `json:"filename"`
	ContentType string        `json:"content_type"`
	Ttl         time.Duration `json:"ttl"`
}

type PresignedPostUrl struct {
	Url         string        `json:"url"`
	Ttl         time.Duration `json:"ttl"`
	Key         string        `json:"key"`
	Id          string        `json:"id"`
	Filename    string        `json:"filename"`
	ContentType string        `json:"content_type"`
}

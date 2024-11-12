package book

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Book struct {
	ID            uuid.UUID `gorm:"primarykey"`
	Title         string
	Author        string
	PublishedDate time.Time
	ImageURL      string
	Description   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}

type Books []*Book

type DTO struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	PublishedDate string `json:"published_date"`
	ImageURL      string `json:"image_url"`
	Description   string `json:"description"`
}

type Form struct {
	Title         string `json:"title" validate:"required,max=255"`
	Author        string `json:"author" validate:"required,alphaspace,max=255"`
	PublishedDate string `json:"published_date" validate:"required,datetime=2006-01-02"`
	ImageURL      string `json:"image_url" validate:"url"`
	Description   string `json:"description"`
}

func (f *Form) ToModel() *Book {
	pubDate, _ := time.Parse("2006-01-02", f.PublishedDate)
	return &Book{
		Title:         f.Title,
		Author:        f.Author,
		PublishedDate: pubDate,
		ImageURL:      f.ImageURL,
		Description:   f.Description,
	}
}

func (b *Book) ToDTO() *DTO {
	return &DTO{
		ID:            b.ID.String(),
		Title:         b.Title,
		Author:        b.Author,
		PublishedDate: b.PublishedDate.Format("2006-01-02"),
		ImageURL:      b.ImageURL,
		Description:   b.Description,
	}
}

func (bs Books) ToDTO() []*DTO {
	dtos := make([]*DTO, len(bs))
	for _, b := range bs {
		dtos = append(dtos, b.ToDTO())
	}
	return dtos
}

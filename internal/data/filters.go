package data

import (
	"math"

	"github.com/Nana-Seyram/crest-countries/internal/validator"
)

type Filters struct {
	Page     int
	PageSize int
}

type Metadata struct {
	CurrentPage  int `json:"currentPage"`
	PageSize     int `json:"pageSize"`
	FirstPage    int `json:"firstPage"`
	LastPage     int `json:"lastPage"`
	TotalRecords int `json:"totalRecords"`
}

func ValidateFilters(v *validator.Validator, f Filters) {
	v.Check(f.Page != 0, "page", "must be greater than 0")
	v.Check(f.Page <= 250, "page", "must be a maximum of 250")
	v.Check(f.PageSize != 0, "page_size", "must be greater than 0")
	v.Check(f.PageSize <= 250, "page_size", "must be a maximum of 250")
}

func calculateMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}

func (f Filters) limit() int {
	return f.PageSize
}

func (f Filters) offset() int {
	return (f.Page - 1) * f.PageSize
}

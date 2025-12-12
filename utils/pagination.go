package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Total    int `json:"total"`
	Pages    int `json:"pages"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

func GetPaginationParams(c *gin.Context) (page, pageSize int) {
	// Default values
	page = 1
	pageSize = 1

	// parse query parameters
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	return page, pageSize
}

func CalculateOffset(page, pageSize int) int {
	return (page - 1) * pageSize
}

func CalculateTotalPages(total, pageSize int) int {
	if total == 0 {
		return 1
	}

	pages := total / pageSize
	if total%pageSize > 0 {
		pages++
	}
	return pages
}

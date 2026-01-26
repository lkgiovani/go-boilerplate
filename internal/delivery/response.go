package delivery

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lkgiovani/go-boilerplate/pkg/utils"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Count   int         `json:"count,omitempty"`
}

type ListRequest struct {
	Filters map[string]interface{} `json:"filters"`
	Page    int                    `json:"page"`
	Limit   int                    `json:"limit"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}

func (h *Handler) parsePagination(c *fiber.Ctx) *PaginatedResponse {
	page := c.Query("page", "1")
	limit := c.Query("limit", "10")

	var pageNum, limitNum int
	if p, err := utils.ParseInt(page); err == nil && p > 0 {
		pageNum = p
	} else {
		pageNum = 1
	}

	if l, err := utils.ParseInt(limit); err == nil && l > 0 {
		limitNum = l
	} else {
		limitNum = 10
	}

	return &PaginatedResponse{
		Page:  pageNum,
		Limit: limitNum,
	}
}

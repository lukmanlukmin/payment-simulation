// Package http ...
package http

import (
	"net/http"

	payload "payment-simulation/model/http_payload"

	"github.com/gofiber/fiber/v2"
)

// Transfer ...
// @Summary Create Transfer
// @Description Create Transfer With Hardcoded Merchant
// @Tags Transaction
// @Accept json
// @Produce json
// @Router /transfer [post]
// @Param payload body payload.TransferRequest true "Create Transaction"
// @Success 200 {object} payload.BaseResponse{data=payload.TransferResponse{}}
// @Response 400
// @Response 500
func (h *Handler) Transfer(c *fiber.Ctx) error {
	req := payload.TransferRequest{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(payload.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
	}
	if err := h.Validate.Struct(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(payload.BaseResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	// hardcoded merchant id
	req.MerchantID = 1

	res, err := h.Service.TransactionService.SubmitTransfer(c.Context(), req)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}
	return c.Status(http.StatusAccepted).JSON(payload.BaseResponse{
		Success: true,
		Data:    res,
	})
}

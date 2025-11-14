package handler

import (
	"errors"
	"net/http"

	"github.com/SwanPoi/bmstu_rsoi_lab2/src/rental/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *RentalHandler) GetUserRentals(ctx *gin.Context) {
	username := ctx.GetHeader("X-User-Name")
	if username == "" {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "X-User-Name header is required"})
		return
	}

	rentals, err := h.services.GetUserRentals(username)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, rentals)
}

/**
* Информация об оплате по идентификатору
 */
func (h *RentalHandler) GetUserRentalByUid(ctx *gin.Context) {
	username := ctx.GetHeader("X-User-Name")
	if username == "" {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "X-User-Name header is required"})
		return
	}

	rentalUid := ctx.Param("uid")

	if _, err := uuid.Parse(rentalUid); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "RentalUid must be valid"})
		return
	}

	rental, err := h.services.GetUserRentalByUid(rentalUid, username)

	if err != nil {
		if errors.Is(err, models.ErrorNotFound) {
			message := "Rental with rental_uid = " + rentalUid + " is not found"
			ctx.JSON(http.StatusNotFound, models.ErrorResponse{Message: message})
		} else if errors.Is(err, models.Forbidden) {
			message := "Rental with rental_uid = " + rentalUid + " is not for user with username " + username
			ctx.JSON(http.StatusNotFound, models.ErrorResponse{Message: message})
		} else {
			ctx.JSON(http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, rental)
}
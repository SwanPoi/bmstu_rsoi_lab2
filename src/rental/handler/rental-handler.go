package handler

import (
	"errors"
	"net/http"
	"time"

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

func (h *RentalHandler) CreateRental(ctx *gin.Context) {
	var req models.RentCreation

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "Bad Rental Creation body"})
		return
	}

	validationErr := models.ValidationErrorResponse{
		Message: "Validation Error",
		Errors: make(map[string]string),
	}

	if _, err := time.Parse(time.RFC3339, req.DateFrom); err != nil {
        validationErr.Errors["date-from"] = "Error with parsing time from date-from"
    }

    if _, err := time.Parse(time.RFC3339, req.DateTo); err != nil {
       	validationErr.Errors["date-to"] = "Error with parsing time from date-from"
    }

	if _, err := uuid.Parse(req.CarUID); err != nil {
		validationErr.Errors["car_uid"] = "Car Uid must be valid"
	}

	if _, err := uuid.Parse(req.PaymentUID); err != nil {
		validationErr.Errors["payment_uid"] = "Payment Uid must be valid"
	}

	if len(validationErr.Errors) != 0 {
		ctx.JSON(http.StatusBadRequest, validationErr)
		return
	}

	rental, err := h.services.CreateRental(req)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, rental)
}

func (h *RentalHandler) UpdateRental(ctx *gin.Context) {
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

	var req models.RentalUpsert

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "Bad Rental Upsert body"})
		return
	}

	if err := h.services.UpdateRental(req, rentalUid, username); err != nil {
		if err == models.InvalidStatus {
			ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
		} else if err == models.ErrorNotFound {
			message := "Rental with rental_uid = " + rentalUid + " is not found"
			ctx.JSON(http.StatusNotFound, models.ErrorResponse{Message: message})
		} else {
			ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
		}

		return
	}

	ctx.Status(http.StatusNoContent)
}
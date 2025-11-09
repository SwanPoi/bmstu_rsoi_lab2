package handler

import (
	"net/http"
	"strconv"

	"github.com/SwanPoi/bmstu_rsoi_lab2/src/car/models"
	"github.com/gin-gonic/gin"
)

func (h *CarHandler) GetCars(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "0")
	sizeStr := ctx.DefaultQuery("size", "1")
	showAll := ctx.Query("showAll") == "true"

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
		return
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
		return
	}

	if page < 0 {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "Page must be not less than 0"})
		return
	}

	if size < 1 || size > 100 {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "Size must be greater than 0 but smaller than 101"})
		return
	}

	carsResponse, err := h.services.GetCars(page, size, showAll)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, carsResponse)
}
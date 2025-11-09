package handler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/SwanPoi/bmstu_rsoi_lab2/src/gateway/models"
	"github.com/gin-gonic/gin"
)

func forwardRequest(c *gin.Context, method, targetURL string, headers map[string]string, body []byte) (int, []byte, http.Header, error) {
	if len(c.Request.URL.RawQuery) > 0 {
		targetURL = fmt.Sprintf("%s?%s", targetURL, c.Request.URL.RawQuery)
	}

	req, err := http.NewRequest(method, targetURL, bytes.NewReader(body))
	if err != nil {
		return 0, nil, nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	if c.Request.Header.Get("Content-Type") != "" {
		req.Header.Set("Content-Type", c.Request.Header.Get("Content-Type"))
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, resp.Header, err
	}

	return resp.StatusCode, respBody, resp.Header, nil
}

func (h *GatewayHandler) GetCars(ctx *gin.Context) {
	status, body, headers, err := forwardRequest(ctx, "GET", "http://cars:8070/cars", nil, nil)

	if err != nil {
		ctx.JSON(status, models.ErrorResponse{ Message: err.Error() })
		return
	}

	ctx.Data(status, headers.Get("Content-Type"), body)
}

func (h *GatewayHandler) GetUserRentals(ctx *gin.Context) {

}

func (h *GatewayHandler) GetRentalById(ctx *gin.Context) {

}

func (h *GatewayHandler) RentCar(ctx *gin.Context) {

}

func (h *GatewayHandler) FinishCarRent(ctx *gin.Context) {

}

func (h *GatewayHandler) RevokeRent(ctx *gin.Context) {

}
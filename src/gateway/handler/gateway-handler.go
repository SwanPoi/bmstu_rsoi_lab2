package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/SwanPoi/bmstu_rsoi_lab2/src/gateway/converters"
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
	status, body, headers, err := forwardRequest(ctx, "GET", "http://cars:8070/api/v1/cars", nil, nil)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, models.ErrorResponse{ Message: err.Error() })
		return
	}

	ctx.Data(status, headers.Get("Content-Type"), body)
}

func (h *GatewayHandler) GetUserRentals(ctx *gin.Context) {
	username := ctx.GetHeader("X-User-Name")
	if username == "" {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "X-User-Name header is required"})
		return
	}

	headers := map[string]string{"X-User-Name": username}

	// 1. Получить аренды
	status, body, _, err := forwardRequest(ctx, "GET", "http://rental/api/v1/rentals", headers, nil)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, models.ErrorResponse{Message: err.Error()})
		return
	}

	if status != http.StatusOK {
		ctx.Data(status, "application/json", body)
		return
	}

	var rentals []models.RentalInfo
	if err := json.Unmarshal(body, &rentals); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Rental parsing error"})
		return
	}

	carUIDs := make([]string, len(rentals))
	paymentUIDs := make([]string, len(rentals))

	for i, rental := range rentals {
		carUIDs[i] = rental.CarUID
		paymentUIDs[i] = rental.PaymentUID
	}
	// 2. Получить автомобили
	carUrl := "http://cars:8070/api/v1/cars/query"
	carsRequest := models.CarsRequest{ UIDs: carUIDs }
	carReqBody, _ := json.Marshal(carsRequest)

	carStatus, carBody, _, err := forwardRequest(ctx, "POST", carUrl, nil, carReqBody)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, models.ErrorResponse{Message: err.Error()})
		return
	}

	if carStatus != http.StatusOK {
		ctx.Data(carStatus, "application/json", carBody)
		return
	}

	var cars []models.CarInfo
	if err := json.Unmarshal(carBody, &cars); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Car parsing error"})
		return
	}
	
	// 3. Получить оплаты
	paymentUrl := "http://payment:8050/api/v1/payment/query"
	paymentsRequest := models.PaymentsRequest{ UIDs: paymentUIDs }
	paymentsReqBody, _ := json.Marshal(paymentsRequest)

	paymentStatus, paymentBody, _, err := forwardRequest(ctx, "POST", paymentUrl, nil, paymentsReqBody)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, models.ErrorResponse{Message: err.Error()})
		return
	}

	if paymentStatus != http.StatusOK {
		ctx.Data(paymentStatus, "application/json", paymentBody)
		return
	}

	var payments []models.PaymentInfo
	if err := json.Unmarshal(paymentBody, &payments); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Payment parsing error"})
		return
	}
	
	// 4. Смэтчить в массив RentalResponse
	carMap := make(map[string]models.CarInfo)
    for _, car := range cars {
        carMap[car.CarUID] = car
    }

    paymentMap := make(map[string]models.PaymentInfo)
    for _, payment := range payments {
        paymentMap[payment.PaymentUID] = payment
    }

	rentalsResponse := make([]models.RentalResponse, len(rentals))

	for i, rental := range rentals {
		rentalsResponse[i] = converters.ConvertToRentalResponse(rental, carMap[rental.CarUID], paymentMap[rental.PaymentUID])
	}

	ctx.JSON(http.StatusOK, rentalsResponse)
}

func (h *GatewayHandler) GetRentalById(ctx *gin.Context) {
	username := ctx.GetHeader("X-User-Name")
	if username == "" {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "X-User-Name header is required"})
		return
	}

	headers := map[string]string{"X-User-Name": username}

	rentalUid := ctx.Param("rentalUid")

	if rentalUid == "" {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "RentalUid is required"})
	}

	rentalUrl := "http://rental/api/v1/rental/" + rentalUid
	
	status, body, _, err := forwardRequest(ctx, "GET", rentalUrl, headers, nil)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, models.ErrorResponse{Message: err.Error()})
		return
	}

	if status != http.StatusOK {
		ctx.Data(status, "application/json", body)
		return
	}

	var rental models.RentalInfo
	if err := json.Unmarshal(body, &rental); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Rental parsing error"})
		return
	}

	// Получить авто
	carUrl := "http://cars:8070/api/v1/cars/" + rental.CarUID

	carStatus, carBody, _, err := forwardRequest(ctx, "GET", carUrl, nil, nil)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, models.ErrorResponse{Message: err.Error()})
		return
	}

	if carStatus != http.StatusOK {
		ctx.Data(carStatus, "application/json", carBody)
		return
	}

	var car models.CarInfo
	if err := json.Unmarshal(carBody, &car); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Car parsing error"})
		return
	}
	
	// Получить оплату
	paymentUrl := "http://payment:8050/api/v1/payment/" + rental.PaymentUID

	paymentStatus, paymentBody, _, err := forwardRequest(ctx, "GET", paymentUrl, nil, nil)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, models.ErrorResponse{Message: err.Error()})
		return
	}

	if paymentStatus != http.StatusOK {
		ctx.Data(paymentStatus, "application/json", paymentBody)
		return
	}

	var payment models.PaymentInfo
	if err := json.Unmarshal(paymentBody, &payment); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Payment parsing error"})
		return
	}

	response := converters.ConvertToRentalResponse(rental, car, payment)

	ctx.JSON(http.StatusOK, response)

}

func (h *GatewayHandler) RentCar(ctx *gin.Context) {
	username := ctx.GetHeader("X-User-Name")
	if username == "" {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "X-User-Name header is required"})
		return
	}

	bodyBytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "Fail during reading of request body for car rent"})
		return
	}

	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var rentReq models.RentCreationRequest
	if err := json.Unmarshal(bodyBytes, &rentReq); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Rent request parsing error"})
		return
	}

	carStatusUpsert := models.CarStatusUpsert{
		Availability: false,
	}

	carStatusBytes, err := json.Marshal(carStatusUpsert)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Rent request parsing error"})
		return
	}

	carUrl := "http://cars:8070/api/v1/car/" + rentReq.CarUID

	status, body, _, err := forwardRequest(ctx, "PATCH", carUrl, nil, carStatusBytes)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, models.ErrorResponse{Message: err.Error()})
		return
	}

	if status != http.StatusOK {
		ctx.Data(status, "application/json", body)
		return
	}

	payCreateReq := models.PaymentCreateRequest{
		DateFrom: rentReq.DateFrom,
		DateTo: rentReq.DateTo,
	}

	payCreateBytes, err := json.Marshal(payCreateReq)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Payment Creation request marshaling error"})
		return
	}

	payStatus, payBody, _, err := forwardRequest(ctx, "POST", "http://payment/api/v1/payment", nil, payCreateBytes)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, models.ErrorResponse{Message: err.Error()})
		return
	}

	if payStatus != http.StatusOK {
		ctx.Data(payStatus, "application/json", payBody)
		return
	}

	var paymentResponse models.PaymentCreationResponse

	if err := json.Unmarshal(payBody, &paymentResponse); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Payment response parsing error"})
		return
	}

	rentCreation := models.RentCreation{
		DateFrom: rentReq.DateFrom,
		DateTo: rentReq.DateTo,
		CarUID: rentReq.CarUID,
		PaymentUID: paymentResponse.PaymentUID,
		Username: username,
	}

	rentBytes, err := json.Marshal(rentCreation)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Rental Creation request marshaling error"})
		return
	}

	rentStatus, rentBody, _, err := forwardRequest(ctx, "POST", "http://rental/api/v1/rental", nil, rentBytes)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, models.ErrorResponse{Message: err.Error()})
		return
	}

	if rentStatus != http.StatusOK {
		ctx.Data(rentStatus, "application/json", rentBody)
		return
	}

	var rentalCreationResponse models.RentalInfo

	if err := json.Unmarshal(rentBody, &rentalCreationResponse); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Rental Creation response parsing error"})
		return
	}

	rentResponse := converters.ConvertToCreateRentalResponse(rentalCreationResponse, paymentResponse)

	ctx.JSON(http.StatusOK, rentResponse)
}

func (h *GatewayHandler) FinishCarRent(ctx *gin.Context) {

}

func (h *GatewayHandler) RevokeRent(ctx *gin.Context) {

}
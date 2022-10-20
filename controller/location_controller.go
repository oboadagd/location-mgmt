// Package controller implements api layer for external clients.
// Through implementation of the LocationControllerInterface methods,
// it is possible to define validation and management of parameters
// it also invokes service layer.
package controller

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	respKit "github.com/oboadagd/kit-go/middleware/responses"
	"github.com/oboadagd/location-common/dto"
	"github.com/oboadagd/location-common/enums"
	grpcclient "github.com/oboadagd/location-mgmt/userlocation/client"
	"strconv"
	"strings"
)

const (
	defaultPage       = 1  // default page
	defaultItemsLimit = 10 // default quantity of items per page
)

// LocationControllerInterface is the interface of Location controller layer. Contains definition of
// methods to manage the microservice api.
type LocationControllerInterface interface {
	LocationControllerHelperInterface
	Save(c echo.Context) error
	GetUsersByLocationAndRadius(c echo.Context) error
}

// LocationControllerHelperInterface is the interface helper of Location controller layer. Contains definition of
// methods to validate and manage parameters of microservices api.
type LocationControllerHelperInterface interface {
	SaveValidate(c echo.Context) (*dto.SaveLocationRequest, error)
	GetUsersByLocationAndRadiusValidate(c echo.Context) (*dto.GetUsersByLocationAndRadiusRequest, error)
}

// LocationController represents the Location controller layer.
type LocationController struct{}

// NewLocationController initializes Location controller layer.
func NewLocationController() LocationControllerInterface {
	return &LocationController{}
}

// SaveValidate is a helper method that validates parameters of Save service
func (ctr *LocationController) SaveValidate(c echo.Context) (*dto.SaveLocationRequest, error) {

	req := dto.SaveLocationRequest{}
	if err := c.Bind(&req); err != nil {
		return &req, respKit.GenericBadRequestError(enums.ErrorRequestBodyCode, err.Error())
	}

	vtr := validator.New()
	if err := vtr.RegisterValidation("patternazAZ09", dto.IsPatternUserName); err != nil {
		return &req, respKit.GenericBadRequestError(enums.ErrorRequestBodyCode, err.Error())
	}

	if err := vtr.RegisterValidation("maxDecimals", dto.IsMaxDecimals); err != nil {
		return &req, respKit.GenericBadRequestError(enums.ErrorRequestBodyCode, err.Error())
	}

	cvt := &dto.CustomValidatorSaveLoc{Validator: vtr}

	if err := cvt.Validate(req); err != nil {
		return &req, respKit.GenericBadRequestError(enums.ErrorRequestBodyCode, err.Error())
	}

	return &req, nil
}

// Save invoke remote functionality of Save Location model through a grpc client.
func (ctr *LocationController) Save(c echo.Context) error {

	req, err := ctr.SaveValidate(c)

	if err != nil {
		return err
	}

	return grpcclient.SaveLocation(c, *req)
}

// GetUsersByLocationAndRadiusValidate is a helper method that validates parameters of GetUsersByLocationAndRadius service
func (ctr *LocationController) GetUsersByLocationAndRadiusValidate(c echo.Context) (*dto.GetUsersByLocationAndRadiusRequest, error) {
	var lt, lg, rd float64
	var page, itemsLimit uint64
	req := dto.GetUsersByLocationAndRadiusRequest{}

	if f, err := strconv.ParseFloat(c.Param("latitude"), 64); err != nil {
		return &req, respKit.GenericBadRequestError(enums.ErrorRequestBodyCode, err.Error())
	} else {
		lt = f
	}

	if f, err := strconv.ParseFloat(c.Param("longitude"), 64); err != nil {
		return &req, respKit.GenericBadRequestError(enums.ErrorRequestBodyCode, err.Error())
	} else {
		lg = f
	}

	if f, err := strconv.ParseFloat(c.Param("radius"), 64); err != nil {
		return &req, respKit.GenericBadRequestError(enums.ErrorRequestBodyCode, err.Error())
	} else {
		rd = f
	}

	if strings.TrimSpace(c.QueryParam("page")) == "" {
		page = defaultPage
	} else if ui, err := strconv.ParseUint(c.QueryParam("page"), 10, 64); err != nil {
		return &req, respKit.GenericBadRequestError(enums.ErrorRequestBodyCode, err.Error())
	} else if ui == 0 {
		itemsLimit = defaultPage
	} else {
		page = ui
	}

	if strings.TrimSpace(c.QueryParam("itemsLimit")) == "" {
		itemsLimit = defaultItemsLimit
	} else if ui, err := strconv.ParseUint(c.QueryParam("itemsLimit"), 10, 64); err != nil {
		return &req, respKit.GenericBadRequestError(enums.ErrorRequestBodyCode, err.Error())
	} else if ui == 0 {
		itemsLimit = defaultItemsLimit
	} else {
		itemsLimit = ui
	}

	req = dto.GetUsersByLocationAndRadiusRequest{
		Latitude:   lt,
		Longitude:  lg,
		Radius:     rd,
		Page:       page,
		ItemsLimit: itemsLimit,
	}

	vtr := validator.New()
	if err := vtr.RegisterValidation("patternazAZ09", dto.IsPatternUserName); err != nil {
		return &req, respKit.GenericBadRequestError(enums.ErrorRequestBodyCode, err.Error())
	}

	if err := vtr.RegisterValidation("maxDecimals", dto.IsMaxDecimals); err != nil {
		return &req, respKit.GenericBadRequestError(enums.ErrorRequestBodyCode, err.Error())
	}

	cvt := &dto.CustomValidatorSaveLoc{Validator: vtr}

	if err := cvt.Validate(req); err != nil {
		return &req, respKit.GenericBadRequestError(enums.ErrorRequestBodyCode, err.Error())
	}

	return &req, nil
}

// GetUsersByLocationAndRadius invoke remote functionality of GetUsersByLocationAndRadius through a grpc client.
func (ctr *LocationController) GetUsersByLocationAndRadius(c echo.Context) error {
	req, err := ctr.GetUsersByLocationAndRadiusValidate(c)

	if err != nil {
		return err
	}

	return grpcclient.GetUsersByLocationAndRadius(c, *req)
}

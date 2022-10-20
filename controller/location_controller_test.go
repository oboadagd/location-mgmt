package controller

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/oboadagd/location-common/dto"
	"github.com/oboadagd/location-mgmt/testutils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSaveValidate2(t *testing.T) {
	nameTest := "TestSaveValidate"

	var err error
	var b []byte
	var req *http.Request
	var rec *httptest.ResponseRecorder
	var ctx echo.Context

	locationController := NewLocationController()
	e := echo.New()

	type test struct {
		dataStr        string
		dataNumb       []float64
		resultValidate []string
		answer         string
	}

	tests := []test{
		{"usernamesample", []float64{10, 10}, []string{""}, "success"},
		{"", []float64{10, 10}, []string{"username", "required"}, "userName required failed"},
		{"usr", []float64{10, 10}, []string{"username", "min"}, "userName min failed"},
		{"username123456789", []float64{10, 10}, []string{"username", "max"}, "userName max failed"},
		{"username_1", []float64{10, 10}, []string{"username", "pattern"}, "userName pattern failed"},
		{"usernamesample", []float64{-91, 10}, []string{"latitude", "min"}, "latitude min failed"},
		{"usernamesample", []float64{91, 10}, []string{"latitude", "max"}, "latitude max failed"},
		{"usernamesample", []float64{10.12345678, 10}, []string{"latitude", "maxdecimals"}, "latitude max decimals failed"},
		{"usernamesample", []float64{10, -181}, []string{"longitude", "min"}, "longitude min failed"},
		{"usernamesample", []float64{10, 181}, []string{"longitude", "max"}, "longitude max failed"},
		{"usernamesample", []float64{10, 10.12345678}, []string{"longitude", "maxdecimals"}, "longitude max decimals failed"},
	}

	for _, v := range tests {
		l := dto.SaveLocationRequest{
			UserName:  v.dataStr,
			Latitude:  v.dataNumb[0],
			Longitude: v.dataNumb[1],
		}
		b, err = json.Marshal(l)

		req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(b)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec = httptest.NewRecorder()
		ctx = e.NewContext(req, rec)
		ctx.SetPath("/location-history-mmg/locations")

		_, err = locationController.SaveValidate(ctx)

		if err != nil && testutils.EvaluateErrConditions(err.Error(), v.resultValidate) {
			t.Errorf("%s: Expected %v but got %v", nameTest, v.answer, err.Error())
			return
		}
	}

	t.Logf("%s Success", nameTest)
}

func TestGetUsersByLocationAndRadiusValidate(t *testing.T) {
	nameTest := "TestGetUsersByLocationAndRadiusValidate"
	var err error

	locationController := NewLocationController()

	type test struct {
		data           []string
		resultValidate []string
		answer         string
	}

	tests := []test{
		{[]string{"10", "10", "10"}, []string{""}, "success"},
		{[]string{"10", "10", "0"}, []string{"radius", "gt"}, "radius greater than zero failed"},
		{[]string{"10", "10", ""}, []string{"parsefloat"}, "radius parse float failed"},
		{[]string{"", "10", "10"}, []string{"parsefloat"}, "latitude parse float failed"},
		{[]string{"-91", "10", "10"}, []string{"latitude", "min"}, "latitude min failed"},
		{[]string{"91", "10", "10"}, []string{"latitude", "max"}, "latitude max failed"},
		{[]string{"10.123456789", "10", "10"}, []string{"latitude", "maxdecimals"}, "latitude max decimals failed"},
		{[]string{"10", "", "10"}, []string{"parsefloat"}, "longitude parse float failed"},
		{[]string{"10", "-181", "10"}, []string{"longitude", "min"}, "longitude min failed"},
		{[]string{"10", "181", "10"}, []string{"longitude", "max"}, "longitude max failed"},
		{[]string{"10", "10.123456789", "10"}, []string{"longitude", "maxdecimals"}, "longitude max decimals failed"},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/location-history-mmg/locations/users")

	for _, v := range tests {
		ctx.SetParamNames("latitude", "longitude", "radius")
		ctx.SetParamValues(v.data[0], v.data[1], v.data[2])

		_, err = locationController.GetUsersByLocationAndRadiusValidate(ctx)

		if err != nil && testutils.EvaluateErrConditions(err.Error(), v.resultValidate) {
			t.Errorf("%s: Expected %v but got %v", nameTest, v.answer, err.Error())
			return
		}
	}

	t.Logf("%s Success", nameTest)
}

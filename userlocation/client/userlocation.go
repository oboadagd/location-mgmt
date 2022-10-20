package grpcclient

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/oboadagd/location-common/dto"
	"net/http"

	pb "github.com/oboadagd/location-mgmt/userlocation/proto"
)

// doSaveLocation implements the response management sent from
// SaveLocation.
func doSaveLocation(ctx echo.Context, c pb.UserLocationServiceClient, req *pb.SaveLocationRequest) error {
	log.Info("doSaveLocation was invoked")
	resp, err := c.SaveLocation(context.Background(), req)

	if err != nil {
		log.Fatalf("doSaveLocation error, %v", err)
		return err
	}

	log.Infof("doSaveLocation finished, %v", resp)
	return ctx.JSON(http.StatusCreated, dto.Response{
		Message: resp.Message,
	})
}

// doGetUsersByLocationAndRadius implements the response management sent from
// GetUsersByLocationAndRadius.
func doGetUsersByLocationAndRadius(ctx echo.Context, c pb.UserLocationServiceClient, req *pb.GetUsersByLocationAndRadiusRequest) error {
	log.Info("doGetUsersByLocationAndRadius was invoked")

	pbResp, err := c.GetUsersByLocationAndRadius(context.Background(), req)

	if err != nil {
		log.Fatalf("doGetUsersByLocationAndRadius error, %v", err)
		return err
	}

	log.Infof("doGetUsersByLocationAndRadius finished, %v", pbResp)
	return ctx.JSON(http.StatusOK, pbResp)
}

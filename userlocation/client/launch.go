// Package grpcclient implements routines to reach grpc server out, wait the response,
// then manage the response.
package grpcclient

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/oboadagd/location-common/dto"
	"github.com/oboadagd/location-common/recordtype"
	pb "github.com/oboadagd/location-mgmt/userlocation/proto"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var addr string = ":50061" // url grpc server connection

// SaveLocation implements grpc client connection with grpc server for
// saving the Location model.
func SaveLocation(ctx echo.Context, req dto.SaveLocationRequest) error {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if addr == ":50061" {
		if err := envconfig.Process("LIST", &recordtype.Cfg); err != nil {
			err = errors.Wrap(err, "parse environment variables")
			return err
		}
		addr = fmt.Sprintf("%s%s", recordtype.Cfg.GRPCHost, addr)
	}

	conn, err := grpc.Dial(addr, opts...)

	if err != nil {
		log.Error(err)
		return err
	}

	defer conn.Close()
	c := pb.NewUserLocationServiceClient(conn)

	r := pb.SaveLocationRequest{
		UserName:  req.UserName,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	}

	return doSaveLocation(ctx, c, &r)
}

// GetUsersByLocationAndRadius implements grpc client connection with grpc server for
// getting a list of Location models in a radius from a given location point.
func GetUsersByLocationAndRadius(ctx echo.Context, req dto.GetUsersByLocationAndRadiusRequest) error {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(addr, opts...)

	if err != nil {
		log.Error(err)
		return err
	}

	defer conn.Close()
	c := pb.NewUserLocationServiceClient(conn)

	r := pb.GetUsersByLocationAndRadiusRequest{
		Latitude:   req.Latitude,
		Longitude:  req.Longitude,
		Radius:     req.Radius,
		Page:       req.Page,
		ItemsLimit: req.ItemsLimit,
	}

	return doGetUsersByLocationAndRadius(ctx, c, &r)
}

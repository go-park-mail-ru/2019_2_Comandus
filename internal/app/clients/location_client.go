package clients

import (
	"context"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/location/delivery/grpc/location_grpc"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type LocationClient struct {
	conn *grpc.ClientConn
}

func (c *LocationClient) Connect() error {
	conn, err := grpc.Dial(LOCATION_PORT, grpc.WithInsecure())
	if err != nil {
		return errors.Wrap(err, "grpc.Dial()")
	}
	c.conn = conn
	return nil
}

func (c *LocationClient) GetCountry(id int64) (*location_grpc.Country, error) {
	client := location_grpc.NewLocationHandlerClient(c.conn)
	req := &location_grpc.CountryID{
		ID: id,
	}

	country, err := client.GetCountry(context.Background(), req)
	if err != nil {
		return nil, errors.Wrap(err, "client.GetCountry()")
	}

	return country, nil
}

func (c *LocationClient) GetCity(id int64) (*location_grpc.City, error) {
	client := location_grpc.NewLocationHandlerClient(c.conn)
	req := &location_grpc.CityID{
		ID: id,
	}

	city, err := client.GetCity(context.Background(), req)
	if err != nil {
		return nil, errors.Wrap(err, "client.GetCountry()")
	}

	return city, nil
}

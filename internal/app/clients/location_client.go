package clients

import (
	"context"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/location/delivery/grpc/location_grpc"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"log"
)

func GetCountry(id int64) (*location_grpc.Country, error) {
	conn, err := grpc.Dial(":8087", grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrap(err, "grpc.Dial()")
	}

	defer func(){
		if err := conn.Close(); err != nil {
			// TODO: use zap logger
			log.Println("conn.Close()", err)
		}
	}()

	client := location_grpc.NewLocationHandlerClient(conn)
	req := &location_grpc.CountryID{
		ID:		id,
	}

	country, err := client.GetCountry(context.Background(), req)
	if err != nil {
		return nil, errors.Wrap(err, "client.GetCountry()")
	}

	return country, nil
}

func GetCity(id int64) (*location_grpc.City, error) {
	conn, err := grpc.Dial(":8087", grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrap(err, "grpc.Dial()")
	}

	defer func(){
		if err := conn.Close(); err != nil {
			// TODO: use zap logger
			log.Println("conn.Close()", err)
		}
	}()

	client := location_grpc.NewLocationHandlerClient(conn)
	req := &location_grpc.CityID{
		ID:		id,
	}

	city, err := client.GetCity(context.Background(), req)
	if err != nil {
		return nil, errors.Wrap(err, "client.GetCountry()")
	}

	return city, nil
}

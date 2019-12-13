package logrpc

import (
	"context"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/location"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/location/delivery/grpc/location_grpc"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type LocationServer struct {
	Ucase location.Usecase
}

func NewLocationServerGrpc(gserver *grpc.Server, ucase location.Usecase) {
	server := &LocationServer{
		Ucase: ucase,
	}
	location_grpc.RegisterLocationHandlerServer(gserver, server)
	reflection.Register(gserver)
}

func (s *LocationServer) TransformCountryRPC(location *model.Country) *location_grpc.Country {
	if location == nil {
		return nil
	}

	res := &location_grpc.Country{
		ID:                   location.ID,
		Name:                 location.Name,
	}
	return res
}

func (s *LocationServer) TransformCityRPC(location *model.City) *location_grpc.City {
	if location == nil {
		return nil
	}

	res := &location_grpc.City{
		ID:                   location.ID,
		CountryID:            location.CountryID,
		Name:                 location.Name,
	}
	return res
}

func (s *LocationServer) TransformCountryData(location *location_grpc.Country) *model.Country {
	res := &model.Country{
		ID:   location.ID,
		Name: location.Name,
	}
	return res
}

func (s *LocationServer) TransformCityData(location *location_grpc.City) *model.City {
	res := &model.City{
		ID:        location.ID,
		CountryID: location.CountryID,
		Name:      location.Name,
	}
	return res
}

func (s *LocationServer) GetCountry(context context.Context,req *location_grpc.CountryID) (*location_grpc.Country, error) {
	country, err := s.Ucase.GetCountry(req.ID)
	if err != nil {
		return nil, err
	}
	res := s.TransformCountryRPC(country)
	return res, nil
}

func (s *LocationServer) GetCity(context context.Context, req *location_grpc.CityID) (*location_grpc.City, error) {
	city, err := s.Ucase.GetCity(req.ID)
	if err != nil {
		return nil, err
	}
	res := s.TransformCityRPC(city)
	return res, nil
}
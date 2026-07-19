package model

import "github.com/kotafan1rich/geo_logic_api_go/internal/errors"

type InfraObject struct {
	ID       uint64
	GeoPoint GeoPoint
	Address  string
	Name     *string

	Type InfraType
}

func NewInfraObject(geopoint GeoPoint, address string, name *string, infraType InfraType) (*InfraObject, error) {
	if address == "" {
		return nil, errors.ErrInvalidAddress
	}
	return &InfraObject{
		GeoPoint: geopoint,
		Address:  address,
		Name:     name,
		Type:     infraType,
	}, nil
}

func (o *InfraObject) UpdateAddress(address string) error {
	if address == "" {
		return errors.ErrInvalidAddress
	}
	o.Address = address
	return nil
}

func (o *InfraObject) UpdateName(name *string) {
	o.Name = name
}

func (o *InfraObject) UpdateType(infraType InfraType) {
	o.Type = infraType
}

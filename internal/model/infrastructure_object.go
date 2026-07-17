package model

import "github.com/kotafan1rich/geo_logic_api_go/internal/errors"

type InfrastructureObject struct {
	ID       uint64
	GeoPoint GeoPoint
	Address  string
	Name     *string

	Type InfrastructureType
}

func NewInfrastructureObject(geopoint GeoPoint, address string, name *string, infraType InfrastructureType) (*InfrastructureObject, error) {
	if address == "" {
		return nil, errors.ErrInvalidAddress
	}
	return &InfrastructureObject{
		GeoPoint: geopoint,
		Address:  address,
		Name:     name,
		Type:     infraType,
	}, nil
}

func (o *InfrastructureObject) UpdateAddress(address string) error {
	if address == "" {
		return errors.ErrInvalidAddress
	}
	o.Address = address
	return nil
}

func (o *InfrastructureObject) UpdateName(name *string) {
	o.Name = name
}

func (o *InfrastructureObject) UpdateType(infraType InfrastructureType) {
	o.Type = infraType
}

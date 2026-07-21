package model

import (
	"github.com/kotafan1rich/geo_logic_api_go/internal/errors"
)

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

func (o *InfraObject) UpdateTypeID(typeID uint64) {
	o.Type.ID = typeID
}

func (o *InfraObject) Update(lat, long *float64, address, name *string, typeID *uint64) error {
	if lat != nil {
		err := o.GeoPoint.UpdateLat(*lat)
		if err != nil {
			return err
		}
	}
	if long != nil {
		err := o.GeoPoint.UpdateLong(*long)
		if err != nil {
			return err
		}
	}

	if address != nil {
		err := o.UpdateAddress(*address)
		if err != nil {
			return err
		}
	}

	if name != nil {
		o.UpdateName(name)
	}

	if typeID != nil {
		o.UpdateTypeID(*typeID)
	}
	return nil
}

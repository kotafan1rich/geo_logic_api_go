package model

import "github.com/kotafan1rich/geo_logic_api_go/internal/errors"

type InfrastructureType struct {
	ID        uint64
	Slug      string
	Name      string
	Weight    float64
	MaxRadius uint16
}

func NewInfrastructureType(slug, name string, weight float64, maxRadius uint16) (*InfrastructureType, error) {
	if slug == "" {
		return nil, errors.ErrInvalidSlug
	}
	if name == "" {
		return nil, errors.ErrInvalidName
	}
	if weight <= 0 {
		return nil, errors.ErrInvalidWeight
	}
	return &InfrastructureType{
		Slug:      slug,
		Name:      name,
		Weight:    weight,
		MaxRadius: maxRadius,
	}, nil
}

func (t *InfrastructureType) UpdateSlug(slug string) error {
	if slug == "" {
		return errors.ErrInvalidSlug
	}
	t.Slug = slug
	return nil
}

func (t *InfrastructureType) UpdateName(name string) error {
	if name == "" {
		return errors.ErrInvalidName
	}
	t.Name = name
	return nil
}

func (t *InfrastructureType) UpdateWeight(weight float64) error {
	if weight <= 0 {
		return errors.ErrInvalidWeight
	}
	t.Weight = weight
	return nil
}

func (t *InfrastructureType) UpdateMaxRadius(maxRadius uint16) {
	t.MaxRadius = maxRadius
}

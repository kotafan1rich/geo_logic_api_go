package model

import "github.com/kotafan1rich/geo_logic_api_go/internal/errors"

type InfraType struct {
	ID        uint64
	Slug      string
	Name      string
	Weight    float64
	MaxRadius uint16
}

func NewInfraType(slug, name string, weight float64, maxRadius uint16) (*InfraType, error) {
	if slug == "" {
		return nil, errors.ErrInvalidSlug
	}
	if name == "" {
		return nil, errors.ErrInvalidName
	}
	if weight <= 0 {
		return nil, errors.ErrInvalidWeight
	}
	return &InfraType{
		Slug:      slug,
		Name:      name,
		Weight:    weight,
		MaxRadius: maxRadius,
	}, nil
}

func (t *InfraType) UpdateSlug(slug string) error {
	if slug == "" {
		return errors.ErrInvalidSlug
	}
	t.Slug = slug
	return nil
}

func (t *InfraType) UpdateName(name string) error {
	if name == "" {
		return errors.ErrInvalidName
	}
	t.Name = name
	return nil
}

func (t *InfraType) UpdateWeight(weight float64) error {
	if weight <= 0 {
		return errors.ErrInvalidWeight
	}
	t.Weight = weight
	return nil
}

func (t *InfraType) UpdateMaxRadius(maxRadius uint16) {
	t.MaxRadius = maxRadius
}

func (t *InfraType) Update(slug, name *string, weight *float64, maxRadius *uint16) error {
	if slug != nil {
		err := t.UpdateSlug(*slug)
		if err != nil {
			return err
		}
	}

	if name != nil {
		err := t.UpdateName(*name)
		if err != nil {
			return err
		}
	}

	if weight != nil {
		err := t.UpdateWeight(*weight)
		if err != nil {
			return err
		}
	}

	if maxRadius != nil {
		t.UpdateMaxRadius(*maxRadius)
	}
	return nil
}

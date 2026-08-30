package model

type TrackedLocation struct {
	ID       uint64
	UserID   uint64
	GeoPoint GeoPoint
}

func NewTrackedLocation(userID uint64, geopoint *GeoPoint) *TrackedLocation {
	return &TrackedLocation{
		UserID:   userID,
		GeoPoint: *geopoint,
	}
}

func (t *TrackedLocation) UpdateUserID(userID uint64) {
	t.UserID = userID
}

func (t *TrackedLocation) UpdateGeoPoint(geopoint GeoPoint) {
	t.GeoPoint = geopoint
}

func (t *TrackedLocation) Update(userID *uint64, lat, long *float64) error {
	if userID != nil {
		t.UpdateUserID(*userID)
	}

	if lat != nil {
		err := t.GeoPoint.UpdateLat(*lat)
		if err != nil {
			return err
		}
	}

	if long != nil {
		err := t.GeoPoint.UpdateLong(*long)
		if err != nil {
			return err
		}
	}
	return nil
}

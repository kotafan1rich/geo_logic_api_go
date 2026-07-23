package model

type TrackedLocation struct {
	ID       uint64
	UserID   uint64
	GeoPoint GeoPoint
}

func NewTrackedLocation(userID uint64, geopoint GeoPoint) *TrackedLocation {
	return &TrackedLocation{
		UserID:   userID,
		GeoPoint: geopoint,
	}
}

func (t *TrackedLocation) UpdateUserID(userID uint64) {
	t.UserID = userID
}

func (t *TrackedLocation) UpdateGeoPoint(geopoint GeoPoint) {
	t.GeoPoint = geopoint
}

func (t *TrackedLocation) Update(userID *uint64, geopoint *GeoPoint) {
	if userID != nil {
		t.UpdateUserID(*userID)
	}
	if geopoint != nil {
		t.UpdateGeoPoint(*geopoint)
	}
}

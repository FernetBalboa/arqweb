package domain

type POIFilter struct {
	Title    string  `form:"title"`
	Category string  `form:"category"`
	Limit    int64   `form:"limit"`
	MaxLat   float64 `form:"north_lat"`
	MaxLong  float64 `form:"east_long"`
	MinLat   float64 `form:"south_lat"`
	MinLong  float64 `form:"west_long"`
	Bound    bool    `form:"bound"`
	Hidden   bool    `form:"hidden"`
}

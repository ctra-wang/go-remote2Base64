package r2b

type PointInfo struct {
	X, Y, With float64
	Align      int `json:"align" default:"1|2|3"`
	Info       string
}

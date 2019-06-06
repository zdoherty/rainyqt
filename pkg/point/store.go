package point

type PointStore interface {
	All() ([]Point, error)
	GetByRadius(point Point, radius Distance) ([]Point, error)
	GetByName(name string) (Point, error)
	Put(Point) error
	Delete(Point) error
}

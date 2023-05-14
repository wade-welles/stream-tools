package obs

// TODO: Probably will eventually need a grid to make local checks without
// contacting OBS by simply modeling it well on our side (and makes it easier to
// get rid of OBS later)

type Position struct {
	X, Y float64
}

type Rectangle struct {
	Position
	Height, Width float64
}

// TODO:
//   Maybe perimeter, cross length, overlaps with

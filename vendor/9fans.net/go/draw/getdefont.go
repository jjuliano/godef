package draw

import "bytes"

func getdefont(d *Display) (*subfont, error) {
	return d.readSubfont("*default*", bytes.NewReader(defontdata), nil)
}

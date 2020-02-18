package main

import "golang.org/x/tour/pic"

func Pic(dx, dy int) [][]uint8 {


	v := make([][]uint8, dy)

	for i := 0; i < len(v); i++ {

		v[i] = make([]uint8,dx)

	}

	for i := 0; i < len(v); i++ {

		for j := 0; j < len(v[i]); j++ {

			if i < j {
				v[i][j] = 125
			}	else{
				v[i][j] = 255
			}
		}
	}

	return v
}

func main() {
	pic.Show(Pic)
}

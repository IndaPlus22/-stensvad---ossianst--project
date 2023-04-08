package planet

import (
	"github.com/go-gl/mathgl/mgl32"
)

func Gen(res uint32) ([]float32, []uint32) {
	corners := []mgl32.Vec3{
		{0.0, 1.0, 0.0},
		{0.0, -1.0, 0.0},
		{1.0, 0.0, 0.0},
		{0.0, 0.0, 1.0},
		{-1.0, 0.0, 0.0},
		{0.0, 0.0, -1.0},
	}

	return genDividedTriangle(corners[0], corners[2], corners[3], res)
}

/*
	i
	0		   B
	|		  / \
	1		 o---o
	|		/ \ / \
	2	  BA-BABC--BC
	|     / \ / \ / \
	3    o---o---o---o
	|	/ \ / \ / \ / \
	4  A---o---o---o---C
	|
	   0___1___2___3___4   k

	i
	0		   0
	|		  / \
	1		 1---2
	|		/ \ / \
	2	   3---4---5
	|     / \ / \ / \
	3    6---7---8---9
	|	/ \ / \ / \ / \
	4  10--11--12--13--14
	|
	   0___1___2___3___4   k
*/

func genDividedTriangle(a, b, c mgl32.Vec3, res uint32) ([]float32, []uint32) {
	vertices := []float32{}
	indices := []uint32{}
	index := uint32(0)

	for i := uint32(0); i <= res; i++ {
		BA := lerp(b, c, float32(i)/float32(res))
		BC := lerp(b, a, float32(i)/float32(res))

		for k := uint32(0); k <= i; k++ {
			fraction := float32(0.0)
			if i != 0 {
				fraction = float32(k) / float32(i)
			}

			BABC := lerp(BA, BC, fraction)
			vertices = append(vertices, BABC.X(), BABC.Y(), BABC.Z(), 0, 0, 0)

			println(BABC.X(), BABC.Y())
			if i < res {
				indices = append(indices, index, index+i+2, index+i+1)
				if k < i {
					indices = append(indices, index, index+1, index+i+2)
				}
			}

			index++
		}
	}

	return vertices, indices
}

func lerp(v1, v2 mgl32.Vec3, t float32) mgl32.Vec3 {
	return v1.Mul(1.0 - t).Add(v2.Mul(t))
}

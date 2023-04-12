package planet

import (
	"github.com/go-gl/mathgl/mgl32"
)

/*
GenPlanet generates the vertices and indices of a planet with the specified radius and vertex resolution.

Parameters:
- radius: the radius of the planet
- res: the resolution of the planet, in number of vertices per face

Returns:
- vertices: the vertices of the planet, as a float32 array
- indices: the indices of the vertices that form the triangles of the planet

Example usage:

	vertices, indices := GenPlanet(1.0, 10)
*/
func GenPlanet(radius float32, res uint32) ([]float32, []uint32) {
	points, indices := genOctahedron(res)

	normalizePointDistances(points, radius)

	// TODO: cool noise algorithms for mountains and valleys

	normals := calculateVertexNormals(points, indices)

	// add points and normals together as vertices in float32 array
	vertices := []float32{}
	for i := 0; i < len(points); i++ {
		vertices = append(vertices,
			points[i][0],
			points[i][1],
			points[i][2],
			normals[i][0],
			normals[i][1],
			normals[i][2])
	}

	return vertices, indices
}

/*
// Points:

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

// Indices:

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

func genOctahedron(res uint32) ([]mgl32.Vec3, []uint32) {
	// Points and indices for octahedron:
	corners := []mgl32.Vec3{
		{0.0, 1.0, 0.0},  // 0
		{1.0, 0.0, 0.0},  // 1
		{0.0, 0.0, 1.0},  // 2
		{-1.0, 0.0, 0.0}, // 3
		{0.0, 0.0, -1.0}, // 4
		{0.0, -1.0, 0.0}, // 5
	}

	faces := []float32{
		// Top faces:
		0, 1, 2,
		0, 2, 3,
		0, 3, 4,
		0, 4, 1,
		// Bottom faces:
		5, 2, 1,
		5, 3, 2,
		5, 4, 3,
		5, 1, 4,
	}

	points := []mgl32.Vec3{}
	indices := []uint32{}

	index := uint32(0)

	for i := 0; i < 8; i++ {
		triPoints, triInds := genDividedTriangle(
			corners[uint32(faces[i*3])],
			corners[uint32(faces[i*3+1])],
			corners[uint32(faces[i*3+2])],
			res,
			&index)

		points = append(points, triPoints...)
		indices = append(indices, triInds...)
	}

	mergeDuplicateVertices(points, indices)

	return points, indices
}

func lerp(v1, v2 mgl32.Vec3, t float32) mgl32.Vec3 {
	return v1.Mul(1.0 - t).Add(v2.Mul(t))
}
func genDividedTriangle(a, b, c mgl32.Vec3, res uint32, index *uint32) ([]mgl32.Vec3, []uint32) {
	points := []mgl32.Vec3{}
	indices := []uint32{}

	for i := uint32(0); i <= res; i++ {
		BA := lerp(b, c, float32(i)/float32(res))
		BC := lerp(b, a, float32(i)/float32(res))

		for k := uint32(0); k <= i; k++ {
			fraction := float32(0.0)
			if i != 0 {
				fraction = float32(k) / float32(i)
			}
			BABC := lerp(BA, BC, fraction)

			points = append(points, BABC)

			// Add the two new triangles formed by adding the point
			if i < res {
				indices = append(indices, *index, *index+i+2, *index+i+1)
				if k < i {
					indices = append(indices, *index, *index+1, *index+i+2)
				}
			}

			*index++
		}
	}

	return points, indices
}

func mergeDuplicateVertices(vertices []mgl32.Vec3, indices []uint32) {
	uniqueVertices := make(map[mgl32.Vec3]uint32)
	mergedIndices := make([]uint32, len(indices))

	for i := 0; i < len(indices); i++ {
		vertexPos := vertices[indices[i]]

		// check if the vertex is already in the uniqueVertices map
		mergedIndex, ok := uniqueVertices[vertexPos]

		if !ok {
			// the vertex is not in the map, add it to uniqueVertices
			mergedIndex = uint32(len(uniqueVertices))
			uniqueVertices[vertexPos] = mergedIndex
		}

		// update the index in the mergedIndices array
		mergedIndices[i] = mergedIndex
	}

	// modify vertices to remove duplicates
	newVertices := make([]mgl32.Vec3, len(uniqueVertices))
	for vertex, index := range uniqueVertices {
		newVertices[index] = vertex
	}

	copy(vertices, newVertices)
	copy(indices, mergedIndices)
}

func normalizePointDistances(points []mgl32.Vec3, radius float32) {
	for i := 0; i < len(points); i++ {
		points[i] = points[i].Normalize().Mul(radius)
	}
}

func calculateVertexNormals(points []mgl32.Vec3, indices []uint32) []mgl32.Vec3 {
	normals := make([]mgl32.Vec3, len(points))

	// iterate over each triangle and calculate its normal
	for i := 0; i < len(indices); i += 3 {
		v1 := points[indices[i]]
		v2 := points[indices[i+1]]
		v3 := points[indices[i+2]]

		// calculate triangle normal using cross product
		normal := v2.Sub(v1).Cross(v3.Sub(v1))

		// add triangle normal to each vertex normal
		normals[indices[i]] = normals[indices[i]].Add(normal)
		normals[indices[i+1]] = normals[indices[i+1]].Add(normal)
		normals[indices[i+2]] = normals[indices[i+2]].Add(normal)
	}

	// normalize each vertex normal
	for i := 0; i < len(normals); i++ {
		normals[i] = normals[i].Normalize()
	}

	return normals
}

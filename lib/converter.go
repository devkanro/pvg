package lib

import (
	"fmt"
	"image"
	"image/color"
	"strings"
)

func Convert(img image.Image, transparent color.Color) string {
	transparentValue := colorToInt(transparent)
	size := img.Bounds().Size()
	data := getColorData(img)
	delete(data, transparentValue)

	paths := make([]string, 0)

	for pixel, points := range data {
		paths = append(paths, path(pixel, points))
	}

	return fmt.Sprintf(`
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 -0.5 %d %d" shape-rendering="crispEdges">
%s
</svg>`,
		size.X, size.Y, strings.Join(paths, "\n"))
}

func path(pixel uint32, points []image.Point) string {
	pointsData := make([]string, 0)
	start := image.Point{X: -1, Y: -1}
	l := 0

	for _, point := range points {
		if start.X < 0 {
			start = point
			l = 1
			continue
		}
		if point.Y > start.Y || start.X+l != point.X {
			pointsData = append(pointsData, pathData(start, l))
			start = point
			l = 1
			continue
		}
		l++
	}

	if start.X >= 0 && l > 0 {
		pointsData = append(pointsData, pathData(start, l))
	}

	return fmt.Sprintf(`<path stroke="#%06X" d="%s"/>`, pixel, strings.Join(pointsData, ""))
}

func pathData(point image.Point, len int) string {
	return fmt.Sprintf("M%d %dh%d", point.X, point.Y, len)
}

func getColorData(img image.Image) map[uint32][]image.Point {
	result := map[uint32][]image.Point{}
	bounds := img.Bounds()
	size := bounds.Size()
	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			pixel := img.At(x, y)
			value := colorToInt(pixel)

			points, ok := result[value]
			if !ok {
				points = make([]image.Point, 0)
			}
			points = append(points, image.Point{
				X: x,
				Y: y,
			})

			result[value] = points
		}
	}
	return result
}

func colorToInt(color color.Color) uint32 {
	r, g, b, _ := color.RGBA()
	r = r >> 8
	g = g >> 8
	b = b >> 8
	return r<<16 | g<<8 | b
}

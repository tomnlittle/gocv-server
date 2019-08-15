package server

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"strconv"

	"gocv.io/x/gocv"
)

type encodingDefault struct {
	minQuality     uint8
	maxQuality     uint8
	DesiredQuality int
	qualityParam   int
	format         gocv.FileExt
}

// EncodingDefaults are used to easily define the default values for a given encoding
var EncodingDefaults = map[string]encodingDefault{
	"jpeg": {
		minQuality:     0,
		maxQuality:     100,
		DesiredQuality: 85,
		qualityParam:   gocv.IMWriteJpegQuality,
		format:         gocv.JPEGFileExt,
	},
	"jpg": {
		minQuality:     0,
		maxQuality:     100,
		DesiredQuality: 85,
		qualityParam:   gocv.IMWriteJpegQuality,
		format:         gocv.JPEGFileExt,
	},
	"png": {
		minQuality:     0,
		maxQuality:     9,
		DesiredQuality: 3,
		qualityParam:   gocv.IMWritePngCompression,
		format:         gocv.PNGFileExt,
	},
}

func parseStringAsInt(s string) (int, error) {
	s64, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return -1, fmt.Errorf("Cannot cast string as int")
	}

	// TODO: Handle Overflow
	return int(s64), nil
}

// Encode .
func Encode(mat gocv.Mat, format, quality string) ([]byte, error) {

	// format needs to be defined as well if quality is
	if format == "" && quality != "" {
		return nil, fmt.Errorf("Format needs to be specified with quality")
	}

	if format == "" {
		format = "png"
	}

	encoding, ok := EncodingDefaults[format]
	if !ok {
		return nil, fmt.Errorf("Unsupported format '%v'", format)
	}

	if quality != "" {

		quality, err := parseStringAsInt(quality)
		if err != nil {
			return nil, fmt.Errorf("Quality should be an integer")
		}

		encoding.DesiredQuality = quality
	}

	if encoding.DesiredQuality < int(encoding.minQuality) || encoding.DesiredQuality > int(encoding.maxQuality) {
		return nil, fmt.Errorf("Quality %v should be between %v and %v", encoding.DesiredQuality, encoding.minQuality, encoding.maxQuality)
	}

	params := []int{encoding.qualityParam, int(encoding.DesiredQuality)}
	return gocv.IMEncodeWithParams(encoding.format, mat, params)
}

// ------------------------------------------------ Standardised Functions ---------------------------------------------

// CVFunction .
type CVFunction func(mat gocv.Mat, parameters map[string]string) (*gocv.Mat, error)

// FunctionMappings maps an incoming id to the desired function
var FunctionMappings = map[string]CVFunction{
	"resize": Resize,
	"rotate": Rotate,
}

// Resize .
func Resize(mat gocv.Mat, parameters map[string]string) (*gocv.Mat, error) {

	width := parameters["width"]
	height := parameters["height"]

	if height == "" || width == "" {
		return nil, fmt.Errorf("'height' and 'width' are required fields for Resizing")
	}

	dst := gocv.NewMat()

	heightI, err := parseStringAsInt(height)
	if err != nil {
		return nil, fmt.Errorf("Height cannot be cast as int")
	}

	widthI, err := parseStringAsInt(width)
	if err != nil {
		return nil, fmt.Errorf("Width cannot be cast as int")
	}

	point := image.Point{widthI, heightI}

	gocv.Resize(mat, &dst, point, 0, 0, gocv.InterpolationArea)

	return &dst, nil
}

// Rotate Note: Resultant Image may be larger than the original
func Rotate(mat gocv.Mat, parameters map[string]string) (*gocv.Mat, error) {
	angle := parameters["angle"]

	if angle == "" {
		return nil, fmt.Errorf("'angle' is a required filed for performing a rotation")
	}

	angleI, err := parseStringAsInt(angle)
	if err != nil {
		return nil, fmt.Errorf("Angle cannot be cast as Int")
	}

	angleF := float64(angleI)

	width := float64(mat.Cols())
	height := float64(mat.Rows())

	// // calculate the new size
	newWidth := int(width*math.Sin(angleF) + width)
	newHeight := int(height*math.Cos(angleF) + height)

	paddedMat := gocv.NewMat()
	colour := color.RGBA{0, 0, 0, 0}
	gocv.CopyMakeBorder(mat, &paddedMat, newHeight/2, newHeight/2, newWidth/2, newWidth/2, gocv.BorderConstant, colour)

	// // TODO: Allow Float Angles
	dst := gocv.NewMat()
	newSize := image.Point{newWidth, newHeight}
	centre := image.Point{newWidth / 2, newHeight / 2}
	rMat := gocv.GetRotationMatrix2D(centre, float64(angleI), 1.0)
	gocv.WarpAffineWithParams(paddedMat, &dst, rMat, newSize, gocv.InterpolationLinear, gocv.BorderConstant, colour)

	return &dst, nil
}

package server

import (
	"fmt"
	"image"
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

	return int(s64), nil
}

// EncodeMatrix .
func EncodeMatrix(mat gocv.Mat, format, quality string) ([]byte, error) {

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
	"resize": ResizeMatrix,
	"clahe":  CLAHE,
}

// ResizeMatrix simple resizes an image to the given width and height
func ResizeMatrix(mat gocv.Mat, parameters map[string]string) (*gocv.Mat, error) {

	width := parameters["width"]
	height := parameters["height"]

	if height == "" || width == "" {
		return nil, fmt.Errorf("Height and Width are required fields for Resizing")
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

// CLAHE - Contrast Limited Adaptive Histogram Equalisation - https://docs.opencv.org/3.1.0/d5/daf/tutorial_py_histogram_equalization.html
func CLAHE(mat gocv.Mat, parameters map[string]string) (*gocv.Mat, error) {

	// convert image to a LAB encoding before apply clahe

	lab := gocv.NewMat()
	gocv.CvtColor(mat, &lab, gocv.ColorRGBToLab)
	mv := gocv.Split(lab)

	clahe := gocv.NewCLAHEWithParams(2.0, image.Point{8, 8})
	clahe.Apply(mv[0], &mv[0])

	dst := gocv.NewMat()
	gocv.Merge(mv, &dst)

	gocv.CvtColor(dst, &dst, gocv.ColorLabToRGB)

	return &dst, nil
}

package util

import (
	"sort"

	"github.com/signintech/gopdf"
)

//BuildPDFFromImages build a pdf from multiple images (reports)
func BuildPDFFromImages(outputPath string, filenames []string) error {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: gopdf.Rect{W: 841.89, H: 595.28}}) // Landscape A4

	//sort pdf
	filenames = sort.StringSlice(filenames)
	for _, file := range filenames {
		pdf.AddPage()
		err := pdf.Image(file, 0, 50, nil)
		if err != nil {
			return err
		}
	}

	if err := pdf.WritePdf(outputPath); err != nil {
		return err
	}

	return nil
}

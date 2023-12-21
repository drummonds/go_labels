// Create a page of labels with a given image
package main

import (
	"fmt"

	"github.com/jung-kurt/gofpdf"
)

const (
	thin             = 0.2
	thick            = 3.0
	LABEL_HEIGHT     = 67.7
	LABEL_HEIGHT_GAP = 0
	LABEL_WIDTH      = 99.1
	LABEL_WIDTH_GAP  = 2
	LABLES_PER_PAGE  = 8
)

type Label struct {
	Include  bool // should this label be included in the mail merge
	Title    string
	Name     string
	Surname  string
	Address  string
	Address2 string
	Town     string
	PostCode string
	Country  string
}

// Start with 1 page model
type Pdf struct {
	Pdf *gofpdf.Fpdf
	Row int
	Col int
}

//Formatting is done with no margins so positioning is absolute from top
//left corner of page
func New() *Pdf {
	// Create new page
	result := new(Pdf)
	result.Pdf = gofpdf.New("Portrait", "mm", "A4", "")
	p := result.Pdf
	p.SetMargins(0, 0, 0)
	p.SetAutoPageBreak(false, 0) // manage pages manually and bottom margin to 0
	p.SetFont("Helvetica", "", 10)
	p.SetFillColor(200, 200, 220)
	result.AddBlankPage()

	// report on size of page and margins
	left, top, right, bottom := p.GetMargins()
	fmt.Printf("Margins left %6.1f top %6.1f right %6.1f bottom %6.1f\n", left, top, right, bottom)
	width, height := p.GetPageSize()
	fmt.Printf("Page Size  width %6.1f height %6.1f\n", width, height)
	return result
}

// Print text at a point with sensible default
// Also increment lines
// func print(p *gofpdf.Fpdf, x, y float64, s string) float64 {
// 	if s == "" {
// 		return y // do nothing
// 	}
// 	p.Text(x, y, s)
// 	return y + 5.0
// }

// Set background colour so that if there is a little misalignment with the labels
// and you are using a none white background you don't have an edge problem
func (pdf *Pdf) AddBlankPage() {
	p := pdf.Pdf
	p.AddPage()
	p.SetDrawColor(218, 213, 124)
	p.SetFillColor(218, 213, 124)
	p.Rect(0, 0, 221, 297, "F")
}

// Chrome and edge print 1 mm further to left and 1mm higher
// top of text is 2.5
func (pdf *Pdf) Add(ad *Label) {
	y := 10.0 + LABEL_HEIGHT*float64(pdf.Row)
	x := 5.5 + (LABEL_WIDTH+LABEL_WIDTH_GAP)*float64(pdf.Col)
	p := pdf.Pdf
	p.SetXY(x, y)

	options := gofpdf.ImageOptions{
		ReadDpi:   false,
		ImageType: "",
	}
	p.ImageOptions("label.png", x, y, 99.1, 67.7, false, options, 0, "")
	// p.Image(, x+50, y+27, 4, 4*1190/1080, false, "", 0, "")
	if pdf.Col >= 1 && pdf.Row >= 4 {
		p.AddPage()
		pdf.Col = 0
		pdf.Row = 0
	} else {
		if pdf.Col >= 1 {
			pdf.Col = 0
			pdf.Row += 1
		} else {
			pdf.Col += 1
		}
	}
}

func (pdf *Pdf) Write() {
	fileStr := "labels.pdf"
	pdf.Pdf.OutputFileAndClose(fileStr)
}

func getDefault(row []string, index int, defaultValue string) string {
	if index >= len(row) {
		return defaultValue
	} else {
		return row[index]
	}
}

// This is a framework to read in a mailmerge list or a
func readImages(inputs chan *Label) {
	// f, err := excelize.OpenFile("xmasCards2021.xlsx")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// defer func() {
	// 	// Close the spreadsheet.
	// 	if err := f.Close(); err != nil {
	// 		fmt.Println(err)
	// 	}
	// }()
	// // Get value from cell by given worksheet name and axis.
	// cell, err := f.GetCellValue("Live", "B2")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(cell)
	// // Get all the rows in the Live.
	// // row, err := f.GetRows("Live")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	var ad *Label
	for i := 0; i < LABLES_PER_PAGE; i++ {
		// if i > 0 && i < 999 && getDefault(row, 2, "") == "1" {
		ad = new(Label)
		ad.Include = true
		inputs <- ad
		// }
	}
	close(inputs)
}

func writeLabels(inputs chan *Label, p *Pdf) {
	for ad := range inputs {
		if ad.Include {
			// fmt.Printf("Still to do %s %s\n", ad.Name, ad.Surname)
			p.Add(ad)
		}
	}
}

func main() {
	inputs := make(chan *Label, 200) // only expect 150 cards
	readImages(inputs)
	pdf := New()
	writeLabels(inputs, pdf)
	pdf.Write()
}

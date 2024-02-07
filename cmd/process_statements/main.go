package main

// TODOs
// 0. Port PyMuPDF, PyPdf, PdfPlumber to Go?
// 1. Figure out which file actually contains the statement (-0 vs -1)
// Approach:
// 2. Figure out if the year matches the statement title
// 3. Extract the following data from the pdfs
// - financial
// - management
// - text and tidbits (sentiment analysis)
// 4. Build intelligence through analysis of past data & news reports & any other source

func main() {
	// unipdf.
}

// 	// pdf.DebugOn = true
// 	content, err := readPdf("./statements/ASIANPAINT/2022-2023.pdf") // Read local pdf file
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(content)
// 	return
// }

// func readPdf(path string) (string, error) {
// 	r, err := pdf.Open(path)
// 	if err != nil {
// 		return "", err
// 	}
// 	totalPage := r.NumPage()

// 	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
// 		p := r.Page(pageIndex)
// 		if p.V.IsNull() {
// 			continue
// 		}

// 		rows, _ := p.GetTextByRow()
// 		for _, row := range rows {
// 			println(">>>> row: ", row.Position)
// 			for _, word := range row.Content {
// 				fmt.Println(word.S)
// 			}
// 		}
// 	}
// 	return "", nil
// }

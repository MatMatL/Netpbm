package Netpbm

//all pre-draw functions work except Read/Save P6 (not done)

//import needed library
import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// Define the PPM structure and Pixel struct
type PPM struct {
	data          [][]Pixel
	width, height int
	magicNumber   string
	max           uint8
}

type Pixel struct {
	R, G, B uint8
}

// ReadPPM reads a PPM image from a file and returns a struct that represents the image.
func ReadPPM(filename string) (*PPM, error) {
	//Part 1 : This part manage the file  open/close

	// Open file and check for error
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture du fichier :", err)
		return nil, err
	}
	// Close the file at the end of the function
	defer file.Close()

	//Part 2 : Defining of variables

	//define the structure used for all the functions
	FilledPPM := PPM{}
	// Create new scanner
	scanner := bufio.NewScanner(file)
	// Defines variables used during the function
	var magickeycheck, sizecheck, maxcheck = false, false, false
	var ligne, emplacement, pixelemp int = 0, 0, 0

	//Part 3 : Reading of the file line by line and processing them

	//this 'for' will scan the entire file and execute the code for each line of the file (line will be recovered using the 'scanner.Text()' fonction)
	for scanner.Scan() {
		//checking for '#' at the beginning of the line to skip comments
		if strings.HasPrefix(scanner.Text(), "#") {
			continue
			//if there's no comments, then :
		} else {
			//check if magicNumber is already found
			if magickeycheck == false {
				magickeycheck = true
				FilledPPM.magicNumber = scanner.Text()
				// If magicNumber found, check if file size is already found
			} else if sizecheck == false {
				sizecheck = true
				// Splits the line with the sizes using strings.Split() using space as splitting argument
				size := strings.Split(scanner.Text(), " ")

				// Converts the string containing the number into int using strconv.Atoi() to assign first part(number) of the string to the width
				FilledPPM.width, err = strconv.Atoi(size[0])
				// Checks for error during conversion
				if err != nil {
					fmt.Println("Erreur lors de la conversion 1")
				}

				// Convert the string containing the number into int using strconv.Atoi() to assigns second part(number) of the string to the height
				FilledPPM.height, err = strconv.Atoi(size[1])
				// Check for error during conversion
				if err != nil {
					fmt.Println("Erreur lors de la conversion 2")
				}

				// Create the array with the right dimensions using width and height we just found out
				FilledPPM.data = make([][]Pixel, FilledPPM.height)
				for i := range FilledPPM.data {
					FilledPPM.data[i] = make([]Pixel, FilledPPM.width)
				}

				// If sizes are found, check if file max is already found
			} else if maxcheck == false {
				maxcheck = true
				// Convert the string containing the number into int using strconv.Atoi() to assign number to temporary variable 'tempomax' (because max is Uint8)
				tempomax, err := strconv.Atoi(scanner.Text())
				// Check for error during conversion
				if err != nil {
					fmt.Println("Erreur lors du split 2")
				}
				// Fills the struct with 'max' value
				FilledPPM.max = uint8(tempomax)

			} else if magickeycheck == true && sizecheck == true && maxcheck == true {
				// Checks for magicNumber to diffenretiate P3 and P6 encryption
				if FilledPPM.magicNumber == "P3" {
					// Splits the current line using strings.Split() using space as splitting argument
					Currentline := strings.Split(scanner.Text(), " ")
					// Define pixel var using Pixel struct
					Pixel := Pixel{}

					// Range the current line to fill the array with values
					for i := 0; i < len(Currentline); i++ {
						// Convert value (in type string) into int using strconv.Atoi
						nombre, _ := strconv.Atoi(Currentline[i])

						// Fill the pixel according to the pixelemp (first fill R then G, then B and when it reaches B it fills the data with the pixel values and start again for next data)
						switch pixelemp {
						case 0:
							Pixel.R = uint8(nombre)
							pixelemp++
						case 1:
							Pixel.G = uint8(nombre)
							pixelemp++
						case 2:
							Pixel.B = uint8(nombre)
							pixelemp = 0

							//P6 is working when i print values but when ppm go into test file, values 255 changes into 253 (u can try it by removing '//' before the print above)

							//fmt.Println(Pixel)
							FilledPPM.data[ligne][emplacement] = Pixel
							//fmt.Println(FilledPPM.data[ligne][emplacement])

							emplacement++
						}
					}
					// Check for magicNumber to differentiate P1 and P4 encryption
				} else if FilledPPM.magicNumber == "P6" {
					// Define pixel var using Pixel struct
					Pixel := Pixel{}
					// Filling of the array
					var emplacement, compressedline int = 0, 0
					for _, number := range scanner.Text() {

						// Condition to go to the next line when the end of one is reached
						if emplacement == FilledPPM.width {
							emplacement = 0
							// Used += instead of incrementation because incrementation doesn't work here (i don't know why...)
							compressedline += 1
						}
						// Fill the pixel according to the pixelemp (first fill R then G, then B and when it reaches B it fills the data with the pixel values and start again for next data)
						switch pixelemp {
						case 0:
							Pixel.R = uint8(number)
							pixelemp++
						case 1:
							Pixel.G = uint8(number)
							pixelemp++
						case 2:
							Pixel.B = uint8(number)
							pixelemp = 0
							FilledPPM.data[compressedline][emplacement] = Pixel
							emplacement++
						}
					}

					//Error if the magic number isn't valid
				} else {
					fmt.Println("Erreur, magic number pas reconnue")
				}
				// Add a line each time a line is read
				ligne++
				// Reset emplacement and pixelemp for next line
				emplacement = 0
				pixelemp = 0
			}
		}
	}
	// Returns the PGM struct filled with file data
	return &PPM{FilledPPM.data, FilledPPM.width, FilledPPM.height, FilledPPM.magicNumber, FilledPPM.max}, nil
}

// Size returns the width and height of the image.
func (ppm *PPM) Size() (int, int) {
	return ppm.width, ppm.height
}

// At returns the value of the pixel at (x, y).
func (ppm *PPM) At(x, y int) Pixel {
	return ppm.data[y][x]
}

// Set sets the value of the pixel at (x, y).
func (ppm *PPM) Set(x, y int, value Pixel) {
	ppm.data[x][y] = value
}

// Save saves the PPM image to a file and returns an error if there was a problem.
func (ppm *PPM) Save(filename string) error {
	//Manage file creating, opening and closing of the file
	file, err := os.Create(filename)
	defer file.Close()

	//write magicNumber, width, height and max in the file
	fmt.Fprintf(file, "%s\n%d %d\n%d\n", ppm.magicNumber, ppm.width, ppm.height, ppm.max)

	// Range all the data of the struct and write value of the pixel of the location
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			CurrentPixel := ppm.data[i][j]
			fmt.Fprint(file, CurrentPixel.R, " ", CurrentPixel.G, " ", CurrentPixel.B, " ")
		}
		// Skip to the next line
		fmt.Fprintln(file)
	}

	return err
}

// Invert inverts the colors of the PPM image.
func (ppm *PPM) Invert() {
	//browse all the array
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			// Invert the value of each r, g and b using the max value
			ppm.data[i][j].R = uint8(ppm.max) - ppm.data[i][j].R
			ppm.data[i][j].G = uint8(ppm.max) - ppm.data[i][j].G
			ppm.data[i][j].B = uint8(ppm.max) - ppm.data[i][j].B
		}
	}
}

// Flip flips the PPM image horizontally.
func (ppm *PPM) Flip() {
	// This function is going to invert 1rst line with the last, then second with second-to-last ect
	// So here we define number of changes
	var nb_change int = (ppm.width / 2)
	//tempo is used as a stocker
	var tempo Pixel
	// Range all the array line by line and then exchange
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < nb_change; j++ {
			tempo = ppm.data[i][j]
			ppm.data[i][j] = ppm.data[i][ppm.width-j-1]
			ppm.data[i][ppm.width-j-1] = tempo
		}
	}
}

// Flop flops the PPM image vertically.
func (ppm *PPM) Flop() {
	// Same logic as above but in the other direction
	var nb_change int = (ppm.height / 2)
	var tempo Pixel
	for i := 0; i < ppm.width; i++ {
		for j := 0; j < nb_change; j++ {
			tempo = ppm.data[j][i]
			ppm.data[j][i] = ppm.data[ppm.height-j-1][i]
			ppm.data[ppm.height-j-1][i] = tempo
		}
	}
}

// SetMagicNumber sets the magic number of the PPM image.
func (ppm *PPM) SetMagicNumber(magicNumber string) {
	ppm.magicNumber = magicNumber
}

// SetMaxValue sets the max value of the PPM image.
func (ppm *PPM) SetMaxValue(maxValue uint8) {
	// Browse all the array
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			//re-calculate the value depending on the new max value (using cross product)
			ppm.data[i][j].R = uint8(float64(ppm.data[i][j].R) * float64(maxValue) / float64(ppm.max))
			ppm.data[i][j].G = uint8(float64(ppm.data[i][j].G) * float64(maxValue) / float64(ppm.max))
			ppm.data[i][j].B = uint8(float64(ppm.data[i][j].B) * float64(maxValue) / float64(ppm.max))
		}
	}
	//Change max to new max value
	ppm.max = maxValue
}

// Rotate90CW rotates the PPM image 90° clockwise.
func (ppm *PPM) Rotate90CW() {
	// Tempo is used as a stocker
	var tempo int

	// Create new array just like the old one to stock new data
	rotatedData := make([][]Pixel, ppm.width)
	for i := range rotatedData {
		rotatedData[i] = make([]Pixel, ppm.height)
	}

	//assigns the value of a pixel to its equivalent after rotation
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			rotatedData[i][j] = ppm.data[(ppm.width-1)-j][i]
		}
	}

	// Actualise the data
	tempo = ppm.width
	ppm.width = ppm.height
	ppm.height = tempo
	ppm.data = rotatedData
}

// ToPGM converts the PPM image to PGM.
func (ppm *PPM) ToPGM() *PGM {
	// Create new array just like the old one but of uint8 instead of Pixel to stock new data
	data := make([][]uint8, ppm.width)
	for i := range data {
		data[i] = make([]uint8, ppm.height)
	}

	// Range data and fill the uint8 array according to the pixel value (making a average value of all R, G and B)
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			R := int(ppm.data[i][j].R)
			G := int(ppm.data[i][j].G)
			B := int(ppm.data[i][j].B)
			average := (R + G + B) / 3
			data[i][j] = uint8(average)
		}
	}

	// Return a pbm struct filled with p2 values
	return &PGM{data, ppm.width, ppm.height, "P2", ppm.max}
}

// ToPBM converts the PPM image to PBM.
func (ppm *PPM) ToPBM() *PBM {
	// Create new array just like the old one but it uses boolean instead of Pixel to stock new data
	data := make([][]bool, ppm.width)
	for i := range data {
		data[i] = make([]bool, ppm.height)
	}

	// Range data and fill the boolean array depending on the pixel value (if value closer of 0 than max then false, and if closer max than 0 then true)
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			if uint8((int(ppm.data[i][j].R)+int(ppm.data[i][j].G)+int(ppm.data[i][j].B))/3) < ppm.max/2 {
				data[i][j] = true
			} else {
				data[i][j] = false
			}
		}
	}

	// Return a pbm struct filled with p1 values
	return &PBM{data, ppm.width, ppm.height, "P1"}
}

// Define the Point structure
type Point struct {
	X, Y int
}

// DrawLine draws a line between two points.
func (ppm *PPM) DrawLine(p1, p2 Point, color Pixel) {

	// Define the number of width and height between the 2 points
	dx := Abs(p2.X - p1.X)
	dy := Abs(p2.Y - p1.Y)

	// I applied the Bresenham algorithm
	sx := -1
	if p1.X < p2.X {
		sx = 1
	}

	sy := -1
	if p1.Y < p2.Y {
		sy = 1
	}

	err := dx - dy

	for {
		// Prevent from out of range pixels
		if (p1.Y < ppm.height) && (p1.X < ppm.width) {
			ppm.data[p1.Y][p1.X] = color
		}

		if p1.X == p2.X && p1.Y == p2.Y {
			break
		}

		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			p1.X += sx
		}
		if e2 < dx {
			err += dx
			p1.Y += sy
		}
	}
}

// Small fuction just to return absolute value of an int
func Abs(nb int) int {
	if nb < 0 {
		nb = -nb
		return nb
	} else {
		return nb
	}
}

// DrawRectangle draws a rectangle.
func (ppm *PPM) DrawRectangle(p1 Point, width, height int, color Pixel) {
	// Define all 4 corners of the rectangle
	TopLeft := Point{
		X: p1.X,
		Y: p1.Y,
	}

	TopRight := Point{
		X: (p1.X + width),
		Y: p1.Y,
	}

	BottomLeft := Point{
		X: p1.X,
		Y: (p1.Y + height),
	}

	BottomRight := Point{
		X: (p1.X + width),
		Y: (p1.Y + height),
	}

	ppm.DrawLine(TopLeft, TopRight, color)
	ppm.DrawLine(BottomLeft, BottomRight, color)
	ppm.DrawLine(TopLeft, BottomLeft, color)
	ppm.DrawLine(TopRight, BottomRight, color)
}

// DrawFilledRectangle draws a filled rectangle.
func (ppm *PPM) DrawFilledRectangle(p1 Point, width, height int, color Pixel) {
	ppm.DrawRectangle(p1, width, height, color)
	if width == 0 && height == 0 {
		return
	} else if width == 0 && height > 0 {
		ppm.DrawFilledRectangle(p1, width, height-1, color)
	} else if width > 0 && height == 0 {
		ppm.DrawFilledRectangle(p1, width-1, height, color)
	} else if width > 0 && height > 0 {
		ppm.DrawFilledRectangle(p1, width-1, height-1, color)
	}
}

// DrawCircle draws a circle.
func (ppm *PPM) DrawCircle(center Point, radius int, color Pixel) {

	for x := 0; x < ppm.height; x++ {
		for y := 0; y < ppm.width; y++ {
			dx := float64(x) - float64(center.X)
			dy := float64(y) - float64(center.Y)
			distance := math.Sqrt(dx*dx + dy*dy)

			if math.Abs(distance-float64(radius)) < 1.0 && distance < float64(radius) {
				ppm.Set(x, y, color)
			}
		}
	}
	ppm.Set(center.X-(radius-1), center.Y, color)
	ppm.Set(center.X+(radius-1), center.Y, color)
	ppm.Set(center.X, center.Y+(radius-1), color)
	ppm.Set(center.X, center.Y-(radius-1), color)
}

// DrawFilledCircle draws a filled circle.
func (ppm *PPM) DrawFilledCircle(center Point, radius int, color Pixel) {
	ppm.DrawCircle(center, radius, color)
	if radius != 0 {
		ppm.DrawFilledCircle(center, radius-1, color)
	}
}

// DrawTriangle draws a triangle.
func (ppm *PPM) DrawTriangle(p1, p2, p3 Point, color Pixel) {
	ppm.DrawLine(p1, p2, color)
	ppm.DrawLine(p1, p3, color)
	ppm.DrawLine(p2, p3, color)
}

// DrawFilledTriangle draws a filled triangle.
func (ppm *PPM) DrawFilledTriangle(p1, p2, p3 Point, color Pixel) {
	ppm.DrawTriangle(p1, p2, p3, color)
	if p1 != p2 {
		NewPoint := Point{
			X: p1.X + 1,
			Y: p1.Y,
		}
		ppm.DrawFilledTriangle(NewPoint, p2, p3, color)
	}
}

// DrawPolygon draws a polygon.
func (ppm *PPM) DrawPolygon(points []Point, color Pixel) {
	// ...
}

// DrawFilledPolygon draws a filled polygon.
func (ppm *PPM) DrawFilledPolygon(points []Point, color Pixel) {
	// ...
}

// DrawKochSnowflake draws a Koch snowflake.
func (ppm *PPM) DrawKochSnowflake(n int, start Point, width int, color Pixel) {
	// N is the number of iterations.
	// Koch snowflake is a 3 times a Koch curve.
	// Start is the top point of the snowflake.
	// Width is the width all the lines.
	// Color is the color of the lines.
	// ...
}

// DrawSierpinskiTriangle draws a Sierpinski triangle.
func (ppm *PPM) DrawSierpinskiTriangle(n int, start Point, width int, color Pixel) {
	// N is the number of iterations.
	// Start is the top point of the triangle.
	// Width is the width all the lines.
	// Color is the color of the lines.
	// ...
}

// DrawPerlinNoise draws perlin noise.
// this function Draw a perlin noise of all the image.
func (ppm *PPM) DrawPerlinNoise(color1 Pixel, color2 Pixel) {
	// Color1 is the color of 0.
	// Color2 is the color of 1.
}

// KNearestNeighbors resizes the PPM image using the k-nearest neighbors algorithm.
func (ppm *PPM) KNearestNeighbors(newWidth, newHeight int) {
	// ...
}

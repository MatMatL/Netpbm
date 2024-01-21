package Netpbm

//all functions work except Read/Save P5 (not done)

//import needed library
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// define the PGM structure
type PGM struct {
	data          [][]uint8
	width, height int
	magicNumber   string
	max           uint8
}

// This fuction is used to read a .pgm file
func ReadPGM(filename string) (*PGM, error) {
	//Part 1 : This part manages the file  open/close

	// Opens the file and checks for error
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture du fichier :", err)
		return nil, err
	}
	// Close the file at the end of the function
	defer file.Close()

	//Part 2 : Defining of variables

	// Define the structure used all function long
	FilledPGM := PGM{}
	// Create new scanner
	scanner := bufio.NewScanner(file)
	// Define variables used during function
	var magickeycheck, sizecheck, maxcheck = false, false, false
	var line int = 0

	// Part 3 : Reading of the file line by line and processing them

	// This 'for' will scan the entire file and execute the code for each line of the file (line will be recovered using the 'scanner.Text()' function)
	for scanner.Scan() {
		// Checking for '#' at the beginning of the line to skip comments
		if strings.HasPrefix(scanner.Text(), "#") {
			continue
			// If there's no comments, then :
		} else {
			// Checks if magicNumber is already found
			if magickeycheck == false {
				magickeycheck = true
				FilledPGM.magicNumber = scanner.Text()

				// If magicNumber found, check if file size is already found
			} else if sizecheck == false {
				sizecheck = true
				// Split the line with the size using strings.Split() using space as splitting argument
				size := strings.Split(scanner.Text(), " ")

				// Converts the string containing the number into an int using strconv.Atoi() to assign first part(number) of the string to the width
				FilledPGM.width, err = strconv.Atoi(size[0])
				// Checks for error during conversion
				if err != nil {
					fmt.Println("Erreur lors de la conversion 1")
				}

				// Converts the string containing the number into int using strconv.Atoi() to assign second part(number) of the string to the height
				FilledPGM.height, err = strconv.Atoi(size[1])
				// Check for error during conversion
				if err != nil {
					fmt.Println("Erreur lors de la conversion 2")
				}

				// Creates the array with the right dimensions using width and height we just found out
				FilledPGM.data = make([][]uint8, FilledPGM.height)
				for i := range FilledPGM.data {
					FilledPGM.data[i] = make([]uint8, FilledPGM.width)
				}

				// If sizes are found, check if file max is already found
			} else if maxcheck == false {
				maxcheck = true
				// Converts the string containing the number into int using strconv.Atoi() to assigns number to temporary variable 'tempomax' (because max is Uint8)
				tempomax, err := strconv.Atoi(scanner.Text())
				// Check for error during conversion
				if err != nil {
					fmt.Println("Erreur lors de la conversion 3")
				}
				// Fills the struct with 'max' value
				FilledPGM.max = uint8(tempomax)

			} else if magickeycheck == true && sizecheck == true && maxcheck == true {
				// Check for magicNumber to diffenretiate P2 and P5 encryption
				if FilledPGM.magicNumber == "P2" {
					// Splits the current line using strings.Split() using space as splitting argument
					Currentline := strings.Split(scanner.Text(), " ")
					// Range the current line to fill the array with values
					for i := 0; i < FilledPGM.width; i++ {
						// Converts value (in type string) into int using strconv.Atoi
						nombre, err := strconv.Atoi(Currentline[i])
						// Checks for error during conversion
						if err != nil {
							fmt.Println("Erreur lors de la conversion 4")
						}

						// Filling of the array
						FilledPGM.data[line][i] = uint8(nombre)
					}
					// Checks for magicNumber to diffenretiate P1 and P4 encryption
				} else if FilledPGM.magicNumber == "P5" {
					// Filling of the array
					var emplacement, compressedline int = 0, 0
					for _, number := range scanner.Text() {

						// Condition to go to the next line when the end of one is reached
						if emplacement == FilledPGM.width {
							emplacement = 0
							// Used += instead of incrementation because incrementation doesn't work here (i don't know why...)
							compressedline += 1
						}
						// Condition to deal with to much data (not needed if file is correct)
						if compressedline != FilledPGM.width {
							FilledPGM.data[compressedline][emplacement] = uint8(number)
							emplacement++
						}
					}

					//Error if the magic number isn't valid
				} else {
					fmt.Println("Erreur, magic number pas reconnue")
				}
				// Add a line each time a line is read
				line++
			}
		}
	}
	// Returns PGM struct filled with file data
	return &PGM{FilledPGM.data, FilledPGM.width, FilledPGM.height, FilledPGM.magicNumber, FilledPGM.max}, nil
}

// Size returns the width and height of the image.
func (pgm *PGM) Size() (int, int) {
	return pgm.width, pgm.height
}

// At returns the value of the pixel at (x, y).
func (pgm *PGM) At(x, y int) uint8 {
	return pgm.data[x][y]
}

// Set sets the value of the pixel at (x, y).
func (pgm *PGM) Set(x, y int, value uint8) {
	pgm.data[x][y] = value
}

// Save saves the PGM image to a file and returns an error if there was a problem.
func (pgm *PGM) Save(filename string) error {
	// Manage file creating, opening and closing
	file, err := os.Create(filename)
	defer file.Close()

	//write magicNumber, width, height and max in the file
	fmt.Fprintf(file, "%s\n%d %d\n%d\n", pgm.magicNumber, pgm.width, pgm.height, pgm.max)

	//range all the data of the struct and write value of the location
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			fmt.Fprint(file, pgm.data[i][j], " ")
		}
		//skip to next line
		fmt.Fprintln(file)
	}

	return err
}

// Invert inverts the colors of the PGM image.
func (pgm *PGM) Invert() {
	//Browse all the array
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			// Inverts the value using the max value
			pgm.data[i][j] = uint8(pgm.max) - pgm.data[i][j]
		}
	}
}

// Flip flips the PGM image horizontally.
func (pgm *PGM) Flip() {
	// Function is going to invert 1rst line with the last, then second with second-to-last and so on
	//so here we are defining the number of changes
	var nb_change int = (pgm.width / 2)
	// Tempo is used as a stocker
	var tempo uint8
	// Range all the array line by line and exchange
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < nb_change; j++ {
			tempo = pgm.data[i][j]
			pgm.data[i][j] = pgm.data[i][pgm.width-j-1]
			pgm.data[i][pgm.width-j-1] = tempo
		}
	}
}

// Flop flops the PGM image vertically.
func (pgm *PGM) Flop() {
	// Same logic as above but in the other direction
	var nb_change int = (pgm.height / 2)
	var tempo uint8
	for i := 0; i < pgm.width; i++ {
		for j := 0; j < nb_change; j++ {
			tempo = pgm.data[j][i]
			pgm.data[j][i] = pgm.data[pgm.height-j-1][i]
			pgm.data[pgm.height-j-1][i] = tempo
		}
	}
}

// SetMagicNumber sets the magic number of the PGM image.
func (pgm *PGM) SetMagicNumber(magicNumber string) {
	pgm.magicNumber = magicNumber
}

// SetMaxValue sets the max value of the PGM image.
func (pgm *PGM) SetMaxValue(maxValue uint8) {
	//browse all the array
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			//re-calculate the value depending on the new max value (using cross product)
			pgm.data[i][j] = uint8(float64(pgm.data[i][j]) * float64(maxValue) / float64(pgm.max))
		}
	}
	//change max to new max value
	pgm.max = maxValue
}

// Rotate90CW rotates the PGM image 90° clockwise.
func (pgm *PGM) Rotate90CW() {
	//tempo is used as a stocker
	var tempo int

	// Creates a new array just like the old one to stock new data
	rotatedData := make([][]uint8, pgm.width)
	for i := range rotatedData {
		rotatedData[i] = make([]uint8, pgm.height)
	}

	// Assigns the value of a pixel to its equivalent after rotation
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			rotatedData[i][j] = pgm.data[(pgm.width-1)-j][i]
		}
	}

	// Actualise the data
	tempo = pgm.width
	pgm.width = pgm.height
	pgm.height = tempo
	pgm.data = rotatedData
}

// ToPBM converts the PGM image to PBM.
func (pgm *PGM) ToPBM() *PBM {
	// Creates a new array just like the old one but it uses boolean instead of uint8 to stock new data
	data := make([][]bool, pgm.width)
	for i := range data {
		data[i] = make([]bool, pgm.height)
	}

	// Ranges data and fills the boolean array depending on the pixel value (if value closer of 0 than max then false, and if closer max than 0 then true)
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			if pgm.data[i][j] >= (uint8(pgm.max) / 2) {
				data[i][j] = false
			} else if pgm.data[i][j] < (uint8(pgm.max) / 2) {
				data[i][j] = true
			}
		}
	}

	//return a pbm struct filled with p1 values
	return &PBM{data, pgm.width, pgm.height, "P1"}
}

package Netpbm

// All functions work except save P4 (not done)

//Import needed library
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// define the PBM structure
type PBM struct {
	data          [][]bool
	width, height int
	magicNumber   string
}

// Fuction used to read a .pbm file
func ReadPBM(filename string) (*PBM, error) {
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

	// Define the structure used all fonction long
	FilledPBM := PBM{}
	// Create new scanner
	scanner := bufio.NewScanner(file)
	// Define variables used during fuction
	var magickeycheck, sizecheck bool = false, false
	var line int = 0

	//Part 3 : Reading of the file line by line and processing them

	// This 'for' will scan the entire file and execute the code for each line of the file (line will be recovered using the 'scanner.Text()' fonction)
	for scanner.Scan() {
		// Checking for '#' at the begining of the line to skip comments
		if strings.HasPrefix(scanner.Text(), "#") {
			continue
			// If there's no comments, then :
		} else {
			// Check if magicNumber is already found
			if magickeycheck == false {
				magickeycheck = true
				FilledPBM.magicNumber = scanner.Text()

				// If magicNumber found, check if file size is already found
			} else if sizecheck == false {
				sizecheck = true
				// Splits the line with the sizes using strings.Split() using space as splitting argument
				size := strings.Split(scanner.Text(), " ")

				// Convert the string containing the number into int using strconv.Atoi() to assigns first part(number) of the string to the width
				FilledPBM.width, err = strconv.Atoi(size[0])
				// Check for error during conversion
				if err != nil {
					fmt.Println("Erreur lors de la conversion 1")
				}

				// Convert the string containing the number into int using strconv.Atoi() to assigns second part(number) of the string to the height
				FilledPBM.height, err = strconv.Atoi(size[1])
				// Check for error during conversion
				if err != nil {
					fmt.Println("Erreur lors de la conversion 2")
				}

				// Create the array with the right dimensions using width and height we just found out
				FilledPBM.data = make([][]bool, FilledPBM.height)
				for i := range FilledPBM.data {
					FilledPBM.data[i] = make([]bool, FilledPBM.width)
				}

			} else if magickeycheck == true && sizecheck == true {
				// Check for magicNumber to deffenretiate P1 and P4 encryption
				if FilledPBM.magicNumber == "P1" {
					emplacement := 0
					// Range the current line to fill the array with true and false depending on 0 and 1
					for _, number := range scanner.Text() {
						// Condition to avoid spaces and only keep 0 and 1
						if number != ' ' {
							// Fills the array using 'ToBool()' fonction to translate number into boolean
							FilledPBM.data[line][emplacement] = ToBool(number)
							emplacement++
						}
					}
					// Check for magicNumber to diffenretiate P1 and P4 encryption
				} else if FilledPBM.magicNumber == "P4" {
					// Define the number of bytes used to contain 1 single line
					var bytesnumber int
					// Condition used to have the right amount of bytes
					if FilledPBM.width%8 == 0 {
						bytesnumber = (FilledPBM.width / 8)

					} else {
						bytesnumber = (FilledPBM.width / 8) + 1
					}
					// Define the number of bit(s) of padding
					padding := (bytesnumber * 8) - FilledPBM.width

					// Send the current line to the fonction 'ToBinary' to process the character and translate them to binary while removing the padding
					binaire := ToBinary(scanner.Text(), bytesnumber, padding)

					// Filling of the array
					var emplacement, compressedline int = 0, 0
					for _, number := range binaire {

						// Condition to go to the next line when the end of one is reached
						if emplacement == FilledPBM.width {
							emplacement = 0
							// Used += instead of incrementation because incrementation doesn't work here (i don't know why...)
							compressedline += 1
						}
						// Condition to deal with to much data (not needed if file is correct)
						if compressedline != FilledPBM.width {
							FilledPBM.data[compressedline][emplacement] = ToBool(number)
							emplacement++
						}
					}
					// Error if the magic number isn't valid
				} else {
					fmt.Println("Erreur, magic number pas reconnue")
				}
				// Add a line each time a line is read
				line++
			}
		}
	}
	// Return PBM struct filled with file data
	return &PBM{FilledPBM.data, FilledPBM.width, FilledPBM.height, FilledPBM.magicNumber}, nil
}

// Function used to translate a number into a boolean. It takes as argument a number (0 or 1) and returns the corresponding boolean
func ToBool(nb rune) bool {
	var Boobool bool
	if nb == '0' {
		Boobool = false
	} else if nb == '1' {
		Boobool = true
	}
	return Boobool
}

// The goal of the fonction 'ToBinary'is to process the character and translate them to binary while removing padding. It takes as argument a string containing characters, the number of bytes per line and the padding then it returns a string of 0 and 1 containing all data of the input string
func ToBinary(input string, bytesnumber, padding int) string {
	var binary string = ""

	// Continue as long as all the strings aren't all ranged
	for i := 0; i < len(input); i++ {
		//"translate" the char into it's binary form
		bin := fmt.Sprintf("%08b", input[i])
		// Remove padding from last bytes of each line
		if bytesnumber != 1 {
			if i != 0 && (i+1)%bytesnumber == 0 {
				bin = bin[:len(bin)-padding]
			}
		} else if bytesnumber == 1 {
			bin = bin[:len(bin)-padding]
		}
		// Add the char converted in binary to final string
		binary += bin
	}
	// Return the final string
	return binary
}

// Size returns the width and height of the image.
func (pbm *PBM) Size() (int, int) {
	return pbm.width, pbm.height
}

// At returns the value of the pixel at (x, y).
func (pbm *PBM) At(x, y int) bool {
	return pbm.data[x][y]
}

// Set sets the value of the pixel at (x, y).
func (pbm *PBM) Set(x, y int, value bool) {
	pbm.data[x][y] = value
}

// Save saves the PBM image to a file and returns an error if there was a problem. (only P1 made)
func (pbm *PBM) Save(filename string) error {
	//Manage file creating, opening and closing of the file
	file, err := os.Create(filename)
	defer file.Close()

	//write magicNumber, width and height in the file
	fmt.Fprintf(file, "%s\n%d %d\n", pbm.magicNumber, pbm.width, pbm.height)

	// Range all the data of the struct and write 0 or 1 according to the data (0 or 1)
	for _, i := range pbm.data {
		for _, j := range i {
			if j {
				fmt.Fprint(file, "1 ")
			} else {
				fmt.Fprint(file, "0 ")
			}
		}
		// Skip to next line
		fmt.Fprintln(file)
	}
	return err
}

// Invert inverts the colors of the PBM image.
func (pbm *PBM) Invert() {
	// Browse all the array
	for i := 0; i < pbm.height; i++ {
		for j := 0; j < pbm.width; j++ {
			// Change true to false and false to true for each location
			if pbm.data[i][j] {
				pbm.data[i][j] = false
			} else {
				pbm.data[i][j] = true
			}
		}
	}
}

// Flip flips the PBM image horizontally.
func (pbm *PBM) Flip() {
	// The fuction is going to invert 1rst line with the last, then second with second-to-last and so on
	// So here we are defining the number of changes
	var nb_change int = (pbm.width / 2)
	//tempo is used as stock
	var tempo bool
	//range all the array line by line and exchange
	for i := 0; i < pbm.height; i++ {
		for j := 0; j < nb_change; j++ {
			tempo = pbm.data[i][j]
			pbm.data[i][j] = pbm.data[i][pbm.width-j-1]
			pbm.data[i][pbm.width-j-1] = tempo
		}
	}
}

// Flop flops the PBM image vertically.
func (pbm *PBM) Flop() {
	// Same logic as above but in the other direction
	var nb_change int = (pbm.height / 2)
	var tempo bool
	for i := 0; i < pbm.width; i++ {
		for j := 0; j < nb_change; j++ {
			tempo = pbm.data[j][i]
			pbm.data[j][i] = pbm.data[pbm.height-j-1][i]
			pbm.data[pbm.height-j-1][i] = tempo
		}
	}
}

// SetMagicNumber sets the magic number of the PBM image.
func (pbm *PBM) SetMagicNumber(magicNumber string) {
	pbm.magicNumber = magicNumber
}

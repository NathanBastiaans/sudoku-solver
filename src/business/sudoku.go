package business

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
)

// Field holds all the data for a single cell
type Field struct {
	Row             uint
	Column          uint
	Number          int
	PossibleNumbers []uint
}

// Board hold all the field data
type Board struct {
	Fields []*Field
}

// Box is the data for the 3x3 boxes on the field
type Box struct {
	Rows []uint
	Cols []uint
}

// Populate the boxes data
var boxes = []Box{
	{Rows: []uint{1, 2, 3}, Cols: []uint{1, 2, 3}},
	{Rows: []uint{1, 2, 3}, Cols: []uint{4, 5, 6}},
	{Rows: []uint{1, 2, 3}, Cols: []uint{7, 8, 9}},
	{Rows: []uint{4, 5, 6}, Cols: []uint{1, 2, 3}},
	{Rows: []uint{4, 5, 6}, Cols: []uint{4, 5, 6}},
	{Rows: []uint{4, 5, 6}, Cols: []uint{7, 8, 9}},
	{Rows: []uint{7, 8, 9}, Cols: []uint{1, 2, 3}},
	{Rows: []uint{7, 8, 9}, Cols: []uint{4, 5, 6}},
	{Rows: []uint{7, 8, 9}, Cols: []uint{7, 8, 9}},
}

// Init vars
var iterationsSinceLastSolve = 0
var solvedNumberCount = 0
var maxIterationsSinceLastSolve = 100

// GenerateBoardFromFile generates a board from the given file
func GenerateBoardFromFile(filePath string) (Board, error) {
	// create the board
	var board = Board{}

	// Open the given file
	file, err := os.Open(filePath)
	if err != nil {
		return board, err
	}

	// close the file when done. Fatal errors when fails
	defer func() {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// row should be 1
	var row uint = 1
	// col should be 0 because it gets increased at the start of the loop
	var col uint = 0

	// Init a new scanner
	scanner := bufio.NewScanner(file)

	// Loop over all lines
	for scanner.Scan() {
		line := scanner.Text()

		// If the length of a line isn't 9 the line is incorrect for the sudoku solver
		if len(line) != 9 {
			return board, errors.New(fmt.Sprintf("line is not 9. Instead got %v", line))
		}

		// Loop over all numbers in the line
		for _, number := range line {
			col = col + 1

			// if col is greater than 9 we're at a new row
			if col > 9 {
				row = row + 1
				col = 1
			}

			// parse the sting number to an integer
			nr, err := strconv.ParseInt(string(number), 10, 0)
			if err != nil {
				return board, err
			}

			// Append the field to the board
			board.Fields = append(board.Fields, &Field{
				Row:    row,
				Column: col,
				Number: int(nr),
			})
		}
	}

	// Quick check to validate the length of the field
	if len(board.Fields) != 81 {
		return board, errors.New("wrong field count")
	}

	return board, nil
}

// HasNumberInRow checks if the given number exists in the row of the given field
func (board Board) HasNumberInRow(field Field, number uint) bool {
	for _, f := range board.Fields {
		if f.Row != field.Row {
			continue
		}

		if f.Number == int(number) {
			return true
		}
	}

	return false
}

// HasNumberInColumn checks if the given number exists in the column of the given field
func (board Board) HasNumberInColumn(field Field, number uint) bool {
	for _, f := range board.Fields {
		if f.Column != field.Column {
			continue
		}

		if f.Number == int(number) {
			return true
		}
	}

	return false
}

// GetBoxForField returns the correct box data for the given field
func GetBoxForField(field Field) (Box, error) {

	for _, box := range boxes {
		if contains(box.Cols, int(field.Column)) && contains(box.Rows, int(field.Row)) {
			return box, nil
		}
	}

	return Box{}, errors.New("could not find correct box for field")
}

// HasNumberInBox checks if the given number exists in the box of the given field
func (board Board) HasNumberInBox(field Field, number uint) bool {

	box, err := GetBoxForField(field)

	if err != nil {
		log.Fatal(err)
	}

	for _, f := range board.Fields {

		if !contains(box.Rows, int(f.Row)) || !contains(box.Cols, int(f.Column)) {
			continue
		}

		if f.Number == int(number) {
			return true
		}
	}

	return false
}

// Solve solves the given board and returns the data
func Solve(board Board) (Board, error) {

	// infinite loops for infinite loops
	for {
		// reset the solvedNumberCount when a new iteration over the board starts
		solvedNumberCount = 0

		// Loop over all the fields in the board
		for _, field := range board.Fields {

			// number is already existing
			if field.Number != 0 {
				solvedNumberCount++

				// There are 81 fields, if all numbers are solved we're done
				if solvedNumberCount == 81 {
					return board, nil
				}

				// Skip solved numbers
				continue
			}

			// Reset the possible numbers for the current field
			field.PossibleNumbers = []uint{}

			// Loop over all numbers to see if the number can fit in the given field
			for i := 1; i <= 9; i++ {
				if !board.HasNumberInBox(*field, uint(i)) &&
					!board.HasNumberInColumn(*field, uint(i)) &&
					!board.HasNumberInRow(*field, uint(i)) {
					field.PossibleNumbers = append(field.PossibleNumbers, uint(i))
				}
			}

			// Increase the last solve iteration count
			iterationsSinceLastSolve++

			// If there's only one possible number it's the solution
			if len(field.PossibleNumbers) == 1 {
				field.Number = int(field.PossibleNumbers[0])
				field.PossibleNumbers = []uint{}
				iterationsSinceLastSolve = 0
			}

			// Check to see if we don't have too many iterations since the last solve.
			// This is to avoid infinite loops for unsolvable boards
			if iterationsSinceLastSolve > maxIterationsSinceLastSolve {
				return board, errors.New("too many iterations since last solve. Might be an unsolvable board")
			}
		}
	}
}

// PrintBoard prints the given board data in a readable sudoku format
func PrintBoard(board Board) {
	for i, field := range board.Fields {
		if i % 3 == 0 && i != 0 {
			fmt.Print("|")
		}
		if i % 9 == 0  {
			fmt.Println("\r\n|---+---+---|")
			fmt.Print("|")
		}
		fmt.Print(field.Number)
	}
	fmt.Print("|")
	fmt.Println("\r\n|---+---+---|")

}

// contains is a little helper method to check if the given number is in the given slice
func contains(slice []uint, nr int) bool {
	for _, a := range slice {
		if int(a) == nr {
			return true
		}
	}
	return false
}

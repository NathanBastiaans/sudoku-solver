package main

import (
	"log"
	"os"
	"sudoku-solver/src/business"
)

// main starts the solver
func main() {

	// Get the file from the command line arguments, if not present this will fatally crash the program
	fileName := os.Args[1]

	// Generate the board
	board, err := business.GenerateBoardFromFile(fileName)

	if err != nil {
		log.Fatal(err)
	}

	// Solve the board
	board, err = business.Solve(board)

	if err != nil {
		log.Fatal(err)
	}

	// Print the board
	business.PrintBoard(board)

}

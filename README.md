# Sudoku solver #
This is a tool to solve sudoku puzzles

## How do I get set up? ##
To run the tool you can run the following command
```
go run main.go <file>
```

To run the tool with the demo sudoku run the following command 
```
go run main.go demo.txt
```

This will run the tool and will print the solution to the puzzle in the commandline. 
If a file is not passed to the commandline it will fatally crash the program.

## Sudoku files 
The given files have to be in the specified format: 
 - There should only be 9 characters per line 
 - All characters should be numbers 
 - Empty fields should be 0
 - The total characters should be 81 (not counting new lines)

Also check the `demo.txt` file to see what a correct file looks like
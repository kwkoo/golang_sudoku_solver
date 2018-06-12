# Go Sudoku Solver

This app loads a puzzle from <http://davidbau.com/generated/sudoku.txt>.

It then solves the puzzle using backtracking.

There's also a web interface which listens on port 8080.

![screenshot](images/solver.gif)

The Gorilla WebSocket library is used to send update events from the server to the web client as the puzzle is being solved.

Don't forget to include the `--recurse-submodules` option when cloning the repository.

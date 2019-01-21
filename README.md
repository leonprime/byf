# DLX Beat Your Father

An application of Knuth's Dancing Links paper to a particular 2D and 3D tiling
game.

## The paper

https://arxiv.org/abs/cs/0011047

## The game

**Gagne Ton Papa**, or **Beat Your Father** is a French board game that includes
a number of pentomino-like pieces.  The manual lists a number of fun games that
can be played with the pieces.  Among them are 2D and 3D tiling puzzles that
challenge the player to fill a space with various selections of the pieces.

The primary game is a two player race to see who can solve their 2D tiling game
first.  To scale for age and skill differences, a player can choose a larger
board with a more gnarly piece selection.  The game comes with a 5x13 board and
dividers and cards that provide piece selections for games up to 6x5 in size.

## The motivation

What the manual doesn't list are all the possible tiling games and what pieces
would form a viable solution set for them.  One notable omission is that it
lacks a reference solution to the entire 5x13 game board, which can be handy
when putting the game away.  Another glaring omission is the lack of a 3D tiling
game that uses as many as the pieces as possible.  There are 67 unit cubes of
pieces available, but the largest 3D game is only 45 unit cubes.

This program is an implementation of Knuth's techniques to generate solutions
for all 2D and 3D tiling games over a given set of pieces.

## Pieces

To see the available pieces,

    ./byf -show

The game comes with the pieces specified by `ooOvVzZiiIlLnpstrY`.

Use the `-pieces` argument to point to a different data file.  Pentominos are
provided in `data/pentominoes.txt`.

## Examples

### Put the game away

Solving the 5x13 tiling is straightforward.  We choose 65 unit cubes by tossing
out the `i` piece.  The number of solutions is huge, so we limit the number of
solutions to the first 10 found,

    ./byf -max 10 5 13 ooOvVzZiIlLnpstrY

### Largest cube

64 unit cubes can be used to build a 4x4x4 cube.  We just need to toss out 3
unit cubes of pieces, which can be the `o` and `i` pieces. Then we can generate
a sample of solutions as follows,

    ./byf -max 10 4 4 4 oOvVzZiIlLnpstrY


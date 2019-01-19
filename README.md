# DLX Beat Your Father

An application of Knuth's Dancing Links paper to a particular 2D and 3D tiling game.

## The paper

https://arxiv.org/abs/cs/0011047

## The game

**Gagne Ton Papa**, or **Beat Your Father** is a French board game that includes
a number of pentomino-like pieces.  The manual lists a number of fun games that
can be played with the pieces.  Among them are 2D and 3D tiling puzzles that
challenge the player to fill a space with various selections of the pieces.

## The motivation

What the manual doesn't list are all the possible tiling games and what pieces
would form a viable solution set for them.  One notable omission is that it
lacks a reference solution to the entire 5x13 game board, which can be handy
when putting the game away.  Another glaring omission is the lack of a 3D tiling
game that uses as many as the pieces as possible.  There are 67 unit cubes of
pieces available, but the largest 3D game is only 45 unit cubes.

This program uses Knuth's techniques to generate all 2D and 3D tiling games.

## The program

TODO: words about the program

## Pieces

TODO: words about the pieces and describe them somehow

## Examples

### Put the game away

Solving the 5x13 tiling is straightforward.  We choose 65 unit cubes by tossing
out the `i` piece.  The number of solutions is huge, so we only pick a small
sample,

    ./byf -max 10 5 13 ooOvVzZiIlLnpstrY

### Largest cube

64 unit cubes can be used to build a 4x4x4 cube.  We just need to toss out 3
unit cubes of pieces, which can be the `o` and `i` pieces. Then we can generate
a sample of solutions as follows,

    ./byf -max 10 4 4 4 oOvVzZiIlLnpstrY

### Torture your friend

Want to test the strength of a friendship?  Let's find some set of pieces for
a tiling puzzle that doesn't have a solution.

TODO: mechanism

### Hardest puzzle

Bigger puzzles are more likely to have a solution for any subset of pieces, but
we can still find the *hardest* puzzle simply by looking for a set of pieces
that produce the smallest number of solutions.

TODO: this

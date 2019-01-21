# Beat Your Father with DLX

An application of Knuth's Dancing Links algorithm to a particular 2D and 3D tiling
game named **Beat Your Father**.

## The technique

* https://arxiv.org/abs/cs/0011047
* https://www.youtube.com/watch?v=t9OcDYfHqOk

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
would form a viable solution set for them.  One glaring omission is that it
lacks a reference solution to the entire 5x13 game board, which can be handy
when putting the game away.  Another is the lack of a 3D tiling
game that uses as many of the pieces as possible.  There are 67 unit cubes of
pieces available, but the largest 3D game is only 45 unit cubes.

This program is an implementation of Knuth's techniques to generate solutions
for all 2D and 3D tiling games over a given set of pieces.  It is then used to
correct these omissions for fun.

## Pieces

To see the available pieces,

    ./byf -show

The default pieces are for the **Beat Your Father** game, which comes with the
pieces `ooOvVzZiiIlLnpstrY`.  Most games are played with some subset of this.

Pentominos are provided in `data/pentominoes.txt` and are numbered `FILNPTUVWXYZ`.

Use the `-pieces` argument to point to a different data file.

## Examples

### Put the game away

Solving the 5x13 tiling is straightforward.  We choose 65 unit cubes by tossing
out one of the the `i` piece.  The number of solutions is huge, so we limit the
number of solutions to the first one found,

    ./byf -max 1 13 5 ooOvVzZiIlLnpstrY

![A solution to 13x5 Beat Your Father](./docs/13x5_ooOvVzZiIlLnpstrY/0.png "Logo Title Text 1")

### Largest cube

64 unit cubes can be used to build a 4x4x4 cube.  We just need to toss out 3
unit cubes of pieces, which can be the duplicate `o` and `i` pieces.  To see
if it's possible, let's print some solutions,

    ./byf 4 4 4 oOvVzZiIlLnpstrY

There's a lot of solutions, so this will take awhile.  TODO XXX specify how
much.  Press `CTRL-C` to interrupt and print the solutions it found so far,

### Verify pentominoes

On the [wiki for pentominoes](https://en.wikipedia.org/wiki/Pentomino), a number
of 2D and 3D solutions are provided.  Let's verify that this program can
generate the same results.

#### 2D pentominoes

We can count and verify the 2 solutions to the 3x20 case.  We'll lay it out as
20x3 so the output is horizontal.

    ./byf -pieces data/pentominoes.txt 20 3 FILNPTUVWXYZ
    game "20x3_FILNPTUVWXYZ" has 8 solutions
        time taken: 105.598042ms
        steps: 461658
    wrote all 8 solutions to ./solutions/20x3_FILNPTUVWXYZ

There are 8 solutions because the chiral pieces were included, meaning we count the mirror images of piece F, P, L, etc.  To leave them out, we can use the `-nochiral` argument,

    ./byf -pieces data/pentominoes.txt -nochiral 20 3 FILNPTUVWXYZ
    game "20x3_FILNPTUVWXYZ" has 4 solutions
        time taken: 11.440136ms
        steps: 58022
    wrote all 4 solutions to ./solutions/20x3_FILNPTUVWXYZ

Visual inspection will confirm that there are 2 unique solutions and each has a duplicate rotated 180 degrees.

Here, we render one solution for 5x12 and visually confirm it,

    ./byf -pieces data/pentominoes.txt -max 1 12 5 FILNPTUVWXYZ


![A solution to 12x5 pentominoes](./docs/12x5_FILNPTUVWXYZ/0.png "Logo Title Text 1")

#### 3D pentominoes

We can verify the solution counts match by dividing by the 8-fold symmetry.

    ./byf -print 0 -pieces data/pentominoes.txt 2 3 10 FILNPTUVWXYZ
    game "2x3x10_FILNPTUVWXYZ" has 96 solutions
        time taken: 1.601036413s
        steps: 7089848

    ./byf -print 0 -pieces data/pentominoes.txt 2 5 6 FILNPTUVWXYZ
    game "2x5x6_FILNPTUVWXYZ" has 2112 solutions
        time taken: 13.341344262s
        steps: 47928520

    ./byf -print 0 -pieces data/pentominoes.txt 3 4 5 FILNPTUVWXYZ
    game "3x4x5_FILNPTUVWXYZ" has 31520 solutions
        time taken: 2m27.858245257s
        steps: 647787028

This is my 3rd time implementing a mastermind solver, and my first real project
in go. It answers a very narrow problem: given a current game, what is/are the
next optimal guess/es?

Warning: For some settings, the runtime can be significantly long.

Currently the game to analyze is hardcoded. The original idea was that the
program would take the path to a JSON file as its main argument.

Setting up a game looks something like this:

```go
g := game{
  colors: chars("bdglpry"),
  slots:  6,
  rounds: []round{
    round{chars("bbbddd"), 0, 1},
    round{chars("ggglll"), 2, 1},
    round{chars("ppprrr"), 0, 2}}}
```

and produces something that looks a bit like this:

```
After guess "bbbddd" 18750/117649 combinations remain.
After guess "ggglll" 1734/117649 combinations remain.
After guess "ppprrr" 180/117649 combinations remain.

There are 18 optimal guesses: gddplp grrbbl dgdlpp rrgbbl rgrlbb dgdplp rrglbb
rgrbbl ddgppl ddglpp rrgblb gddppl ddgplp rgrblb dgdppl grrlbb grrblb gddlpp
```

I'm using `b`, `d`, `g`, `l`, `p`, `r`, and `y` for `blue`, `dark`, `green`,
`light`, `purple`, `red`, and `yellow`.

The `nextGuesses()` func is a beautiful opportunity for some go-style
concurrency/parallelism magic... but I haven't done that yet. So far waiting
has been fine, and anyway this go version beats the pants off the python
version in the speed department.

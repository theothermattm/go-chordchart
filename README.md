# Go ChordChart

A Simple CLI Utility to output guitar chord fingerings.

```
 ./go-chordchart -c C
Searching chord(s): {C}

Chord:  C         or   Cmaj
Type:   major triad
Notes:  C  E  G

          E  A  D  G  B  e
          ----------------
          0  3  2  0  1  0
          0  3  5  5  5  3
          3  3  2  0  1  0
          3  x  2  0  1  0
          x  3  2  0  1  0
          x  3  5  5  5  0

```

Written in Go as a learning project.

Chord source is [http://gospelmusic.org.uk/resources/chord_chart_big.htm](http://gospelmusic.org.uk/resources/chord_chart_big.htm)

# Downloading and Running (Mac)

```
curl -L https://github.com/theothermattm/go-chordchart/raw/main/dist/go-chordchart-macos > go-chordchart && chmod +x go-chordchart
./go-chordchart
```

This will show you help. An example command is above.

# Building

```
go get
rice embed-go
go build
```

package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	log	"github.com/sirupsen/logrus"
	flag "github.com/ogier/pflag"
)

type chord struct {
	name      string
	notes     string
	fingering string
	chordType string
}

var (
	inputChord   string
	mappedChords map[string]chord
)

func main() {
	flag.Parse()
	// if user does not supply flags, print usage
	// we can clean this up later by putting this into its own function
	if flag.NFlag() == 0 {
		fmt.Printf("Usage: %s [options]\n", os.Args[0])
		fmt.Println("Options:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Printf("Searching chord(s): %s\n", inputChord)
	fmt.Print(mappedChords[inputChord])
}

func parseChordFile() {
	file, err := os.Open("./chords.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var text = scanner.Text()
		// Sample line: `A         or   Amaj      [0 0 2 2 2 0] (Db E  A) : major triad`
		var re = regexp.MustCompile(`^(.*)\[(.*)\].*\((.*)\).*:(.*)`)
		result := re.FindAllStringSubmatch(text, -1)
		// TODO good lord clean this up!
		for _, m := range result {
			if m[1] != "" || m[2] != "" {
					var chordName = m[1]
					if strings.Index(m[1], "or") > -1 {
						var reOr = regexp.MustCompile(`^(.*)or(.*)`)
						reOrResult := reOr.FindAllStringSubmatch(chordName, -1)
						for _, n := range reOrResult{
							var chordKey = strings.Trim(strings.ToUpper(n[1]), " ")
							newChord := chord{name: n[1], fingering: m[2], notes: m[3], chordType: m[4]}
							log.Debug("Keying new chord with OR: {" + chordKey + "}\n")
							mappedChords[chordKey] = newChord
						}
					} else {
						var chordKey = strings.Trim(strings.ToUpper(m[1]), " ")
						log.Debug("Keying new chord: {" + chordKey + "}\n")
						newChord := chord{name: m[1], fingering: m[2], notes: m[3], chordType: m[4]}
						mappedChords[chordKey] = newChord
					}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}

func init() {
	// TODO add a debug flag
	log.SetLevel(log.DebugLevel)
	flag.StringVarP(&inputChord, "chord", "c", "", "Search Chords")
	mappedChords = make(map[string]chord)
	parseChordFile()
}

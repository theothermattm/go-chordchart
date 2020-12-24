package main

import (
	"bufio"
	"fmt"
	flag "github.com/ogier/pflag"
	log "github.com/sirupsen/logrus"
	"os"
	"regexp"
	"strings"
)

type chord struct {
	name      string
	notes     string
	fingering string
	chordType string
}

var (
	inputChord   string
	debugMode    string
	mappedChords map[string][]chord
)

func main() {

	// some of this taken from https://www.freecodecamp.org/news/writing-command-line-applications-in-go-2bc8c0ace79d/
	flag.Parse()
	// if user does not supply flags, print usage
	// we can clean this up later by putting this into its own function
	if flag.NFlag() == 0 {
		fmt.Printf("Usage: %s [options]\n", os.Args[0])
		fmt.Println("Options:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if len(debugMode) > 0 {
		log.SetLevel(log.DebugLevel)
	}

	fmt.Printf("Searching chord(s): %s\n\n", inputChord)
	var transformedInput = strings.Trim(strings.ToUpper(inputChord), " ")
	fmt.Println(formatOutput(mappedChords[transformedInput]))
}

func formatOutput(chords []chord) string {

	var output = "No chords found"
	if len(chords) > 1 {
		output = ""
		output = output + "Chord: \t\t" + chords[1].name + "\n"
		output = output + "Type: \t\t" + chords[1].chordType + "\n"
		output = output + "Notes: \t\t" + chords[1].notes + "\n\n"
		for _, chord := range chords {
			output = output + formatFingering(chord.fingering) + "\n"
		}
	}

	return output
}

func formatFingering(fingering string) string {
	var values = strings.Split(fingering, " ")
	var newResult = ""
	for _, value := range values {
		if len(value) > 1 {
			newResult = newResult + " " + value
		} else {
			newResult = newResult + "  " + value
		}

	}

	return newResult
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
		// Sample line from file:
		// A         or   Amaj      [0 0 2 2 2 0] (Db E  A) : major triad
		var re = regexp.MustCompile(`^(.*)\[(.*)\].*\((.*)\).*:(.*)`)
		result := re.FindAllStringSubmatch(text, -1)
		// TODO good lord clean this up!
		for _, m := range result {
			var chordName = m[1]
			if strings.Index(m[1], "or") > -1 {
				var reOr = regexp.MustCompile(`^(.*)or(.*)`)
				reOrResult := reOr.FindAllStringSubmatch(chordName, -1)
				for _, n := range reOrResult {
					var chordKey = strings.Trim(strings.ToUpper(n[1]), " ")
					newChord := chord{name: strings.Trim(m[1], " "),
						fingering: strings.Trim(m[2], " "), notes: strings.Trim(m[3], " "), chordType: strings.Trim(m[4], " ")}
					log.Debug("Keying new chord with OR: {" + chordKey + "}\n")
					mappedChords[chordKey] = append(mappedChords[chordKey], newChord)
				}
			} else {
				var chordKey = strings.Trim(strings.ToUpper(m[1]), " ")
				log.Debug("Keying new chord: {" + chordKey + "}\n")
				newChord := chord{name: m[1], fingering: m[2], notes: m[3], chordType: m[4]}
				mappedChords[chordKey] = append(mappedChords[chordKey], newChord)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}

func init() {
	// TODO add a debug flag
	flag.StringVarP(&inputChord, "chord", "c", "", "Search Chords")
	flag.StringVarP(&debugMode, "debug", "d", "", "Debug Mode")
	mappedChords = make(map[string][]chord)
	parseChordFile()
}

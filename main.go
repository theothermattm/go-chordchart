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

const NoteSpacing = "  "

type chord struct {
	name      string
	notes     string
	fingering string
	chordType string
}

var (
	inputChord   string
	debugMode    bool
	mappedChords map[string] []chord
)

func main() {
	var transformedInput = strings.Trim(strings.ToUpper(inputChord), " ")
	fmt.Printf("Searching chord(s): {%s}\n\n", transformedInput)
	fmt.Println(formatOutput(mappedChords[transformedInput]))
}

func formatOutput(chords []chord) string {

	var output = "No chords found"
	if len(chords) >= 1 {
		output = ""
		output += "Chord: \t" + chords[0].name + "\n"
		output += "Type: \t" + chords[0].chordType + "\n"
		output += "Notes: \t" + chords[0].notes + "\n\n"
		notes := "EADGBe"
		notesArray := strings.Split(notes, "")

		output += "\t"
		for _, note := range notesArray {
			output += NoteSpacing + note
		}
		output += "\n"

		output += "\t" + NoteSpacing

		dividerString := ""
		for i := 1; i <= len(notes) ; i++ {
			 dividerString += " " + NoteSpacing
		}
		dividerString = strings.ReplaceAll(dividerString, " ", "-")
		dividerString = dividerString[:len(dividerString)-len(NoteSpacing)]
		output += dividerString + "\n"

		for _, chord := range chords {
			output += "\t" + formatFingering(chord.fingering) + "\n"
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
		// TODO clean this up with some functions
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
				newChord := chord{name: strings.Trim(m[1], " "), fingering: strings.Trim(m[2], " "), 
					notes: strings.Trim(m[3], " "), chordType: strings.Trim(m[4], " ")}
				mappedChords[chordKey] = append(mappedChords[chordKey], newChord)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}

func init() {
	flag.StringVarP(&inputChord, "chord", "c", "", "Search Chords")
	flag.BoolVarP(&debugMode, "debug", "d", false, "Debug Mode")

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

	if debugMode {
		log.Info("Setting debug mode")
	  log.SetLevel(log.DebugLevel)
	}

	mappedChords = make(map[string][]chord)
	parseChordFile()
}

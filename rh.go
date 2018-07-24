package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type sentenceBuzzwordCount struct {
	sentence      string
	buzzwordCount int
}

//byCount implements sort.Interface for []sentenceBuzzwordCount based on the count
//of buzzwords
type byCount []sentenceBuzzwordCount

func (a byCount) Len() int {
	return len(a)
}
func (a byCount) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a byCount) Less(i, j int) bool {
	return a[i].buzzwordCount < a[j].buzzwordCount
}

var buzzwords = []string{"clean", "training", "personal", "hours", "weights", "cardio",
	"price", "parking", "crowd", "locker", "sauna", "steam", "room",
	"pros", "staff", "nice", "location", "pool", "time", "weight",
	"muscle", "professional", "convenient", "fee", "goal", "member",
	"towel", "equipment", "jacuzzi", "space", "great"}

func main() {

	filename := os.Args[1]
	max, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var reviews []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		reviews = append(reviews, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	highlights := review_highlights(reviews, max)

	for _, highlight := range highlights {
		fmt.Println(highlight)
		fmt.Println()
	}
}

func review_highlights(reviews []string, max int) (highlights []string) {

	var potentialHighlights []sentenceBuzzwordCount

	//iterate through reviews
	for _, review := range reviews {
		//split each review into individual sentences
		sentences := strings.FieldsFunc(review, splitOnPunctuation)

		for _, sentence := range sentences {
			numOfBuzzwords := getBuzzwordCount(sentence, buzzwords)

			//disregard sentences with no buzzwords
			if numOfBuzzwords > 0 {
				potential := sentenceBuzzwordCount{
					sentence:      sentence,
					buzzwordCount: numOfBuzzwords,
				}

				potentialHighlights = append(potentialHighlights, potential)
			}
		}
	}

	highlights = getHighlightsFromPotentials(potentialHighlights, max)
	return highlights
}

func getBuzzwordCount(sentence string, buzzwords []string) int {
	count := 0

	sentence = strings.ToLower(sentence)
	for _, buzzword := range buzzwords {
		if strings.Contains(sentence, buzzword) {
			count++
		}
	}
	return count
}

func getHighlightsFromPotentials(potentialHighlights []sentenceBuzzwordCount, max int) (highlights []string) {
	//sort array of potential highlights so those with most buzzwords come first
	sort.Sort(sort.Reverse(byCount(potentialHighlights)))

	//if the number of highlights requested is greater than the number of potential
	//highlights, set the counter = the number of potential highlights
	var counter = max
	if max > len(potentialHighlights) {
		counter = len(potentialHighlights)
	}

	for i := 0; i < counter; i++ {
		highlights = append(highlights, format(potentialHighlights[i].sentence))
	}
	return highlights
}

//trim leading space from sentence and add period
func format(sentence string) string {
	return strings.TrimSpace(sentence) + "."
}

//helper func to split reviews into sentences
func splitOnPunctuation(c rune) bool {
	return c == '.' || c == '!' || c == '?'
}

package streamer

import (
	"os"
	"bufio"
)

/**
Read lines from a file specified by its path.
It may return a list of strings (one line for each position) and an error that indicates if a problem has occurred.
 */
func LoadTextFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

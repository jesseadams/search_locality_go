package main

import(
  "fmt"
  "path/filepath"
  "os"
  "github.com/alexflint/go-arg"
  "bufio"
  "regexp"
)

var args struct {
  MaxDistance int `arg:"required"`
  TermOne string `arg:"positional,required"`
  TermTwo string `arg:"positional,required"`
}

var matches []string

func searchFile(path string, termOne string, termTwo string) bool {
  matchTermOne := false
  matchTermTwo := false
  regexOne := regexp.MustCompile("(?i)" + termOne)
  regexTwo := regexp.MustCompile("(?i)" + termTwo)

  file, err := os.Open(path)
  if err != nil {
    return false
  }

  scanner := bufio.NewScanner(file)
  scanner.Split(bufio.ScanWords)

  count := 1

  for scanner.Scan() {
    word := scanner.Text()

    if matchTermOne && count <= args.MaxDistance {
      matchResult := regexTwo.MatchString(word)
      if matchResult {
        return true
      }
      count++
    } else if matchTermTwo && count <= args.MaxDistance {
      matchResult := regexOne.MatchString(word)

      if matchResult {
        return true
      }
      count++
    } else {
      matchTermOne = false
      matchTermTwo = false

      count = 1
      matchResult := regexOne.MatchString(word)
      if matchResult {
        matchTermOne = true
      } else {
        matchResult := regexTwo.MatchString(word)
        if matchResult {
          matchTermTwo = true
        }
      }
    }
  }

  return false
}

func searchFiles(path string, file os.FileInfo, err error) error {

  if !file.IsDir() {
    match := searchFile(path, args.TermOne, args.TermTwo)
    if match {
      matches = append(matches, path)
    }
  }
  return nil
}

func main() {
  arg.MustParse(&args)
  fmt.Printf("Max Distance: %d, Terms: %s, %s\n", args.MaxDistance, args.TermOne, args.TermTwo)

  err := filepath.Walk("documents", searchFiles)
  if err != nil {
    fmt.Printf("Error! %s" , err)
  }

  fmt.Printf("Matches: %d\n", len(matches))
  for _,match := range(matches) {
    fmt.Printf("Matched: %s\n", match)
  }
}

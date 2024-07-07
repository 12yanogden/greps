/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/12yanogden/cat"
	"github.com/12yanogden/errors"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "greps",
	Short: "Global Regular Expression Print Section",
	Long: `A shell command that prints the lines between two regular expressions given.
The command can search a file or read a stirng from a pipe.`,
	Args: cobra.MinimumNArgs(2),
	Run:  greps,
}

func isInputFromPipe() bool {
	fileInfo, _ := os.Stdin.Stat()
	return fileInfo.Mode()&os.ModeCharDevice == 0
}

func getSearchTextFromPipe(cmd *cobra.Command) string {
	searchText, err := bufio.NewReader(cmd.InOrStdin()).ReadString('\n')
	errors.Peek(err)

	return strings.TrimSpace(searchText)
}

func getSearchTextFromFile(file string) string {
	searchText, err := cat.Cat(file)
	errors.Peek(err)

	return searchText
}

func getSearchText(cmd *cobra.Command) string {
	searchText := ""
	isInputFromPipe := isInputFromPipe()
	file, err := cmd.Flags().GetString("file")
	errors.Peek(err)

	// Validate input
	if !isInputFromPipe && file == "" {
		panic("greps: no search text given")
	} else if isInputFromPipe && file != "" {
		panic("greps: cannot search pipe and file input")
	}

	// Get search text from input
	if isInputFromPipe {
		searchText = getSearchTextFromPipe(cmd)
	} else {
		searchText = getSearchTextFromFile(file)
	}

	return searchText
}

func compileRegex(regex string) *regexp.Regexp {
	compiledRegex, err := regexp.Compile(regex)
	errors.Peek(err)

	return compiledRegex
}

func search(searchText string, regex *regexp.Regexp) []int {
	location := regex.FindIndex(([]byte(searchText)))

	if len(location) == 0 {
		panic("greps: no match found for regex: " + regex.String())
	}

	return location
}

func greps(cmd *cobra.Command, args []string) {
	regex1 := compileRegex(args[0])
	regex2 := compileRegex(args[1])
	searchText := getSearchText(cmd)
	originalSearchText := searchText

	// Search for first regex
	location1 := search(searchText, regex1)

	// Update search text
	searchText = searchText[location1[1]:]

	// Search for second regex
	location2 := search(searchText, regex2)

	// Align location2 with original search text
	location2[0] += location1[1]
	location2[1] += location1[1]

	// Print result
	fmt.Println(originalSearchText[location1[0]:location2[1]])
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("file", "f", "", "Specify a file to search")
}

package main

import (
	"encoding/json"
	"fmt"
	"github.com/gernest/front"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"github.com/spf13/viper"
)

type Config struct {
	SourcePath string
	DestPath   string
}

type SnippetInfo struct {
	DirName     string
	SnippetName string
	FullPath    string
	Content     string
	Data        map[string]interface{}
}

type SnippetJSON struct {
	Prefix      string   `json:"prefix"`
	Scope       string   `json:"scope"`
	Body        []string `json:"body"`
	Description string   `json:"description"`
}

func main() {
	// Initialize the frontmatter parser
	frontMatter := front.NewMatter()
	frontMatter.Handle("---", front.YAMLHandler)

	config := loadConfig()
	snippets := loadSnippets(config.SourcePath, frontMatter)
	writeSnippets(snippets, config.DestPath)
}

func loadConfig() Config {
	// Look for config in the home directory and the current
	// directory
	viper.SetConfigName("vscode-snippets")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")

	// Read the config file
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Error reading config file.\n%s\n\n"+
			"Please create a config file (vscode-snippets.yaml) in your home directory", err))
	}

	sourcePath := viper.GetString("source_path")
	if !pathExists(sourcePath) {
		panic(fmt.Errorf("source_path not specified in config file"))
	}

	destPath := viper.GetString("dest_path")
	if !pathExists(destPath) {
		panic(fmt.Errorf("source_path not specified in config file"))
	}

	return Config{
		SourcePath: sourcePath,
		DestPath:   destPath,
	}
}

func loadSnippets(snippetsPath string, frontMatter *front.Matter) map[string][]SnippetInfo {
	snippets := make(map[string][]SnippetInfo)

	filepath.Walk(snippetsPath, func(path string, f os.FileInfo, err error) error {
		// Ignore directories
		if f.IsDir() {
			return nil
		}

		// Derive the folder name and snippet name from the snippet's
		// relative path
		//
		// Relative path will always be something like "ruby/rails.snippets"
		// The dirName then becomes "ruby", and the snippetName "rails.snippets"
		relPath, _ := filepath.Rel(snippetsPath, path)
		dirName := filepath.Dir(relPath)
		snippetName := filepath.Base(relPath)

		// Ignore files that are not snippets
		if strings.HasPrefix(snippetName, ".") {
			return nil
		}

		// Read the snippet file and get the parsed results
		contentBytes, _ := ioutil.ReadFile(path)
		contentString := string(contentBytes)
		metadata, snippetBody, _ := frontMatter.Parse(strings.NewReader(contentString))

		// Create the snippet
		snippetInfo := SnippetInfo{DirName: dirName, SnippetName: snippetName,
			FullPath: path,
			Content: snippetBody,
			Data: metadata,
		}

		// Append the snippet to the list
		snippets[dirName] = append(snippets[dirName], snippetInfo)

		return nil
	})

	return snippets
}

func writeSnippets(snippets map[string][]SnippetInfo, snippetsWritePath string) {
	for folderName, folderSnippets := range snippets {
		jsonFileData := make(map[string]SnippetJSON)

		for _, snippet := range folderSnippets {
			// Prepare the various fields of the snippet
			description, _ := snippet.Data["desc"].(string)
			snippetIdentifier := "vs/" + folderName + "/" + snippet.SnippetName
			snippetScope := ""
			snippetPrefix := "."

			// Create the snippetJSON object
			snippetJson := SnippetJSON{
				Prefix:      snippetPrefix,
				Scope:       snippetScope,
				Body:        strings.Split(snippet.Content, "\n"),
				Description: description,
			}

			jsonFileData[snippetIdentifier] = snippetJson
		}

		// Marshall the data with two-space indents
		jsonDataMarshalled, _ := json.MarshalIndent(jsonFileData, "", "  ")

		// Write the file
		jsonOutputFileName := folderName + ".code-snippets"
		ioutil.WriteFile(
			filepath.Join(snippetsWritePath, jsonOutputFileName),
			jsonDataMarshalled,
			0644,
		)

		// Print it out
		fmt.Println(filepath.Join(snippetsWritePath, jsonOutputFileName))
	}
}

func pathExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}

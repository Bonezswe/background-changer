package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var (
	version string
)

type Mode int

const (
	Center Mode = iota
	Crop
	Fit
	Span
	Stretch
	Tile
)

func init() {
	setCmd.Flags().String("url", "", "URL of image")
	setCmd.Flags().Bool("save", false, "Bool to save image to default path")
}

func main() {
	var rootCmd = &cobra.Command{
		Use: "cbg",
	}

	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(versionCmd)

	configCmd.AddCommand(getConfigCmd)
	configCmd.AddCommand(setPathCmd)
	configCmd.AddCommand(addUrlCmd)
	configCmd.AddCommand(removeUrlCmd)

	rootCmd.Execute()

}

// Get, set, config, search, list, version

type Picture struct {
	Name string
	Path string
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get current background",
	Run: func(cmd *cobra.Command, args []string) {
		current, err := Get()

		if err != nil {
			panic("Get current failed")
		}

		fmt.Printf("Current Background: %s", current)
	},
}

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set background image",
	Run: func(cmd *cobra.Command, args []string) {
		url, _ := cmd.Flags().GetString("url")
		save, _ := cmd.Flags().GetBool("save")

		if url == "" {
			panic("No URL provided")
		}

		if save {
			file, err := DownloadAndSave(url)

			if err != nil {
				panic(err)
			}

			config := ReadCofnig()
			defaultPath := config.DefaultPath
			filename := file.Name()

			SetFromFile(
				fmt.Sprintf("%s/%s", defaultPath, filename),
			)
		} else {
			SetFromFile(url)
		}
	},
}

var configCmd = &cobra.Command{
	Use:   "config [sub]",
	Short: "Config commands",
}

var getConfigCmd = &cobra.Command{
	Use:   "show",
	Short: "Prints current config settings",
	Run: func(cmd *cobra.Command, args []string) {
		config := ReadCofnig()
		path := config.DefaultPath
		urls := config.URLStore

		fmt.Printf("Path: %s\n Urls: %s\n", path, urls)
	},
}

var setPathCmd = &cobra.Command{
	Use:   "set",
	Short: "update default path",
	Run: func(cmd *cobra.Command, args []string) {
		setDefaultPath(args[0])
	},
}

var addUrlCmd = &cobra.Command{
	Use:   "addUrl",
	Short: "Add a website to search list",
	Run: func(cmd *cobra.Command, args []string) {
		AddUrlToStore(args[0])
	},
}

var removeUrlCmd = &cobra.Command{
	Use:   "removeUrl",
	Short: "Remove an existing url from url store",
	Run: func(cmd *cobra.Command, args []string) {
		RemoveUrlFromStore(args[0])
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List images",
	Run: func(cmd *cobra.Command, arg []string) {
		pics := []Picture{}

		files, _ := os.Open("C:\\Users\\bonezswe\\Pictures\\backgrounds")

		defer files.Close()

		fileInfo, _ := files.ReadDir(-1)

		for _, file := range fileInfo {
			name := file.Name()
			fmt.Println(name)

			pics = append(
				pics,
				Picture{
					Name: name,
					Path: fmt.Sprintf("C:\\Users\\bonezswe\\Pictures\\backgrounds\\%s", name),
				},
			)
		}

		searcher := func(input string, index int) bool {
			pic := pics[index]
			name := strings.Replace(strings.ToLower(pic.Name), " ", "", -1)
			path := strings.Replace(strings.ToLower(input), " ", "", -1)

			return strings.Contains(name, path)
		}

		templates := &promptui.SelectTemplates{
			Label:    "{{ .Name }}?",
			Active:   "{{ .Name | cyan }}",
			Inactive: "{{ .Name | white }}",
			Selected: "{{ .Name | green }}",
			Details: `
			---Pictures---
			{{ "Name:" | faint }}    {{ .Name }}
			`,
		}

		prompt := promptui.Select{
			Label:     "Select image",
			Items:     pics,
			Searcher:  searcher,
			Templates: templates,
			Size:      4,
		}

		i, _, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		fmt.Printf("Selected %q: ", pics[i].Name)

		SetFromFile(pics[i].Path)
	},
}

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search websites for image",
	Run: func(cmd *cobra.Command, args []string) {
		// resolution := GetScreenResolution()
		config := ReadCofnig()
		urls := config.URLStore

		prompt := promptui.Select{
			Label: "Select a site: ",
			Items: urls,
		}

		_, result, err := prompt.Run()

		if err != nil {
			return
		}

		OpenBrowser(result)

		urlPrompt := promptui.Prompt{
			Label: "Enter url: ",
		}

		imageUrl, err := urlPrompt.Run()

		if err != nil {
			return
		}

		SetFromUrl(imageUrl)
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print current version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Cbg version: %s\n", version)
	},
}

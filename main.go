package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/ncruces/zenity"
)

func SortFileNameAscend(files []fs.DirEntry) {
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})
}

func SortFileNameDescend(files []fs.DirEntry) {
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() > files[j].Name()
	})
}

func main() {
	inDir, err := zenity.SelectFile(
		zenity.Filename(""),
		zenity.Directory(),
		zenity.DisallowEmpty(),
		zenity.Title("Select input directory"),
	)
	if err != nil {
		zenity.Error(err.Error(),
			zenity.Title("Error"),
			zenity.ErrorIcon,
		)
		log.Fatal(err)
	}

	// ask for a name (must have consecutive #)
	replaceFileName, err := zenity.Entry("Enter new formatted name (must have consecutive # symbols):",
		zenity.Title("Add a new entry"))
	if err != nil {
		zenity.Error(err.Error(),
			zenity.Title("Error"),
			zenity.ErrorIcon,
		)
		log.Fatal(err)
	}

	// get hash and format for renaming
	var replaceHashStr string
	for h := 0; h < strings.Count(replaceFileName, "#"); h++ {
		replaceHashStr = replaceHashStr + "#"
	}
	leadZero := "%0" + strconv.Itoa(strings.Count(replaceFileName, "#")) + "d"

	// get list of file names in folder
	files, err := os.ReadDir(inDir)
	SortFileNameAscend(files)
	if err != nil {
		log.Fatal(err)
	}

	i := 1
	ignore := []string{".DS_Store", "._.DS_Store"}
	for _, oldFileName := range files {
		if !oldFileName.IsDir() && !contains(ignore, oldFileName.Name()) && !strings.HasPrefix(oldFileName.Name(), ".") {
			newFileName := strings.Replace(replaceFileName, replaceHashStr, fmt.Sprintf(leadZero, i), -1) + filepath.Ext(oldFileName.Name())

			fmt.Println("----------------------------------------------------------------------------------------")
			fmt.Println("old:", filepath.Join(inDir, oldFileName.Name()))
			fmt.Println("new:", filepath.Join(inDir, newFileName))
			fmt.Println("----------------------------------------------------------------------------------------")
			fmt.Println()

			e := os.Rename(filepath.Join(inDir, oldFileName.Name()), filepath.Join(inDir, newFileName))
			if e != nil {
				log.Fatal(e)
			}
			i += 1
		}
	}

	zenity.Info("Renaming complete.",
		zenity.Title("Information"),
		zenity.InfoIcon)

}

func contains(list []string, target string) bool {
	for _, str := range list {
		if str == target {
			return true
		}
	}
	return false
}

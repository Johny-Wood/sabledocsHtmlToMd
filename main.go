package main

import (
  "fmt"
  "log"
  "os"
  // "path"
  "path/filepath"
  "regexp"
  "strings"

  "github.com/BurntSushi/toml"
  md "github.com/JohannesKaufmann/html-to-markdown"
  "github.com/JohannesKaufmann/html-to-markdown/plugin"
  "github.com/PuerkitoBio/goquery"
  "golang.org/x/text/cases"
  "golang.org/x/text/language"
)

type Formatter struct {
  Settings Settings
  Translation Translation 
  inputPath string
  outputPath string
  configPath string
}

type Settings struct {
  ExcludeInputFiles []string
}

type Translation struct {
  TablesT map[string]string
  EntitiesT map[string]string
  ReqResT map[string]string
}


func resolveFilePath(filename string) (string, error) {
  // Get the directory of the current Go file
	goFileDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return "", err
	}

	// Define the filename you want to read
	// filename := "example.txt"

	// Construct the file path in the current directory
	filePath := filepath.Join(goFileDir, filename)

	// Check if the file exists in the current directory
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		// If file doesn't exist in current directory, try to use os.Executable()
		execPath, err := os.Executable()
		if err != nil {
			fmt.Println("Error getting executable path:", err)
			return "", err
		}

		execDir := filepath.Dir(execPath)
		filePath = filepath.Join(execDir, filename)

		// Check again if the file exists using the executable directory
		_, err = os.Stat(filePath)
		if os.IsNotExist(err) {
			fmt.Println("File does not exist in the current directory or in the directory of the executable.")
			return "", err
		}
	}
  
  return filePath, nil
 //  // Get the path of the running executable
	// execPath, err := os.Executable()
	// if err != nil {
	// 	fmt.Println("Error getting executable path:", err)
	// 	return "", err
	// }
	//
	// // Derive the directory of the executable
	// execDir := filepath.Dir(execPath)
	//
	// // Define the filename you want to read
	// // filename := "example.txt"
	//
	// // Construct the full file path relative to the executable's directory
	// filePath := filepath.Join(execDir, filename)
	//
	// // Check if the file exists
	// _, err = os.Stat(filePath)
	// if os.IsNotExist(err) {
	// 	fmt.Println("File does not exist in the directory of the executable.")
	// 	return "", err
	// }
	//
 //  return filePath, nil
 // //  // Get the current working directory
	// // workingDir, err := os.Getwd()
	// // if err != nil {
	// // 	fmt.Println("Error getting working directory:", err)
	// // 	return "", err
	// // }
	// //
	// // // Define the filename you want to read
	// // // filename := "example.txt"
	// //
	// // // Check if the file exists in the current directory
	// // filePath := filepath.Join(workingDir, filename)
	// // _, err = os.Stat(filePath)
	// //
	// // if os.IsNotExist(err) {
	// // 	fmt.Println("File does not exist in the current directory.")
	// // 	return "", err
	// // }
 // //  
	// // // Open and read the file
 // //  return filePath, nil 
	// // // file, err := os.Open(filePath)
	// // // if err != nil {
	// // // 	fmt.Println("Error opening file:", err)
	// // // 	return
	// // // }
	// // // defer file.Close()
	// //
	// // // Read file contents
	// // // ... (Your file reading logic here)
	// // // fmt.Println("Reading file:", filename)
 // //  // ex, err := os.Executable()
 // //  // if err != nil {
 // //  //   return "", err
 // //  // }
 // //  //
 // //  // dir := filepath.Dir(ex)
 // //  // exPath := filepath.Dir(ex)
 // //  // if strings.Contains(dir, "go-build") {
 // //  //   return filename, nil 
 // //  // } else {
 // //  //   filePath := path.Join(exPath, filename)
 // //  //   return filePath, nil
 // //  // }
}

func SplitAny(s string, seps string) []string {
  splitter := func(r rune) bool {
    return strings.ContainsRune(seps, r)
  }
  return strings.FieldsFunc(s, splitter)
}


func (f Formatter) RemoveEmptyLines(html string) string {
  re := regexp.MustCompile(`(?m)^\s*$[\r\n]*|(\r\n|\n)`)
  return re.ReplaceAllString(html, "")
}


func (f Formatter) RemoveCodeTagsAndFormatLinks(htmlContent string) string {
  doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
  if err != nil {
    log.Fatal(err)
  }

  doc.Find("td code").Each(func(i int, codeSelection *goquery.Selection) {
    aTag := codeSelection.Find("a")
    if aTag.Length() > 0 {
      href, exists := aTag.Attr("href")
      if exists {
        // Modify href if it follows the specified pattern
        parts := SplitAny(href, ".-_#")
        entityPattern := `\b\w+-\b`
        re := regexp.MustCompile(entityPattern)
        matches := re.FindAllString(href, -1)

        if len(parts) > 1 {
          serviceName := strings.ToLower(parts[len(parts)-1])
          servicePrefix := strings.TrimSuffix(matches[0], "-")

          // Translate heading
          if value, ok := f.Translation.EntitiesT[cases.Title(language.English, cases.NoLower).String(servicePrefix)]; ok {
            servicePrefix = strings.ToLower(value)
          }

          aTag.SetAttr("href", fmt.Sprintf("#%s-%s", servicePrefix, serviceName))
        }
      }
    }

    // Replace <code> tag with its contents
    codeSelection.ReplaceWithSelection(codeSelection.Contents())
  })

  doc.Find("div div code").Each(func(i int, codeSelection *goquery.Selection) {
    aTag := codeSelection.Find("a")
    if aTag.Length() > 0 {
      href, exists := aTag.Attr("href")
      if exists {
        // Modify href if it follows the specified pattern
        parts := SplitAny(href, ".-#")

        entityPattern := `\b\w+-\b`
        re := regexp.MustCompile(entityPattern)
        matches := re.FindAllString(href, -1)

        if len(parts) > 1 {
          serviceName := strings.ToLower(parts[len(parts)-1])
          servicePrefix := strings.TrimSuffix(matches[0], "-")
          // Translate heading
          if value, ok := f.Translation.EntitiesT[cases.Title(language.English, cases.NoLower).String(servicePrefix)]; ok {
            servicePrefix = strings.ToLower(value)
          }


          if len(matches) > 0 {
            aTag.SetAttr("href", fmt.Sprintf("#%s-%s", servicePrefix, serviceName))
          }
        }
      }
    }

    // Replace <code> tag with its contents
    codeSelection.ReplaceWithSelection(codeSelection.Contents())
  })

  modifiedHTML, err := doc.Html()
  if err != nil {
    log.Fatal(err)
  }

  return modifiedHTML
}

func (f Formatter) FormatHeading(htmlContent string) string {
  doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
  if err != nil {
    log.Fatal(err)
  }

  doc.Find("h4 span").Each(func(i int, codeSelection *goquery.Selection) {
    aTag := codeSelection.Find("a")
    if aTag.Length() > 0 {
      serviceStr := strings.Replace(aTag.Text(), aTag.Children().Text(), "", -1)
      serviceSlice := strings.Split(serviceStr, " ") 
      if len(serviceSlice) > 1 {
        servicePrefix := serviceSlice[0]
        serviceName := serviceSlice[1]

        // Translate heading
        if value, ok := f.Translation.EntitiesT[cases.Title(language.English, cases.NoLower).String(servicePrefix)]; ok {
          servicePrefix = value
        }


        // Replace <code> tag with its contents
        newHeadingStr := cases.Title(language.English, cases.NoLower).String(servicePrefix) + " " + serviceName
        newHeadingHtml := fmt.Sprintf("<h4>%s</h4>", newHeadingStr)
        codeSelection.ReplaceWithHtml(newHeadingHtml)
      }
    }
  })

  // h4 a
  doc.Find("h4").Each(func(i int, codeSelection *goquery.Selection) {
    aTag := codeSelection.Find("a")
    if aTag.Length() > 0 {
      serviceStr := aTag.First().Text()
      serviceSlice := strings.Split(serviceStr, " ") 
      if len(serviceSlice) > 1 {
        servicePrefix := serviceSlice[0]
        serviceName := serviceSlice[1]

        // Translate heading
        if value, ok := f.Translation.EntitiesT[cases.Title(language.English, cases.NoLower).String(servicePrefix)]; ok {
          servicePrefix = value
        }

        newHeadingStr := cases.Title(language.English, cases.NoLower).String(servicePrefix) + " " + serviceName
        newHeadingHtml := fmt.Sprintf("<h4>%s</h4>", newHeadingStr)
        codeSelection.ReplaceWithHtml(newHeadingHtml)
      }
    }
  })



  modifiedHTML, err := doc.Html()
  if err != nil {
    log.Fatal(err)
  }

  return modifiedHTML
}

func (f Formatter) TransformAnchors(htmlContent []byte) string {
  doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(htmlContent)))
  if err != nil {
    log.Fatal(err)
  }

  doc.Find("p").Each(func(i int, codeSelection *goquery.Selection) {
    aTag := codeSelection.Find("a")
    if aTag.Length() > 0 {
      href, exists := aTag.Attr("href")
      if exists {
        parts := SplitAny(href, ".-#")
        entityPattern := `\b\w+-\b`
        re := regexp.MustCompile(entityPattern)
        matches := re.FindAllString(href, -1)

        if len(parts) > 1 {
          if len(matches) > 0 {
            serviceName := strings.ToLower(parts[len(parts)-1])
            servicePrefix := strings.TrimSuffix(matches[0], "-")

            // Translate heading
            if value, ok := f.Translation.EntitiesT[cases.Title(language.English, cases.NoLower).String(servicePrefix)]; ok {
              servicePrefix = strings.ToLower(value)
            }

            aTag.SetAttr("href", fmt.Sprintf("#%s-%s", servicePrefix, serviceName))
          }
        }
      }
    }
  })

  doc.Find("p a span").Each(func(i int, codeSelection *goquery.Selection) {
    aTag := codeSelection
    if aTag.Length() > 0 {
      servicePrefix := strings.Replace(aTag.Text(), aTag.Children().Text(), "", -1)

      // Translate heading
      if value, ok := f.Translation.EntitiesT[cases.Title(language.English, cases.NoLower).String(servicePrefix)]; ok {
        servicePrefix = value
      }


      // Replace <code> tag with its contents
      codeSelection.ReplaceWithHtml(cases.Title(language.English, cases.NoLower).String(servicePrefix))
    }
  })

  modifiedHTML, err := doc.Html()
  if err != nil {
    log.Fatal(err)
  }
  return modifiedHTML 
}

func (f Formatter) TranslateReqRes(htmlContent string) string {
  if len(f.Translation.ReqResT) > 0 {
    req := "Request:"
    res := "Response:"

    htmlReqTransalte := strings.ReplaceAll(htmlContent, req, f.Translation.ReqResT["Request"] + ":")
    htmlResTransalte := strings.ReplaceAll(htmlReqTransalte, res, f.Translation.ReqResT["Response"] + ":")

    return htmlResTransalte
  } else {
    return htmlContent
  }
}

func (f Formatter) TranslateTableHead(htmlContent string) string {
  if len(f.Translation.TablesT) > 0 {
    doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
    if err != nil {
      log.Fatal(err)
    }

    doc.Find("thead tr th").Each(func(i int, codeSelection *goquery.Selection) {
      if codeSelection.Length() > 0 {
        title := codeSelection.Text()


        // Translate heading
        if value, ok := f.Translation.TablesT[title]; ok {
          title = value
        }

        codeSelection.SetText(title)
      }
    })

    modifiedHTML, err := doc.Html()
    if err != nil {
      log.Fatal(err)
    }
    return modifiedHTML 
  } else {
    return htmlContent
  }
}


func isExcludedInputFile(filePath string, exclusions []string) bool {
  fileName := filepath.Base(filePath)
  for _, exclusion := range exclusions {
    if exclusion == fileName {
      return true
    }
  }
  return false
}

func main() {
  // Initialize settings
  formatter := &Formatter{
    inputPath: "file.html",
    outputPath: "output.md",
    configPath: "config.toml",
  }

  // Get config file path
  configFilePath, err := resolveFilePath(formatter.configPath)
  if err != nil {
    panic(err)
  }


  // Try to read and decode CONFIG file
  // cPath, err := filepath.Glob("./*.toml")
  if _, err := os.Stat(configFilePath); !os.IsNotExist(err) {

  // if _, err := os.Stat(cPath[0]); !os.IsNotExist(err) {
    fileContents, err := os.ReadFile(configFilePath)
    // fileContents, err := os.ReadFile(cPath[0])
    if err != nil {
      fmt.Printf("Error while reading config.toml file. Error - %s. Continue working with default setting.", err)
    } else {

      // _, err := toml.Decode(string(fileContents), &formatter.Translation)
      _, err := toml.Decode(string(fileContents), &formatter)
      if err != nil {
        fmt.Printf("Error while decoding config.toml file. Error - %s. Continue working with default setting.", err)
      } 
    }
  }


  // Read all HTML files in the folder
  htmlFiles, err := filepath.Glob("*.html")
  if err != nil {
    panic(err)
  }

  for _, htmlFile := range htmlFiles {
    // Check exclusion criteria from config file
    fileName := strings.Split(htmlFile, ".")[0]

    shouldExclude := isExcludedInputFile(htmlFile, formatter.Settings.ExcludeInputFiles)
    if shouldExclude {
      fmt.Printf("Skipping file %s based on exclusion criteria\n", htmlFile)
      continue
    }

    // Read HTML content
    htmlContent, err := os.ReadFile(htmlFile)
    if err != nil {
      panic(err)
    }


    // Apply HTML transformations
    modifiedHTML := formatter.RemoveCodeTagsAndFormatLinks(string(htmlContent))
    modifiedHTML = formatter.RemoveEmptyLines(modifiedHTML)
    modifiedHTML = formatter.TransformAnchors([]byte(modifiedHTML))
    modifiedHTML = formatter.FormatHeading(modifiedHTML)
    modifiedHTML = formatter.TranslateReqRes(modifiedHTML)
    modifiedHTML = formatter.TranslateTableHead(modifiedHTML)

    // Convert HTML to Markdown
    opt := &md.Options{
      EscapeMode: "disabled",
    }
    converter := md.NewConverter("", true, opt)
    converter.Use(plugin.GitHubFlavored())
    mdContent, err := converter.ConvertString(modifiedHTML)
    if err != nil {
      log.Fatal(err)
    }

    // Write Markdown content to a file with the same name
    err = os.WriteFile(fmt.Sprintf("%s.md", fileName), []byte(mdContent), 0644)
    if err != nil {
      panic(err)
    }
  }
}

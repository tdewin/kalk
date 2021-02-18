package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"gopkg.in/yaml.v2"
)

type General struct {
	UniqueID    string `yaml:"id"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}
type Input struct {
	Name        string `yaml:"name"`
	Default     string `yaml:"default"`
	HumanName   string `yaml:"humanName"`
	Type        string `yaml:"type"`
	Description string `yaml:"description"`
}
type Output struct {
	Name        string `yaml:"name"`
	HumanName   string `yaml:"humanName"`
	Type        string `yaml:"type"`
	Description string `yaml:"description"`
	Formula     string `yaml:"formula"`
}
type Calculator struct {
	General General  `yaml:"general"`
	Inputs  []Input  `yaml:"input"`
	Outputs []Output `yaml:"output"`
}

func isLetterOrNumber(test byte) bool {
	if (test >= 'a' && test <= 'z') || (test >= 'A' && test <= 'Z') || (test >= '0' && test <= '9') {
		return true
	}
	return false
}

func main() {
	fmt.Println("Parsing YAML file")

	var directory string
	var fileNameIn string
	var htmlNameIn string
	flag.StringVar(&directory, "d", "default", "Directory for input & output")
	flag.StringVar(&fileNameIn, "i", "index.yaml", "YAML file to parse.")
	flag.StringVar(&htmlNameIn, "o", "index.html", "Output file")
	flag.Parse()

	fileName := filepath.Join(directory, fileNameIn)
	htmlName := filepath.Join(directory, htmlNameIn)

	if fileName == "" {
		fmt.Println("Please provide yaml file by using -i option")
		return
	}

	if htmlName == "" {
		fmt.Println("Please provide index file by using -o option")
		return
	}

	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Error reading YAML file: %s\n", err)
		return
	}

	htmlFile, err := os.Create(htmlName)
	if err != nil {
		fmt.Printf("Error opening stream out file: %s\n", err)
		return
	}
	defer htmlFile.Close()
	htmlFile.WriteString(`
	<!doctype html>
	<html lang="en">
	  <head>
		<!-- Required meta tags -->
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1">
	
		<!-- Bootstrap CSS -->
		<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.0.0-beta1/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-giJF6kkoqNQ00vy+HMDP7azOuL0xtbfIcaT9wjKHr8RbDVddVHyTfAAsrekwKmP1" crossorigin="anonymous">
	
		<title>Calculation Page</title>
		<script>
			function roundFloat(floatnum,precision) {
				pexp = Math.pow(10,precision)
				return Math.round(floatnum*pexp)/pexp
			}
			round1 = (floatnum) => roundFloat(floatnum,1)
			round2 = (floatnum) => roundFloat(floatnum,2)
			round3 = (floatnum) => roundFloat(floatnum,3)
			round4 = (floatnum) => roundFloat(floatnum,4)

			function selectfunc(el) {
				el.select();
			}
		</script>
	  </head>
	  <body>
	  <div class="container-fluid mainpage">
`)

	r := bytes.NewReader(yamlFile)
	dec := yaml.NewDecoder(r)

	calculator := Calculator{}
	calculators := []Calculator{}
	for dec.Decode(&calculator) == nil {
		calculators = append(calculators, calculator)
		//make sure it is clean
		calculator = Calculator{}
	}

	for _, calculator := range calculators {

		id := calculator.General.UniqueID
		if id == "" {
			id = strings.Replace(uuid.New().String(), "-", "", -1)
		}

		fmt.Fprintf(htmlFile, "    <div class='calculator mb-5'>\n")
		fmt.Fprintf(htmlFile, "    <h1>%s</h1>\n", calculator.General.Name)
		fmt.Fprintf(htmlFile, "    <hr>\n")

		fmt.Fprint(htmlFile, "        <script>\n")

		getters := make(map[string]string)

		for _, input := range calculator.Inputs {
			inputtypelower := strings.ToLower(input.Type)
			if inputtypelower != "divider" {
				localid := fmt.Sprintf("%s_%s", id, input.Name)
				itype := "parseFloat"

				if inputtypelower == "number" {
					itype = "parseInt"
				} else if len(inputtypelower) > 5 && inputtypelower[0:5] == "round" {
					itype = "parseFloat"
				} else if inputtypelower == "ceil" {
					itype = "parseFloat"
				} else if inputtypelower == "floor" {
					itype = "parseFloat"
				}
				getters[input.Name] = fmt.Sprintf("get%s()", localid)
				fmt.Fprintf(htmlFile, "\nfunction get%s() { return %s(eval(document.getElementById('%s').value)) }", localid, itype, localid)
			}
		}
		for _, output := range calculator.Outputs {
			outputtypelower := strings.ToLower(output.Type)
			if outputtypelower != "divider" {
				localid := fmt.Sprintf("%s_%s", id, output.Name)

				parsetype := "parseFloat"
				showtype := ""
				if strings.ToLower(output.Type) == "number" {
					parsetype = "parseInt"
					showtype = "Math.round"
				} else if len(outputtypelower) > 5 && outputtypelower[0:5] == "round" {
					parsetype = "parseFloat"
					showtype = fmt.Sprintf("round%c", outputtypelower[5])
				} else if outputtypelower == "ceil" {
					parsetype = "parseFloat"
					showtype = "Math.ceil"
				} else if outputtypelower == "floor" {
					parsetype = "parseFloat"
					showtype = "Math.floor"
				}

				getters[output.Name] = fmt.Sprintf("get%s()", localid)
				fmt.Fprintf(htmlFile, "\nfunction get%s() { return %s(eval(document.getElementById('%s').value)) }", localid, parsetype, localid)
				fmt.Fprintf(htmlFile, "\nfunction set%s(newVal) { document.getElementById('%s').value = %s(newVal) }", localid, localid, showtype)
			}
		}

		fmt.Fprintf(htmlFile, "\nfunction calculate%s() {", id)
		for _, output := range calculator.Outputs {
			outputtypelower := strings.ToLower(output.Type)
			if outputtypelower != "divider" {
				formula := []byte{}

				word := []byte{}

				for i := 0; i < len(output.Formula); i++ {
					b := output.Formula[i]
					if isLetterOrNumber(b) {
						word = append(word, b)
					} else {
						if len(word) > 0 {
							strcheck := string(word)
							if getter, ok := getters[strcheck]; ok {
								for _, b := range []byte(getter) {
									formula = append(formula, b)
								}
							} else {
								for _, b := range word {
									formula = append(formula, b)
								}
							}
							word = []byte{}
						}
						formula = append(formula, b)
					}
				}
				if len(word) > 0 {
					strcheck := string(word)
					if getter, ok := getters[strcheck]; ok {
						for _, b := range []byte(getter) {
							formula = append(formula, b)
						}
					} else {
						for _, b := range word {
							formula = append(formula, b)
						}
					}
					word = []byte{}
				}

				localid := fmt.Sprintf("%s_%s", id, output.Name)
				fmt.Fprintf(htmlFile, "\n set%s(eval('%s'))", localid, formula)
			}
		}
		fmt.Fprintf(htmlFile, "\n}\n")
		fmt.Fprint(htmlFile, "        </script>\n")
		fmt.Fprintf(htmlFile, "<form id='%s'>", id)
		for _, input := range calculator.Inputs {
			celltypelower := strings.ToLower(input.Type)
			if celltypelower == "divider" {
				fmt.Fprint(htmlFile, "<div class='form-group row mt-5'><div class='col-sm-12'></div></div>")
			} else {

				localid := fmt.Sprintf("%s_%s", id, input.Name)
				description := ""
				if input.Description != "" {
					description = fmt.Sprintf("<small id='emailHelp' class='form-text text-muted'>%s</small>", input.Description)
				}
				fmt.Fprintf(htmlFile, `
			<div class="form-group row mt-2">
			 <label for="%s" class="col-sm-2 col-form-label">%s</label>
			 <div class="col-sm-10">
 			 	<input type="text" class="form-control" id="%s" placeholder="%s" value="%s">
				%s
			  </div>
			</div>
			
 			`, localid, input.HumanName, localid, input.Default, input.Default, description)
			}
		}
		fmt.Fprintf(htmlFile, `
			<div class="form-group row mt-3 mb-3"><div class="col-sm-12"><button type="button" onclick="calculate%s();return false" class="btn btn-primary">Calculate</button></div></div>
		`, id)
		for _, output := range calculator.Outputs {
			celltypelower := strings.ToLower(output.Type)
			if celltypelower == "divider" {
				fmt.Fprint(htmlFile, "<div class='form-group row mt-5'><div class='col-sm-12'></div></div>")
			} else {

				localid := fmt.Sprintf("%s_%s", id, output.Name)
				description := ""
				if output.Description != "" {
					description = fmt.Sprintf("<small id='emailHelp' class='form-text text-muted'>%s</small>", output.Description)
				} else {
					description = fmt.Sprintf("<small id='emailHelp' class='form-text text-muted'>%s = %s</small>", output.Name, output.Formula)
				}

				fmt.Fprintf(htmlFile, `
			<div class="form-group row">
			 <label for="%s" class="col-sm-2 col-form-label">%s</label>
			 <div class="col-sm-10">
 			 	<input type="text" class="form-control" id="%s" placeholder="0" value="0" oncontextmenu="selectfunc(this);">
				  %s
			  </div>
			</div>
 			`, localid, output.HumanName, localid, description)
			}
		}
		fmt.Fprint(htmlFile, "</form>")
		fmt.Fprintf(htmlFile, "</div>\n")

	}

	htmlFile.WriteString(`
		</div>
		<!-- Option 1: Bootstrap Bundle with Popper -->
		<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.0.0-beta1/dist/js/bootstrap.bundle.min.js" integrity="sha384-ygbV9kiqUc6oa4msXn9868pTtWMgiQaeYH7/t7LECLbyPA2x65Kgf80OJFdroafW" crossorigin="anonymous"></script>
	  </body>
	</html>
	`)

	fmt.Printf("Done %s", htmlName)

}

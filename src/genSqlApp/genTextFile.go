// vi:nu:et:sts=4 ts=4 sw=4
// See License.txt in main repository directory

// Generate SQL Application programs for GO

package genSqlApp

import (
	"../shared"
	"../util"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func GenTextFile(mdl string, outPath string, data interface{}) error {
	var err	    error
	var tmpl	*template.Template

	log.Printf("\tGenTextFile mdl:%s fn:%s ...", mdl, outPath)

	outData := strings.Builder{}

	// Parse and execute the template.
	name := filepath.Base(mdl)
	tmpl, err = template.New(name).Delims("[[", "]]").Funcs(sharedData.Funcs()).ParseFiles(mdl)
	if err != nil {
		return err
	}
	if sharedData.Debug() {
		log.Println("\t\t\t input data to template:", data)
		log.Println("\t\texecuting template...")
	}
	err = tmpl.ExecuteTemplate(&outData, name, data)
	if err != nil {
		return err
	}

	// Save the generated file to the output file path.
	if !sharedData.Noop() {
		// Delete existing file.
		if outPath, err = util.IsPathRegularFile(outPath); err == nil {
			if sharedData.Replace() {
				if err = os.Remove(outPath); err != nil {
					return errors.New(fmt.Sprint("Error - could not delete:", outPath, err))
				}
			} else {
				return errors.New(fmt.Sprint("Error - overwrite error of:", outPath))
			}
		}
		// Write the file to disk
		err := ioutil.WriteFile(outPath, []byte(outData.String()), 0664)
		if err != nil {
			return errors.New(fmt.Sprint("Error:", outPath, err))
		}
	} else {
		log.Println("<<<<<<<<<<<<<<<<<<<<<<<<", outPath, ">>>>>>>>>>>>>>>>>>>>>>>>>")
		log.Println(outData.String())
		log.Println("<<<<<<<<<<<<<<<<<<<<<<<< End of", outPath, ">>>>>>>>>>>>>>>>>>>>>>>>>>")
	}

	return nil
}

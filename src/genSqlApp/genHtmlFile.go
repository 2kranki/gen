// vi:nu:et:sts=4 ts=4 sw=4
// See License.txt in main repository directory

// Generate SQL Application programs for GO

package genSqlApp

import (
	"../shared"
	"../util"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var htmlTmpls template.Template

func GenHtmlFile(mdl string, fn string, data interface{}) error {
	var outPath string
	var err error
	var tmpl *template.Template

	log.Printf("\tGenHtmlFile mdl:%s fn:%s ...", mdl, fn)

	outData := strings.Builder{}
	if sharedData.Debug() {
		log.Println("\t\texecuting template...")
		log.Println("\t\tdata:", nil)
	}

	name := filepath.Base(mdl)
	tmpl, err = template.New(name).Delims("[[", "]]").Funcs(sharedData.Funcs()).ParseFiles(mdl)
	if err != nil {
		return err
	}
	err = tmpl.ExecuteTemplate(&outData, name, data)
	if err != nil {
		return err
	}

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
		log.Println("<<<<<<<<<<<<<<<<<<<<<<<<", fn, ">>>>>>>>>>>>>>>>>>>>>>>>>")
		log.Println(outData.String())
		log.Println("<<<<<<<<<<<<<<<<<<<<<<<<", fn, ">>>>>>>>>>>>>>>>>>>>>>>>>>")
	}

	return err
}

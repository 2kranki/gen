// vi:nu:et:sts=4 ts=4 sw=4
// See License.txt in main repository directory

// Generate HTML Files

package genCmn

import (
	"fmt"
	"genapp/pkg/sharedData"
	"html/template"
	"io/ioutil"
	"log"
	"strings"

	"github.com/2kranki/go_util"
)

var htmlTmpls template.Template

func GenHtmlFile(mdl *util.Path, outPath *util.Path, data interface{}) error {
	var err error
	var tmpl *template.Template

	log.Printf("\tGenHtmlFile mdl:%s fn:%s ...", mdl.String(), outPath.String())

	outData := strings.Builder{}
	if sharedData.Debug() {
		log.Println("\t\texecuting template...")
		log.Println("\t\tdata:", nil)
	}

	name := mdl.Base()
	tmpl, err = template.New(name).Delims("[[", "]]").Funcs(sharedData.Funcs()).ParseFiles(mdl.String())
	if err != nil {
		return err
	}
	err = tmpl.ExecuteTemplate(&outData, name, data)
	if err != nil {
		return err
	}

	if !sharedData.Noop() {
		// Delete existing file.
		if outPath.IsPathRegularFile() {
			if !sharedData.Replace() {
				return fmt.Errorf("Error - overwrite error of %s\n", outPath)
			}
		}
		// Write the file to disk replacing an existing file.
		err := ioutil.WriteFile(outPath.String(), []byte(outData.String()), 0664)
		if err != nil {
			return fmt.Errorf("Error: I/O error for %s: %s\n", outPath.String(), err.Error())
		}
	} else {
		log.Println("<<<<<<<<<<<<<<<<<<<<<<<<", outPath.String(), ">>>>>>>>>>>>>>>>>>>>>>>>>")
		log.Println(outData.String())
		log.Println("<<<<<<<<<<<<<<<<<<<<<<<<", outPath.String(), ">>>>>>>>>>>>>>>>>>>>>>>>>>")
	}

	return err
}

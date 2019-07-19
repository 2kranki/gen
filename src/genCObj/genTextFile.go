// vi:nu:et:sts=4 ts=4 sw=4
// See License.txt in main repository directory

// Generate C Object - text file

package genCObj

import (
	"../shared"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"text/template"
)

func GenTextFile(mdl *util.Path, outPath *util.Path, data interface{}) error {
	var err	    error
	var tmpl	*template.Template

	log.Printf("\tGenTextFile mdl:%s fn:%s ...", mdl.String(), outPath.String())

	outData := strings.Builder{}

	// Parse and execute the template.
	name := mdl.Base()
	tmpl, err = template.New(name).Delims("[[", "]]").Funcs(sharedData.Funcs()).ParseFiles(mdl.String())
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
		// Write the file to disk
		err := ioutil.WriteFile(outPath.String(), []byte(outData.String()), 0664)
		if err != nil {
			return errors.New(fmt.Sprint("Error:", outPath.String(), err.Error()))
		}
	} else {
		log.Println("<<<<<<<<<<<<<<<<<<<<<<<<", outPath, ">>>>>>>>>>>>>>>>>>>>>>>>>")
		log.Println(outData.String())
		log.Println("<<<<<<<<<<<<<<<<<<<<<<<< End of", outPath, ">>>>>>>>>>>>>>>>>>>>>>>>>>")
	}

	return nil
}

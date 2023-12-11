/*
Copyright 2021 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"log"
	"os"

	"github.com/spf13/pflag"

	"vitess.io/vitess/go/tools/asthelpergen"

	"vitess.io/vitess/go/tools/goimports"
)

func main() {
	var options asthelpergen.Options
	var verify bool

	pflag.StringSliceVar(&options.Packages, "in", nil, "Go packages to load the generator")
	pflag.StringVar(&options.RootInterface, "iface", "", "Root interface generate rewriter for")
	pflag.StringSliceVar(&options.Clone.Exclude, "clone_exclude", nil, "don't deep clone these types")
	pflag.StringSliceVar(&options.Equals.AllowCustom, "equals_custom", nil, "generate custom comparators for these types")
	pflag.BoolVar(&verify, "verify", false, "ensure that the generated files are correct")
	pflag.Parse()

	result, err := asthelpergen.GenerateASTHelpers(&options)
	if err != nil {
		log.Fatal(err)
	}

	if verify {
		for _, err := range asthelpergen.VerifyFilesOnDisk(result) {
			log.Fatal(err)
		}
		log.Printf("%d files OK", len(result))
	} else {
		for fullPath, file := range result {
			content, err := goimports.FormatJenFile(file)
			if err != nil {
				log.Fatalf("failed to apply goimport to '%s': %v", fullPath, err)
			}
			err = os.WriteFile(fullPath, content, 0664)
			if err != nil {
				log.Fatalf("failed to save file to '%s': %v", fullPath, err)
			}
			log.Printf("saved '%s'", fullPath)
		}
	}
}

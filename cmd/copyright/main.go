/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/blend/go-sdk/ansi"
	"github.com/blend/go-sdk/copyright"
)

type flagStrings []string

func (fs flagStrings) String() string {
	return strings.Join(fs, ", ")
}

func (fs *flagStrings) Set(flagValue string) error {
	if flagValue == "" {
		return fmt.Errorf("invalid flag value; is empty")
	}
	*fs = append(*fs, flagValue)
	return nil
}

var (
	flagNotice  string
	flagCompany string
	flagYear    int

	flagRestrictions           string
	flagRestrictionsOpenSource bool
	flagRestrictionsInternal   bool

	flagInject bool
	flagRemove bool

	flagIncludeFiles = flagStrings{}
	flagExcludeFiles = flagStrings{}
	flagIncludeDirs  = flagStrings{}
	flagExcludeDirs  = flagStrings{}

	flagExitFirst bool
	flagVerbose   bool
	flagDebug     bool
)

func init() {
	flag.BoolVar(&flagVerbose, "verbose", false, "If verbose output should be shown")
	flag.BoolVar(&flagDebug, "debug", false, "If verbose output should be shown")

	flag.BoolVar(&flagExitFirst, "exit-first", false, "If the program should exit on the first verification error")

	flag.StringVar(&flagNotice, "notice", copyright.DefaultNoticeBodyTemplate, "The notice body template; use '-' to read from standard input")
	flag.StringVar(&flagCompany, "company", "", "The company name to use in the notice body template")
	flag.IntVar(&flagYear, "year", time.Now().UTC().Year(), "The year to use in the notice body template")

	flag.StringVar(&flagRestrictions, "restrictions", "", "The restrictions to use in the notice body template.")
	flag.BoolVar(&flagRestrictionsOpenSource, "restrictions-open-source", false, "The restrictions should be the open source default")
	flag.BoolVar(&flagRestrictionsInternal, "restrictions-internal", false, "The restrictions should be the internal default")

	flag.BoolVar(&flagInject, "inject", false, "If we should inject the notice")
	flag.BoolVar(&flagRemove, "remove", false, "If we should remove the notice")

	flag.Var(&flagIncludeFiles, "include-file", "Files to include via glob match")
	flag.Var(&flagExcludeFiles, "exclude-file", "Files to exclude via glob match")
	flag.Var(&flagIncludeDirs, "include-dir", "Directories to include via glob match")
	flag.Var(&flagExcludeDirs, "exclude-dir", "Directories to exclude via glob match")
	flag.Parse()
}

func main() {
	ctx := context.Background()

	if flagNotice == "" {
		fmt.Fprintln(os.Stderr, "--notice was provided an empty string; cannot continue")
		os.Exit(1)
	}

	if strings.TrimSpace(flagNotice) == "-" {
		notice, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%+v\n", err)
			os.Exit(1)
		}
		flagNotice = string(notice)
	}

	var roots []string
	if args := flag.Args(); len(args) > 0 {
		roots = args
	} else {
		roots = []string{"."}
	}

	var restrictions string
	if flagRestrictions != "" {
		restrictions = flagRestrictions
	} else if flagRestrictionsOpenSource {
		restrictions = copyright.DefaultRestrictionsOpenSource
	} else if flagRestrictionsInternal {
		restrictions = copyright.DefaultRestrictionsInternal
	}

	engine := copyright.Copyright{
		Config: copyright.Config{
			NoticeBodyTemplate: flagNotice,
			Company:            flagCompany,
			Restrictions:       restrictions,
			Year:               flagYear,
			IncludeFiles:       flagStringsWithDefault(flagIncludeFiles, copyright.DefaultIncludeFiles),
			ExcludeFiles:       flagStringsWithDefault(flagExcludeFiles, copyright.DefaultExcludeFiles),
			IncludeDirs:        flagStringsWithDefault(flagIncludeDirs, copyright.DefaultIncludeDirs),
			ExcludeDirs:        flagStringsWithDefault(flagExcludeDirs, copyright.DefaultExcludeDirs),
			ExitFirst:          &flagExitFirst,
			Verbose:            &flagVerbose,
			Debug:              &flagDebug,
		},
	}

	var action func(context.Context) error
	var actionLabel string

	if flagInject {
		action = engine.Inject
		actionLabel = "inject"
	} else if flagRemove {
		action = engine.Remove
		actionLabel = "remove"
	} else {
		action = engine.Verify
		actionLabel = "verify"
	}

	var didFail bool
	for _, root := range roots {
		engine.Root = root
		maybeFail(ctx, action, &didFail)
	}
	if didFail {
		fmt.Printf("copyright %s %s!\n", actionLabel, ansi.Red("failed"))
		os.Exit(1)
	}
	fmt.Printf("copyright %s %s!\n", actionLabel, ansi.Green("ok"))
}

func flagStringsWithDefault(flagPointer flagStrings, defaultValues []string) []string {
	if flagPointer != nil && len(flagPointer) > 0 {
		return flagPointer
	}
	return defaultValues
}

func maybeFail(ctx context.Context, action func(context.Context) error, didFail *bool) {
	err := action(ctx)
	if err != nil {
		if err == copyright.ErrFailure {
			*didFail = true
			return
		}
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}

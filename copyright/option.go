/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package copyright

// Option is a function that modifies a config.
type Option func(*Copyright)

// OptVerbose sets if we should show verbose output.
func OptVerbose(verbose bool) Option {
	return func(p *Copyright) {
		p.Config.Verbose = &verbose
	}
}

// OptDebug sets if we should show debug output.
func OptDebug(debug bool) Option {
	return func(p *Copyright) {
		p.Config.Debug = &debug
	}
}

// OptExitFirst sets if we should stop after the first failure.
func OptExitFirst(exitFirst bool) Option {
	return func(p *Copyright) {
		p.Config.ExitFirst = &exitFirst
	}
}

// OptRoot sets the root directory to start the profanity check.
func OptRoot(root string) Option {
	return func(p *Copyright) {
		p.Config.Root = root
	}
}

// OptIncludeFiles sets the include glob filter for files.
func OptIncludeFiles(includeGlobs ...string) Option {
	return func(p *Copyright) {
		p.Config.IncludeFiles = includeGlobs
	}
}

// OptExcludeFiles sets the exclude glob filter for files.
func OptExcludeFiles(excludeGlobs ...string) Option {
	return func(p *Copyright) {
		p.Config.ExcludeFiles = excludeGlobs
	}
}

// OptIncludeDirs sets the include glob filter for files.
func OptIncludeDirs(includeGlobs ...string) Option {
	return func(p *Copyright) {
		p.Config.IncludeDirs = includeGlobs
	}
}

// OptExcludeDirs sets the exclude glob filter for directories.
func OptExcludeDirs(excludeGlobs ...string) Option {
	return func(p *Copyright) {
		p.Config.ExcludeDirs = excludeGlobs
	}
}

// OptNoticeBodyTemplate sets the notice body template.
func OptNoticeBodyTemplate(noticeBodyTemplate string) Option {
	return func(p *Copyright) {
		p.Config.NoticeBodyTemplate = noticeBodyTemplate
	}
}

// OptCompany sets the notice body template company.
func OptCompany(company string) Option {
	return func(p *Copyright) {
		p.Config.Company = company
	}
}

// OptRestrictions sets the notice body template restrictions.
func OptRestrictions(restrictions string) Option {
	return func(p *Copyright) {
		p.Config.Restrictions = restrictions
	}
}

// OptConfig sets the config in its entirety.
func OptConfig(cfg Config) Option {
	return func(p *Copyright) {
		p.Config = cfg
	}
}

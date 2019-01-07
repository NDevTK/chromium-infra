// Copyright 2018 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package crauditcommits

import (
	"fmt"
	"strings"

	"golang.org/x/net/context"
)

const (
	dirLayoutTests = "third_party/blink/web_tests"
	dirSKCMS       = "third_party/skcms"
	dirSkiaAPIDocs = "site/user/api"

	fileAFDO            = "chrome/android/profiles/newest.txt"
	fileDEPS            = "DEPS"
	fileFuchsiaSDKLinux = "build/fuchsia/linux.sdk.sha1"
	fileFuchsiaSDKMac   = "build/fuchsia/mac.sdk.sha1"
	fileSkiaManifest    = "manifest/skia"
	fileSkiaTasks       = "infra/bots/tasks.json"
)

// AutoRollRulesForFilesAndDirs returns an AccountRules instance for an account
// which should only modify the given set of files and directories.
func AutoRollRulesForFilesAndDirs(account string, files, dirs []string) AccountRules {
	return AccountRules{
		Account: account,
		Funcs: []RuleFunc{
			func(ctx context.Context, ap *AuditParams, rc *RelevantCommit, cs *Clients) *RuleResult {
				ruleName := fmt.Sprintf("OnlyModifies_%s", strings.Join(append(files, dirs...), "+"))
				return OnlyModifiesFilesAndDirsRule(ctx, ap, rc, cs, ruleName, files, dirs)
			},
		},
		notificationFunction: fileBugForAutoRollViolation,
	}
}

// AutoRollRulesForDirList returns an AccountRules instance for an account
// which should only modify the given set of directories.
func AutoRollRulesForDirList(account string, dirs []string) AccountRules {
	return AutoRollRulesForFilesAndDirs(account, []string{}, dirs)
}

// AutoRollRulesForFileList returns an AccountRules instance for an account
// which should only modify the given set of files.
func AutoRollRulesForFileList(account string, files []string) AccountRules {
	return AutoRollRulesForFilesAndDirs(account, files, []string{})
}

// AutoRollRulesDEPS returns an AccountRules instance for an account which should
// only modify the ``DEPS`` file.
func AutoRollRulesDEPS(account string) AccountRules {
	return AutoRollRulesForFileList(account, []string{fileDEPS})
}

// AutoRollRulesDEPSAndTasks returns an AccountRules instance for an account
// which should only modify the ``DEPS`` and ``infra/bots/tasks.json`` files.
func AutoRollRulesDEPSAndTasks(account string) AccountRules {
	return AutoRollRulesForFileList(account, []string{fileDEPS, fileSkiaTasks})
}

// AutoRollRulesFuchsiaSDKVersion returns an AccountRules instance for an
// account which should only modifiy ``build/fuchsia/sdk.sha1``.
func AutoRollRulesFuchsiaSDKVersion(account string) AccountRules {
	return AutoRollRulesForFileList(account, []string{fileFuchsiaSDKLinux, fileFuchsiaSDKMac})
}

// AutoRollRulesSKCMS returns an AccountRules instance for an account which
// should only modify ``third_party/skcms``.
func AutoRollRulesSKCMS(account string) AccountRules {
	return AutoRollRulesForDirList(account, []string{dirSKCMS})
}

// AutoRollRulesLayoutTests returns an AccountRules instance for an account
// which should only modify ``third_party/blink/web_tests``.
func AutoRollRulesLayoutTests(account string) AccountRules {
	return AutoRollRulesForDirList(account, []string{dirLayoutTests})
}

// AutoRollRulesAPIDocs returns an AccountRules instance for an account which
// should only modify ``site/user/api``.
func AutoRollRulesAPIDocs(account string) AccountRules {
	return AutoRollRulesForDirList(account, []string{dirSkiaAPIDocs})
}

// AutoRollRulesSkiaAssets returns an AccountRules instance for an account which
// should only modify Skia assets.
func AutoRollRulesSkiaAssets(account string, assets []string) AccountRules {
	files := make([]string, 0, len(assets)+1)
	for _, asset := range assets {
		files = append(files, fmt.Sprintf("infra/bots/assets/%s/VERSION", asset))
	}
	files = append(files, fileSkiaTasks)
	return AutoRollRulesForFileList(account, files)
}

// AutoRollRulesSkiaManifest returns an AccountRules instance for an account
// which should only modify ``manifest/skia``.
func AutoRollRulesSkiaManifest(account string) AccountRules {
	return AutoRollRulesForFileList(account, []string{fileSkiaManifest})
}

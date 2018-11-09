// Copyright 2018 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package crauditcommits

import (
	"golang.org/x/net/context"
)

// AutoRollRulesAFDOVersion returns an AccountRules instance for an account
// which should only modify ``chrome/android/profiles/newest.txt``.
func AutoRollRulesAFDOVersion(account string) AccountRules {
	return AccountRules{
		Account: account,
		Funcs: []RuleFunc{
			OnlyModifiesAFDOVersion,
		},
		notificationFunction: fileBugForAutoRollViolation,
	}
}

// AutoRollRulesDEPS returns an AccountRules instance for an account which should
// only modify the ``DEPS`` file.
func AutoRollRulesDEPS(account string) AccountRules {
	return AccountRules{
		Account: account,
		Funcs: []RuleFunc{
			OnlyModifiesDEPSFile,
		},
		notificationFunction: fileBugForAutoRollViolation,
	}
}

// AutoRollRulesDEPSAndTasks returns an AccountRules instance for an account
// which should only modify the ``DEPS`` and ``infra/bots/tasks.json`` files.
func AutoRollRulesDEPSAndTasks(account string) AccountRules {
	return AccountRules{
		Account: account,
		Funcs: []RuleFunc{
			OnlyModifiesDEPSAndTasks,
		},
	}
}

// AutoRollRulesFuchsiaSDKVersion returns an AccountRules instance for an
// account which should only modifiy ``build/fuchsia/sdk.sha1``.
func AutoRollRulesFuchsiaSDKVersion(account string) AccountRules {
	return AccountRules{
		Account: account,
		Funcs: []RuleFunc{
			OnlyModifiesFuchsiaSDKVersions,
		},
		notificationFunction: fileBugForAutoRollViolation,
	}
}

// AutoRollRulesSKCMS returns an AccountRules instance for an account which
// should only modify ``third_party/skcms``.
func AutoRollRulesSKCMS(account string) AccountRules {
	return AccountRules{
		Account: account,
		Funcs: []RuleFunc{
			OnlyModifiesSKCMS,
		},
		notificationFunction: fileBugForAutoRollViolation,
	}
}

// AutoRollRulesLayoutTests returns an AccountRules instance for an account
// which should only modify ``third_party/WebKit/LayoutTests``.
func AutoRollRulesLayoutTests(account string) AccountRules {
	return AccountRules{
		Account: account,
		Funcs: []RuleFunc{
			OnlyModifiesLayoutTests,
		},
		notificationFunction: fileBugForAutoRollViolation,
	}
}

// OnlyModifiesDEPSFile is a RuleFunc that verifies that the only file
// modified by the audited CL is ``DEPS``.
func OnlyModifiesDEPSFile(ctx context.Context, ap *AuditParams, rc *RelevantCommit, cs *Clients) *RuleResult {
	return OnlyModifiesFileRule(ctx, ap, rc, cs, "OnlyModifiesDEPSFile", "DEPS")
}

// OnlyModifiesDEPSAndTasks is a RuleFunc that verifies that the only files
// modified by the audited CL are ``DEPS`` and ``infra/bots/tasks.json``.
func OnlyModifiesDEPSAndTasks(ctx context.Context, ap *AuditParams, rc *RelevantCommit, cs *Clients) *RuleResult {
	files := []string{
		"DEPS",
		"infra/bots/tasks.json",
	}
	return OnlyModifiesFilesRule(ctx, ap, rc, cs, "OnlyModifiesDEPS+tasks.json", files)
}

// OnlyModifiesAFDOVersion is a RuleFunc which verifies that the only file
// modified by the audited CL is ``chrome/android/profiles/newest.txt``.
func OnlyModifiesAFDOVersion(ctx context.Context, ap *AuditParams, rc *RelevantCommit, cs *Clients) *RuleResult {
	return OnlyModifiesFileRule(ctx, ap, rc, cs, "OnlyModifiesAFDOVersion", "chrome/android/profiles/newest.txt")
}

// OnlyModifiesFuchsiaSDKVersions is a RuleFunc which verifies that the only
// files modified by the audited CL is ``build/fuchsia/linux.sdk.sha1`` and
// ``build/fuchsia/mac.sdk.sha1``.
func OnlyModifiesFuchsiaSDKVersions(ctx context.Context, ap *AuditParams, rc *RelevantCommit, cs *Clients) *RuleResult {
	files := []string{
		"build/fuchsia/linux.sdk.sha1",
		"build/fuchsia/mac.sdk.sha1",
	}
	return OnlyModifiesFilesRule(ctx, ap, rc, cs, "OnlyModifiesFuchsiaSDKVersions", files)
}

// OnlyModifiesSKCMS is a RuleFunc which verifies that the audited CL only
// modifies files in the ``third_party/skcms`` directory.
func OnlyModifiesSKCMS(ctx context.Context, ap *AuditParams, rc *RelevantCommit, cs *Clients) *RuleResult {
	return OnlyModifiesDirRule(ctx, ap, rc, cs, "OnlyModifiesSKCMS", "third_party/skcms")
}

// OnlyModifiesLayoutTests is a RuleFunc which verifies that the audited CL
// only modifies files in the ``third_party/WebKit/LayoutTests`` directory.
func OnlyModifiesLayoutTests(ctx context.Context, ap *AuditParams, rc *RelevantCommit, cs *Clients) *RuleResult {
	return OnlyModifiesDirRule(ctx, ap, rc, cs, "OnlyModifiesLayoutTests", "third_party/WebKit/LayoutTests")
}

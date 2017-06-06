// Copyright 2015 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package markdown

import "sort"

type registeredCoreRule struct {
	id   int
	rule CoreRule
}

var registeredCoreRules []registeredCoreRule

type coreRulesById []registeredCoreRule

func (r coreRulesById) Len() int           { return len(r) }
func (r coreRulesById) Less(i, j int) bool { return r[i].id < r[j].id }
func (r coreRulesById) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }

type registeredBlockRule struct {
	id         int
	rule       BlockRule
	terminates []int
}

var registeredBlockRules []registeredBlockRule

type blockRulesById []registeredBlockRule

func (r blockRulesById) Len() int           { return len(r) }
func (r blockRulesById) Less(i, j int) bool { return r[i].id < r[j].id }
func (r blockRulesById) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }

type registeredInlineRule struct {
	id   int
	rule InlineRule
}

var registeredInlineRules []registeredInlineRule

type inlineRulesById []registeredInlineRule

func (r inlineRulesById) Len() int           { return len(r) }
func (r inlineRulesById) Less(i, j int) bool { return r[i].id < r[j].id }
func (r inlineRulesById) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }

func indexInt(a []int, n int) int {
	for i, m := range a {
		if m == n {
			return i
		}
	}
	return -1
}

func RegisterCoreRule(id int, rule CoreRule) {
	registeredCoreRules = append(registeredCoreRules, registeredCoreRule{
		id:   id,
		rule: rule,
	})
	sort.Sort(coreRulesById(registeredCoreRules))

	coreRules = coreRules[:0]
	for _, r := range registeredCoreRules {
		coreRules = append(coreRules, r.rule)
	}
}

func RegisterBlockRule(id int, rule BlockRule, terminates []int) {
	registeredBlockRules = append(registeredBlockRules, registeredBlockRule{
		id:         id,
		rule:       rule,
		terminates: terminates,
	})
	sort.Sort(blockRulesById(registeredBlockRules))

	blockRules = blockRules[:0]
	blockquoteTerminatedBy = blockquoteTerminatedBy[:0]
	listTerminatedBy = listTerminatedBy[:0]
	referenceTerminatedBy = referenceTerminatedBy[:0]
	paragraphTerminatedBy = paragraphTerminatedBy[:0]
	for _, r := range registeredBlockRules {
		blockRules = append(blockRules, r.rule)
		if indexInt(r.terminates, 300) != -1 {
			blockquoteTerminatedBy = append(blockquoteTerminatedBy, r.rule)
		}
		if indexInt(r.terminates, 500) != -1 {
			listTerminatedBy = append(listTerminatedBy, r.rule)
		}
		if indexInt(r.terminates, 600) != -1 {
			referenceTerminatedBy = append(referenceTerminatedBy, r.rule)
		}
		if indexInt(r.terminates, 1100) != -1 {
			paragraphTerminatedBy = append(paragraphTerminatedBy, r.rule)
		}
	}
}

func RegisterInlineRule(id int, rule InlineRule) {
	registeredInlineRules = append(registeredInlineRules, registeredInlineRule{
		id:   id,
		rule: rule,
	})
	sort.Sort(inlineRulesById(registeredInlineRules))

	inlineRules = inlineRules[:0]
	for _, r := range registeredInlineRules {
		inlineRules = append(inlineRules, r.rule)
	}
}

func init() {
	RegisterCoreRule(100, ruleInline)
	RegisterCoreRule(200, ruleLinkify)
	RegisterCoreRule(300, ruleReplacements)
	RegisterCoreRule(400, ruleSmartQuotes)

	RegisterBlockRule(100, ruleCode, nil)
	RegisterBlockRule(200, ruleFence, []int{1100, 600, 300, 500, 1200})
	RegisterBlockRule(300, ruleBlockQuote, []int{1100, 600, 500, 1200})
	RegisterBlockRule(400, ruleHR, []int{1100, 600, 300, 500, 1200})
	RegisterBlockRule(500, ruleList, []int{1100, 600, 300, 1200})
	RegisterBlockRule(600, ruleReference, nil)
	RegisterBlockRule(700, ruleHeading, []int{1100, 600, 300, 1200})
	RegisterBlockRule(800, ruleLHeading, nil)
	RegisterBlockRule(900, ruleHTMLBlock, []int{1100, 600, 300, 1200})
	RegisterBlockRule(1000, ruleTable, []int{1100, 600, 1200})
	RegisterBlockRule(1100, ruleParagraph, nil)

	RegisterInlineRule(100, ruleText)
	RegisterInlineRule(200, ruleNewline)
	RegisterInlineRule(300, ruleEscape)
	RegisterInlineRule(400, ruleBackticks)
	RegisterInlineRule(500, ruleStrikeThrough)
	RegisterInlineRule(600, ruleEmphasis)
	RegisterInlineRule(700, ruleLink)
	RegisterInlineRule(800, ruleImage)
	RegisterInlineRule(900, ruleAutolink)
	RegisterInlineRule(1000, ruleHTMLInline)
	RegisterInlineRule(1100, ruleEntity)
	RegisterInlineRule(1200, ruleSupScript)
}

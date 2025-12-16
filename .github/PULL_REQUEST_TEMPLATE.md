## ✨ Summary

> What is the change? Why is it needed? Describe why this PR is necessary.

## 🏷️ Type

- [ ] 🐞 Bugfix (non-breaking change that fixes an issue)
- [ ] 💚 Improvement (non-breaking change that adds/modifies features to existing functionality)
- [ ] ⚡️ New feature (non-breaking change that adds functionality)
- [ ] ⚠️ Breaking change (change that is not backwards compatible and/or changes current functionality)

## 📦 Changes

> Describe the layers, classes, entities, etc. that the PR modifies.
> If possible, add details about how it was implemented, attaching screenshots or designs for better understanding.

---

## 🎯 Motivation and Context _(optional)_

> What is the reason? Related to which team? To solve which use case?

## 📷 Screenshots _(optional - required for UI changes)_

<table>
<tr>
<td></td>
<td>Before</td>
<td>After</td>
</tr>
<tr>
<td>Mobile</td>
<td>

_paste image_

</td>
<td>

_paste image_

</td>
</tr>
<tr>
<td>Desktop</td>
<td>

_paste image_

</td>
<td>

_paste image_

</td>
</tr>
</table>

## 📈 Test Cases _(optional)_

> Test links + users, CURLs, scopes and/or necessary data

## 🔗 Reference Links _(optional)_

> Documentation links, Figma, Jira tasks, etc.

## 🚨 Priority _(optional)_

- [ ] ⚪️ Very Low
- [ ] 🟢 Low
- [ ] 🟡 Medium
- [ ] 🟠 High
- [ ] 🔴 Very High

## Dependencies _(optional)_

#### Are there other apps whose release depends on this PR being deployed? (or vice versa)

---

## ✅ Before merging...

<!-- - [ ] PR has an assigned milestone (if targeting `develop`) -->
- [ ] PR is assigned to the appropriate person(s)
<!-- - [ ] PR was sent for review to the Slack channel `b2b-pull-request` -->
- [ ] PR was approved by the owner
- [ ] PR was tested and does not generate failures for the impacted use cases
- [ ] PR has the applicable labels

## Merge Rules

|              Pull Request              |   Merge Type   |
| :------------------------------------: | :------------: |
|    **feature/fix/etc** → `develop`     | _SQUASH MERGE_ |
| **feature/hotfix/fix/etc** → `release` | _SQUASH MERGE_ |
|         **hotfix** → `master`          | _SQUASH MERGE_ |
|         **release** → `master`         | _MERGE COMMIT_ |
|        **backport** → `develop`        | _MERGE COMMIT_ |

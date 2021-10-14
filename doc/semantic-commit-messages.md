---
type: "docs"
weight: 1
title: "语义化 Commit 消息"
---

# 语义化 Commit 消息规范

[原文](https://gist.github.com/joshbuchea/6f47e86d2510bce28f8e7f42ae84c716)

See how a minor change to your commit message style can make you a better programmer.

Format: `<type>(<scope>): <subject>`

`<scope>` is optional

## 例子

```
feat: add hat wobble
^--^  ^------------^
|     |
|     +-> Summary in present tense.
|
+-------> Type: chore, docs, feat, fix, refactor, style, or test.
```

More Examples:

- `feat`: (new feature for the user, not a new feature for build script) (新特性，但不是脚本的新特性)
- `fix`: (bug fix for the user, not a fix to a build script) (bug 修复，但不是脚本的修复)
- `docs`: (changes to the documentation) (文档改动)
- `style`: (formatting, missing semi colons, etc; no production code change) (代码风格修改，无代码逻辑变动)
- `refactor`: (refactoring production code, eg. renaming a variable) (重构代码，如变量重命名等等)
- `test`: (adding missing tests, refactoring tests; no production code change) (测试用例的增删改)
- `chore`: (updating grunt tasks etc; no production code change) (更新脚本、配置、依赖等等与代码逻辑无关的内容)

References:

- https://www.conventionalcommits.org/
- https://seesparkbox.com/foundry/semantic_commit_messages
- http://karma-runner.github.io/1.0/dev/git-commit-msg.html

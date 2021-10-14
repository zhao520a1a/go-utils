---
type: "docs"
weight: 7
title: "方案&代码审查"
---

# 新功能方案审查

有时候我们会需要修改一个项目某个模块的整体实现方案，方案本身的规模大不到开评审会，但也小不到可以无需审查，这时候就适用「新功能方案审查」。

## 步骤一：提 issue

在 Gitlab 项目下创建一个新的 issue，描述方案的**背景**和**解决方案**，如果有多备选方案也请直接注明。[示例](https://gitlab.pri.ibanyu.com/server/alertmanager/service/issues/137)如下：

![issue](./issue.png)

## 步骤二：邀请同事审查

找到项目 owner、mentor 或负责该项目的同事，邀请他 (们) 一起审查方案，过程可以用讨论的形式留存在 issue 下方的评论区中，并注明**最终方案**以及**这么选择的原因**。

## 步骤三：实现并关联 issue

使用 Gitlab 的 issue 关联功能将你的 Merge Requests 与 issue 进行关联，并在所有 Merge Requests 合并后关闭 issue。


# 代码审查规范 (试运行)

* 原则上，每个 MR 至少经过团队内一位工程师审查，但 MR 发起者有权直接合并。
* 每个 MR 有一个 LGTM 即表示审查通过。审查通过后，由 MR 发起者自行执行合并操作。
* MR 大小与请求他人审查的诚意成反比，对于过大的 MR，审查人员可拒绝审查。

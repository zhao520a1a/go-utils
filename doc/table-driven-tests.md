---
type: "docs"
weight: 3
title: "表驱动测试"
---

# 优先考虑表驱动测试

[原文](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)发表于 2019 年 5 月 7 日，作者为 DaveCheney。

我是测试的拥趸，特别是[单元测试](https://dave.cheney.net/2019/04/03/absolute-unit-test)和（[正确实践的](https://www.youtube.com/watch?v=EZ05e7EMOLM)）TDD。围绕
Go 项目发展起来的一个实践是表驱动测试。这篇文章探讨了编写表驱动测试的方式和原因。

假设我们有一个拆分字符串的函数：

```go
// Split 将 s 切分为由 sep 分隔的子字符串，并返回分隔符之间的子字符串的 slice
func Split(s, sep string) []string {
    var result []string
    i := strings.Index(s, sep)
    for i > -1 {
        result = append(result, s[:i])
        s = s[i+len(sep):]
        i = strings.Index(s, sep)
    }
    return append(result, s)
}
```

在 Go 中，单元测试只是（具备一定规则的）普通 Go 函数。所以我们在同目录下的文件中，用相同的包名 `strings`，为这个函数写单元测试。

```go
package strings

import (
	"reflect"
	"testing"
)

func TestSplit(t *testing.T) {
    got := Split("a/b/c", "/")
    want := []string{"a", "b", "c"}
    if !reflect.DeepEqual(want, got) {
         t.Fatalf("expected: %v, got: %v", want, got)
    }
}
```

测试只是有些规则的普通 Go 函数：

1. 测试函数的名字必须以 `Test` 开头。
2. 测试函数的参数必须是 `*testing.T` 类型。此类型由 `testing`
   包自身注入，用于提供打印（print）、跳过（skip）测试和使测试失败（fail）的方法。

在测试中，我们用一些输入调用 `Split`，然后将其与期望的结果进行比较。

## 代码覆盖率

下一个问题是，这个包的覆盖率是多少？幸运的是 go test 支持分支覆盖分析。我们可以这样调用：

```shell
% go test -coverprofile=c.out
PASS
coverage: 100.0% of statements
ok      split   0.010s
```

这告诉我们分支覆盖率为 100%，并不奇怪，毕竟代码中只有一个分支。

如果我们想钻研下覆盖率报告，go tool 有几个选项可以打印覆盖率报告。可以用`go tool cover -func`细化出各个函数的覆盖率：

```shell
% go tool cover -func=c.out
split/split.go:8:       Split          100.0%
total:                  (statements)   100.0%
```

不必为此感到激动，因为这个包只有一个函数，但我相信你会找更「刺激」的包来测试。

### 包装到 `.bashrc` 中

这俩命令对我很有用，我有一个 shell alias，可以用一条命令运行测试覆盖率和报告：

```shell
cover () {
    local t=$(mktemp -t cover)
    go test $COVERFLAGS -coverprofile=$t $@ \
        && go tool cover -func=$t \
        && unlink $t
}
```

## 超越 100% 的覆盖率

所以，我们写了一个测试用例，得到了 100%
的覆盖率，但这并不是故事的结局。我们有很好的分支覆盖率，但可能需要测试一些边界条件。比如，如果用逗号分隔会发生什么？

```go
func TestSplitWrongSep(t *testing.T) {
    got := Split("a/b/c", ",")
    want := []string{"a/b/c"}
    if !reflect.DeepEqual(want, got) {
        t.Fatalf("expected: %v, got: %v", want, got)
    }
}
```

或者，如果源字符串中没有分隔符呢？

```go
func TestSplitNoSep(t *testing.T) {
    got := Split("abc", "/")
    want := []string{"abc"}
    if !reflect.DeepEqual(want, got) {
        t.Fatalf("expected: %v, got: %v", want, got)
    }
}
```

我们开始构建一套测试边界条件的用例，这很好。

## 介绍表驱动测试

然而，在我们的测试中有很多重复。对于每个测试用例，只有测试用例的输入、期望输出和名字在变化，其他一切都是重复代码 (boilerplate
code)。我们只需要设置好所有输入和期望输出，就像关系型数据表中的行数据，然后将它们传给单个测试函数，由这个函数来遍历这张表，这便是我们今天要介绍的表驱动测试。

```go
func TestSplit(t *testing.T) {
    type test struct {
        input string
        sep   string
        want  []string
    }

    tests := []test{
        {input: "a/b/c", sep: "/", want: []string{"a", "b", "c"}},
        {input: "a/b/c", sep: ",", want: []string{"a/b/c"}},
        {input: "abc", sep: "/", want: []string{"abc"}},
    }

    for _, tc := range tests {
        got := Split(tc.input, tc.sep)
        if !reflect.DeepEqual(tc.want, got) {
            t.Fatalf("expected: %v, got: %v", tc.want, got)
        }
    }
}
```

我们声明了一个结构体来存放测试输入和期望输出，这就是表。`tests`
结构体通常是一个局部声明，因为我们想在这个包的其他测试中复用这个名字。

事实上，我们甚至不需要为类型命名，可以使用匿名结构体减少重复代码：

```go
func TestSplit(t *testing.T) {
    tests := []struct {
        input string
        sep   string
        want  []string
    }{
        {input: "a/b/c", sep: "/", want: []string{"a", "b", "c"}},
        {input: "a/b/c", sep: ",", want: []string{"a/b/c"}},
        {input: "abc", sep: "/", want: []string{"abc"}},
    }

    for _, tc := range tests {
        got := Split(tc.input, tc.sep)
        if !reflect.DeepEqual(tc.want, got) {
            t.Fatalf("expected: %v, got: %v", tc.want, got)
        }
    }
}
```

现在，增加新测试的方法很直观，向 `tests` 结构体增加一行即可。比如，如果我们的输入字符串以分隔符结尾会发生什么？

```go
{input: "a/b/c", sep: "/", want: []string{"a", "b", "c"}},
{input: "a/b/c", sep: ",", want: []string{"a/b/c"}},
{input: "abc", sep: "/", want: []string{"abc"}},
{input: "a/b/c/", sep: "/", want: []string{"a", "b", "c"}}, // trailing sep7&%
```

但我们运行 go test 时会得到

```shell
% go test
--- FAIL: TestSplit (0.00s)
    split_test.go:24: expected: [a b c], got: [a b c ]
```

抛开测试失败，还有几个问题需要谈论。

第一个是，将每个测试从函数重写为表中的一行后，我们无法知道失败的是哪个测试。虽然测试文件中有注释，但在 go test 输出中没法获取到它。

有几个方法可以解决这个问题。你会在 Go 代码库中看到各种风格，因为随着人们持续试验表测试这种形式，它的惯用法也在演化。

### 枚举测试用例

因为测试存在 slice 中，我们可以在失败消息中打印测试用例的序号：

```go
func TestSplit(t *testing.T) {
    tests := []struct {
        input string
        sep   string
        want  []string
    }{
        {input: "a/b/c", sep: "/", want: []string{"a", "b", "c"}},
        {input: "a/b/c", sep: ",", want: []string{"a/b/c"}},
        {input: "abc", sep: "/", want: []string{"abc"}},
        {input: "a/b/c/", sep: "/", want: []string{"a", "b", "c"}},
    }

    for i, tc := range tests {
        got := Split(tc.input, tc.sep)
        if !reflect.DeepEqual(tc.want, got) {
            t.Fatalf("test %d: expected: %v, got: %v", i+1, tc.want, got)
        }
    }
}
```

现在当我们运行 `go test` 会得到

```shell
% go test
--- FAIL: TestSplit (0.00s)
    split_test.go:24: test 4: expected: [a b c], got: [a b c ]
```

这样看起来好了一些，我们能够知道第 4 个测试失败了。但程序遍历的序号是从 0 开始，这需要在测试用例间保持一致；如果一些基于
0，另一些基于 1，会产生歧义。并且，如果测试用例的列表很长，只能手动数数来确定失败的测试用例是哪个 fixture 构成了 4 号测试用例。

### 给你的测试用例命名

另一个常见模式是在测试 fixture 中包含 `name` 字段。

```go
func TestSplit(t *testing.T) {
    tests := []struct {
        name  string
        input string
        sep   string
        want  []string
    }{
        {name: "simple", input: "a/b/c", sep: "/", want: []string{"a", "b", "c"}},
        {name: "wrong sep", input: "a/b/c", sep: ",", want: []string{"a/b/c"}},
        {name: "no sep", input: "abc", sep: "/", want: []string{"abc"}},
        {name: "trailing sep", input: "a/b/c/", sep: "/", want: []string{"a", "b", "c"}},
    }

    for _, tc := range tests {
        got := Split(tc.input, tc.sep)
        if !reflect.DeepEqual(tc.want, got) {
            t.Fatalf(" expected: %v, got: %v", tc.name, tc.want, got)
        }
    }
}
```

现在当测试失败时，我们可以通过描述性的名字知道测试在做什么，不再需要根据输出去分析，而且这些名字可以用来搜索。

```shell
% go test
--- FAIL: TestSplit (0.00s)
    split_test.go:25: trailing sep: expected: [a b c], got: [a b c ]
```

我们可以用 map 进一步 dry：

```go
func TestSplit(t *testing.T) {
    tests := map[string]struct {
        input string
        sep   string
        want  []string
    }{
        "simple":       {input: "a/b/c", sep: "/", want: []string{"a", "b", "c"}},
        "wrong sep":    {input: "a/b/c", sep: ",", want: []string{"a/b/c"}},
        "no sep":       {input: "abc", sep: "/", want: []string{"abc"}},
        "trailing sep": {input: "a/b/c/", sep: "/", want: []string{"a", "b", "c"}},
    }

    for name, tc := range tests {
        got := Split(tc.input, tc.sep)
        if !reflect.DeepEqual(tc.want, got) {
            t.Fatalf("%s: expected: %v, got: %v", name, tc.want, got)
        }
    }
}
```

我们将测试用例定义为，测试名到测试 fixture 的 map，而不是结构体 slice。用 map 有个附带好处，可能改进测试的有效性。

Map 迭代序是*未定义的* [^1]。这意味着我们每次运行 `go test`，测试都可能以不同的顺序运行。

[^1]: 请不要发邮件给我，争论 map 的迭代序是随机的，[它是未定义的](https://golang.org/ref/spec#For_statements)。

这点是很有帮助的，对于测试只有在按语句顺序运行时才会通过的情况。这可以是因为存在全局状态，而且后续的测试会依赖前面测试对全局状态的修改。

## 介绍子测试（sub tests）

修复失败的测试之前，我们的表驱动测试方案中有几个其他的问题需要解决。

第一个是，我们在测试用例失败时调用
`t.Fatalf`。这意味着一个测试用例失败后，就不再测试其他用例。因为测试用例以未定义的顺序运行，我们不知道这个用例是唯一会失败的，还是第一个失败的。

`testing` 包可以解决这个问题，如果我们费力的将每个测试用例都写成函数，但那太冗长了。好消息是 Go 1.7
增加了一个新功能——[子测试](https://blog.golang.org/subtests)，能轻松的为表驱动测试解决这个问题。

```go
func TestSplit(t *testing.T) {
    tests := map[string]struct {
        input string
        sep   string
        want  []string
    }{
        "simple":       {input: "a/b/c", sep: "/", want: []string{"a", "b", "c"}},
        "wrong sep":    {input: "a/b/c", sep: ",", want: []string{"a/b/c"}},
        "no sep":       {input: "abc", sep: "/", want: []string{"abc"}},
        "trailing sep": {input: "a/b/c/", sep: "/", want: []string{"a", "b", "c"}},
    }

    for name, tc := range tests {
        t.Run(name, func(t *testing.T) {
            got := Split(tc.input, tc.sep)
            if !reflect.DeepEqual(tc.want, got) {
                t.Fatalf("expected: %v, got: %v", tc.want, got)
            }
        })
    }
}
```

每个子测试都有名字，而且会在测试运行时自动打印出来。

```shell
% go test
--- FAIL: TestSplit (0.00s)
    --- FAIL: TestSplit/trailing_sep (0.00s)
        split_test.go:25: expected: [a b c], got: [a b c ]
```

每个子测试都是匿名函数，因此我们可以在保持表驱动测试紧凑的同时，使用 `t.Fatalf`、`t.Skipf` 和所有其他 `testing.T` 提供的辅助方法。

### 可以直接执行单个子测试用例

因为子测试有名字，你可以在运行子测试时用 `go test -run` 参数选择一部分名字。

```shell
% go test -run=.*/trailing -v
=== RUN   TestSplit
=== RUN   TestSplit/trailing_sep
--- FAIL: TestSplit (0.00s)
    --- FAIL: TestSplit/trailing_sep (0.00s)
        split_test.go:25: expected: [a b c], got: [a b c ]
```

## 比较所得与所想

现在我们准备好修复测试用例了，先看下错误。

```shell
--- FAIL: TestSplit (0.00s)
    --- FAIL: TestSplit/trailing_sep (0.00s)
        split_test.go:25: expected: [a b c], got: [a b c ]
```

你能指出问题吗？很明显 `reflect.DeepEqual` 在抱怨两个 slice 不同。但指出具体的不同并不容易，你必须发现 `c`
后面有个额外的空格。在这个简单的示例中还算容易发现，但当你比较两个复杂的深层嵌套 gRPC 结构体时，可有你受的。

如果我们切换到 `%#v` 语法，从而以 Go 声明的样式查看值，可以改进输出。

```go
got := Split(tc.input, tc.sep)
if !reflect.DeepEqual(tc.want, got) {
    t.Fatalf("expected: %#v, got: %#v", tc.want, got)
}
```

现在当我们运行测试时，很明显问题是 slice 中有个额外的空元素。

```shell
% go test
--- FAIL: TestSplit (0.00s)
    --- FAIL: TestSplit/trailing_sep (0.00s)
        split_test.go:25: expected: []string{"a", "b", "c"}, got: []string{"a", "b", "c", ""}
```

但在我们修复测试用例之前，我想再谈下选择正确的方式呈现测试失败。`Split`
函数很简单，接收一个字符串并返回一个字符串 slice，但如果需要打印结构体指针呢？

这个示例中 `%#v` 并不好用：

```go
func main() {
    type T struct {
        I int
    }
    x := []*T{{1}, {2}, {3}}
    y := []*T{{1}, {2}, {4}}
    fmt.Printf("%v %v\n", x, y)
    fmt.Printf("%#v %#v\n", x, y)
}
```

不出意料的，第一个 `fmt.Printf` 打印了没用的地址 slice：

```
[0xc000096000 0xc000096008 0xc000096010] [0xc000096018 0xc000096020 0xc000096028]
```

然而我们的 `%#v` 版本并没更好，打印了转型为 `*main.T` 的地址 slice：

```
[]*main.T{(*main.T)(0xc000096000), (*main.T)(0xc000096008), (*main.T)(0xc000096010)} []*main.T{(*main.T)(0xc000096018), (*main.T)(0xc000096020), (*main.T)(0xc000096028)}
```

因为使用任何 `fmt.Printf` verb 都有限制，我想介绍谷歌的 [go-cmp](https://github.com/google/go-cmp) 库。

cmp 库的目标是专门比较两个值。这类似 `reflect.DeepEqual`，但它有更多的能力。使用 cmp 包，你可以理所当然的写：

```go
func main() {
    type T struct {
        I int
    }
    x := []*T{{1}, {2}, {3}}
    y := []*T{{1}, {2}, {4}}
    fmt.Println(cmp.Equal(x, y)) // false
}
```

但对我们的测试函数更有用的是 `cmp.Diff` 函数，它会为两个值的不同之处递归的生成文本描述。

```go
func main() {
    type T struct {
        I int
    }
    x := []*T{{1}, {2}, {3}}
    y := []*T{{1}, {2}, {4}}
    diff := cmp.Diff(x, y)
    fmt.Printf(diff)
}
```

产生：

```shell
% go run
{[]*main.T}[2].I:
         -: 3
         +: 4
```

这告诉我们，`T` 的 slice 的元素 2，`I` 字段期望为 3，但实际为 4。

把这些东西整合起来，我们得到了表驱动 go-cmp 测试。

```go
func TestSplit(t *testing.T) {
    tests := map[string]struct {
        input string
        sep   string
        want  []string
    }{
        "simple":       {input: "a/b/c", sep: "/", want: []string{"a", "b", "c"}},
        "wrong sep":    {input: "a/b/c", sep: ",", want: []string{"a/b/c"}},
        "no sep":       {input: "abc", sep: "/", want: []string{"abc"}},
        "trailing sep": {input: "a/b/c/", sep: "/", want: []string{"a", "b", "c"}},
    }

    for name, tc := range tests {
        t.Run(name, func(t *testing.T) {
            got := Split(tc.input, tc.sep)
            diff := cmp.Diff(tc.want, got)
            if diff != "" {
                t.Fatalf(diff)
            }
        })
    }
}
```

运行这个我们得到：

```shell
% go test
--- FAIL: TestSplit (0.00s)
    --- FAIL: TestSplit/trailing_sep (0.00s)
        split_test.go:27: {[]string}[?->3]:
                -: <non-existent>
                +: ""
FAIL
exit status 1
FAIL    split   0.006s
```

使用 `cmp.Diff`，测试结果不单告诉我们得到的不同于想要的，还告诉我们 fixture
的第三个索引不该存在，但实际得到的输出是空字符串 `""`。知道这个，修复测试失败就容易了。

## 相关文章

1. **[Writing table driven tests in Go](https://dave.cheney.net/2013/06/09/writing-table-driven-tests-in-go)**
2. **[Why bother writing tests at all?](https://dave.cheney.net/2019/05/14/why-bother-writing-tests-at-all)**
3. **[Internets of Interest #7: Ian Cooper on Test Driven Development](https://dave.cheney.net/2018/10/15/internets-of-interest-7-ian-cooper-on-test-driven-development)**
4. **[Automatically run your package’s tests with inotifywait](https://dave.cheney.net/2016/06/21/automatically-run-your-packages-tests-with-inotifywait)**

## 译者注

GoLand 支持[生成](https://www.jetbrains.com/help/go/2021.2/create-tests.html)表驱动测试风格的测试文件，在测试文件中一键[运行单个表测试](https://www.jetbrains.com/help/go/2021.2/performing-tests.html#run-individual-table-tests)，用并排视图、直观清晰的 [diff Testify assertion](https://www.jetbrains.com/help/go/2021.2/using-the-testify-toolkit.html#compare-expected-and-actual-values) 的期望值和实际值。Map
类型的 `tests` 可以通过 [Live template](https://www.jetbrains.com/help/go/2021.2/using-live-templates.html) 轻松实现。

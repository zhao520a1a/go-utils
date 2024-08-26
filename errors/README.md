---
type: "docs"
weight: 5
title: "错误处理"
---

# 引
错误即 error，在 Go 程序中通常以 err 之名存在；错误处理即 error handling，指我们如何在 Go 程序中合理地处理错误。为了避免歧义，本文剩余内容将直接使用错误的英文单词 error 代替。

关于本话题的详细调研内容请参考这篇[博客](https://zhenghe-md.github.io/blog/2020/10/05/Go-Error-Handling-Research/)，本文的重心在于阐述最终选择的**最佳实践**。

> 希望以此此文建立更加扎实的工程实践方法论，进一步提升交付项目质量。

- 前端错误信息可读性差
- 后端服务识别报错困难
```go
if strings.Contains(err.Error(), m) {
    tracelog.Info(ctx, fmt.Sprintf("ignore err :%v", err))
    return 
}

ErrEmptyResult = errors.New(`[scanner]: empty result`)
if err == scanner.ErrEmptyResult {
    err = nil
}
```
# 道 
相关术语说明：

| **英文** | **中文** |
| --- | --- |
| error-code-based | 基于错误码 |
| exception-based | 基于异常 |
| error wrapping/unwrapping | 包装错误/解包装 |
| error inspection | 错误检查 |
| error formatting | 错误格式化 |
| error chain | 错误链表，即通过包装将错误组织成链表结构 |
| error class | 错误类别、类型 |

## 错误处理方式？
通常错误处理方案分为两种：error-code-based 和 exception-based。很早就有人 [指出 exception-based 错误处理更不利于写出优质的代码](https://devblogs.microsoft.com/oldnewthing/20050114-00/?p=36693)，也更难辨别优质和劣质的代码； Go 在设计时选择了 error-code-based error handling 方案，许多来自 Java、Python 等语言的工程师习惯了 exception-based 的方案，遇到 Go 时感到十分不习惯，本质原因是不了解语法后面的设计理念。
## Go 语法的设计理念？
### "happy path" 与 "sad path" 地位相同
如果我们将函数的正常逻辑路径称为 "happy path"，异常逻辑路径称为 "sad path"。在使用 exception-based error handling 的编程语言时，工程师认为 "sad path" 是一种需要额外考虑的特殊情况，需要特殊对待；而在 Go 开发者眼里，"happy path" 和 "sad path" 都是一般的情况，二者应该同样重要，被同等对待。
```go
// Go 鼓励工程师将逻辑的 "happy path" 留在函数缩进的最外层，而把 "sad path" 放到第二级缩进中；
// 采用 fail-fast 的策略结束执行，不要使用其它返回值的特征作为调用成功与否的依据
func Do() (ret interface{}, err error) {
    // happy path
    v1, err := A()
    if err != nil {
      // sad path 1
  }

    v2, err := B(v1)
    if err != nil {
      // sad path 2
  }

    ret = process(v1, v2)
    return
}
```
### errors are values
任何实现了 Error 接口的数据类型都是 error，它们与字符串、整数、结构体相比并没有特别之处。Go 官博中提出了[Errors are values](https://blog.golang.org/errors-are-values)的理念，鼓励开发者显式地在 error 出现的地方直接处理，任何实现了 Error 接口的数据类型都是 error，它们与字符串、整数、结构体相比并没有特别之处。
## 都有谁关心 error？
在任意一个服务的生命周期中，通常至少有 3 个角色关心 error：
### 应用程序与error
error 类型检查 面向的是应用程序，用于逻辑判断；应用程序可能拥有各种各样的外部依赖，比如第三方服务、内部 RPC 服务、数据库服务、消息队列服务等。应用程序需要能够准确、方便、健壮地获取 error 特征，从容地根据 error 的特点处理。
### 程序维护者与error
遇到线上问题时，服务的维护者接到报警后，需要根据详细的 error 信息做根源分析，这时信息越多越好，当然更高的可读性能够帮助维护人员更快地定位问题，解决问题。
PS：一个设计精良的 errors package 要能够让工程师自如地处理 error 与各个角色之间的信息传递。
### 用户 与 error
当服务运行遇到 error 时，需要向普通的 C 端用户提供友好、明确的消息提示，让他明白系统正处于异常状态，可以稍后重试或联系客服、技术人员。消息应该是对人类友好的自然语言。除此之外，系统内部的细节，如错误栈信息，不应当直接暴露给 C 端用户，对于未明确定义的 error 更应如此。主要原因在于：

- 用户不应该关心服务的实现细节
- 暴露不必要的细节可能会降低系统安全性
# 术
## error handling 涉及那些方面？

- checking：判断 error 发生与否
- inspection：检查 error 类型
- formatting ：打印 error 上下文。

很多大佬针对上面的环节，提出了自己的优化见解，如：Russ Cox 早在 2018 年末发布了两个新提议：

- [Error Handling](https://go.googlesource.com/proposal/+/master/design/go2draft-error-handling-overview.md) - 尝试解决 checking 代码冗长的问题；
- [Error Values](https://go.googlesource.com/proposal/+/master/design/go2draft-error-values-overview.md) - 尝试解决 inspection 的信息丢失以及 formatting 的上下文信息不足问题;
## error 上下文 != 堆栈信息？
为什么采用堆栈信息不行？这里引用 [GopherCon 2019](https://about.sourcegraph.com/blog/go/gophercon-2019-handling-go-errors) 中工程师的观点，原因如下：

- 它们很难阅读
- 它们很难解析
- 它们说的是哪里出错了，而不是为什么
## error 的逻辑调用栈是啥？
Ben Johnson 提出 [Failure is your Domain](https://middlemost.com/failure-is-your-domain/) 的观点，认为每个项目应当构建特有的 error handling package，并提出逻辑调用栈 (logical stack) 的概念；针对 error 多层嵌套调用场景下的上下文注入问题，采用 error wrapping 的方式给程序提供额外的信息，用于后续决策。之后我们还能做什么？

- 按严重性对错误进行分类。
- 按类型对错误进行分类。
- 添加特定于应用程序的数据。
- 查询以上所有内容
## 方案调研
> 抄作业，先看看别人咋做的？

### 官方SDK

- Go 1.13 以前：提供 Error 接口及 errors.New、fmt.Errorf 两个构建 error 的方法，建议创建一个实现错误接口的自定义结构，并将原始错误作为该结构上的字段，例如：[PathError](https://pkg.go.dev/os#PathError)。
- Go 1.13 以后：提供类似的原生解决方案，支持利用 %w 格式化符号实现 [error wrapping](https://go.dev/blog/go1.13-errors)，并提供 Unwrap、errors.Is 以及 errors.As 来解决 error inspection 的问题。
### 社区轮子
> Go 社区的开发者们因为语言本身对 error handling 的支持不足，因此创造各种各样的轮子；

| 代码库 | 特点 |
| --- | --- |
| [juju errors](https://github.com/juju/errors) | 在 wrap error 时，可以选择保留或隐藏 error 产生的原因 (cause)，但它的 Cause 方法仅 unwrap 一层。 |
| [pkg/errors](https://github.com/pkg/errors) | 由 Go 核心工程师 R Dave Cheney 的开发被广泛使用。
提供 wrapping 和调用栈捕获的功能，并利用 %+v 格式化 error，会递归地遍历 error chain，展示更多的细节，它认为只有整个 error chain 最末端的 error 最有价值。 |
| [upspin err](https://github.com/upspin/upspin/tree/master/errors) | 定制化 error 的实践范本，同时引入了 errors.Is 和 errors.Match 用于辅助检查 error 类型。 |
| [hashicorp/errwrap](https://github.com/hashicorp/errwrap) | 支持将 errors 组织成树状结构，并提供 Walk 方法遍历这棵树。

 |
| [pingcap parser errors](https://github.com/pingcap/parser/blob/release-3.0/terror/terror.go) | 在 pkg/errors 的基础上二次开发，并增加了 error 类的概念。 |
| [cockroachdb/errors](https://github.com/cockroachdb/errors) | 考虑了 error 在进程间传递的场景，让 error handling 具备网络传播兼容能力。 |

# 器
## 当前现状
目前，内部生产环境使用 Go 1.16 +，通查项目中都有这样一个基本的 error handling 工具类(虽然你有，我也不一定用)，此外服务研发团队内部并没有统一 error handling 规范与方案。
```go
type XError struct {
    code int
    err  error
}

// code 与 msg 是一一绑定的
errMsg = map[errorCode]string{
    CommonErrCode: "",
    MysqlErrCode:  "mysql",
    RedisErrCode:  "redis",
}
```

下面我整几段代码，先抛开代码逻辑的正确性，只谈 error handling 逻辑，看看错误是被如何处理的？
```go
// RPC 请求
func baseRequest(ctx context.Context, url string, params map[string]string) (interface{}, error) {
    b, err := client.RequestWithContext(ctx, conf.URL+url, client.BuildOptions(&client.Options{
      Method:           "GET",
      Headers:          map[string][]string{},
      Params:           params,
      Timeout:          conf.Timeout,
      HostWithVip:      conf.VipHost,
      RecordWithParams: true,
  }))
    if err != nil {
      // ① - 纯通过手动拼接函数名的方式，便于生成日志，格式混乱，容易遗忘；
      traceerror.Error(ctx, fmt.Sprintf("#coupon#baseRequest##error# res=%v,err=%v", string(b), err))
      // ② - 不包装 error 上下文信息，同一错误在日志展现上是不同的；
      return nil, err
  }

    var resp *Resp
    if err = json.Unmarshal(b, &resp); err != nil {
      return nil, err
  }
    if resp.Code != "1" {
      return nil, xerror.New(resp.Message)
  }
    return resp.Data, nil
}
```
```go
func OnlineStrategyFlow(ctx context.Context, flowID int64) (map[int64]*UpFlowResult, int64, error) {
    // ③ - 方便复用对，函数名做了声明，但日志拼接格式还是很混乱；
    fun := "#STRATEGY#OnlineStrategyFlow#"
    logger.Infof(ctx, "%s 策略流上线开始执行.flowID:%d", fun, flowID)
    strategies, err := service3.QueryStrategyInFlow(ctx, flowID)
    if err != nil {
      return nil, 0, err
  }

    if len(strategies) == 0 {
      logger.Errorf(ctx, "%s 策略流下策略信息不存在.flow_id:%d", fun, flowID)
      // ④ - 使用了特定 error，但却丢失了错误上下文；
      return nil, 0, xerror.ErrDBResEmpty
  }

}
```
```go
func (o *Crowd) ScanRecordSyncData(ctx context.Context, bind xgin.Bind) (ret interface{}, err error) {
    var req crowd_service.ScanParam
    if err = bind(&req); err != nil {
      // ⑤ - 采用了 error code ，但没有当前函数信息；
      return nil, xerror.WrapCode(xerror.CommonErrCode, err)
  }
    return o.CrowdModule.ScanRecordSyncData(ctx, req)
}
```
造成的后果：

- error 被多次日志打印，重复冗余干扰问题定位，并影响监控配置；
- error 判断困难，无法区分特定error；
- error 上下文信息丢失，没有当前函数信息，更无法串联出逻辑调用栈；

==> **前端报错没有可读性，全靠 traceId 和个人经验来定位问题！！**




## error 的消费者

在写出合理的 error handling 代码前，我们有必要思考：谁是 error 的消费者？在任意一个服务的生命周期中，通常至少会有 3 类人关心 error：

* 应用程序本身 (application)
* 服务的用户 (end user)
* 服务的维护者 (operator)
 
## error handling 的最佳实践

接下来我们阐述如何利用 dry 的 errors 模块完成我们日常项目开发的 error handling 需求。

### 替换 errors

在很多地方，我们会用到 Go 自带的 errors  package，比如：

```go
import (
	"errors"
)

func GetOneWorkflow(ctx context.Context, db manager.XDB, where map[string]interface{}) (*bpm.Workflow, error) {
	if nil == db {
		return nil, errors.New("manager.XDB object couldn't be nil")
	}
	//...
}
```

这里用到了 `errors.New`，在 dry 的 errors package 中，我们同样提供了 `New` 函数，因此只需替换依赖即可：

```go
func GetOneWorkflow(ctx context.Context, db manager.XDB, where map[string]interface{}) (*bpm.Workflow, error) {
	if nil == db {
		return nil, errors.New("manager.XDB object couldn't be nil")
	}
	//...
}
```

### 用 errors.Op 替换 fun

只要阅读过老项目的代码，你就不可避免地会遇到 `fun` 这个变量，比如：

```go
func CheckIAMAuth(ctx context.Context, header, uri, method string) bool {
	fun := "CheckIAMAuth -->"
  // ...
}
```

`fun` 存在的意义就是打日志的时候将当前的函数名称记录下来，方便排查问题时从 Kibana 上快速检索。我们可以利用 `errors.Op` 来代替 `fun`：

```go
func CheckIAMAuth(ctx context.Context, header, uri, method string) bool {
  op := errors.Op("CheckIAMAuth")
  // ...
}
```

下文我们会继续介绍 op 如何替代 fun 的功能，并提供更简洁、完整的错误信息。

### 指定 error 的类型

尽管在应用程序中遇到的具体 error 种类可能很多，但仔细梳理可以发现它们可以被划归到有限的几类：

```go
package errors

const (
	Conflict   Class = "conflict"          // Action cannot be performed
	Internal   Class = "internal"          // Internal error
	Invalid    Class = "invalid"           // Validation failed
	NotFound   Class = "not_found"         // Entity does not exist
	PermDenied Class = "permission_denied" // Does not have permission
	Other      Class = "other"             // Unclassified error
)
```

在使用 dry 的 `errors.E` 方法时，我们可以直接指定 error 的类型，比如：

```go
errors.E(op, errors.Invalid, "message")
errors.E(op, err, errors.Internal, "message")
//...
```

### 指定 error code

如果你的应用会直面 C 端用户，暴露 error code 通常能减少更多的沟通成本，通过 `errors.E` 可以指定 error code：

```go
const (
  ErrCodeTrafficLimited int = 100001 // 流量限制
  ErrCodeInvalidToken int32 = 100002 // 无效 token
  //...
)

errors.E(op, err, errors.Conflict, ErrCodeTrafficLimited)
errors.E(op, err, errors.Invalid, ErrCodeInvalidToken)
```

任意整型 (int, int32) 参数传入 `errors.E` 都将被认为是 error code，因此你可以在应用中自行定义 error code。

### 指定 error msg

不论面对 C 端用户还是运维人员，精简、完备的文案可以大大提高问题定位的速度，通过 `errors.E` 可以制定 error msg：

```go
errors.E(op, err, "this is message")
errors.E(op, err, fmt.Sprintf("paramA %v paramB %v", paramA, paramB))
//...
```

由于 `errors.E` 函数因为其功能和使用的便利性需要，会根据参数的类型决定它的用途，因此如果同一类型的参数被传入两次，只有最后的参数会覆盖生效 (相关讨论见 [issue](https://gitlab.pri.ibanyu.com/quality/dry/issues/5))。

### 无中生有的 error

在一些情况下，如参数不符合要求、找不到数据、与已知数据冲突、没有权限等，我们需要创建新的 error，这时候优先使用 `errors.E`

```go
func (m *handlerGroupAgent) List(ctx context.Context, db manager.XDB, ...) (records []*bpm.HandlerGroup, err error) {
	op := errors.Op("DefaultHandlerGroupAgent.List")

	if query == nil {
		err = errors.E(op, errors.Invalid, "query condition is nil")
		return
	}
  //...
}
```

当然，你也可以使用 `errors.New`，但会失去当前执行的函数名称 (即 op)、error 类型 (即 errors.Invalid) 等信息。

### 包装下层 error

在访问下层函数、RPC、DB 时，会返回 error，使用 `errors.E` 来对它合理包装：

```go
func (m *WorkflowInstanceService) Add(ctx context.Context, ins *bpm.WorkflowInstance) (lastInsertID int64, err error) {
	op := errors.Op("WorkflowInstanceService.Add")
  // ...
  rawOpContext, err := to.JSON(ins.OpContext)
	if err != nil {
    err = errors.E(op, err, errors.Invalid)
		return
	}
  // ...
}
```

如果你发现单个函数中会重复出现多次类似的代码，可以在 `defer` 中统一处理，让代码更加精简：

```go
func (m *WorkflowInstanceService) Add(ctx context.Context, ins *bpm.WorkflowInstance) (lastInsertID int64, err error) {
	op := errors.Op("WorkflowInstanceService.Add")
	
	defer func() {
		if err != nil {
			err = errors.E(op, err)
		}
	}()
 
  rawOpContext, err := to.JSON(ins.OpContext)
  if err != nil {
		return
	}
  // ...
}
```

在同一个函数中，可以对同一个 errors 执行多次包装操作，但实际上只有第一次包装会生效，如：

```go
func (m *WorkflowInstanceService) Add(ctx context.Context, ins *bpm.WorkflowInstance) (lastInsertID int64, err error) {
	op := errors.Op("WorkflowInstanceService.Add")

	defer func() {
		if err != nil {
			err = errors.E(op, err)
		}
	}()
 
  rawOpContext, err := to.JSON(ins.OpContext)
  if err != nil {
    err = errors.E(op, err, errors.Invalid)
		return
	}
  // ...
}
```

这时候只有 `err = errors.E(op, err, errors.Invalid)` 会生效，`defer` 中的包装不会生效。

需要特别注意的是：**理论上，每个错误在向上抛出的过程中，应当只被指定一次类型，且应该在最接近跨进程交互的地方指定，中间层仅做简单包装和透传**，即 `errors.E(op, err)`。

### 满足应用的消费需求：利用 `errors.Is` 判断 error 类型

在上层处理 error 时，如 controller，可以利用 `errors.Is` 判断 error 类型：

```go
func (m *WorkflowController) AddWorkflow(ctx context.Context, req *AddWorkflowReq) (res *AddWorkflowRes, err error) {
	op := errors.Op("ProcessController.AddWorkflow")
  // ...
  last, err := m.workflowService.Get(ctx, map[string]interface{}{
		"name":     name,
		"_orderby": "created_at DESC",
	})
	if err != nil && !errors.Is(err, errors.NotFound) {
		return res, errors.E(op, err)
	}
  // ...
}
```

### 满足维护者的消费需求：在最上层逻辑中打印 error 信息

在整个 error chain 上，error 应当只被打印一次，且一般在最上层打印。下面的示例代码中有三个方法：GrpcServiceImpl.HandleDelUser、UserService.DelUser 以及 DelUser，前面的依赖后面的，形成依赖链条 (chain)。

```go
// grpc/user.go
func (m *GrpcServiceImpl) HandleDelUser(ctx context.Context, id int64) {
  op := errors.Op("HandleDelUser")
  //...
  err = m.userService.DelUser(ctx, id)
  if err != nil {
    xlog.Errorf(ctx, errors.E(op, err))
    return
  }
}

// service/user.go
func (m *UserService) DelUser(ctx context.Context, id int64) (err error) {
  op := errors.Op("UserService.DelUser")
  //...
  _, err = tidb.DelUser(where)
  if err != nil {
   	err = errors.E(op, err, errors.Internal)
    return
  }
  //...
}

// tidb/user.go
func DelUser(ctx context.Context, db manager.XDB, where map[string]interface{}) (n int64, err error) {
  op := errors.Op("tidb.DelUser")
  //...
  if err != nil {
    err = errors.E(op, err, errors.NotFound, fmt.Sprintf("where %v", where))
    return
  }
  //...
}
```

最底层的 DelUser 将 error 往上抛，在 handler 上打印错误信息，会得到：

```sh
HandleDelUser: UserService.DelUser: tidb.DelUser [not found] where ...
```

我们可以看到**逻辑调用栈**和实际 error 发生地上下文参数信息。

### 满足用户的消费需求：`errors.Code` 和 `errors.Msg`

如果你利用「指定 error code」一节给出的方案在 error chain 上记录了 error code，就可以利用 `errors.Code` 将其取出；类似地，`errors.Msg` 会从 error chain 上最近的包含 error msg 的节点中取出 error msg。我们可以利用这两个函数在 http/grpc 响应中填入相应的信息：

```go
res.Errinfo = &grpcutil.ErrInfo{
  Code: int32(errors.ErrCode(err)),
  Msg:  errors.ErrMsg(err),
}
```

## 其它功能函数

### 合并多个 errors

有时候你的逻辑可能是并发地访问外部服务，不论单次访问是否有问题，都不想停止其它访问的正常进行，这时可能会产生一个或多个 error。这些 errors 之间是并列关系，这时我们可以使用 `errors.Combine` 来处理：

```go
// ...
var errs []error

for _, req := range reqs {
  data, err := req.Do()
  if err != nil {
    errs = append(errs, err)
  }
}

err = errors.Combine(errs)
// ...
```

`errors.Combine` 会妥善处理好 errs 为空的情况，因此不必担心特殊情况影响代码书写逻辑。

---
type: "docs"
weight: 5
title: "错误处理"
---

# 错误处理

错误即 error，在 Go 程序中通常以 err 之名存在；错误处理即 error handling，指我们如何在 Go 程序中合理地处理错误。为了避免歧义，本文剩余内容将直接使用错误的英文单词 error 代替。

关于本话题的详细调研内容请参考这篇[博客](https://zhenghe-md.github.io/blog/2020/10/05/Go-Error-Handling-Research/)，本文的重心在于阐述最终选择的**最佳实践**。

## error 的消费者

在写出合理的 error handling 代码前，我们有必要思考：谁是 error 的消费者？在任意一个服务的生命周期中，通常至少会有 3 类人关心 error：

* 应用程序本身 (application)
* 服务的用户 (end user)
* 服务的维护者 (operator)

我们的 dry (don't repeat yourself) [项目](https://gitlab.pri.ibanyu.com/quality/dry)中的[errors](https://gitlab.pri.ibanyu.com/quality/dry/tree/master/errors) 模块，正是为这个目的而诞生。

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
import (
	"gitlab.pri.ibanyu.com/quality/dry.git/errors"
)

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
---
type: "docs"
weight: 2
title: "项目布局"
---

# 项目布局规范

通常当工程师进入新环境，会从维护老项目开始切入，逐渐熟悉公司的技术栈和效率工具。这时候老项目的一些习惯，如：代码风格、项目布局、合作开发等等，不论是好的还是坏的，都会不自觉地影响新人，最后形成路径依赖。在这个过程中，如果没有人主动去思考为什么，那么这些习惯也将被无理由地继承下去。

本节要谈的就是其中很重要的一点：**项目布局**。

## 有缺陷的布局方案

在正式介绍最佳实践之前，有必要先了解常见的有缺陷的布局方案。

### 方案一：Monolithic package

Monolithic package 就是把项目的所有文件放在同一个 package 内部，这种方案的优势就是简单，永远不会有循环依赖；缺陷也很明显，当项目规模变大时，将变得难以维护：单个文件内代码量过大、文件之间的依赖变混乱，甚至对于 IDE 来说也是不小的挑战。

在公司内部的一些早期小型项目，就存在 monolithic package 的影子，比如 `base/configuration` 服务的布局方案：

```shell
.
├── Dockerfile
├── docker
├── go.mod
├── logic
│   ├── config.go
│   ├── dbitem.go
│   └── logic.go
├── main.go
└── processor
    └── proc_thrift.go
```

所有的代码逻辑都在 `logic` 文件夹中。

### 方案二：Rails-style layout

Rails 风格的布局方案将代码按照功能类型 (functional type) 拆分。比如将 handlers、controllers 以及 models 分别放到不同的 package 中，这种方案我们也曾经使用过，比如 `base/uauth` 服务的布局：

```shell
.
├── Dockerfile
├── Makefile
├── Makefile~
├── README.md
├── build.sh
├── common
│   └── env.go
├── controller
│   ├── httpcontroller
│   └── thriftcontroller
├── docker
├── go.mod
├── go.sum
├── main.go
├── model
│   ├── dao
│   ├── daoimpl
│   ├── domain
│   ├── error.go
│   └── model.go
├── pkg
│   ├── cache
│   └── config
├── router
│   ├── httprouter
│   └── thriftrouter
└── sql
    └── uauth.sql
```

这种方案有两个问题：**命名啰嗦**和**循环依赖**。假如你在 router 中引用 controller，就可能出现 `controller.UserController` 这种名称；在 controller 中引用 model 就可能出现 `dao.UserDao` 这种名称。啰嗦的名称会给代码的阅读者带来一定的负担，对于有洁癖的开发者来说也是难以忍受的存在，但这个问题并不致命，Rails 风格项目布局更大的问题在于循环依赖。功能之间很容易产生循环依赖，除非维护者能够在划分功能模块时保证单向依赖，但通常当项目变得复杂时这将很难做到。

### 方案三：Group by module

正如公司的组织架构，既存在按照功能划分成行政、人力、设计、产品、研发等职能部门，也存在按照模块划分成不同的 BU，每个 BU 内部再划分出各自的职能部门，项目的布局亦是如此。如：你可以建一个 users package，负责用户相关的所有逻辑，从 model 到 controller。

这种方案的缺陷与 Rails 风格方案类似，会产生啰嗦的命名，如 `users.User`，同时不同的模块之间极易产生循环依赖。

## 理想的布局方案

理想的布局方案应当满足哪些要求？

* 容易理解、上手、可维护性好
* 命名精简、可读性好
* 避免循环依赖
* 方便构建单元测试和集成测试
* ...

阅读 Ben Johnson 提出的 Standard Package Layout (SPL) 后，我根据公司当前的基础设施和发布系统特点，提出一个改良版的 SPL。改良版的布局方案可以用 4 句话概括：

* 将领域类型 (domain types) 放在与项目同名的 subpackage 中
* 按照依赖关系来组织 subpackages
* 使用共享的 mock subpackage
* 利用每个 subpackage 的 init 函数注入依赖

以下是 `bpm` 的基本布局：

```shell
.
├── Dockerfile
├── Makefile
├── README.md
├── bpm
├── bpmn
├── config
├── doc
├── docker
├── engine
├── go.mod
├── go.sum
├── grpc
├── http
├── main.go
├── mock
├── notifier
├── oss
├── rpc
├── sql
└── tidb
```

接下来，我们利用 `bpm` 项目来解析上面 4 句话的含义。

### 将领域类型 (domain types) 放在与项目同名的 subpackage 中

每个应用都会有独特的语言来描述数据和数据交互，这种语言便是领域 (domain)。如果你有一个电子商务应用，那么领域就可能包含顾客、账号、信用卡、库存等等；如果你是 Facebook，那么领域就可能包含用户、赞以及关注关系。领域是建立在实现层之上的抽象，它与所使用的技术无关。

bpm 是流程管理系统，其领域中包含的两个核心概念是工作流 (workflow) 及其实例 (workflow instance)。以 workflow 为例，我们可以在 bpm 项目中的 bpm subpackage 中定义它的数据结构及服务接口：

```go
// bpm/workflow.go
type Workflow struct {
	ID           int64                    `json:"id" bdb:"id"`
	Name         string                   `json:"name" bdb:"name"`
	Version      int64                    `json:"version" bdb:"version"`
	ProjectID    int64                    `json:"project_id" bdb:"project_id"`
	ProjectName  string                   `json:"project_name" bdb:"project_name"`
	DeployStatus bpm.WorkflowDeployStatus `json:"deploy_status" bdb:"deploy_status"`
	XMLUri       string                   `json:"xml_uri" bdb:"xml_uri"`
	CreatedBy    string                   `json:"created_by" bdb:"created_by"`
	UpdatedBy    string                   `json:"updated_by" bdb:"updated_by"`
	CreatedAt    time.Time                `json:"created_at" bdb:"created_at"`
	UpdatedAt    time.Time                `json:"updated_at" bdb:"updated_at"`
	Def          *bpmn.Def                `json:"-" bdb:"-"`
}

type WorkflowService interface {
	Add(ctx context.Context, workflow *Workflow) (lastInsertID int64, err error)
	Get(ctx context.Context, where map[string]interface{}) (workflow *Workflow, err error)
	Set(ctx context.Context, workflow *Workflow) (rowsAffected int64, err error)
	Del(ctx context.Context, where map[string]interface{}) (rowsAffected int64, err error)
	List(ctx context.Context, where map[string]interface{}) (workflows []*Workflow, total int64, err error)
}
```

bpm subpackage 仅包含像 `Workflow` 这样的简单的领域数据结构以及 `WorkflowService` 这样的服务接口定义，我们称之为领域类型。你也可以往这里添加 `Workflow` 的其它方法，但前提是**这些方法不能依赖项目中的其它任何 package，也不能依赖外部服务或数据库**。原因在于它将成为所有其它 subpackages 互相依赖的**支点**，这么做可以保证消除循环依赖。

为了方便理解下文，我们接着在 bpm subpackage 中定义 `XMLStorageService`，负责存储和读取 xml 文件：

```go
// bpm/xml_storage.go
type XMLStorageService struct {
  UploadXML(ctx context.Context, data []byte) (uri string, err error)
	LoadXML(ctx context.Context, uri string) (data []byte, err error)
}
```

以及 `DBManager`，负责管理数据库连接：

```go
type DBManager interface {
	Begin(ctx context.Context) (*manager.Tx, error)
	GetDB(ctx context.Context) (*manager.DB, error)
}
```

### 按照依赖关系来组织 subpackages

既然 bpm subpackage 不能依赖外部服务、数据库，我们就必须将这些依赖推到其它 subpackages 中，这些 subpackage 将作为领域类型具体实现的适配器 (adapter)。

举例如下：假设 `WorkflowService` 背后的持久化存储是 TiDB，我们就可以引入 tidb subpackage，后者负责提供 `bpm.WorkflowService` 的具体实现，由于 workflow 的详细配置信息是一个 xml 文件，它不会被存放在关系型数据库中，因此 `WorkflowService` 还将依赖 `bpm.XMLStorageService`：

```go
// tidb/workflow.go
package tidb

import (/*...*/)

type WorkflowService struct {
  db 							  *sql.DB
  xmlStorageService bpm.XMLStorageService
}

func (m *WorkflowService) Add(ctx context.Context, wf *bpm.Workflow) (lastInsertID int64, err error) {/*...*/}
func (m *WorkflowService) Get(ctx context.Context, where map[string]interface{}) (workflow *Workflow, err error) {/*...*/}
// ...
```

当我们有一天想更换持久化层的实现，比如使用 MongoDB，BoltDB，就可以类似地再引入一个 mongo subpackage 或者 bolt subpcakge。

此外，我们还可以利用这种方式来引入 subpackages 之间的层级依赖。假设你想要在 TiDB 层前面添加缓存层，那么可以引入一个 inmemory subpackage，后者以 tidb subpackage 为后端，实现基于内存的 LRU 缓存逻辑：

```go
// inmemory/user.go
package inmemory

import (/**/)

type WorkflowCache struct {
  cache map[int]*Workflow
  service bpm.WorkflowService
}

func (m *WorkflowCache) Add(ctx context.Context, wf *bpm.Workflow) (lastInsertID int64, err error) {/*...*/}
func (m *WorkflowCache) Get(ctx context.Context, where map[string]interface{}) (workflow *Workflow, err error) {/*...*/}
```

这里的关键点在于：

* 其它 subpackages 作为 bpm subpackage 的适配器
* 其它 subpackages 之间的依赖都通过 bpm subpackage 来中转

这样就能有效地消除 subpackages 之间的循环依赖。我们也可以从 Go 的标准 library 中看到这种布局，如：`io.Reader` 是 io 的领域知识，`tar.Reader`、`gzip.Reader` 以及 `multipart.Reader` 这些都是 `io.Reader` 的实现，同时这些实现之间也存在层级 (layered) 依赖关系，我们会看到 `os.File` 被包裹在 `bufio.Reader`中、`bufio.Reader` 被包裹在 `gzip.Reader` 中、`gzip.Reader` 被包裹在 `tar.Reader` 中。

#### 控制 subpackages 之间的依赖

subpackages 之间不仅只存在线性的层次依赖，即 A 依赖 B、B 依赖 C，还可能存在嵌套依赖，如 A 依赖 B 和 C，如上文中的 `WorkflowService` 同时依赖 `DBManager` 以及 `XMLStorageService`。其中 `XMLStorageService` 通过 OSS 来实现。当我们想要更换 `XMLStorageService` 实现时，无需修改任何 `WorkflowService` 的实现代码逻辑；当我们想要更换 `WorkflowService` 实现时，无需修改任何 `XMLStorageService` 的实现，二者之间的依赖关系仅靠 bpm subpackage 定义的领域类型维系，耦合度很低。

#### 用 subpackage 控制对标准包的依赖

上述这种技巧并不局限于控制外部依赖，我们也可以用它来控制对标准包的依赖。比如，net/http package 属于标准包，我们也可以在项目中引入 http subpackage，来控制对 net/http 的依赖。

```go
// http/handler.go
package http

import (/*...*/)

type Handler struct {
  WorkflowService bpm.WorkflowService
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  // handle request
}
```

这总做法粗看起来很奇怪，为什么要取一个和标准包一样的名字，如果某个地方需要同时引用 http 和 net/http，岂不是很尴尬？实际上这种设计是有意而为之，只要你不允许项目的其它地方引用 net/http，问题就不存在了，而这种限制恰恰能够帮助你从源头上将所有对 net/http 的依赖控制在 http subpackage 中，这样项目的依赖关系将变得更加清晰。

现在，`http.Handler` 就成为领域类型与 HTTP 协议之间的适配器。

### 使用共享的 mock subpackage

现在，所有的 subpackages 之间都依靠 bpm subpackage 中的定义的领域类型作为沟通的桥梁，我们就很容易通过依赖注入的方式实现 mock。

假设我们希望利用本地的数据库来做简单的 end-to-end test，就可以引入公共的 mock subpackage，在里面实现简单的 mock，同样以 `WorkflowService` 为例，引入 DBManager 的 mock：

```go
// mock/db_manager.go
type DBManager struct {
	BeginFn      func(ctx context.Context) (*manager.Tx, error)
	BeginInvoked bool

	GetDBFn      func(ctx context.Context) (*manager.DB, error)
	GetDBInvoked bool
}

func (m *DBManager) Begin(ctx context.Context) (*manager.Tx, error) {
	m.BeginInvoked = true
	return m.BeginFn(ctx)
}

func (m *DBManager) GetDB(ctx context.Context) (*manager.DB, error) {
	m.GetDBInvoked = true
	return m.GetDBFn(ctx)
}
```

剩下的工作就是在测试时，将数据库本地化的实现注入到 `BeginFn` 和 `GetDBFn` 中，然后在初始化时将 `mockDBManager` 传递给 `WorkflowService` 即可。

### 利用每个 subpackage 的 init 函数注入依赖

设计好整体布局后，只需要一根线将它们串联起来。这根线就是每个 subpackage 的 init 函数，以 grpc subpackage 中的 init 函数为例：

```go
// grpc/init.go
import (
	"gitlab.pri.ibanyu.com/server/bpm/service.git/engine"
	"gitlab.pri.ibanyu.com/server/bpm/service.git/notifier"
	"gitlab.pri.ibanyu.com/server/bpm/service.git/oss"
	"gitlab.pri.ibanyu.com/server/bpm/service.git/rpc"
	"gitlab.pri.ibanyu.com/server/bpm/service.git/tidb"
)

var HandleGrpcBPM *GrpcBPM

func init() {
	var workflowCtl = NewWorkflowController(tidb.DefaultWorkflowService, oss.DefaultXMLStorageService)
	var workflowInstanceCtl = NewWorkflowInstanceController(tidb.DefaultWorkflowService, tidb.DefaultWorkflowInstanceService)

	HandleGrpcBPM = &GrpcBPM{
		workflowCtl:            workflowCtl,
		workflowInstanceCtl:    workflowInstanceCtl,
	}
}
```

利用 Go 的启动机制串联依赖关系，这很 Go。

## 代码生成

为了更好的实施这一项目布局方案，我们搭建了代码生成 CLI [pf](https://gitlab.pri.ibanyu.com/quality/pf)，点击链接了解详情。

## 参考文献

* [Standard Package Layout --- Ben Johnson](https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1)
* [Practical Go: package design --- Dave Cheney](https://dave.cheney.net/practical-go/presentations/gophercon-singapore-2019.html#_package_design)
* [Gitlab: server/bpm/service](https://gitlab.pri.ibanyu.com/server/bpm/service/commit/e231bcd3032b46f896a09c6ecd9d9ae36133adc1)
* [Building WTF Dial](https://medium.com/wtf-dial/wtf-dial-domain-model-9655cd523182)
* [WTF Dial: Data storage with BoltDB](https://medium.com/wtf-dial/wtf-dial-boltdb-a62af02b8955)

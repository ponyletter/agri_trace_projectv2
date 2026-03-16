

# 农产品溯源系统 - 部署与使用指南

> **项目名称**：基于国密区块链的农产品溯源系统（以阿克苏苹果为例）
> **技术栈**：Go 1.22 + Gin + GORM + MySQL 8.0 + Redis + Hyperledger Fabric v2.2（国密版）+ IPFS + Vue3 + 微信小程序
> **论文支撑**：《基于区块链的农产品溯源关键技术研究》第五章

---

## 目录结构

```
agri_trace_project/
├── backend/                    # Go 后端项目
│   ├── main.go                 # 主程序入口
│   ├── go.mod / go.sum         # Go 模块依赖
│   ├── config/
│   │   ├── config.go           # 配置加载模块
│   │   └── config.yaml         # 配置文件（数据库/Redis/区块链/IPFS）
│   ├── controller/
│   │   ├── auth.go             # 用户认证控制器（登录/注册/获取信息）
│   │   └── trace.go            # 溯源核心控制器（批次/记录/查询）
│   ├── middleware/
│   │   ├── jwt.go              # JWT 鉴权中间件
│   │   └── cors.go             # CORS 跨域中间件
│   ├── model/
│   │   └── model.go            # GORM 数据模型定义
│   ├── pkg/
│   │   ├── blockchain/
│   │   │   └── client.go       # Fabric 链码调用（含优雅降级 Mock）
│   │   ├── ipfs/
│   │   │   └── client.go       # IPFS 文件上传（含优雅降级 Mock）
│   │   └── gmsm/
│   │       └── sm3.go          # 国密 SM3 哈希工具
│   ├── router/
│   │   └── router.go           # Gin 路由注册
│   ├── utils/
│   │   └── database.go         # 数据库初始化工具
│   └── Dockerfile              # 后端 Docker 镜像构建文件
├── frontend/                   # Vue3 前端项目
│   ├── src/
│   │   ├── main.ts             # 应用入口
│   │   ├── App.vue             # 根组件
│   │   ├── router/index.ts     # 路由配置
│   │   ├── api/
│   │   │   ├── request.ts      # Axios 请求封装（含 JWT 拦截器）
│   │   │   └── trace.ts        # API 接口定义
│   │   └── views/
│   │       ├── login/          # 登录页面
│   │       ├── DashboardView.vue # 系统概览
│   │       ├── batch/          # 批次管理页面
│   │       ├── trace/          # 溯源记录与公开查询页面
│   │       └── admin/          # 用户管理页面（管理员专用）
│   └── vite.config.ts          # Vite 构建配置
├── miniprogram/                # 微信小程序
│   ├── app.js / app.json       # 小程序全局配置
│   ├── utils/api.js            # 网络请求封装
│   └── pages/
│       ├── index/              # 溯源查询首页（扫码/手动输入）
│       └── result/             # 溯源结果展示页（时间轴）
├── db/
│   └── init.sql                # MySQL 初始化脚本（含模拟数据）
├── scripts/                    # Python 图表生成脚本
│   ├── run_all_charts.py       # 一键生成所有图表
│   ├── gen_usecase_diagram.py  # 用例图（图5-1）
│   ├── gen_architecture_diagram.py # 架构图（图5-2）
│   ├── gen_er_diagram.py       # 数据库ER图（图5-3）
│   ├── gen_trace_flow_diagram.py # 流程图（图5-4/5-5）
│   └── gen_performance_charts.py # 性能测试图（图5-8~5-11）
├── pic/                        # 生成的论文图片
├── data/                       # 性能测试 CSV 数据
├── docker-compose.yml          # Docker Compose 编排文件
├── nginx.conf                  # Nginx 反向代理配置
└── 第五章_系统实现与测试评估.md  # 第五章论文正文
```

---

## 快速启动（推荐：Docker Compose 一键部署）

### 前置条件

```bash
# 确保已安装以下工具
docker --version       # >= 24.0
docker compose version # >= 2.0
go version             # >= 1.22（本地开发时需要）
node --version         # >= 18（前端构建时需要）
python3 --version      # >= 3.9（生成图表时需要）
```

### 步骤一：构建前端静态文件

```bash
cd frontend
pnpm install
pnpm build
# 构建产物在 frontend/dist/ 目录
```

### 步骤二：修改后端配置

编辑 `backend/config/config.yaml`，根据实际环境修改以下配置：

```yaml
database:
  dsn: "root:你的密码@tcp(mysql:3306)/agri_trace?charset=utf8mb4&parseTime=True&loc=Local"

blockchain:
  mock_mode: true  # 无 Fabric 网络时设为 true，系统自动降级
  # mock_mode: false  # 连接真实 Fabric 网络时设为 false

ipfs:
  mock_mode: true  # 无 IPFS 节点时设为 true
```

### 步骤三：启动所有服务

```bash
# 在项目根目录执行
# 基础模式（不启动 IPFS，适合快速体验）
docker compose up -d

# 完整模式（含 IPFS 节点）
docker compose --profile full up -d

# 查看服务状态
docker compose ps

# 查看后端日志
docker compose logs -f backend
```

### 步骤四：访问系统

| 服务 | 地址 | 说明 |
| :--- | :--- | :--- |
| Web 管理端 | http://localhost:80 | Vue3 前端界面 |
| 后端 API | http://localhost:8080 | Go 后端服务 |
| 健康检查 | http://localhost:8080/health | 服务状态检查 |
| 溯源查询（示例） | http://localhost/trace/BATCH-APPLE-20251025-001 | 无需登录 |

**默认账号**（密码均为 `123456`，需先运行 init.sql 以实际mysql数据为准）：

| 账号 | 角色 | 说明 |
| :--- | :--- | :--- |
| admin | 管理员 | 全部权限 |
| farmer01 | 种植户 | 创建批次、录入种植/采收信息 |
| inspector01 | 质检员 | 录入质检结果 |
| transporter01 | 物流商 | 录入运输信息 |
| retailer01 | 销售商 | 录入上架信息 |

> **注意**：init.sql 中的密码哈希为占位符，首次使用请通过管理员账号的注册接口重置密码，或直接修改 SQL 中的 `password_hash` 字段为 bcrypt 哈希值。

---

## 本地开发模式（无 Docker）

### 启动后端

```bash
# 1. 确保 MySQL 和 Redis 已启动
# 2. 导入数据库
mysql -u root -p agri_trace < db/init.sql

# 3. 启动后端
cd backend
export PATH=$PATH:/usr/local/go/bin
go run main.go
# 后端监听 :8080
```

### 启动前端

```bash
cd frontend
pnpm install
pnpm dev
# 前端监听 :9528，API 请求自动代理至 :8080
```

---

## 连接真实 Fabric 网络

当 `config.yaml` 中 `blockchain.mock_mode` 设为 `false` 时，系统将尝试连接真实 Fabric 网络。

### 使用 fabric-gm（国密版）

本项目参考 [fabric-gm](https://github.com/ponyletter/fabric-gm) 进行国密改造。启动步骤如下：

```bash
# 1. 克隆国密 Fabric 项目
git clone https://github.com/ponyletter/fabric-gm.git
cd fabric-gm

# 2. 启动测试网络（需要 Docker）
cd test-network
./network.sh up createChannel -c agrichannel -ca

# 3. 部署链码（链码路径需根据实际调整）
./network.sh deployCC -ccn agritrace -ccp ../chaincode/agritrace -ccl go

# 4. 将 Fabric 连接配置文件复制到后端
cp organizations/peerOrganizations/org1.example.com/connection-org1.yaml \
   ../agri_trace_project/backend/config/fabric/config.yaml

# 5. 修改 backend/config/config.yaml
#    blockchain.mock_mode: false
#    blockchain.config_path: ./config/fabric/config.yaml
```

---

## 生成论文图表

```bash
# 安装 Python 依赖
pip3 install matplotlib numpy

# 一键生成所有图表
cd scripts
python3 run_all_charts.py

# 图片保存在 ../pic/ 目录
# CSV 数据保存在 ../data/ 目录
```

---

## API 接口文档

### 认证接口

| 方法 | 路径 | 说明 | 鉴权 |
| :--- | :--- | :--- | :--- |
| POST | `/api/v1/auth/login` | 用户登录，返回 JWT Token | 否 |
| GET | `/api/v1/auth/profile` | 获取当前用户信息 | 是 |
| POST | `/api/v1/auth/register` | 注册新用户（仅管理员） | 是（admin） |

### 批次管理接口

| 方法 | 路径 | 说明 | 鉴权 |
| :--- | :--- | :--- | :--- |
| POST | `/api/v1/batches` | 创建农产品批次 | 是（farmer/admin） |
| GET | `/api/v1/batches` | 获取批次列表 | 是 |

### 溯源接口

| 方法 | 路径 | 说明 | 鉴权 |
| :--- | :--- | :--- | :--- |
| POST | `/api/v1/trace/records` | 添加溯源节点记录（触发上链） | 是 |
| GET | `/api/v1/trace/records` | 获取溯源记录列表 | 是 |
| GET | `/api/v1/trace/:trace_code` | 根据溯源码查询完整溯源链 | **否**（供小程序使用） |
| GET | `/api/v1/block/info?height=N` | 查询区块信息 | 否 |

### 请求示例（添加溯源记录）

```bash
curl -X POST http://localhost:8080/api/v1/trace/records \
  -H "Authorization: Bearer <your_jwt_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "batch_id": 1,
    "node_type": "harvesting",
    "operation_time": "2025-10-25 09:00:00",
    "location": "优质苹果种植示范基地A区",
    "env_data": {
      "weather": "晴",
      "method": "人工采摘",
      "sugar_content": "18%"
    }
  }'
```

---

## 微信小程序配置

1. 使用微信开发者工具打开 `miniprogram/` 目录
2. 修改 `app.js` 中的 `globalData.baseUrl` 为后端实际域名（需已备案 HTTPS）
3. 在微信公众平台配置合法域名
4. 编译并上传审核

---

## 常见问题

**Q: 启动后端报 "数据库连接失败"？**
A: 检查 MySQL 是否已启动，且 `config.yaml` 中的 DSN 配置正确。

**Q: 上链操作返回 "Fabric 网络未连接"？**
A: 将 `config.yaml` 中 `blockchain.mock_mode` 设为 `true`，系统将使用模拟数据，不影响业务功能演示。

**Q: 如何生成真实的 bcrypt 密码哈希？**
A: 使用以下 Go 代码生成：
```go
import "golang.org/x/crypto/bcrypt"
hash, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
fmt.Println(string(hash))
```

---

```
frontend/src/views/login/LoginView.vue        ← 覆盖此文件
frontend/src/views/trace/TracePublicView.vue  ← 覆盖此文件

1. 纯前端溯源系统：https://github.com/sugar2026dev/agricultural-trace-systemgit  可以参考
2. 电子证照(Fabric+IPFS)：https://gitee.com/surgar2022/fabric-eCert-trace-ipfs_ok  这个是目前go vue Fabric ipfs可以直接启动的 完整成功项目,但缺少国密，也不是阿克苏苹果溯源系统。
3. 国密底层基础库：https://github.com/tjfoc/gmsm
4. Fabric 国密改造版：https://github.com/ponyletter/fabric-gm  把这个当作你写好的，已经有的。 
我的相关资料
https://gitee.com/surgar2022/grad-proj-block-chain-agri_md
```



```
很好，目前代码可以跑，没问题了，最新版本我上传到了github上，https://github.com/ponyletter/agri_trace_project。 但还需要增加元素和功能，数据库表可能还需要多一些，表和表属性结构可能需要增加。还有sql 数据也尽可能多一些，包括用户，图片也需要实际给出，图片要实际的图片，如果有地图，位置也需要理论可行的。 
例如：UI元素
表单组件：输入框、下拉选择框、数字加减器（如采摘数量）、日期/时间选择器。

多媒体组件：图片上传控件（用于上传现场采摘照片）、验证码图片。

导航与布局组件：左侧多级菜单侧边栏、顶部面包屑导航、数据统计卡片（仪表盘）。

数据可视化组件：折线图（预期与实际数据对比）、列表表格、操作日志卡片视图。

弹窗/模态框组件：电子合格证弹窗、区块链交易详情弹窗、地理位置定位弹窗。

特色业务组件：二维码（用于扫码溯源） 一物一码溯源 这个需要有 、防伪印章（政府部门盖章）

例如：功能模块

用户与权限管理：系统管理员登录、溯源参与者（农户、加工商等）注册与登录、消费者独立溯源入口等 ，目前微信小程序我已经写好了，只需要留出接口即可，和之前一样。登录界面  需要有注册，注册也需要验证码，web端溯源应该是单独的一个界面，而不是在登录界面，登录密码也是之前的github上的，123456。 

数据大屏展示（控制台）：宏观展示平台访客、消息、交易金额、订单量及趋势走势

种植环节管理：记录作物类型、施肥、灌溉、除草周期及详细地址

数据采集与上传：参与者上传原材料信息、地理位置、生产日期，并支持上传现场实拍图片佐证

加工环节管理：精细化记录每一次加工操作（如：清洗、烘干、包装），并绑定操作人ID、电话和时间，责任到人

物流追踪定位：记录产品流转的部门、时间，并提供具体的地理位置定位信息，可以的话，增加地图功能，可视化

防伪与认证：系统自动生成带有二维码、防伪承诺、检测打勾和官方印章的电子合格证

区块链底层存证：展示Hyperledger Fabric的底层交易日志、区块哈希、以及联盟链的节点架构。 

请你参考
https://github.com/ponyletter/agri_trace_project，按照之前的思路，重新给出项目压缩包。

```





```
你现在是一位资深的 Go 语言全栈架构师、区块链国密底层专家，同时也是一位精通农业信息化方向硕士学位论文写作的学术指导专家。
我正在撰写一篇题为《基于区块链的农产品溯源关键技术研究》的硕士论文，以“阿克苏苹果”作为具体溯源实例。本项目参考并融合了以下开源项目的架构与技术思想：
1. 纯前端溯源系统：https://github.com/sugar2026dev/agricultural-trace-systemgit  可以参考
2. 电子证照(Fabric+IPFS)：https://github.com/ponyletter/fabric-eCert-trace-ipfs_ok  这个是目前go vue Fabric ipfs可以直接启动的 完整成功项目,但缺少国密，也不是阿克苏苹果溯源系统。可能需要先启动下这个。 
3. 国密底层基础库：https://github.com/tjfoc/gmsm
4. Fabric 国密改造版：https://github.com/ponyletter/fabric-gm  把这个当作你写好的，已经有的。 
我的相关资料
https://github.com/ponyletter/gradProjBlockChainAgri_md
硕士学位论文v2.md是目前的 pic 第五章参考其他硕士论文中相关web图片
技术栈为：后端 Go 1.22 (Gin + GORM)；数据库 MySQL 8.0 + Redis；底层基于 Hyperledger Fabric v2.2（国密版）+ IPFS；前端为 Vue3 + 微信小程序。前四章已定稿，现在需要你帮我完成代码落地支撑、图表生成以及《第五章》的撰写。
请严格按照以下 4 个模块的任务要求，一步步输出完整的成果：
### 任务一：设计符合国标的 MySQL 数据库 (输出 SQL 脚本) sql尽量通用，不要出现品牌名字，不要出现地区等相关。
1. 参考《GB/T 29373-2012 农产品追溯要求 果蔬》，设计规范的表结构：用户角色表、农产品批次表(agri_batches)、溯源节点流转表(trace_records，包含种植、质检装箱、运输、终端上架等，参考pic下的图片)、IPFS文件关联表等。
2. 提供丰富的 `INSERT` 初始模拟数据（以2025年12月阿克苏苹果为例，注意实际情况，实际的苹果成熟季节，必须能完美支撑后续的系统查询测试。包含详尽的字段中文注释。
### 任务二：核心 Go 后端代码与“优雅降级”机制 (输出 Go 代码)
1. 提供基于 Gin + GORM 的核心架构代码（路由、JWT鉴权中间件、核心控制器）。
2. 提供微信小程序溯源查询 API（如 `/api/v1/trace/{trace_code}`），返回时间轴数据、区块哈希、高度、上链时间等。
3. **重点（Mock模式）**：在 Fabric/IPFS 调用层实现降级机制。当检测到区块链未启动时，系统不报错，而是静默读取 MySQL 数据，并自动生成模拟的国密区块哈希（0x开头）和 IPFS CID，确保前端和小程序在无链环境下也能完美运行并展示界面。
### 任务三：系统架构 UML 与 自动化生成脚本 (输出 Python 代码)
为了方便我直接放入 Word，如果可以不要只输出 Mermaid 文本，而是提供一个可以直接运行的 **Python 脚本 (`generate_uml.py`)**。
1. 该脚本内部包含 多个Mermaid 代码字符串：系统用例图、核心ER图、溯源时序交互图等。尽量多一些。 
2. 可以的话 把相关图表自动渲染并下载保存为 `pic/系统用例图.png`、`pic/核心ER图.png`、`pic/溯源时序图.png`等。
3. 确保脚本会自动创建 `pic` 目录。
### 任务四：撰写论文《第五章 农产品溯源系统实现与测试》 (输出标准 Markdown)
请结合上述所有代码与设计，撰写直接可用于毕业论文的《第五章》文本。要求：
1. **学术文风**：严谨书面语，杜绝口语化，杜绝滥用解释性括号。术语（如BCCSP、链码）首次说明后统一使用中文或缩写。
2. **内容结构**：
   - 5.1 系统环境搭建（描述服务器、国密Fabric镜像编译、IPFS部署，融入提及 ponyletter/fabric-gm 和 tjfoc/gmsm）。
   - 5.2 核心功能模块实现（描述Gin认证、上链逻辑，结合 github/gitee 参考项目思想）。 注意及Mock降级机制不要说，实际应该是默认已经启动了。是最终的https://github.com/ponyletter/gradProjBlockChainAgri_md/硕士学位论文v2.md论文
   - 5.3 系统功能展示例如`![图5-1 溯源系统用例图](pic/图5-1  溯源系统用例图.jpg)`、`![图5-3 系统登录页面](pic/图5-3  系统登录页面.png)`、`![图5-4 溯源系统工作界面](pic/图5-4  溯源系统工作界面.png)`、`![图5-5 区块链浏览器界面](pic/区块链浏览器界面.png)`、`![图5-6 移动端溯源查询首页](pic/图101 溯源查询首页小程序.png)`、`![图5-7 移动端溯源查询结果页](pic/图100 溯源查询结果页小程序.png)`）。
   还有之前生成的，也可以插入，目前的pic只是参考，因为我是先写代码再写程序的，实际可以替换
   - 5.4 系统性能初步评估（结合你生成的性能图表和数据，客观分析系统性能与国密开销），可以的话用caliper测试，得出，目前的grad-proj-block-chain-agri_md pic下可以参考，也可以测试下其他的得出数据，得出数据需要先存入csv，也需要给出python代码，下次我直接运行pyhton即可得到图片。。
最后给出部署教程。 
请深呼吸，发挥你最顶尖的架构和学术能力，后端go 前端vue，可以直接启动项目，无报错。 另外确保 Python 脚本可以直接复制运行并出图，方便我Markdown 文本可以直接复制进 Word。现在开始输出！给出压缩包。 注意这个项目可能比较大，docker需要linux系统。 
```



````
你现在是一位资深的 Go 语言全栈架构师、区块链国密底层专家，同时也是一位精通农业信息化方向硕士学位论文写作的学术指导专家。




我正在撰写一篇题为《基于区块链的农产品溯源关键技术研究》的硕士论文，以“阿克苏苹果”作为具体溯源实例。本项目参考并融合了以下开源项目的架构与技术思想：

1. 纯前端溯源系统：https://github.com/sugar2026dev/agricultural-trace-systemgit  可以参考

2. 电子证照(Fabric+IPFS)：https://gitee.com/surgar2022/fabric-eCert-trace-ipfs_ok

3. 国密底层基础库：https://github.com/tjfoc/gmsm

4. Fabric 国密改造版：https://github.com/ponyletter/fabric-gm  把这个当作你写好的，已经有的。 

我的相关资料

https://gitee.com/surgar2022/grad-proj-block-chain-agri_md

硕士学位论文v2.md是目前的 pic 第五章参考其他硕士论文中相关web图片

技术栈为：后端 Go 1.22 (Gin + GORM)；数据库 MySQL 8.0 + Redis；底层基于 Hyperledger Fabric v2.2（国密版）+ IPFS；前端为 Vue3 + 微信小程序。前四章已定稿，现在需要你帮我完成代码落地支撑、图表生成以及《第五章》的撰写。




请严格按照以下 4 个模块的任务要求，一步步输出完整的成果：




### 任务一：设计符合国标的 MySQL 数据库 (输出 SQL 脚本)

1. 参考《GB/T 29373-2012 农产品追溯要求 果蔬》，设计规范的表结构：用户角色表、农产品批次表(agri_batches)、溯源节点流转表(trace_records，包含种植、质检装箱、运输、终端上架等，参考pic下的图片)、IPFS文件关联表等。

2. 提供丰富的 `INSERT` 初始模拟数据（以2025年12月阿克苏苹果为例，注意实际情况，实际的苹果成熟季节，必须能完美支撑后续的系统查询测试。包含详尽的字段中文注释。




### 任务二：核心 Go 后端代码与“优雅降级”机制 (输出 Go 代码)

1. 提供基于 Gin + GORM 的核心架构代码（路由、JWT鉴权中间件、核心控制器）。

2. 提供微信小程序溯源查询 API（如 `/api/v1/trace/{trace_code}`），返回时间轴数据、区块哈希、高度、上链时间等。

3. **重点（Mock模式）**：在 Fabric/IPFS 调用层实现降级机制。当检测到区块链未启动时，系统不报错，而是静默读取 MySQL 数据，并自动生成模拟的国密区块哈希（0x开头）和 IPFS CID，确保前端和小程序在无链环境下也能完美运行并展示界面。




### 任务三：系统架构 UML 与 自动化生成脚本 (输出 Python 代码)

为了方便我直接放入 Word，如果可以不要只输出 Mermaid 文本，而是提供一个可以直接运行的 **Python 脚本 (`generate_uml.py`)**。

1. 该脚本内部包含 多个Mermaid 代码字符串：系统用例图、核心ER图、溯源时序交互图等。尽量多一些。 

2. 可以的话 把相关图表自动渲染并下载保存为 `pic/系统用例图.png`、`pic/核心ER图.png`、`pic/溯源时序图.png`等。

3. 确保脚本会自动创建 `pic` 目录。





### 任务四：撰写论文《第五章 农产品溯源系统实现与测试》 (输出标准 Markdown)

请结合上述所有代码与设计，撰写直接可用于毕业论文的《第五章》文本。要求：

1. **学术文风**：严谨书面语，杜绝口语化，杜绝滥用解释性括号。术语（如BCCSP、链码）首次说明后统一使用中文或缩写。

2. **内容结构**：

   - 5.1 系统环境搭建（描述服务器、国密Fabric镜像编译、IPFS部署，融入提及 ponyletter/fabric-gm 和 tjfoc/gmsm）。

   - 5.2 核心功能模块实现（描述Gin认证、上链逻辑，结合 github/gitee 参考项目思想）。 注意及Mock降级机制不要说，实际应该是默认已经启动了。是最终的https://gitee.com/surgar2022/grad-proj-block-chain-agri_md/硕士学位论文v2.md论文

   - 5.3 系统功能展示例如`![图5-1 溯源系统用例图](pic/图5-1  溯源系统用例图.jpg)`、`![图5-3 系统登录页面](pic/图5-3  系统登录页面.png)`、`![图5-4 溯源系统工作界面](pic/图5-4  溯源系统工作界面.png)`、`![图5-5 区块链浏览器界面](pic/区块链浏览器界面.png)`、`![图5-6 移动端溯源查询首页](pic/图101 溯源查询首页小程序.png)`、`![图5-7 移动端溯源查询结果页](pic/图100 溯源查询结果页小程序.png)`）。

   还有之前生成的，也可以插入，目前的pic只是参考，因为我是先写代码再写程序的，实际可以替换

   - 5.4 系统性能初步评估（结合你生成的性能图表和数据，客观分析系统性能与国密开销），可以的话用caliper测试，得出，目前的grad-proj-block-chain-agri_md pic下测试，可以参考，也可以测试下其他的得出数据，得出数据需要先存入csv，也需要给出python代码，下次我直接运行pyhton即可得到图片。。

最后给出部署教程。 

请深呼吸，发挥你最顶尖的架构和学术能力，后端go 前端vue，可以直接启动项目，无报错。 另外确保 Python 脚本可以直接复制运行并出图，方便我Markdown 文本可以直接复制进 Word。现在开始输出！



````



````
你现在是一位资深的 Go 语言全栈架构师、区块链专家，同时也是一位精通计算机/农业信息化方向硕士学位论文写作的学术指导专家。
我正在撰写一篇题为《基于区块链的农产品溯源关键技术研究》的硕士论文，以“阿克苏苹果”作为具体溯源实例。技术栈已确定为：后端 Go 1.22 (Gin + GORM)；数据库 MySQL 8.0 + Redis；底层区块链基于 Hyperledger Fabric v2.2（国密改造版）+ IPFS；前端为 Vue3 管理端 + 微信小程序查询端。

现在我已经完成了前四章的理论和设计，需要你帮我完成核心的代码落地支撑、数据库设计以及论文《第五章：系统实现与测试》的撰写。

请严格按照以下 4 个模块的任务要求，一步步输出完整的成果：

### 任务一：设计全面且符合论文要求的 MySQL 数据库 (输出完整 SQL 脚本)
1. 请参考《GB/T 29373-2012 农产品追溯要求 果蔬》国家标准，设计一套完整、规范的关系型数据库表结构。
2. 表必须包含：用户角色表(sys_users, sys_roles)、农产品批次表(agri_batches)、溯源节点流转表(trace_records，需包含种植、质检装箱、冷链运输、终端上架等环节)、IPFS文件关联表(ipfs_files)。
3. **关键要求**：必须提供丰富的 `INSERT` 初始模拟数据（以阿克苏苹果为例，时间设定在2025年12月）。这些数据需要完美支撑论文中的 ER图、实体图以及系统测试。
4. 在 SQL 脚本中加入详尽的字段中文注释。

### 任务二：提供后端核心逻辑代码与“优雅降级”机制 (输出 Go 代码)
1. 给出基于 Gin + GORM 的核心结构代码（如路由定义、鉴权中间件、核心业务 Controller）。
2. 提供微信小程序所需的 API 接口代码（如 `/api/v1/trace/{trace_code}`），返回的数据结构需完全匹配溯源时间轴的展示需求（包含区块哈希、区块高度、上链时间、流转节点信息等）。
3. **重点机制（优雅降级/Mock模式）**：由于本地不一定随时启动 Fabric 国密网络和 IPFS 节点，请在区块链调用层（Service层）实现一个 Mock 机制：当检测到 Fabric/IPFS 未启动或连接失败时，系统**不要报错**，而是静默降级为直接读取 MySQL 数据库中的数据，并自动生成模拟的区块哈希（以 `0x` 开头）和模拟的 CID。确保前端和微信小程序在任何情况下都能正常运行和展示完整界面，方便我截图。

### 任务三：生成辅助论文配图的 UML 代码 (输出 Mermaid 代码)
为了方便我在 Word 中插入正规的系统架构图和设计图，请使用 Mermaid 语法生成以下图表代码（我将在本地渲染成图片）：
1. **系统用例图 (Use Case Diagram)**：包含管理员、种植户、加工商、物流商、消费者等角色的用例。
2. **核心实体关系图 (ER Diagram)**：基于任务一设计的表结构。
3. **农产品供应链溯源流程图 (Sequence Diagram)**：从苹果采摘到消费者扫码查询的时序交互过程。

### 任务四：撰写论文《第五章 农产品溯源系统实现与测试》 (输出标准的 Markdown 文本)
请根据前述代码和本项目的整体设计，为我输出可以直接复制到 Word 中的《第五章》完整内容。要求如下：
1. **语言风格**：严谨的学术论文书面语，杜绝口语化，不要出现多余的解释性括号。技术名词首次出现时定义，后续直接使用缩写。
2. **内容结构**：
   - 5.1 系统环境搭建（描述服务器、Go环境、Docker部署Fabric国密镜像、IPFS部署等）。
   - 5.2 核心功能模块实现（描述Gin后端认证、数据防篡改上链逻辑、小程序接口设计及平滑降级机制）。
   - 5.3 系统功能测试（描述用户管理测试、数据录入测试。**必须精准嵌入以下图片占位符，因为我已经准备好了这些图片**：`![图5-1 溯源系统用例图](pic/图5-1  溯源系统用例图.jpg)`、`![图5-3 系统登录页面](pic/图5-3  系统登录页面.png)`、`![图5-4 溯源系统工作界面](pic/图5-4  溯源系统工作界面.png)`、`![图5-5 区块链浏览器界面](pic/区块链浏览器界面.png)`。重点描述小程序的测试，并插入 `![图5-6 移动端溯源查询首页](pic/图101 溯源查询首页小程序.png)` 和 `![图5-7 移动端溯源查询结果页](pic/图100 溯源查询结果页小程序.png)`）。
   - 5.4 系统性能初步评估（结合国密算法的时间开销，说明高并发下的TPS和响应时间）。

请深呼吸，确保各个模块逻辑连贯。代码要求可运行、无明显Bug，Markdown排版要求精美规范。现在开始你的输出：


这个目前提示词还是需要修改，实际应该以
````



```
结合之前的信息 输入输出回答，还有目前的相关信息
相关代码：
纯农产品溯源部分前端：
https://github.com/sugar2026dev/agricultural-trace-systemgit
电子证照前后端 ipfs 区块链：
https://gitee.com/surgar2022/fabric-eCert-trace-ipfs_ok
苏州同济区块链研究院实现的代码:
https://github.com/tjfoc/gmsm
目前fabric国密改造：
https://github.com/ponyletter/fabric-gm
国家标准 GBT 29373-2012农产品追溯要求 果蔬.md

准备搞定第五章
很好，目前我需要写程序了，写的程序主要有sql，还需要有默认数据，每个表都有默认数据，注意苹果作为具体的溯源应用实例，表尽量全，方便word中介绍表，数据库设计图，方便写 总体框架图，实体图，用例图，er图等，mysql sql数据尽量word都有。 
小程序也需要留下接口
03/10/2026  01:27 PM            76,689 图100 溯源查询结果页小程序.png
03/10/2026  01:27 PM            29,503 图101 溯源查询首页小程序.png
目前有图片 已经写好了。
需要结合之前的
图5-1  溯源系统用例图.jpg
功能模块等 登录工作界面，技术路线图.png
图3-1 农产品中心化溯源模式
图3-2 农产品供应链溯源流程图
图3-2 农产品供应链溯源流程图
还有国密，尽量都有。 
方便的话，sql 实体图er图。
还有截图，界面登录 测试功能等实际的，也可以直接截图或者录制启动视频，方便我直接插入word中。
如果不启动fabric-gm 国密的数据库，则默认应该也会返回值，也就是前端后端 就可以起来整个系统，并且是实际测试环境，最终的，不用在前端或者后端提示。
如果可以生成对应图片，然后在一个markdown形式保存，图片在一个目录中。
请你给出对应的 可以跑的代码，还有第五章内容markdown形式。
请你根据这个需求，给出一个提示词 给manus用的。请你给出对应的提示词，尽量成段，描述清晰。
给出manus 的提示词，我让他来写代码。 


PBFT区块链共识改进
https://github.com/fangvv/SPBFT
Parlia GitHub
BFT-SMaRt GitHub
https://github.com/bft-smart/library

Tendermint GitHub
https://github.com/tendermint/tendermint

https://github.com/apache/incubator-resilientdb/tree/resilientdb-legacy-eurosys23
用
GitHub 地址：https://github.com/hyperledger-caliper/caliper
测压
原始 PBFT 论文：《Practical Byzantine Fault Tolerance》 - 该论文详细描述了 PBFT 的工作原理及其设计。

https://www.scs.stanford.edu/nyu/03sp/sched/bfs.pdf
https://github.com/corgi-kx/blockchain_consensus_algorithm/tree/master/pbft
https://github.com/fares017/T-PBFT/blob/main/T-PBFT_An_EigenTrust-based_practical_Byzantine_fault_tolerance_consensus_algorithm.pdf



```





````
[1] 霍红,钟海岩. 农产品供应链质量安全中区块链技术投入的演化分析[J].运筹与管理,2023,32(01):15-21.
 [2] 王新庄. 食品安全问题探讨及法律规制研究——评《食品安全法原理》[J].食品安全质量检测学报,2022,13(17):5769.
 [3] 陆秋俊. 基于物联网技术构建现代农业种植及食品溯源系统[J].现代农业科技,2019,(22):252-253.
 [4] Lu Yi,Li Peng,Xu He. A Food anti-counterfeiting traceability system based on Blockchain and Internet of Things[J].Procedia Computer Science,2022,199
 [5] Jin Hai,Xiao Jiang. Towards trustworthy blockchain systems in the era of “Internet of value”: development, challenges, and future trends[J].Science China Information Sciences,2021,65(5):
 [6] 倪雪莉,马卓,王群. 区块链P2P网络及安全研究[J].计算机工程与应用,2024,60(05):17-29.
 [7] Tabatabaei Mohammad Hossein,Vitenberg Roman,Veeraragavan Narasimha Raghavan. Understanding blockchain: Definitions, architecture, design, and system comparison[J].Computer Science Review,2023,50
 [8] 司冰茹,肖江,刘存扬,戴小海,金海. 区块链网络综述[J].软件学报,2024,35(02):773-799.
 [9] 曹琪,阮树骅,陈兴蜀,兰晓,张红霞,金泓键. Hyperledger Fabric平台的国密算法嵌入研究[J].网络与信息安全学报,2021,7(01):65-75.
参考文献：
[10] 江巧玲.东源县农产品质量安全监管问题研究[D].   导师：喻国华.   仲恺农业工程学院,   2020.
[11] 朱祉琴. 浅谈食品卫生法与安全现状分析[J].食品安全导刊,2020,(18):38-39.
[12] 赵阳,孟慧敏. 我国重要产品追溯体系建设实践和对策建议[J].轻工标准与质量,2024,(05):131-134.
[13] 柳祺祺,夏春萍. 基于区块链技术的农产品质量溯源系统构建[J].高技术通讯,2019,29(03):240-248.
[14] 雷志军.基于区块链的农产品溯源信息系统研究[D].   导师：陈德海.   江西理工大学,   2022.
[15] Hu Sensen,Huang Shan,Huang Jing,Su Jiafu. Blockchain and Edge Computing Technology Enabling Organic Agricultural Supply Chain: A Framework Solution to Trust Crisis[J].Computers &amp; Industrial Engineering,2020,(prepublish):
[16] Monteiro Emiliano Soares,da Rosa Righi Rodrigo,Barbosa Jorge Luis Victória,Alberti Antônio Marcos. APTM: A Model for Pervasive Traceability of Agrochemicals[J].Applied Sciences,2021,11(17):
[17] Wang Lu,Xu Longqin,Zheng Zhiying,Liu Shuangyin,Li Xiangtong,Cao Liang,Li Jingbin,Sun Chuanheng. Smart Contract-Based Agricultural Food Supply Chain Traceability[J].IEEE ACCESS,2021,9
[18] Yiu Neo C. K.. Decentralizing Supply Chain Anti-Counterfeiting and Traceability Systems Using Blockchain Technology[J].Future Internet,2021,13(4):
[19] Dey Somdip,Saha Suman,Singh Amit Kumar,McDonaldMaier Klaus. FoodSQRBlock: Digitizing Food Production and the Supply Chain with Blockchain and QR Code in the Cloud[J].Sustainability,2021,13(6):
[20] Dey Somdip,Saha Suman,Singh Amit Kumar,McDonaldMaier Klaus. SmartNoshWaste: Using Blockchain, Machine Learning, Cloud Computing and QR Code to Reduce Food Waste in Decentralized Web 3.0 Enabled Smart Cities[J].Smart Cities,2022,5(1):
[21] 王晶宇,马兆丰,徐单恒,段鹏飞. 支持国密算法的区块链交易数据隐私保护方案[J].信息网络安全,2023,23(03):84-95.
[22] 王家峰. 基于混合算法的互联网访问用户身份认证方法[J].齐齐哈尔大学学报(自然科学版),2024,40(03):5-10.
[23] 刘丁宁.基于国密SM2批量验签的区块链系统的研究与应用[D].   导师：袁玉宇.   北京邮电大学,   2022.
[24] 富强,张路. 基于汽车生态圈的区块链应用模型构建研究[J].专用汽车,2024,(11):64-66.
[25] 曹锋林.区块链在跨境供应链中的应用研究[D].   导师：马栋林.   兰州理工大学,   2023.
[26] 蒋司琪.基于区块链技术的有机农产品溯源系统的设计与实现[D].   导师：刘志镜;李小勇.   西安电子科技大学,   2022.
[27] Agricultural product traceability system based on blockchain technology.W Xie, X Zheng, X Lu, X Lin… - … IEEE Intl Conf on Parallel & …, 2019 - ieeexplore.ieee.orgW Xie, X Zheng, X Lu, X Lin, X Fan2019 IEEE Intl Conf on Parallel & Distributed Processing with …, 2019•ieeexplore.ieee.org
[28] 周文欣.基于区块链技术的农产品溯源系统面临的问题及对策.南方农机,2021-12-28
[29] 王蒙蒙.基于区块链的生鲜农产品供应链系统研究.江西理工大学,2023-05-25
[30] 史亮 张复宏 刘文军.基于区块链的果蔬农产品追溯体系研究.农村经济与科技,2019-08-20
[31] 赵超.区块链技术与农产品供应链融合发展研究.新疆财经,2020-10-15
[32] 高潇.基于区块链的产品溯源认证技术研究.临沂大学,2023-05-01
[33] 王亚伟.区块链技术在食品可信溯源中的应用研究.天津科技大学,2022-06-01
[34] 史亮.基于区块链+物联网的果蔬农产品供应链追溯体系研究.山东农业大学,2020-06-05
[35] 于合龙 陈邦越 徐大明 杨信廷 孙传恒.基于区块链的水稻供应链溯源信息保护模型研究.农业机械学报,2020-06-22 14:19
[36] 魏文彦.基于区块链的农产品安全溯源平台系统设计与实现.武汉轻工大学,2022-06-01
[37] 冯家琦.基于区块链的数据溯源方法研究与应用.重庆邮电大学,2021-06-03
[38] 胡俊辉.基于区块链的水产品交易溯源系统研究与实现.上海海洋大学,2022-05-01
[39] 陈锦雯 罗得寸 唐呈俊 唐晨钧 丁勇.基于区块链的农业物联网可信溯源体系.信息安全学报,2022-03-15
[40] 陈少军.基于区块链的供应链信息可溯源方法研究.信息技术与信息化,2020-08-28

````

```
基于区块链的农产品溯源关键技术研究
区块链；农产品溯源；Hyperledger Fabric；IPFS；国密算法
区块链网络
基于区块链的农产品溯源关键技术研究
需要两篇
```

```
目前我需要先写word，然后再写程序。
相关代码：
纯农产品溯源部分前端：
https://github.com/sugar2026dev/agricultural-trace-systemgit
电子证照前后端 ipfs 区块链：
https://gitee.com/surgar2022/fabric-eCert-trace-ipfs_ok
苏州同济区块链研究院实现的代码:
https://github.com/tjfoc/gmsm
目前fabric国密改造：
https://github.com/ponyletter/fabric-gm



还用了小程序 查询，已经写好了，可能需要在前面提及下，并且插入到word中，并且word中也需要体现，最后给出markdown形式即可。 
03/10/2026  01:27 PM            76,689 图100 溯源查询结果页小程序.png
03/10/2026  01:27 PM            29,503 图101 溯源查询首页小程序.png



C:\02media\tarim\grap_paper\grad-proj-block-chain-agri_md\pic>dir
 Volume in drive C is Win11Pro X64
 Volume Serial Number is 20BE-1486

 Directory of C:\02media\tarim\grap_paper\grad-proj-block-chain-agri_md\pic

03/10/2026  04:44 PM    <DIR>          .
03/10/2026  03:44 PM    <DIR>          ..
03/10/2026  02:54 PM           157,170 2014年-2024年中国苹果产量及增速情况.png
03/10/2026  01:27 PM            66,662 BCCSP上的SM2算法接口有用性和有效性验证结果.png
03/10/2026  01:27 PM            73,490 Fabric平台国密算法嵌入设计思路.png
03/10/2026  01:27 PM            73,095 使用hyperleger fabric实现的业务网络.png
03/10/2026  01:27 PM            16,447 农产品中心化溯源模式.png
03/10/2026  01:27 PM           106,564 农产品供应链溯源流程图.png
03/10/2026  01:27 PM           107,488 力软JAVA快速开发平台技术栈.png
03/10/2026  01:27 PM           117,147 区块链浏览器界面.png
03/10/2026  01:27 PM            76,689 图100 溯源查询结果页小程序.png
03/10/2026  01:27 PM            29,503 图101 溯源查询首页小程序.png
03/10/2026  01:27 PM           136,148 图102 普刊录用通知.png
03/10/2026  01:27 PM         1,465,979 图102软著.png
03/10/2026  01:27 PM           208,036 图103_启动区块链网络.png
03/10/2026  01:27 PM           158,480 图4-2  BCCSP 国密算法接口实现示意.png
03/10/2026  01:27 PM            49,743 图4-3  fabric-CA 组件下 util 文件新添方法示意.png
03/10/2026  01:27 PM            53,009 图5-1  溯源系统用例图.jpg
03/10/2026  01:27 PM            17,382 图5-2  系统功能模块.png
03/10/2026  01:27 PM            86,846 图5-3  系统登录页面.png
03/10/2026  01:27 PM            75,871 图5-4  溯源系统工作界面.png
03/10/2026  01:27 PM            46,455 基于区块链的农产品溯源模型.png
03/10/2026  01:27 PM           378,751 密码算法时间开销对比.png
03/10/2026  01:27 PM            94,661 技术路线图.png
              22 File(s)      3,595,616 bytes
               2 Dir(s)  349,404,872,704 bytes free

C:\02media\tarim\grap_paper\grad-proj-block-chain-agri_md\pic>



请你分析目前这个word，是否有严重缺陷，应该以https://gitee.com/surgar2022/fabric-eCert-trace-ipfs_ok技术栈为准，例如为什么出现Hyperledger Fabric Java SDK
，实际上用的是go ，代码实际使用 Go 语言，而且应该是论文格式。 例如
3.2  系统总体架构设计
3.4  应用服务层设计
5.2  核心功能模块实现
还有这个得用了https://github.com/ponyletter/fabric-gm 是否word中写了官网镜像？还是国密改造后的？应该是基于 官网的 源码进行编译，然后这个https://github.com/ponyletter/fabric-gm 应该是自己实现的。 基于苏州同济区块链研究院实现的代码:
https://github.com/tjfoc/gmsm 这个开源的代码基础上。 
还有查询https://github.com/ponyletter/fabric-gm具体版本，v2.2吗？
请你分析这个毕业论文，是否有严重缺陷，具体应该如何替换？
还有一些专用名词 容器化部署 latest 是否可以去掉，不适合在论文中？请你分析，如何解决。如何替换。 

请你先分析，然后给出修改后对应的markdown形式的硕士论文，图片用pic目录下的。说明具体修改替换了哪里。 方便我替换word，
最后给出对应的markdown毕业论文，注意目前word中内容都需要。 
```

```
目前我需要写ppt，预答辩用。 
相关资料我都上传到git了
https://github.com/ponyletter/ppt_makingv2
这个项目是我之前的中期markdown笔记和ppt，学术报告等内容
例如：
03/13/2026  11:20 AM         2,725,167 基于区块链技术的农产品溯源关键技术研究_排版后最终.pptx
03/13/2026  11:16 AM         2,605,334 基于区块链技术的农产品溯源关键技术研究学术汇报.pptx
03/13/2026  11:20 AM             9,578 学术报告 ppt.md
03/13/2026  11:16 AM            11,816 中期ppt.md
03/13/2026  11:16 AM           118,266 硕士学位论文v3.md
可以参考：03/13/2026  02:28 PM             8,679 基于区块链技术的农产品溯源关键技术研究 硕士学位论文预答辩.md，这个
包括但不限于上面内容。 

目前我需要参加学校预答辩，例如插入部分pic下的图片 尽量做到图文并茂，参考之前的中期模板和预答辩等内容
项目目录：
C:\02media\doc\ppt_makingv2>dir
 Volume in drive C is Win11Pro X64
 Volume Serial Number is 20BE-1486

 Directory of C:\02media\doc\ppt_makingv2

03/13/2026  02:28 PM    <DIR>          .
03/13/2026  11:16 AM    <DIR>          ..
03/13/2026  11:16 AM    <DIR>          excel_to_coding
03/13/2026  02:28 PM    <DIR>          pic
03/13/2026  11:16 AM               884 README.en.md
03/13/2026  11:16 AM               974 README.md
03/13/2026  11:16 AM            11,816 中期ppt.md
03/13/2026  11:16 AM            13,793 国家标准 GBT 29373-2012农产品追溯要求 果蔬.md
03/13/2026  02:28 PM             8,679 基于区块链技术的农产品溯源关键技术研究 硕士学位论文预答辩.md
03/13/2026  11:16 AM         1,950,081 基于区块链技术的农产品溯源关键技术研究_v13.docx
03/13/2026  11:20 AM         2,725,167 基于区块链技术的农产品溯源关键技术研究_排版后最终.pptx
03/13/2026  11:16 AM         2,605,334 基于区块链技术的农产品溯源关键技术研究学术汇报.pptx
03/13/2026  11:16 AM            18,224 塔里木大学研究生中期考核登记表（专业型）.md
03/13/2026  11:16 AM            33,751 塔里木大学研究生学位论文开题报告.md
03/13/2026  11:20 AM             9,578 学术报告 ppt.md
03/13/2026  11:16 AM             3,076 提示词_.txt
03/13/2026  11:16 AM            21,512 普刊-安徽农学通报-可信农产品供应链的区块链构建路径研究.md
03/13/2026  11:16 AM           114,395 硕士学位论文.md
03/13/2026  11:16 AM           118,266 硕士学位论文v3.md
03/13/2026  02:28 PM            21,305 预答辩PPT模板草稿 (Markdown形式).md
              16 File(s)      7,656,835 bytes
               4 Dir(s)  340,920,545,280 bytes free

C:\02media\doc\ppt_makingv2>

图片目录：
C:\02media\doc\ppt_makingv2>cd pic

C:\02media\doc\ppt_makingv2\pic>dir
 Volume in drive C is Win11Pro X64
 Volume Serial Number is 20BE-1486

 Directory of C:\02media\doc\ppt_makingv2\pic

03/13/2026  02:28 PM    <DIR>          .
03/13/2026  02:28 PM    <DIR>          ..
03/13/2026  11:16 AM           157,170 2014年-2024年中国苹果产量及增速情况.png
03/13/2026  11:16 AM            66,662 BCCSP上的SM2算法接口有用性和有效性验证结果.png
03/13/2026  11:16 AM            73,490 Fabric平台国密算法嵌入设计思路.png
03/13/2026  11:16 AM           247,805 Hyperledger Fabric 联盟链结构.png
03/13/2026  11:16 AM            83,012 TPS并发曲线.png
03/13/2026  11:16 AM            73,095 使用hyperleger fabric实现的业务网络.png
03/13/2026  11:16 AM            16,447 农产品中心化溯源模式.png
03/13/2026  11:16 AM           106,564 农产品供应链溯源流程图.png
03/13/2026  11:16 AM           107,488 力软JAVA快速开发平台技术栈.png
03/13/2026  11:16 AM           119,845 区块链浏览器界面.png
03/13/2026  11:16 AM            70,274 区块链结构图.png
03/13/2026  11:16 AM           137,922 国密上链时序图.png
03/13/2026  11:16 AM            76,689 图100 溯源查询结果页小程序.png
03/13/2026  11:16 AM            29,503 图101 溯源查询首页小程序.png
03/13/2026  11:16 AM           136,148 图102 普刊录用通知.png
03/13/2026  11:16 AM         1,465,979 图102软著.png
03/13/2026  11:16 AM           208,036 图103_启动区块链网络.png
03/13/2026  11:16 AM           158,480 图4-2  BCCSP 国密算法接口实现示意.png
03/13/2026  11:16 AM            49,743 图4-3  fabric-CA 组件下 util 文件新添方法示意.png
03/13/2026  11:16 AM            53,009 图5-1  溯源系统用例图.jpg
03/13/2026  11:16 AM            17,382 图5-2  系统功能模块.png
03/13/2026  11:16 AM            46,455 基于区块链的农产品溯源模型.png
03/13/2026  02:28 PM            18,434 塔里木大学logo1.png
03/13/2026  02:28 PM            21,261 塔里木大学logo2.png
03/13/2026  11:16 AM           378,751 密码算法时间开销对比.png
03/13/2026  11:16 AM           122,104 延迟并发曲线.png
03/13/2026  11:16 AM            95,297 批次管理界面.png
03/13/2026  11:16 AM            94,661 技术路线图.png
03/13/2026  11:16 AM           158,197 数据大屏.png
03/13/2026  11:16 AM            25,989 新增加工记录表单.png
03/13/2026  11:16 AM            32,439 新增物流记录.png
03/13/2026  11:16 AM            29,615 新增种植记录表单.png
03/13/2026  11:16 AM           146,855 溯源业务流程图.png
03/13/2026  11:16 AM           289,663 溯源查询结果.png
03/13/2026  11:16 AM           174,090 溯源系统用例图.png
03/13/2026  11:16 AM            97,008 溯源记录界面.png
03/13/2026  11:16 AM            82,366 物流追踪.png
03/13/2026  11:16 AM           138,352 系统整体架构图.png
03/13/2026  11:16 AM            67,472 系统注册界面.png
03/13/2026  11:16 AM            52,964 系统登录界面.png
03/13/2026  11:16 AM            33,053 质检记录表单.png
03/13/2026  11:16 AM            80,367 链上链下的存储模型.png
              42 File(s)      5,640,136 bytes
               2 Dir(s)  340,920,016,896 bytes free

C:\02media\doc\ppt_makingv2\pic>

参考其他的markdown文档，例如普刊-安徽农学通报-可信农产品供应链的区块链构建路径研究.md 硕士学位论文v3.md
和
参考 基于区块链技术的农产品溯源关键技术研究成品ppt参考.pptx 类型是偏绿色，农学，需要有logo 塔里木大学logo，
03/13/2026  02:28 PM            18,434 塔里木大学logo1.png
03/13/2026  02:28 PM            21,261 塔里木大学logo2.png

硕士学位论文v3.md是最新的

其他提示：
2．研究生进行汇报，汇报时间在8分钟以内；
3．考核小组成员提问，研究生回答问题。



请你给出一份ppt，农学，20页内，汇报时间尽量在8分钟以内，不超过10分钟，用于预答辩用。另外给出一份逐字稿。 
需要提及：GBT 29373-201，还有图片可能需要子图形式。
请你给出对应的ppt和逐字稿。




请你先给出markdown形式的 ppt模板 草稿 预答辩ppt模板，图片用相对路径
给出预答辩ppt 草稿 markdown形式。 逐字稿，还有五位老师可能提问问题，以及回答 都是markdown形式。



```

```
ppt制作
相关资料我都上传到git了
https://github.com/sugar2026dev/ppt_making.git
这个项目是我之前的中期markdown笔记和ppt，
基于区块链技术的农产品溯源关键技术研究——塔里木大学研究生中期.md
柏小康 基于区块链技术的农产品溯源关键技术研究成品ppt参考.pptx
目前我需要做一次学术报告，例如插入部分pic下的图片 尽量做到图文并茂，参考之前的中期模板
Microsoft Corporation. All rights reserved.
C:\02media\doc\ppt_making\pic>dir
 Volume in drive C is Win11Pro X64
 Volume Serial Number is 20BE-1486
 Directory of C:\02media\doc\ppt_making\pic
03/04/2026  11:16 PM    <DIR>          .
03/04/2026  11:23 PM    <DIR>          ..
09/15/2025  09:56 AM            66,662 BCCSP上的SM2算法接口有用性和有效性验证结果.png
11/07/2024  11:12 AM            73,490 Fabric平台国密算法嵌入设计思路.png
11/17/2025  01:21 PM            73,095 使用hyperleger fabric实现的业务网络.png
12/09/2025  03:19 AM            16,447 农产品中心化溯源模式.png
11/08/2024  05:28 PM           106,564 农产品供应链溯源流程图.png
11/08/2024  02:03 AM           107,488 力软JAVA快速开发平台技术栈.png
11/28/2025  05:04 PM           117,147 区块链浏览器界面.png
12/10/2025  02:56 AM            76,689 图100 溯源查询结果页小程序.png
12/10/2025  02:57 AM            29,503 图101 溯源查询首页小程序.png
12/10/2025  03:05 AM           136,148 图102 普刊录用通知.png
12/11/2025  04:06 PM         1,465,979 图102软著.png
12/10/2025  03:12 AM           208,036 图103_启动区块链网络.png
12/09/2025  03:22 AM           158,480 图4-2  BCCSP 国密算法接口实现示意.png
12/09/2025  03:22 AM            49,743 图4-3  fabric-CA 组件下 util 文件新添方法示意.png
12/09/2025  03:24 AM            53,009 图5-1  溯源系统用例图.jpg
12/09/2025  03:24 AM            17,382 图5-2  系统功能模块.png
12/09/2025  03:25 AM            86,846 图5-3  系统登录页面.png
12/09/2025  03:37 AM            75,871 图5-4  溯源系统工作界面.png
09/15/2025  09:56 AM            46,455 基于区块链的农产品溯源模型.png
03/04/2026  11:16 PM            18,434 塔里木大学logo1.png
03/04/2026  11:16 PM            21,261 塔里木大学logo2.png
12/09/2025  02:53 AM           378,751 密码算法时间开销对比.png
11/08/2024  05:21 PM            94,661 技术路线图.png
              23 File(s)      3,478,141 bytes
               2 Dir(s)  362,484,707,328 bytes free
参考其他的markdown文档，例如普刊-安徽农学通报-可信农产品供应链的区块链构建路径研究.md 硕士学位论文.md
和草稿学术报告 ppt.md
参考 基于区块链技术的农产品溯源关键技术研究成品ppt参考.pptx 类型是偏绿色，农学，需要有logo 塔里木大学logo，请你给出一份ppt，农学，20页内，时间尽量控制在5-10分钟，用于学术报告用。另外给出一份逐字稿。 

```



```
我既需要原来的 使用Hyperledger Fabric实现的业务网络，又需要新增的Hyperledger Fabric 联盟链底层架构 这种可以并存吗 尽量图片与图片隔开，不要连续图片。 


图xx展示了在执行核心交易流转时的国密底层签名与背书通信时序过程。

图5-2 溯源系统Web端身份认证与注册界面



```



```
这个项目
https://github.com/ponyletter/gradProjBlockChainAgri_md
硕士学位论文v2.md 是我目前的论文。
目前我的论文中新增了图片，在pic新增目录下，

C:\02media\tarim\grap_paper\grad-proj-block-chain-agri_md\pic新增>dir
 Volume in drive C is Win11Pro X64
 Volume Serial Number is 20BE-1486

 Directory of C:\02media\tarim\grap_paper\grad-proj-block-chain-agri_md\pic新增

03/12/2026  08:19 PM    <DIR>          .
03/12/2026  08:23 PM    <DIR>          ..
03/11/2026  09:29 PM           247,805 Hyperledger Fabric 联盟链结构.png
03/11/2026  09:26 PM            70,274 区块链结构图.png
03/11/2026  04:58 PM           114,545 国密上链时序图.png
03/11/2026  04:58 PM           209,867 溯源业务流程图.png
03/12/2026  10:11 AM           295,433 溯源系统用例图.png
03/11/2026  04:58 PM           224,367 系统整体架构图.png
03/11/2026  09:28 PM            80,367 链上链下的存储模型.png
               7 File(s)      1,242,658 bytes
               2 Dir(s)  340,190,355,456 bytes free

C:\02media\tarim\grap_paper\grad-proj-block-chain-agri_md\pic新增>


还有在webpic目录下，有最新的web系统展示图片
C:\02media\tarim\grap_paper\grad-proj-block-chain-agri_md\webpic>dir
 Volume in drive C is Win11Pro X64
 Volume Serial Number is 20BE-1486

 Directory of C:\02media\tarim\grap_paper\grad-proj-block-chain-agri_md\webpic

03/12/2026  03:23 PM    <DIR>          .
03/12/2026  08:23 PM    <DIR>          ..
03/12/2026  03:14 PM           112,359 区块链浏览器界面.png
03/12/2026  03:07 PM            94,540 批次管理界面.png
03/12/2026  03:06 PM           143,194 数据大屏.png
03/12/2026  03:13 PM            25,989 新增加工记录表单.png
03/12/2026  03:14 PM            32,439 新增物流记录.png
03/12/2026  03:12 PM            29,615 新增种植记录表单.png
03/12/2026  03:05 PM           245,651 消费者溯源查询.png
03/12/2026  03:05 PM           414,364 溯源查询结果.png
03/12/2026  03:12 PM            96,089 溯源记录界面.png
03/12/2026  03:11 PM            80,816 物流追踪.png
03/12/2026  03:05 PM            35,218 电子合格证.png
03/12/2026  03:04 PM           705,710 系统注册界面.png
03/12/2026  03:04 PM           668,019 系统登录界面.png
03/12/2026  03:06 PM           660,008 系统登录界面_填写密码.png
03/12/2026  03:13 PM            33,053 质检记录表单.png
              15 File(s)      3,377,064 bytes
               2 Dir(s)  340,190,429,184 bytes free

C:\02media\tarim\grap_paper\grad-proj-block-chain-agri_md\webpic>
需要在  第五章体现出来
例如系统登录页面
溯源系统工作界面 a b c d等这种子图，webpic里面很多图片，只需要选一部分即可。 
应该多一些子图描述，这样不会太占格子。
目前pic新增需要插入到对应界面中，然后webpic需要选部分插入到第五章中。 
请你根据新增的图片，参考硕士学位论文v2.md ，给出硕士学位论文v3.md 目前只需要插入图片或者根据实际情况，替换图片。
请你给出markdown文档。 图片用相对路径。 


```

```
都用黑白，不用彩色
图5-2 系统整体架构图.png 标题需要上移动，然后部分文字需要整体左边移动。 
内容分别是：
接入层、应用层、数据存储层、区块链底层、基础设施层
区块链底层 后面的英文不需要。 
图5-4 溯源业务流程图.png
不需要1 2 3 4 5 6【】 也不需要
图5-5 标题上移动
```



## 技术参考

- [Hyperledger Fabric v2.2 官方文档](https://hyperledger-fabric.readthedocs.io/en/release-2.2/)
- [fabric-gm 国密改造版](https://github.com/ponyletter/fabric-gm)
- [tjfoc/gmsm 国密算法库](https://github.com/tjfoc/gmsm)
- [GB/T 29373-2012 农产品追溯要求 果蔬](https://std.samr.gov.cn/)
- [fabric-eCert-trace-ipfs_ok 参考项目](https://github.com/ponyletter/fabric-eCert-trace-ipfs_ok)

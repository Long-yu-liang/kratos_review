# 电商评价微服务系统

## 项目简介
基于 Kratos 框架构建的电商评价微服务系统，覆盖评论创建、回复、审核、申诉、检索全流程，提供完整的商品评价解决方案。

## 技术架构

### 核心技术栈
- **微服务框架**: Kratos + Protobuf
- **服务发现与注册**: Consul
- **数据库**: GORM + MySQL
- **缓存**: Redis
- **搜索**: Elasticsearch
- **消息队列**: Kafka
- **数据同步**: Canal
- **架构模式**: CQRS

### 架构特点
- 采用 CQRS 架构分离读写操作，提升系统性能
- 通过 Protobuf 定义统一接口规范，保证接口一致性
- 微服务架构设计，支持水平扩展和高可用

## 功能特性
- ✅ 评论创建与管理
- ✅ 评论回复功能
- ✅ 审核流程管理
- ✅ 用户申诉处理
- ✅ 智能搜索检索
- ✅ 数据实时同步

## 快速开始

### 环境配置
1. 修改 `config.yaml` 中的配置项：
```yaml
# 数据库配置
database:
  host: "your_mysql_host"
  port: 3306
  username: "your_username"
  password: "your_password"
  dbname: "review_service"

# Redis配置
redis:
  addr: "your_redis_host:6379"
  password: "your_redis_password"

# 其他服务配置...
registry.yaml 服务注册配置
```

### 服务启动
```bash
# 启动评价服务
go run cmd/review-service/main.go

# 或使用 Docker
docker-compose up -d
```

## API 文档

### 在线访问
API Fox 文档: [3ps7nm81x9.apifox.cn](https://3ps7nm81x9.apifox.cn)

### 本地查看
```bash
# 查看 API 规范
cat review-service/openapi.yaml

# 或启动本地 API 文档服务
swagger serve review-service/openapi.yaml
```

## 核心服务

### 服务模块
- review-service: 核心评价服务
- review-b: 业务处理服务
- review-o: 运营管理服务
- review-job: 后台作业服务 用于读取Kafka 写入Elasticsearch

## 开发说明

### 依赖安装
```bash
go mod tidy
```

### 代码生成
```bash
# 生成 Protobuf 代码
make proto

# 生成 Wire 依赖注入
在cmd/文件名/的终端中输出 wire
```

## 部署说明

### 生产环境配置
确保在 `config.yaml` 中配置：
- 正确的数据库连接信息
- Redis 连接配置
- Consul 服务发现地址
- Kafka 消息队列配置
- Elasticsearch 搜索服务配置

### 监控与日志
- 集成 Kratos 框架的日志系统
- 支持 Prometheus 指标收集
- 链路追踪支持

## 注意事项
- 请根据实际环境修改 `config.yaml` 配置文件
- 确保所有依赖服务（MySQL、Redis、Consul等）正常运行
- 生产环境建议配置适当的监控和告警

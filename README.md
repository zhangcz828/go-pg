# Introduction
简易版*金庸群侠传*， 一款对战游戏的Server端。

游戏共设置有两关， 每一关一个Boss. 玩家首次需选定一个Hero。

## Depency package

1. https://github.com/lib/pq
```text
A pure Go postgres driver for Go's database/sql package
```
2. https://github.com/gin-gonic/gin
```text
一个用Go写的web框架
```

## Rest API 介绍

### API 定义
```go
// User rest api
r.GET("/heros", postgres.GetHerosHandler(&store))
r.GET("/session/:id/", postgres.GetSessionHandler(&store)) //获取session
r.PUT("/session/:id", postgres.SelectHeroHandler(&store)) //选择英雄,更新session, PUT /session/:id?hero="张无忌"
r.PUT("/session/:id/fight", postgres.FightHandler(&store))  //打boss, 更新session里面的信息，得分
r.POST("/session/:id/archive", postgres.ArchiveSessionHandler(&store)) // ssMap[:id]存档到db
r.POST("/session/:id/level", postgres.NextLevelHandler(&store)) // 过关，更新session[加血，currentLevel, boss信息]
r.POST("/session/:id/quit", postgres.QuitHandler(&store)) // 玩家下线，自动存档
r.GET("/session/:id/rank", postgres.OnlineRankingHandler) //在线玩家积分排名 ```
```

```go
// Admin rest api
r.POST("/admin/hero", postgres.CreateHeroHandler(&store)) //添加新的hero
r.PUT("/admin/hero", postgres.AdjustHeroHandler(&store)) //节假日对hero进行调整，比如每个战斗力提升20%
```

## 用到的几个关键点
#### 1. Cache
Server实现了cache功能，基于package *github.com/hashicorp/golang-lru*
   > This provides the lru package which implements a fixed-size thread safe LRU cache. It is based on the cache in Groupcache.
   > 
此包实现了一个线程安全的map, 思考： 为什么map不能直接拿来实现cache?

#### 2. Go 单例模式
具体见Cache/cache.go， 思考为什么用Once? 加锁有什么问题？
Reference： https://mp.weixin.qq.com/s/37gV23UVHRA5SYeMEA5Q9w

#### 3. 怎样让程序变的可测？
Mock数据库的各种操作，把所有操作封装到interface， 这样就能在测试中定义一个mockDB

测试的几个注意点：

a) 不要测试你的依赖， mock之

b) 善用interface

c) 测试的套路， table driven
 
## 项目的组织形式

1. cache用于存放LRU cache相关实现和操作
2. modules是根据数据库的table创建的相关struct
3. pkg/connection 里面存放数据库连接的相关操作
4. pkg/postgresql 里面存放的是各种handler/middleware
5. resources 里面存放用到的SQL和一些json

## 项目的测试
1. pkg/postgres/dbstore.go， 定义了interface
```go
type DbStoreInterface interface {
	GetAllHeros() modules.Heros
	CreateHero(h modules.Hero) response
	LoadSessionFromDb(sid string) (*modules.SessionView, error)
	LoadBossFromDB(level int) (modules.Boss, error)
	LoadHeroFromDB(heroName string) (modules.Hero, error)
	Archive(session modules.Session) error
	RemoveSessionFromDB(sid string) error
	AdjustHero() error
}

2. pkg/postgres/mockdb.go, 实现了这个interface
3. Todo： 需要一个server本身的interface， 实现cache和session相关操作，让API变的更可测，测试覆盖率更高。
```
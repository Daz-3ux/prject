/*
负责与外部交互,控制缓存存储和获取的主流程

													是
接收 key --> 检查是否被缓存 -----> 返回缓存值 ⑴

	|  否                         是
	|-----> 是否应当从远程节点获取 -----> 与远程节点交互 --> 返回缓存值 ⑵
	            |  否
	            |-----> 调用`回调函数`，获取值并添加到缓存 --> 返回缓存值 ⑶

geecache/
    |--lru/
        |--lru.go  // lru 缓存淘汰策略
    |--byteview.go // 缓存值的抽象与封装
    |--cache.go    // 并发控制
    |--geecache.go // 负责与外部交互，控制缓存存储和获取的主流程
*/
package geecache

import (
	"fmt"
	"log"
	"sync"
)

/*
	回调 Getter:
	当缓存不存在,用户调用回调函数,得到源数据
	可以支持多种数据源配置,不用一种写一个
*/
// Getter loads data for a key
type Getter interface {
	Get(key string) ([]byte, error)
}

// GetterFunc implements Getter with a function
// 接口型函数: 函数适配器,用于将一个普通的函数转换为 Getter 接口的实现
// 可以将任何符合 func(key string) ([]byte, error) 签名的函数转换为 Getter 接口的实现
type GetterFunc func(key string) ([]byte, error)

// Get implements Getter interface function
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// Group 定义

// A Group is a cache namespace and associated data loaded spread over
/*
		一个 Group 可以看作一个缓存的命名空间
			每个 Group 拥有唯一的名称 name
			例如 [学生成绩]scores [学生信息]info [学生课程]courses
		getter 为缓存未命中时的回调
		mainCache 为并发缓存
*/
type Group struct {
	name				string
	getter			Getter
	mainCache		cache
}

var (
	mu			 sync.RWMutex
	groups = make(map[string]*Group)
)

// NewGroup create a new instance of Group
// 在调用 NewFroup 时使用 GetterFunc 实例化 Getter接口,赋给此处的 getter
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name: 			name,
		getter: 		getter,
		mainCache: 	cache{cacheBytes: cacheBytes},
	}
	// 维护 gourp 哈希表
	groups[name] = g
	return g
}

// GetGroup returns the named group previously created with NewGroup, or
// nil if thers's no such group
func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

// Get value for a key from cache -- key method for GeeCache
/*
		Get 方法实现了 (1) 与 (3)
			(1) 从 mainCache 中查找缓存,如果存在则返回缓存值
			(3) 缓存不存在,调用 load() 方法
*/
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	// 此处的 get 为从 缓存 中获取(缓存存在)
	if v, ok := g.mainCache.get(key); ok {
		log.Println("[GeeCache] hit")
		return v, nil
	}

	// 此处的 get 为实现的接口型函数,无惧类型变化均可获取源数据(缓存不存在)
	// load --> getLocally(分布式对应 getFromPeer) --> g.getter.Get()获取源数据,并将其加入到 mainCache
	return g.load(key)
}

func (g *Group) load(key string) (ByteView, error) {
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	fmt.Println("for test bytes = ", bytes)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bytes)}
	// 加入至缓存
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
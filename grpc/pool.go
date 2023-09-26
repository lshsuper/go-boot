package grpc

import (
	"errors"
	"sync"
	"time"
)

func (h connMap) push(key string, x interface{}) {

	h[key] = x.(*grpcClient)

}
func (h connMap) remove(key string) {

	delete(h, key)

}

func (h connMap) weightRand() *grpcClient {

	var (
		maxKey            string
		maxUseable, index int
	)

	for _, client := range h {
		index++
		l := len(client.usable)
		if index == 1 {
			maxKey = client.tag
			maxUseable = l
			continue
		}

		if l > maxUseable {

			maxUseable = l
			maxKey = client.tag

		}
	}
	if maxUseable == 0 {
		return nil
		//log.Println("引用达到上限:", maxIndex)
	}

	//获取最大索引除值
	return h[maxKey]

}

type connMap map[string]*grpcClient

func (h connMap) len() int { return len(h) }

type grpcPool struct {
	connMap           connMap
	rw                *sync.RWMutex
	addr              string
	maxIdle           int
	maxActive         int
	maxRef            int
	recyTime          time.Duration
	decrTime          time.Duration
	unaryInterceptors []UnaryClientInterceptor
}

type GrpcPoolConf struct {
	Addr              string
	MaxIdle           int
	MaxActive         int
	MaxRef            int
	RecyTime          time.Duration
	DecrTime          time.Duration
	UnaryInterceptors []UnaryClientInterceptor
}

func NewGrpcPool(conf GrpcPoolConf) *grpcPool {

	if conf.MaxActive <= 0 {
		conf.MaxActive = MaxActive
	}

	if conf.MaxRef <= 0 {
		conf.MaxRef = MaxRef
	}

	if conf.MaxIdle <= 0 {
		conf.MaxIdle = MaxIdle
	}

	if conf.RecyTime == 0 {
		conf.RecyTime = RecyTime
	}

	if conf.DecrTime == 0 {
		conf.DecrTime = DecrTime
	}

	pool := &grpcPool{
		rw:                new(sync.RWMutex),
		addr:              conf.Addr,
		maxActive:         conf.MaxActive,
		maxIdle:           conf.MaxIdle,
		maxRef:            conf.MaxRef,
		connMap:           make(connMap),
		decrTime:          conf.DecrTime,
		recyTime:          conf.RecyTime,
		unaryInterceptors: conf.UnaryInterceptors,
	}

	//初始化链接
	for i := 0; i < conf.MaxIdle; i++ {
		//conf.Addr, pool, conf.MaxRef, i
		conn, err := NewGrpcClient(GrpcClientConf{
			Addr:              conf.Addr,
			Pool:              pool,
			Index:             i,
			MaxRef:            pool.maxRef,
			UnaryInterceptors: conf.UnaryInterceptors,
		})
		if err != nil {
			panic("connMap init failed")
		}
		//入队
		pool.connMap.push(conn.tag, conn)
	}

	pool.listenDecrCap()

	return pool
}

func (pool *grpcPool) Get() (*grpcClient, error) {

	for {

		pool.rw.RLock()
		//用最小引用堆来实现（每次获取最小被引用的连客户端连接）
		client := pool.connMap.weightRand()
		pool.rw.RUnlock()
		//消费一下

		if client == nil || client.use() != nil {
			//获取客户端失败
			conn, err := pool.incrCap()
			if err == nil {
				return conn, err
			}
			continue
		}

		return client, nil
	}

	return nil, errors.New("pool is bazy...")

}

func (pool *grpcPool) RWLock() *sync.RWMutex {
	return pool.rw
}

//incrCap 扩容
func (pool *grpcPool) incrCap() (*grpcClient, error) {

	//扩容
	pool.rw.Lock()
	l := pool.connMap.len()
	if l < pool.maxActive {
		client, err := NewGrpcClient(GrpcClientConf{
			Addr:              pool.addr,
			Pool:              pool,
			Index:             l,
			MaxRef:            pool.maxRef,
			UnaryInterceptors: pool.unaryInterceptors,
		})
		if err == nil {
			//扩容
			_ = client.use()
			pool.connMap.push(client.tag, client)
			pool.rw.Unlock()
			return client, err
		}
	}
	pool.rw.Unlock()
	return nil, errors.New("incr  cap  failed...")
}

//decrCap 缩容
func (pool *grpcPool) listenDecrCap() {

	go func() {
		timer := time.NewTimer(pool.decrTime)
		for {

			select {
			case _ = <-timer.C:

				if pool.connMap.len() <= pool.maxIdle {
					timer.Reset(time.Second * 10)
					continue
				}

				pool.rw.Lock()
				keys := make([]string, 0)
				for k, v := range pool.connMap {
					if v.canRecy() {
						keys = append(keys, k)
					}
				}

				for _, v := range keys {
					pool.connMap.remove(v)
				}

				timer.Reset(time.Second * 10)
				pool.rw.Unlock()

			}

		}
	}()

}

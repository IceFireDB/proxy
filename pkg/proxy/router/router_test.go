// Copyright 2014 Wandoujia Inc. All Rights Reserved.
// Licensed under the MIT (MIT-LICENSE.txt) license.

package router


// todo we need to get test back
//
//var (
//	conf       *Conf
//	s          *Server
//	once       sync.Once
//	waitonce   sync.Once
//	conn       zkhelper.Conn
//	redis1     *miniredis.Miniredis
//	redis2     *miniredis.Miniredis
//	proxyMutex sync.Mutex
//)
//
//func InitEnv() {
//	go once.Do(func() {
//		conn = zkhelper.NewConn() // this is a fake customized zk client
//		conf = &Conf{
//			proxyId:         "proxy_test",
//			productName:     "test",
//			coordinatorType: "zookeeper",
//			coordinatorAddr: "localhost:2181",
//			net_timeout:     5,
//			slot_num:        16,
//			// broker:      LedisBroker,
//		}
//
//		// init action path
//		prefix := models.GetWatchActionPath(conf.productName)
//		err := models.CreateActionRootPath(conn, prefix)
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		// init slot
//		err = models.InitSlotSet(conn, conf.productName, conf.slot_num)
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		// init  server group
//		g1 := models.NewServerGroup(conf.productName, 1)
//		g1.Create(conn)
//		g2 := models.NewServerGroup(conf.productName, 2)
//		g2.Create(conn)
//
//		redis1, _ := miniredis.Run()
//		redis2, _ := miniredis.Run()
//
//		s1 := models.NewServer(models.SERVER_TYPE_MASTER, redis1.Addr())
//		s2 := models.NewServer(models.SERVER_TYPE_MASTER, redis2.Addr())
//
//		g1.AddServer(conn, s1)
//		g2.AddServer(conn, s2)
//
//		// set slot range
//		err = models.SetSlotRange(conn, conf.productName, 0, conf.slot_num/2-1, 1, models.SLOT_STATUS_ONLINE)
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		err = models.SetSlotRange(conn, conf.productName, conf.slot_num/2, conf.slot_num-1, 2, models.SLOT_STATUS_ONLINE)
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		go func() { // set proxy online
//			time.Sleep(5 * time.Second)
//			err := models.SetProxyStatus(conn, conf.productName, conf.proxyId, models.PROXY_STATE_ONLINE)
//			if err != nil {
//				log.Fatal(errors.ErrorStack(err))
//			}
//			time.Sleep(2 * time.Second)
//			proxyMutex.Lock()
//			defer proxyMutex.Unlock()
//			pi := s.getProxyInfo()
//			if pi.State != models.PROXY_STATE_ONLINE {
//				log.Fatalf("should be online, we got %s", pi.State)
//			}
//		}()
//
//		proxyMutex.Lock()
//		s = NewServer(":19000", ":11000",
//			conf,
//		)
//		proxyMutex.Unlock()
//		s.Run()
//	})
//
//	waitonce.Do(func() {
//		time.Sleep(10 * time.Second)
//	})
//}
//
//func TestSingleKeyRedisCmd(t *testing.T) {
//	InitEnv()
//	c, err := redis.Dial("tcp", "localhost:19000")
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer c.Close()
//
//	_, err = c.Do("SET", "foo", "bar")
//	if err != nil {
//		t.Error(err)
//	}
//
//	if got, err := redis.String(c.Do("get", "foo")); err != nil || got != "bar" {
//		t.Error("'foo' has the wrong value")
//	}
//
//	_, err = c.Do("SET", "bar", "foo")
//	if err != nil {
//		t.Error(err)
//	}
//
//	if got, err := redis.String(c.Do("get", "bar")); err != nil || got != "foo" {
//		t.Error("'bar' has the wrong value")
//	}
//}
//
//func TestMultiKeyRedisCmd(t *testing.T) {
//	InitEnv()
//	c, err := redis.Dial("tcp", "localhost:19000")
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer c.Close()
//
//	_, err = c.Do("SET", "key1", "value1")
//	if err != nil {
//		t.Fatal(err)
//	}
//	_, err = c.Do("SET", "key2", "value2")
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	var value1 string
//	var value2 string
//	var value3 string
//	reply, err := redis.Values(c.Do("MGET", "key1", "key2", "key3"))
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	if _, err := redis.Scan(reply, &value1, &value2, &value3); err != nil {
//		t.Fatal(err)
//	}
//
//	if value1 != "value1" || value2 != "value2" || len(value3) != 0 {
//		t.Error("value not match")
//	}
//
//	// test del
//	if _, err := c.Do("del", "key1", "key2", "key3"); err != nil {
//		t.Fatal(err)
//	}
//
//	// reset
//	value1 = ""
//	value2 = ""
//	value3 = ""
//	reply, err = redis.Values(c.Do("MGET", "key1", "key2", "key3"))
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	if _, err := redis.Scan(reply, &value1, &value2, &value3); err != nil {
//		t.Fatal(err)
//	}
//
//	if len(value1) != 0 || len(value2) != 0 || len(value3) != 0 {
//		t.Error("value not match", value1, value2, value3)
//	}
//
//	// reset
//	value1 = ""
//	value2 = ""
//	value3 = ""
//
//	_, err = c.Do("MSET", "key1", "value1", "key2", "value2", "key3", "")
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	reply, err = redis.Values(c.Do("MGET", "key1", "key2", "key3"))
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	if _, err := redis.Scan(reply, &value1, &value2, &value3); err != nil {
//		t.Fatal(err)
//	}
//
//	if value1 != "value1" || value2 != "value2" || len(value3) != 0 {
//		t.Error("value not match", value1, value2, value3)
//	}
//}
//
//func TestInvalidRedisCmdUnknown(t *testing.T) {
//	InitEnv()
//	c, err := redis.Dial("tcp", "localhost:19000")
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer c.Close()
//
//	if _, err := c.Do("unknown", "key1", "key2", "key3"); err == nil {
//		t.Fatal(err)
//	}
//}
//
//func TestInvalidRedisCmdPing(t *testing.T) {
//	InitEnv()
//	c, err := redis.Dial("tcp", "localhost:19000")
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer c.Close()
//
//	_, err = c.Do("info")
//	if err != io.EOF {
//		t.Fatal(err)
//	}
//}
//
//func TestInvalidRedisCmdQuit(t *testing.T) {
//	InitEnv()
//	c, err := redis.Dial("tcp", "localhost:19000")
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer c.Close()
//
//	_, err = c.Do("quit")
//	if err != nil {
//		t.Fatal(err)
//	}
//}
//
//func TestInvalidRedisCmdEcho(t *testing.T) {
//	InitEnv()
//	c, err := redis.Dial("tcp", "localhost:19000")
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer c.Close()
//
//	_, err = c.Do("echo", "xx")
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	_, err = c.Do("echo")
//	if err != nil {
//		t.Fatal(err)
//	}
//}
//
//// this should be the last test
//func TestMarkOffline(t *testing.T) {
//	InitEnv()
//
//	suicide := int64(0)
//	proxyMutex.Lock()
//	s.OnSuicide = func() error {
//		atomic.StoreInt64(&suicide, 1)
//		return nil
//	}
//	proxyMutex.Unlock()
//
//	err := models.SetProxyStatus(conn, conf.productName, conf.proxyId, models.PROXY_STATE_MARK_OFFLINE)
//	if err != nil {
//		t.Fatal(errors.ErrorStack(err))
//	}
//
//	time.Sleep(3 * time.Second)
//
//	if atomic.LoadInt64(&suicide) == 0 {
//		t.Error("shoud be suicided")
//	}
//}

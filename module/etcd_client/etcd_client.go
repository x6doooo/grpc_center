package etcd_client

import (
    "github.com/coreos/etcd/clientv3"
    "context"
    "time"
    "errors"
    "sync"
)

type Client struct {
    mu *sync.Mutex
    cfg clientv3.Config
    BaseCli *clientv3.Client
}

const (
    defaultRequestTimeout = 30 * time.Second
)

func NewClient(cfg clientv3.Config) (cli *Client, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = r.(error)
            return
        }
    }()
    cli = &Client{
        mu: &sync.Mutex{},
        cfg: cfg,
    }
    cli.BaseCli, err = clientv3.New(cfg)
    if err != nil {
        panic(err)
    }
    return
}

func (me *Client) Put(key, value string, opts... clientv3.OpOption) (resp *clientv3.PutResponse, err error) {
    ctx, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
    resp, err = me.BaseCli.Put(ctx, key, value, opts...)
    cancel()
    return
}


func (me *Client) Get(key string, opts... clientv3.OpOption) (resp *clientv3.GetResponse, err error) {
    ctx, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
    resp, err = me.BaseCli.Get(ctx, key, opts...)
    cancel()
    return
}

func (me *Client) TxnPut(key, value string) (resp *clientv3.TxnResponse, err error) {
    me.mu.Lock()
    ctx, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
    resp, err = me.BaseCli.Txn(ctx).
        If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).
        Then(clientv3.OpPut(key, value)).Commit()
    cancel()
    me.mu.Unlock()
    return
}

type updateHandlerFunc func(string)(string, error)
func (me *Client) TxnUpdate(key string, updateHandler updateHandlerFunc) (err error) {
    lockPrefixKey := "the_lock_servce/" + key

    var count = 0
    for {
        count++
        err = me.TxnLock(lockPrefixKey)
        if err != nil {
            if count > 5 {
                return errors.New("timeout")
            }
            time.Sleep(1 * time.Second)
        } else {
            break
        }
    }

    defer func() {
        me.TxnUnlock(lockPrefixKey)
        if e := recover(); e != nil {
            err = e.(error)
        }
    }()

    var getResponse *clientv3.GetResponse

    getResponse, err = me.Get(key)
    if err != nil {
        return
    }

    if getResponse.Count != 0 {
        kv := getResponse.Kvs[0]
        value := string(kv.Value)
        value, err = updateHandler(value)
        if err == nil {
            _, err = me.Put(key, value)
        }
    } else {
        err = errors.New("key is not exists")
    }
    return err
}

func (me *Client) TxnLock(key string) error {
    ctx, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
    resp, err := me.BaseCli.Txn(ctx).
        If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).
        Then(clientv3.OpPut(key, "1")).Commit()
    if err != nil {
        return err
    }
    cancel()
    if resp.Succeeded {
        return nil
    }
    return errors.New("the key has locked")
}

func (me *Client) TxnUnlock(key string) error {
    ctx, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
    _, err := me.BaseCli.Delete(ctx, key)
    if err != nil {
        return err
    }
    cancel()
    return nil
}



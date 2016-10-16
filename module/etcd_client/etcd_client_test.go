package etcd_client

import (
    "testing"
    "github.com/coreos/etcd/clientv3"
    "fmt"
    "context"
    "time"
)

func TestClient_Put(t *testing.T) {
    cfg := clientv3.Config{
        Endpoints: []string{
            "127.0.0.1:9537",
            "127.0.0.1:9538",
            "127.0.0.1:9539",
        },
        DialTimeout: defaultRequestTimeout,
    }
    cli, err := NewClient(cfg)
    if err != nil {
        panic(err)
    }

    time.AfterFunc(2 * time.Second, func() {
        r1, err := cli.Put("test/abc", "123")
        fmt.Println(r1, err)

        r2, err := cli.Put("test/abc1", "234")
        fmt.Println(r2, err)

        r3, err := cli.Get("test/", clientv3.WithPrefix())
        fmt.Println(r3, err)

    })


    rch := cli.baseCli.Watch(context.Background(), "test/", clientv3.WithPrefix())
    for wresp := range rch {
        for _, ev := range wresp.Events {
            fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
        }
    }

}

package mcc

import (
    mccPb "grpc_center/pbs/mcc"
    "golang.org/x/net/context"
    "grpc_center/module/etcd_client"
    "github.com/coreos/etcd/clientv3"
    "time"
    "grpc_center/module/auth"
    "strconv"
    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    "encoding/json"
    "grpc_center/module/klog"
    "strings"
    "errors"
)

const (
    defaultRequestTimeout = 10 * time.Second
    service_prefix = "service_grpc/"
    service_lock_prefix = "service_grpc_lock/"
)

type ServiceServer struct {
    Addr   string
    Weight int
    Status int
}

type ServiceInfo struct {
    Name    string
    Desc    string
    Token   string
    Servers []ServiceServer
}

type mcc struct {
    etcdClient *etcd_client.Client
}

// ---- public --------------------

func (me *mcc) Register(ctx context.Context, req *mccPb.RegisterRequest) (resp *mccPb.RegisterResponse, err error) {

    if req.Name == "" {
        err = grpc.Errorf(codes.Internal, "need username")
        return
    }

    token, err := me.register2etcd(req.Name, req.Desc)
    if err != nil {
        err = grpc.Errorf(codes.Internal, "")
        return
    }

    resp.AppToken = token
    return
}

func (me *mcc) Join(ctx context.Context, req *mccPb.JoinRequest) (resp *mccPb.JoinResponse, err error) {

    if req.AppToken == "" {
        err = grpc.Errorf(codes.Internal, "appToken error")
        return
    }

    if req.Addr == "" {
        err = grpc.Errorf(codes.Internal, "server addr error")
        return
    }


    resp.Success, err = me.joinHandler(req.AppToken, req.Addr, req.Weight)
    return
}

func (me *mcc) HeartBeat(ctx context.Context, req *mccPb.HeartBeatRequest) (resp *mccPb.HeartBeatResponse, err error) {
    return
}

func (me *mcc) Lookup(ctx context.Context, req *mccPb.LookupRequest) (resp *mccPb.LookupResponse, err error) {
    return
}


// ---- private --------------------

func (me *mcc) register2etcd(name, desc string) (token string, err error) {

    name = service_prefix + name

    t := time.Now().Unix()
    ts := strconv.Itoa(int(t))

    tokenBase := []byte(name + "|" + ts)
    tokenSrc := auth.TokenService.Encrypt(tokenBase)
    token = string(tokenSrc)

    info := &ServiceInfo{
        Name: name,
        Desc: desc,
        Token: token,
    }

    infoByte, err := json.Marshal(info)
    if err != nil {
        return
    }

    infoStr := string(infoByte)
    _, err = me.etcdClient.TxnPut(name, infoStr)

    if err != nil {
        return
    }
    return
}


func(me *mcc) joinHandler(appToken, addr string, weight int32) (success bool, err error) {
    //me.etcdClient.TxnUpdate()
    tokenSrcByte := auth.TokenService.Decrypt([]byte(appToken))
    tokenSrc := string(tokenSrcByte)
    tokenSrcArr := strings.Split(tokenSrc, "|")
    if len(tokenSrcArr) != 2 || tokenSrcArr[0] == "" {
        err = errors.New("appToken error")
        return
    }
    appNameHasPrefix := tokenSrcArr[0]

    updateHandler := func(value string) (updateValue string, err error) {
        serviceInfo := ServiceInfo{}
        err = json.Unmarshal([]byte(value), &serviceInfo)
        if err != nil {
            return
        }
        // todo: update
    }

    me.etcdClient.TxnUpdate(appNameHasPrefix, updateHandler())
    //name := string()
    return
}


// ---- constructor --------------------

func NewMcc(endpoints []string) *mcc {

    klog.Logger.Info("------------ new mcc server")
    cfg := clientv3.Config{
        Endpoints: endpoints,
        DialTimeout: defaultRequestTimeout,
    }
    klog.Logger.Info("------------ init etcd client")
    cli, err := etcd_client.NewClient(cfg)
    if err != nil {
        panic(err)
    }
    return &mcc{
        etcdClient: cli,
    };
}


// others

func NewError(code codes.Code) error {
    return grpc.Errorf(code, code.String())
}

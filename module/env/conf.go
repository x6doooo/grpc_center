package env

type EnvConf struct {
    Mode         string
    Server_addr  string
    Log_file     string
    Log_max_line int
    Log_backups  int
}

type TokenConf struct {
    Key  string
    Iv   string
    Salt string
}

type EtcdConf struct {
    Addr string
}

type MainConf struct {
    Env  EnvConf
    Token TokenConf
    Etcd []EtcdConf
}




package main

import (
    "bufio"
    "fmt"
    "os"
    "flag"
    "net"
    "sync"
    "runtime"
    . "github.com/hq-cml/reversi"
    "github.com/hq-cml/reversi/ai"
    "github.com/hq-cml/reversi/client/helper"
)

var ip *string = flag.String("h", "127.0.0.1", "ip")
var port *int = flag.Int("p", 9527, "port")
var mode *int = flag.Int("a", 0, "AI")                  //0-手动模式 1-AI自动模式

func checkError(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
        os.Exit(1)
    }
}

//AI自动处理
func ai_auto(conn net.Conn) {
    cmd := "Nhq" //姓名hq
    role := BLACK

    //上报姓名
    _, e := conn.Write([]byte(cmd))
    checkError(e)

    //循环处理
    for {
        //网络
        buf := make([]byte, 128)
        n, e := conn.Read(buf)
        checkError(e)

        if string(buf[0:2]) == "U1" {
            fmt.Println("AI：黑子")
            role = BLACK
        }else if string(buf[0:2]) == "U0" {
            fmt.Println("AI：白子")
            role = WIITE
        }

        do, chess := helper.HandlNetData(buf, n)
        if do {
            //handler_net_data返回true，说明本方需要落子，调用AI落子
            chessBoard := helper.ConverBytesToChessBoard([]byte(chess))
            //首先查看是否可落子
            step , canDown := ai.CheckChessboard(chessBoard, int8(role))
            if step == 0 {
                fmt.Println("目前没有位置可落子，按任意键请对方落子。")
                waitUser()
            } else {
                //Ai落子
                row, col := ai.AiPlayStep(chessBoard, canDown, int8(role))
                cmd := helper.ConvertRowColToServerProtocal(row, col)
                fmt.Printf("AI(%d)落子：%d,%d, cmd:%s\n", role, row, col, cmd)
                _, e := conn.Write([]byte(cmd))
                checkError(e)
            }
        }
        //time.Sleep(time.Second * 5)
    }
}

//接收网络输入
func net_in(conn net.Conn, user_input chan byte) {
    for {
        //网络
        buf := make([]byte, 128)
        n, e := conn.Read(buf)
        checkError(e)

        do, _ := helper.HandlNetData(buf, n)
        if do {
            //handler_net_data返回true，说明本方需要落子，通知显示提示符
            user_input <- byte('c')
        }

        //runtime.Gosched()
    }
    os.Exit(0)
}

//处理用户输入
func std_in(conn net.Conn, user_input chan byte) {
    r := bufio.NewReader(os.Stdin)
    for {
        <-user_input //等待接受用户输入的命令
        fmt.Print("Enter command-> ")
        rawLine, _, _ := r.ReadLine()

        line := string(rawLine)

        if line == "quit" {
            break
        }

        //将用户输入指令转化成Server指令
        cmd := helper.ConvertPlaceToServerProtocal(line)

        _, e := conn.Write([]byte(cmd))
        //fmt.Println("write len:", n)
        checkError(e)
        //runtime.Gosched()
    }
    os.Exit(0)
}

//等待用户按任意键，阻塞程序
func waitUser() {
    r := bufio.NewReader(os.Stdin)
    _, _, _ = r.ReadLine()
}

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())
    //解析参数
    flag.Parse()
    if ip != nil {
        fmt.Println("ip =", *ip, ", port =", *port, ", AI=", *mode)
    }

    //提示输出
    fmt.Println(`
Enter following commands to control:
Nyourname -- report your name. Nhq, eg.
Row&Col -- place pieces in (Row,Col). 3D, eg.
quit -- quit
`)
    address := fmt.Sprintf("%s:%d", *ip, *port)
    //建立连接
    conn, err := net.Dial("tcp", address)
    checkError(err)

    var wg sync.WaitGroup

    if *mode == 0 {
        //手动模式
        user_input := make(chan byte)
        go net_in(conn, user_input);
        go std_in(conn, user_input);
        wg.Add(2);
        //开始
        user_input <- byte('c')
    } else {
        //AI自动模式
        go ai_auto(conn);
        wg.Add(1);
    }

    wg.Wait();
}
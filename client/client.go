package main

import (
    "bufio"
    "fmt"
    "os"
    "flag"
    "net"
    "sync"
    "runtime"
    //"time"
)

var ip *string = flag.String("h", "127.0.0.1", "ip")
var port *int = flag.Int("p", 9527, "port")

func checkError(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
        os.Exit(1)
    }
}

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())
    //解析参数
    flag.Parse()
    if ip != nil {
        fmt.Println("ip =", *ip, ", port =", *port)
    }

    //提示输出
    fmt.Println(`
Enter following commands to control:
Nyourname -- report your name. hq, eg.
Mxy -- place pieces in (x,y)
quit -- quit
`)
    address := fmt.Sprintf("%s:%d", *ip, *port)
    //建立连接
    conn, err := net.Dial("tcp", address)
    checkError(err)


    user_input := make(chan byte)

    var wg sync.WaitGroup
    go net_in(conn, user_input);
    //wg.Add(1);
    go std_in(conn, user_input);
    wg.Add(2);

    //开始
    user_input <- byte('c')
    wg.Wait();
}

//接收网络输入
func net_in(conn net.Conn, user_input chan byte) {
    //ok := true
    for {
        //网络
        buf := make([]byte, 128)
        n, e := conn.Read(buf)
        checkError(e)

        if handler_net_data(buf, n) {
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

        _, e := conn.Write([]byte(line))
        //fmt.Println("write len:", n)
        checkError(e)
        //runtime.Gosched()
    }
    os.Exit(0)
}

//返回true表示轮到己方,展现cmd提示符
//返回false表示继续等待网络数据,cmd提示符收起
func handler_net_data(buf []byte, length int) bool{
    if length == 3 && string(buf[0:2]) == "U1" {
        fmt.Println("Got->",string(buf[0:length-1]), ". [ You are first(BLACK). ]")
        return false
    }
    if length == 3 && string(buf[0:2]) == "U0" {
        fmt.Println("Got->",string(buf[0:length-1]), ". [ You are second(WHITE). ]")
        return false
    }
    if length == 3 && string(buf[0:2]) == "W1" {
        fmt.Println("Got->",string(buf[0:length-1]), ". [ You win! ]")
        return false
    }
    if length == 3 && string(buf[0:2]) == "W0" {
        fmt.Println("Got->",string(buf[0:length-1]), ". [ You lose! ]")
        return false
    }
    if length == 3 && string(buf[0:2]) == "W2" {
        fmt.Println("Got->",string(buf[0:length-1]), ". [ Draw tie! ]")
        return false
    }
    if length == 2 && string(buf[0:1]) == "G" {
        fmt.Println("Got->",string(buf[0:length-1]), ". [ Game over! ]")
        return false
    }
    if length == 66 {
        fmt.Println("Got->",string(buf[0:length]))
        print_board(buf[1:length-1])
        return true
    }
    if length == 69 { //这种情况是U1和棋盘放在一个TCP包中发过来了
        fmt.Println("Got->",string(buf[0:length]))
        print_board(buf[4:length-1])
        return true
    }
    fmt.Println("Got->",string(buf[0:length]), "len:", length)
    fmt.Println("FUCK!")
    os.Exit(1)
    return false
}

//打印棋盘
func print_board(buf []byte) {
    if len(buf) != 64 {
        fmt.Println("Error board!")
        os.Exit(1)
    }
    fmt.Println("  -------------------");
    for y := 7; y >= 0 ; y-- {
        fmt.Printf("%d |", y)
        for x :=0; x <= 7; x++ {
            idx := 8*x + y
            c := string(buf[idx])
            if(c == "1"){
                fmt.Print(" W")
            }else if(c == "2"){
                fmt.Print(" B")
            }else{
                fmt.Print(" *")
            }
        }
        fmt.Print(" |\n")
    }
    fmt.Println("  -------------------");
    fmt.Println("    0 1 2 3 4 5 6 7\n");
}

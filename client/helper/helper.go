package helper

import (
    "fmt"
    "os"
    . "github.com/hq-cml/reversi"
)

/*
 * 处理服务端返回数据
 *
 * 返回：
 *   true  -- 表示轮到己方,展现cmd提示符或者AI自动分析落子
 *   false -- 表示继续等待网络数据
 * TODO：
 *   这里处理不太地道，按理说应该按协议约定的分隔符处理
 */
func HandlNetData(buf []byte, length int)  (bool, string){
    if length == 3 && string(buf[0:2]) == "U1" {
        fmt.Println("Got->",string(buf[0:length-1]), ". [ You are first(BLACK). ]")
        return false, string(buf[0:length-1])
    }
    if length == 3 && string(buf[0:2]) == "U0" {
        fmt.Println("Got->",string(buf[0:length-1]), ". [ You are second(WHITE). ]")
        return false, string(buf[0:length-1])
    }
    if length == 3 && string(buf[0:2]) == "W1" {
        fmt.Println("Got->",string(buf[0:length-1]), ". [ You win! ]")
        return false, string(buf[0:length-1])
    }
    if length == 3 && string(buf[0:2]) == "W0" {
        fmt.Println("Got->",string(buf[0:length-1]), ". [ You lose! ]")
        return false, string(buf[0:length-1])
    }
    if length == 3 && string(buf[0:2]) == "W2" {
        fmt.Println("Got->",string(buf[0:length-1]), ". [ Draw tie! ]")
        return false, string(buf[0:length-1])
    }
    if length == 2 && string(buf[0:1]) == "G" {
        fmt.Println("Got->",string(buf[0:length-1]), ". [ Game over! ]")
        os.Exit(0)
        return false, string(buf[0:length-1])
    }
    if length == 66 {
        fmt.Println("Got->",string(buf[0:length]))
        printBoard(buf[1:length-1])
        return true, string(buf[1:length-1])
    }
    if length == 69 { //这种情况是U1和棋盘放在一个TCP包中发过来了
        fmt.Println("Got->",string(buf[0:length]))
        printBoard(buf[4:length-1])
        return true, string(buf[4:length-1])
    }
    fmt.Println("Got->",string(buf[0:length]), "len:", length)
    fmt.Println("Fuck Something wrong!")
    os.Exit(1)
    return false, ""
}

//打印棋盘
func printBoard(buf []byte) {
    if len(buf) != 64 {
        fmt.Println("Error board!")
        os.Exit(1)
    }
    chessBoard := ConverBytesToChessBoard(buf)
    PrintChessboard(chessBoard)
}

/*
 * 将Java版本的棋盘字符串协议，转化成ChessBoard
 */
func ConverBytesToChessBoard(buf []byte) ChessBoard{
    var chessBoard ChessBoard
    var cnt int8

    for y := 7; y >= 0 ; y-- {
        for x :=0; x <= 7; x++ {
            idx := 8*x + y
            row := cnt/8
            col := cnt%8
            c := string(buf[idx])
            if(c == "1"){
                chessBoard[row][col] = WIITE
            }else if(c == "2"){
                chessBoard[row][col] = BLACK
            }else{
                chessBoard[row][col] = NULL
            }
            cnt++
        }
    }
    return chessBoard
}

/*
 * 将玩家落子行为，转化成Java版server的协议
 * 比如：2F => M56
 */
func ConvertPlaceToServerProtocal(line string) (cmd string){
    var ret [3]byte

    //如果N开头，是报家门，则直接返回
    if []byte(line)[0] == 'N' {
        return line
    }

    if len(line) != 2 {
        return ""
    }

    ret[0] = 'M'
    if line[1] >= byte('a') {
        ret[1] = byte(line[1]) - byte('a')
    } else {
        ret[1] = byte(line[1]) - byte('A')
    }
    ret[2] = LENGTH - (byte(line[0])-byte('0'))

    cmd = fmt.Sprintf("%c%d%d", ret[0], ret[1], ret[2])
    return cmd
}

/*
 * 将AI落子行为，转化成Java版server的协议
 * 比如：row,col(1,6) => Mxy(M66)
 *      row,col(0,7) => Mxy(M47)
 */
func ConvertRowColToServerProtocal(row, col int8) (cmd string){
    //row=>y
    y := 8 - row -1

    //col=>x
    x := col

    cmd = fmt.Sprintf("%c%d%d", 'M', x, y)
    return cmd
}
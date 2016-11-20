package main

import (
    "fmt"
    "bufio"
    "os"
    . "github.com/hq-cml/reversi"
    "github.com/hq-cml/reversi/ai"
)

//等待用户按任意键，用于阻塞程序
func waitUser() {
    r := bufio.NewReader(os.Stdin)
    _, _, _ = r.ReadLine()
}

//统计棋盘战场的得分
func countScoreFinally(chessboard ChessBoard, role int8) (user_score, opponent_score int8){
    var row, col int8

    //确定本方和对方颜色
    self_color, opponent_color := role, -1*role

    //遍历分析棋盘所有位置
    for row =0; row<LENGTH; row++ {
        for col = 0; col < LENGTH; col++ {
            if chessboard[row][col] == opponent_color {
                opponent_score ++
            }
            if chessboard[row][col] == self_color {
                user_score ++
            }
        }
    }
    return
}

//单机版人机对战
func main() {
    var row, col int8
    var chessboard ChessBoard //棋盘
    var cnt int8 //已落子的个数
    var skip_play int8 //当某一方无子可落时，增1，若为2，表示双方都不能落子
    var user_role int8 //玩家颜色：1白色 -1黑色
    var ai_role int8 //玩家颜色：1白色 -1黑色
    var turn int8 //当前轮到哪一方落子
    r := bufio.NewReader(os.Stdin)

    fmt.Println()
    fmt.Println("*************************************")
    fmt.Println("*                                   *")
    fmt.Println("*    ~黑白棋小程序--人机对战AI版~   *")
    fmt.Println("*            作者：HQ               *")
    fmt.Println("*                                   *")
    fmt.Println("*************************************")
    fmt.Println("\n 初始棋盘:")

    //初始化棋盘
    _ = InitChessboard(&chessboard)
    //打印棋盘
    PrintChessboard(chessboard)

    //人机交互
    fmt.Print("\n游戏者执黑先下，输入 0-黑子 1-白子，请选择：")
    rawLine, _, _ := r.ReadLine()
    line := string(rawLine)
    if line == "0" {
        fmt.Println("您执黑先行！")
        user_role = BLACK
    }else if line == "1" {
        fmt.Print("您执白后行！请按任意键，AI落子~")
        waitUser()
        user_role = WIITE
    }else{
        fmt.Println("输入非法，程序退出！")
        os.Exit(1)
    }

    ai_role = user_role * -1 //ai颜色赋值
    turn = BLACK             //黑子先行

    //无限循环，轮流落子，直到分出胜负
    for {
        //如果本轮是黑子（白子），且玩家执黑子（白子），则玩家落子
        if (turn == BLACK && user_role == BLACK) || (turn == WIITE && user_role == WIITE) {
            //玩家落子
            //首先查看是否可落子
            //fmt.Println(chessboard)
            //fmt.Println(user_role)
            step , canDown := ai.CheckChessboard(chessboard, user_role)
            if step == 0 {
                //无子可落
                skip_play ++
                if skip_play == 1 {
                    fmt.Println("你目前没有位置可落子，按回车键让对方下子。")
                    waitUser()
                }else if skip_play == 2 {
                    fmt.Println("双方均没有可落棋子")
                    break
                }
            } else {
                //玩家落子，无限循环等待玩家落下合法子
                for {
                    fmt.Print("\n输入落子的位置(行 列):");
                    line, _, _ := r.ReadLine()
                    if len(line) != 2 {
                        fmt.Println("坐标输入错误，请重新输入~")
                        continue
                    }
                    row = int8(line[0] - byte('1'))
                    if line[1] >= byte('a') {
                        col = int8(line[1] - byte('a'))
                    } else {
                        col =int8( line[1] - byte('A'))
                    }
                    //fmt.Println(x, y, len(line))
                    if row >=0 && row <LENGTH && col >=0 && col <LENGTH && canDown[row][col] {
                        PlacePiece(&chessboard, row, col, user_role) //落子
                        cnt++  //总落子数增加
                        skip_play = 0 //无法落子次数清0
                        break
                    }else{
                        fmt.Println("坐标输入错误，请重新输入~")
                    }
                }
                //fmt.Println(chessboard)
                PrintChessboard(chessboard)
                fmt.Print("请按任意键，AI落子~\n")
                //waitUser() //增加交互性
            }

        } else {
            //AI落子
            //首先查看是否可落子
            step , canDown := ai.CheckChessboard(chessboard, ai_role)
            if step == 0 {
                //无子可落
                skip_play ++
                if skip_play == 1 {
                    fmt.Println("AI目前没有位置可落子，请玩家落子。")
                }else if skip_play == 2 {
                    fmt.Println("双方均没有可落棋子")
                    break
                }
            } else {
                //Ai落子
                x, y := ai.AiPlayStep(chessboard, canDown, ai_role)
                //按上面计算出的最优解，AI最终落子
                PlacePiece(&chessboard, x, y, ai_role)
                skip_play = 0 //无法落子次数清0
                cnt++//总落子数增加
                //fmt.Println(chessboard)
                fmt.Println("\nAI落子(",x+1,",",string(byte(y)+byte('A')),")后：")
                PrintChessboard(chessboard)
            }
        }

        turn *= -1 //下一轮反转
    }

    //统计战场的分数
    user_score, ai_score := countScoreFinally(chessboard, user_role)

    if user_score < ai_score {
        fmt.Println("游戏结束，AI获胜")
    } else if user_score > ai_score {
        fmt.Println("游戏结束，玩家获胜")
    } else {
        fmt.Println("游戏结束，平局")
    }

    fmt.Println("得分情况----", "玩家：", user_score, "; ", "AI：", ai_score)
}

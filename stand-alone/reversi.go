package main

import (
    "fmt"
    "bufio"
    "os"
)

const (
    LENGTH = 8 //棋盘长宽
    WIITE  = 1 //百子颜色
    BLACK  = -1 //黑子颜色
)

//输出棋盘
func PrintChessboard(chessboard [LENGTH][LENGTH] int8) {
    var row, col int
    fmt.Printf("\n  ")

    //输出列号
    for col = 0; col < LENGTH; col++ {
        fmt.Printf("  %c ", int(byte('A')) + col);
    }
    fmt.Println();

    //输出项部横线
    fmt.Print("  ┌");
    //输出一行
    for col = 0; col < (LENGTH-1); col++ {
        fmt.Print("---┬");
    }
    fmt.Print("---┐\n");

    for row = 0; row < LENGTH; row++{
        //输出行号
        fmt.Printf("%2d│", row + 1);
        //输出棋盘各单元格中棋子的状态
        for col = 0; col < LENGTH; col++ {
            if chessboard[row][col] == 1 { //白棋
                fmt.Printf(" ○ │")
            } else if chessboard[row][col] == -1 { //白棋
                fmt.Printf(" ● │")
            } else { //未下子处
                fmt.Printf("   │")
            }
        }
        fmt.Println();
        if row < (LENGTH - 1) {
            fmt.Printf("  ├");  //输出交叉线
            //输出一行
            for col = 0; col < (LENGTH - 1); col++ {
                fmt.Printf("---┼");
            }
            fmt.Printf("---┤\n");
        }
    }
    fmt.Printf("  └");
    //最后一行的横线
    for col = 0; col < (LENGTH - 1); col++ {
        fmt.Printf("---┴");
    }
    fmt.Printf("---┘\n");
}

//分析棋盘，某一方是否还有哪些落子的地方
//参数：
//  chessboard -- 棋盘格局：1白色 -1黑色  0无子
//  role       -- 1表示为白子分析当前情况 -1表示为黑子分析当前情况
//
//返回值:
// setp        -- 可下子的位置的个数
// canDown     -- 可下子的位置
func CheckChessboard(chessboard[LENGTH][LENGTH] int8, role int8) (step int8, canDown [LENGTH][LENGTH] bool) {
    var row_delta, col_delta, row, col, x, y int8
    var self_color, opponent_color int8

    //确定本方和对方颜色
    self_color, opponent_color = role, -1*role

    //遍历分析棋盘所有为空的位置，分析是否符合落子的条件
    for row =0; row < LENGTH; row++ {
        for col =0; col < LENGTH; col++ {
            //跳过已经落子的位置
            if chessboard[row][col] != 0 {
                continue
            }

            //检查当前位置的周围8个方向（如果在边角上，则需要略过）
            for row_delta = -1; row_delta <=1; row_delta++ {
                for col_delta = -1; col_delta <=1; col_delta++ {
                    //忽略当前子和越界子（边角上）
                    if row+row_delta<0 || row+row_delta>=LENGTH || col+col_delta<0 || col+col_delta>=LENGTH || (row_delta == 0 && col_delta == 0) {
                        continue
                    }

                    //若(row,col)四周有对手下的子，即沿着这个方向一直追查是否有本方落子
                    //如果能够找到，说明当前位置是可以落子的，即对地方进行了合围
                    if chessboard[row+row_delta][col+col_delta] == opponent_color {
                        //fmt.Println(row, col, row_delta,  col_delta)
                        //以对手落子为起点
                        x, y = row+row_delta, col+col_delta
                        //沿着这个方向一直找
                        for {
                            x += row_delta
                            y += col_delta
                            //若越界跳出循环
                            if x<0 || x>=LENGTH || y <0 || y>=LENGTH {
                                break
                            }
                            //如果找到了空位置，则也说明无法落子了
                            if chessboard[x][y] == 0 {
                                break
                            }
                            //如果找到了本方棋子，则说明可以形成合围
                            if chessboard[x][y] == self_color {
                                canDown[row][col] = true
                                step ++
                                break
                            }
                        }
                    }
                }
            }
        }
    }

    return step, canDown
}

//指定位置落子
//参数：
//  chessboard -- 棋盘格局：1白色 -1黑色  0无子
//  row, col   -- 行列
//  role       -- 1表示为落白子 -1表示落黑子
//
//注意：
//golang中数字是值属性，所以要改变内容必须传入指针!
func PlacePiece(chessboard *[LENGTH][LENGTH] int8, row, col, role int8) {
    var row_delta, col_delta, x, y int8
    var self_color, opponent_color int8

    //确定本方和对方颜色
    self_color, opponent_color = role, -1*role

    chessboard[row][col] = self_color; //本方落子

    //检查当前位置的周围8个方向（如果在边角上，则需要略过）
    for row_delta = -1; row_delta <= 1; row_delta++ {
        for col_delta = -1; col_delta <= 1; col_delta++ {
            //忽略当前子和越界子（边角上）
            if row+row_delta<0 || row+row_delta>=LENGTH || col+col_delta<0 || col+col_delta>=LENGTH || (row_delta == 0 && col_delta == 0) {
                continue
            }

            //若(row,col)四周有对手下的子，即沿着这个方向一直追查是否有本方落子
            //如果能够找到，则将这中间所有的敌方棋子置换成本方棋子
            if chessboard[row+row_delta][col+col_delta] == opponent_color {
                //以对手落子为起点
                x, y = row+row_delta, col+col_delta
                //沿着这个方向一直找
                for {
                    x += row_delta
                    y += col_delta
                    //若越界跳出循环
                    if x<0 || x>=LENGTH || y <0 || y>=LENGTH {
                        break
                    }
                    //如果找到了空位置，则也说明无法落子了
                    if chessboard[x][y] == 0 {
                        break
                    }
                    //如果找到了本方棋子，则说明可以置换
                    if chessboard[x][y] == self_color {
                        x-=row_delta
                        y-=col_delta
                        for chessboard[x][y] == opponent_color {
                            chessboard[x][y] = self_color //置换
                            x-=row_delta
                            y-=col_delta
                        }
                        break
                    }
                }
            }
        }
    }
}

//分析当前棋局，获得指定方的得分
//参数：
//  chessboard -- 棋盘格局：1白色 -1黑色  0无子
//  role       -- 1表示分析白子得分 -1表示分析黑子得分
func GetScore(chessboard[LENGTH][LENGTH] int8, role int8) (score int8){
    var row, col int8
    var self_color, opponent_color int8

    //确定本方和对方颜色
    self_color, opponent_color = role, -1*role

    //遍历分析棋盘所有位置
    for row =0; row<LENGTH; row++ {
        for col = 0; col < LENGTH; col++ {
            //若棋盘对应位置是地方落的棋子，从总分中减1
            if chessboard[row][col] == opponent_color {
                score -= 1
            }
            //若棋盘对应位置是我方的棋子，总分中加1分
            if chessboard[row][col] == self_color {
                score += 1
            }
            //fmt.Println(row, ":", col, "->", chessboard[row][col], score)
        }
    }
    return
}

//分析所有落子方案，返回最高的一种得分
//参数：
//  chessboard -- 棋盘格局：1白色 -1黑色  0无子
//  canDown    -- 可以落子的位置
//  role       -- 1表示分析白子得分 -1表示分析黑子得分
//返回值：
//  最优方案的得分
func FindBestPlayScore(chessboard [LENGTH][LENGTH] int8, canDown [LENGTH][LENGTH] bool, role int8) (maxScore int8) {
    var row, col, i, j, score int8
    //for i =0; i<LENGTH; i++ {
    //    fmt.Println(chessboard[i])
    //}
    //fmt.Println()
    //for i =0; i<LENGTH; i++ {
    //    fmt.Println(canDown[i])
    //}
    var chessboard_tmp [LENGTH][LENGTH] int8

    maxScore = -128 //因为GetScore可能返回负值，所以如果maxScore采用默认0，会导致削平问题

    //遍历分析棋盘所有位置
    for row =0; row<LENGTH; row++ {
        for col = 0; col < LENGTH; col++ {
            //略过不可落子的位置
            if !canDown[row][col] {
                continue
            }

            //复制棋盘
            for i =0; i<LENGTH; i++ {
                for j = 0; j < LENGTH; j++ {
                    chessboard_tmp[i][j] = chessboard[i][j]
                }
            }

            //在镜像棋盘中落子，然后求出这种方案的得分
            PlacePiece(&chessboard_tmp, row, col, role)
            score = GetScore(chessboard_tmp, role)
            //fmt.Println(score, role)
            if maxScore < score {
                maxScore = score
            }
        }
    }
    //fmt.Println("Max score:", maxScore)
    return
}

//AI落子
//参数：
//  chessboard -- 棋盘格局：1白色 -1黑色  0无子
//  canDown    -- 可以落子的位置
//  role       -- 1表示分析白子落子 -1表示分析黑子落子
//
//核心思想：
// 假设敌方理性落子为前提，穷举本方落子所有可能后，敌方所有的落子可能性的最大值，得到一种对方最差的最大值，即Min-Max算法
func AiPlayStep(chessboard *[LENGTH][LENGTH] int8, canDown [LENGTH][LENGTH] bool, role int8) (row_best, col_best int8){
    var row, col, i, j, score, min_score int8

    var chessboard_snap [LENGTH][LENGTH] int8 //棋盘镜像
    //var canDown_snap [LENGTH][LENGTH] bool   //可落子镜像

    //确定本方和对方颜色
    self_color, opponent_color := role, -1*role
    min_score = 127 //敌方落子的最差的最大得分

    //遍历分析棋盘所有位置
    for row =0; row<LENGTH; row++ {
        for col = 0; col < LENGTH; col++ {
            //略过不可落子的位置
            if !canDown[row][col] {
                continue
            }

            //复制棋盘
            for i =0; i<LENGTH; i++ {
                for j = 0; j < LENGTH; j++ {
                    chessboard_snap[i][j] = chessboard[i][j]
                }
            }

            //试着在镜像棋盘中的一个位子下子
            PlacePiece(&chessboard_snap, row, col, self_color);

            //检查对手是否有地方可下子
            _, canDown_snap := CheckChessboard(chessboard_snap, opponent_color);

            //获得临时棋盘中对方下子的得分情况
            score = FindBestPlayScore(chessboard_snap, canDown_snap, opponent_color);

            //保存对方得分最低的下法
            if score < min_score {
                min_score = score;
                row_best = row;
                col_best = col;
                //fmt.Printf("row:%d, col:%d, score:%d\n", row, col, score)
                //fmt.Println()
            }
        }
    }
    //按上面计算出的最优解，AI最终落子
    PlacePiece(chessboard, row_best, col_best, self_color)
    return
}

//初始化棋盘：
//棋盘坐标：
//  左上角 -- row=0，col=0
//  右上角 -- row=0，col=7
//  左下角 -- row=7，col=0
//  右上角 -- row=7，col=7
func InitChessboard(chessboard *[8][8] int8) int {
    //在棋盘中间位置放置白棋
    chessboard[3][3] = WIITE;
    chessboard[4][4] = WIITE;

    //在棋盘中间位置放置黑棋
    chessboard[3][4] = BLACK;
    chessboard[4][3] = BLACK;

    return 4;
}

//天龙棋局。。。
//func InitChessboard(chessboard *[8][8] int8) int {
//    //在棋盘中间位置放置白棋
//    chessboard[0] = [8]int8{1,1,1,0,0,0,0,0}
//    chessboard[1] = [8]int8{0,1,1,0,0,0,0,0}
//    chessboard[2] = [8]int8{0,-1,1,1,0,0,0,0}
//    chessboard[3] = [8]int8{0,0,-1,1,1,0,0,0}
//    chessboard[4] = [8]int8{0,0,0,-1,-1,1,0,0}
//    chessboard[5] = [8]int8{0,0,0,0,-1,-1,1,1}
//    chessboard[6] = [8]int8{0,0,0,0,0,-1,1,1}
//    chessboard[7] = [8]int8{0,0,0,0,0,0,0,1}
//
//    return 22;
//}

//等待用户按任意键，用于阻塞程序
func waitUser() {
    r := bufio.NewReader(os.Stdin)
    _, _, _ = r.ReadLine()
}

//统计棋盘战场的得分
func countScoreFinally(chessboard[LENGTH][LENGTH] int8, role int8) (user_score, opponent_score int8){
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
    var chessboard [8][8] int8 //棋盘
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
            step , canDown := CheckChessboard(chessboard, user_role)
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
            step , canDown := CheckChessboard(chessboard, ai_role)
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
                x, y := AiPlayStep(&chessboard, canDown, ai_role)
                skip_play = 0 //无法落子次数清0
                cnt++
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

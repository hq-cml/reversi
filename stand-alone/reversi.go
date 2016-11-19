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
func printChessboard(chessboard [LENGTH][LENGTH] int8) {
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

//分析棋牌，某一方是否还有哪些落子的地方
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
    for row =0; row<LENGTH; row++ {
        for col =0; col<LENGTH; col++ {
            //跳过已经落子的位置
            if chessboard[row][col] != 0 {
                continue
            }

            //检查当前位置的周围8个方向（如果在边角上，则需要略过）
            for row_delta = -1; row_delta <=-1; row_delta++ {
                for col_delta = -1; col_delta <=-1; col_delta++ {
                    //忽略当前子和越界子（边角上）
                    if row+row_delta<0 || row+row_delta>=LENGTH || col+col_delta<0 || col+col_delta>=LENGTH || (row_delta == 0 && col_delta == 0) {
                        continue
                    }

                    //若(row,col)四周有对手下的子，即沿着这个方向一致追查是否有本方落子
                    //如果能够找到，说明当前位置是可以落子的，即对地方进行了合围
                    if chessboard[row+row_delta][col+col_delta] == opponent_color {
                        //以对手落子为起点
                        x, y = row+row_delta, col_delta
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
    for row_delta = -1; row_delta <=-1; row_delta++ {
        for col_delta = -1; col_delta <= -1; col_delta++ {
            //忽略当前子和越界子（边角上）
            if row+row_delta<0 || row+row_delta>=LENGTH || col+col_delta<0 || col+col_delta>=LENGTH || (row_delta == 0 && col_delta == 0) {
                continue
            }

            //若(row,col)四周有对手下的子，即沿着这个方向一致追查是否有本方落子
            //如果能够找到，则将这中间所有的敌方棋子置换成本方棋子
            if chessboard[row+row_delta][col+col_delta] == opponent_color {
                //以对手落子为起点
                x, y = row+row_delta, col_delta
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
            //若棋盘对应位置是对手下的棋子，从总分中减1
            if chessboard[row][col] == opponent_color {
                score -= 1
            }

            if chessboard[row][col] == self_color {
                score += 1
            }
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

    var chessboard_tmp [LENGTH][LENGTH] int8

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
            if maxScore < score {
                maxScore = score
            }
        }
    }
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
func AiPlayStep(chessboard *[LENGTH][LENGTH] int8, canDown [LENGTH][LENGTH] bool, role int8) {
    var row, col, row_best, col_best, i, j, score, min_score int8

    var chessboard_snap [LENGTH][LENGTH] int8 //棋盘镜像
    //var canDown_snap [LENGTH][LENGTH] bool   //可落子镜像

    //确定本方和对方颜色
    self_color, opponent_color := role, -1*role
    min_score = 127 //敌方落子的最差的最大值

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
            }
        }
    }
    //按上面计算出的最优解，AI最终落子
    PlacePiece(chessboard, row_best, col_best, self_color)
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

//单机版人机对战
func main() {
    //var row, col, x, y int8
    var chessboard [8][8] int8 //棋盘
    //var cnt int8 //已落子的个数
    //var SkipPlay int8 //当某一方无子可落时，增1，若为2，表示双方都不能落子
    var role int8 //本方颜色：1白色 -1黑色
    r := bufio.NewReader(os.Stdin)

    fmt.Println("\n~黑白棋小程序--人机对战AI版~\n")
    fmt.Println("初始棋盘:")

    //初始化棋盘
    _ = InitChessboard(&chessboard)
    //打印棋盘
    printChessboard(chessboard)

    //人机交互
    fmt.Print("\n游戏者执黑先下，输入 0-黑子 1-白子，请选择：")
    rawLine, _, _ := r.ReadLine()
    line := string(rawLine)
    if line == "0" {
        fmt.Println("您执黑先行！")
        role = -1
    }else if line == "1" {
        fmt.Println("您执白后行！")
        role = 1
    }else{
        fmt.Println("输入非法，程序退出！")
        os.Exit(1)
    }


    _ = role

}

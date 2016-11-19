package main

import "fmt"

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

//检索棋牌，某一方是否还有哪些下子的地方
//参数：
//  chessboard -- 棋盘格局：1白色 -1黑色  0无子
//  role       -- 1表示为白子分析当前情况 -1表示为黑子分析当前情况
//
//返回值:
// setp        -- 可下子的位置的个数
// canDown     -- 可下子的位置
func Check(chessboard[LENGTH][LENGTH] int8, role int8) (step int8, canDown [LENGTH][LENGTH] int8) {
    var row_delta, col_delta, row, col, x, y int8
    var self_color, opponent_color int8

    //确定本方和对方颜色
    self_color, opponent_color = row, -1*role

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
                                canDown[row][col] = 1
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
func PlacePiece(chessboard[LENGTH][LENGTH] int8, row, col, role int8) {
    var row_delta, col_delta, row, col, x, y int8
    var self_color, opponent_color int8

    //确定本方和对方颜色
    self_color, opponent_color = row, -1*role

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

func main() {
    var chessboard [8][8] int8
    printChessboard(chessboard)
}

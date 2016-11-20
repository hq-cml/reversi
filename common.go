package reversi

import "fmt"

const (
    LENGTH = 8    //棋盘长宽
    WIITE  = 1    //白子颜色
    BLACK  = -1   //黑子颜色
)

type ChessBoard [LENGTH][LENGTH]int8
type CanDown [LENGTH][LENGTH]bool


//打印输出棋盘
func PrintChessboard(chessboard ChessBoard) {
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


//初始化棋盘：
//棋盘坐标：
//  左上角 -- row=0，col=0
//  右上角 -- row=0，col=7
//  左下角 -- row=7，col=0
//  右上角 -- row=7，col=7
func InitChessboard(chessboard *ChessBoard) int {
    //在棋盘中间位置放置白棋
    chessboard[3][3] = WIITE;
    chessboard[4][4] = WIITE;

    //在棋盘中间位置放置黑棋
    chessboard[3][4] = BLACK;
    chessboard[4][3] = BLACK;

    return 4;

    //天龙棋局。。。
    //chessboard[0] = [8]int8{1,1,1,0,0,0,0,0}
    //chessboard[1] = [8]int8{0,1,1,0,0,0,0,0}
    //chessboard[2] = [8]int8{0,-1,1,1,0,0,0,0}
    //chessboard[3] = [8]int8{0,0,-1,1,1,0,0,0}
    //chessboard[4] = [8]int8{0,0,0,-1,-1,1,0,0}
    //chessboard[5] = [8]int8{0,0,0,0,-1,-1,1,1}
    //chessboard[6] = [8]int8{0,0,0,0,0,-1,1,1}
    //chessboard[7] = [8]int8{0,0,0,0,0,0,0,1}
    //
    //return 22;
}

//按指定位置落子，并根据规则吃掉地方棋子
//参数：
//  chessboard -- 棋盘格局：1白色 -1黑色  0无子
//  row, col   -- 行列
//  role       -- 1表示为落白子 -1表示落黑子
//
//注意：
//golang中数字是值属性，所以要改变内容必须传入指针!
func PlacePiece(chessboard *ChessBoard, row, col, role int8) {
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
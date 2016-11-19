package main

import "fmt"

//输出棋盘
func print_chessboard(chessboard [8][8] int8) {
    var row, col int
    fmt.Printf("\n  ")

    //输出列号
    for col = 0; col < 8; col++ {
        fmt.Printf("  %c ", int(byte('A')) + col);
    }
    fmt.Println();

    //输出项部横线
    fmt.Print("  ┌");
    //输出一行
    for col = 0; col < 7; col++ {
        fmt.Print("---┬");
    }
    fmt.Print("---┐\n");


    for row = 0; row < 8; row++{
        //输出行号
        fmt.Printf("%2d│", row + 1);
        //输出棋盘各单元格中棋子的状态
        for col = 0; col < 8; col++ {
            if chessboard[row][col] == 1 { //白棋
                fmt.Printf(" ○ │")
            } else if chessboard[row][col] == -1 { //白棋
                fmt.Printf(" ● │")
            } else { //未下子处
                fmt.Printf("   │")
            }
        }
        fmt.Println();
        if row < (8 - 1) {
            fmt.Printf("  ├");  //输出交叉线
            //输出一行
            for col = 0; col < (8 - 1); col++ {
                fmt.Printf("---┼");
            }
            fmt.Printf("---┤\n");
        }
    }
    fmt.Printf("  └");
    //最后一行的横线
    for col = 0; col < (8 - 1); col++ {
        fmt.Printf("---┴");
    }
    fmt.Printf("---┘\n");
}

func main() {
    var chessboard [8][8] int8
    print_chessboard(chessboard)
}

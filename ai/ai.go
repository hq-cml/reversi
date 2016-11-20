package ai

import (
    . "github.com/hq-cml/reversi"
)

/*
 * Mini版AI
 *
 * 核心思想，Min-Max算法：
 *
 * 假设前提为“敌方理性落子”，则：
 *  1. 遍历本方落子所有可能性
 *  2. 针对1中的所有可能性，计算敌方落子的最优解（得分最大值）
 *  3. 在2所有的敌方最优解中，挑出最差的（最小值）一种最优解
 *  4. 在3中最差的一种最优解，就是本方应该落子的位置
 */

//分析棋盘，某一方是否还有哪些落子的地方
//参数：
//  chessboard -- 棋盘格局：1白色 -1黑色  0无子
//  role       -- 1表示为白子分析当前情况 -1表示为黑子分析当前情况
//
//返回值:
// setp        -- 可下子的位置的个数
// canDown     -- 可下子的位置
func CheckChessboard(chessboard ChessBoard, role int8) (step int8, canDown CanDown) {
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
                            if chessboard[x][y] == NULL {
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

//分析当前棋局，获得指定方的得分
//参数：
//  chessboard -- 棋盘格局：1白色 -1黑色  0无子
//  role       -- 1表示分析白子得分 -1表示分析黑子得分
func GetScore(chessboard ChessBoard, role int8) (score int8){
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
func FindBestPlayScore(chessboard ChessBoard, canDown CanDown, role int8) (maxScore int8) {
    var row, col, i, j, score int8
    //for i =0; i<LENGTH; i++ {
    //    fmt.Println(chessboard[i])
    //}
    //fmt.Println()
    //for i =0; i<LENGTH; i++ {
    //    fmt.Println(canDown[i])
    //}
    var chessboard_tmp ChessBoard

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
func AiPlayStep(chessboard ChessBoard, canDown CanDown, role int8) (row_best, col_best int8){
    var row, col, i, j, score, min_score int8

    var chessboard_snap ChessBoard //棋盘镜像
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

            //获得镜像棋盘中对方下子的得分情况
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

    return
}


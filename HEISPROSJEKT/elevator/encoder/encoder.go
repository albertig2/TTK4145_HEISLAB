package encoder

/* 

Elevator State matrix

ID: ->             |        "01"              |   "02"                |      "03"            |            
----------------------------------------------------------------------------------------------
State              |  Idle/Moving/Open_door  |Idle/Moving/Open_door  |Idle/Moving/Open_door  |
Direction          | UP/DOWN/STILL           |UP/DOWN/STILL          |UP/DOWN/STILL          |
Network status     | Online/offline          | Online/offline        | Online/offline        |
Current floor      | 1...N					 | 1...N			     | 1...N			     |
Destination floor  | 1 ..N                   | 1...N			     | 1...N			     |
-----------------------------------------------------------------------------------------------


Order state matrix
ID: ->             |        "01"            |   "02"              |      "03"              |            
--------------------------------------------------------------------------------------------
Floor 1            |  up / down /cab       |  up / down /cab       |  up / down /cab       |
Floor 2            |  up / down /cab       |  up / down /cab       |  up / down /cab       |
Floor 3            |  up / down /cab       |  up / down /cab       |  up / down /cab       |
Floor 4            |  up / down /cab       |  up / down /cab       |  up / down /cab       |
--------------------------------------------------------------------------------------------

*/


import (
	"fmt"
	"strconv"
)

func printMatrix(matrix [][]int) {
	for row := 0; row < len(matrix); row++ {
		rowStr := ""
		for column := 0; column < len(matrix[row]); column++ {
			matrixElement := matrix[row][column]
			matrixElementStr := string(matrixElement)
			rowStr += matrixElementStr + " "

		}
		fmt.Println(rowStr)
		rowStr = ""
	}
}

func printVectorMatrix(matrix [][][]int) {

	for row := 0; row < len(matrix); row++ {
		rowStr := ""
		for column := 0; column < len(matrix[row]); column++ {
			vectorStr := ""
			matrixVector := matrix[row][column]

			vectorStr += "["
			for vectorIndex := 0; vectorIndex < len(matrixVector); vectorIndex++ {
				matrixVectorElement := matrixVector[vectorIndex]
				matrixVectorElementStr := string(matrixVectorElement)

				vectorStr += matrixVectorElementStr + " "
			}
			vectorStr += "]"
			rowStr += vectorStr

			vectorStr = ""
		}
		fmt.Println(rowStr)
		rowStr = ""
	}
}

func encodeMatrix(matrix [][]int) string {

	encodedString := ""
	for row := 0; row < len(matrix); row++ {
		for column := 0; column < len(matrix[row]); column++ {
			matrixElement := matrix[row][column]
			matrixElementStr := strconv.Itoa(matrixElement)

			encodedString += matrixElementStr
		}
		encodedString += "L"
	}
	return encodedString
}

func encodeVectorMatrix(matrix [][][]int) string {
	encodedString := ""
	for row := 0; row < len(matrix); row++ {
		encodedRow := ""
		for column := 0; column < len(matrix[row]); column++ {
			encodedVector := ""
			matrixVector := []int(matrix[row][column])
			for vectorIndex := 0; vectorIndex < len(matrixVector); vectorIndex++ {
				matrixVectorElement := matrixVector[vectorIndex]

				matrixVectorElementStr := strconv.Itoa(matrixVectorElement)
				encodedVector += matrixVectorElementStr
			}
			encodedVector += "V"
			encodedRow += encodedVector
		}
		encodedString += encodedRow
		encodedString += "L"

	}
	return encodedString
}

func encodeData2(senderID string, numFloors int, numButtons int, orderMatrix [][]int, stateMatrix [][][]int) string {

	encodedOrderMatrix := encodeMatrix(orderMatrix)
	encodedStateMatrix := encodeVectorMatrix(stateMatrix)
	fmt.Println(string(numFloors))

	encodedMessage := senderID + strconv.Itoa(numFloors) + strconv.Itoa(numButtons) + encodedOrderMatrix + string('X') + encodedStateMatrix

	return encodedMessage

}

func decodeMatrix(encodedMessage string) [][]int {
	var decodedMatrix [][]int
	newRow := []int{}

	for characterIndex := 0; characterIndex < len(encodedMessage); characterIndex++ {
		if encodedMessage[characterIndex] == 'L' { //Linebreak
			decodedMatrix = append(decodedMatrix, newRow)
			newRow = []int{}
		} else {

			newRow = append(newRow, int(encodedMessage[characterIndex]))
		}
	}
	return decodedMatrix
}

func decodeVectorMatrix(encodedMessage string) [][][]int {
	var decodedMatrix [][][]int
	newRow := [][]int{}
	newStateVector := []int{}

	for characterIndex := 0; characterIndex < len(encodedMessage); characterIndex++ {
		if encodedMessage[characterIndex] == 'L' { // L = Linebreak
			decodedMatrix = append(decodedMatrix, newRow)
			newRow = [][]int{}
		} else if encodedMessage[characterIndex] == 'V' { //V = new vector break
			newRow = append(newRow, newStateVector)
			newStateVector = []int{}
		} else {
			newStateVector = append(newStateVector, int(encodedMessage[characterIndex]))
			// fmt.Println(newStateVector)
		}
	}
	return decodedMatrix
}

func decodeData2(networkMessage string) ([][]int, [][][]int) {

	senderID := string(networkMessage[0]) + string(networkMessage[1])
	numFloors,_ := strconv.Atoi(string(networkMessage[2]))
	numButtons,_ := strconv.Atoi(string(networkMessage[3]))

	matrixBreakIndex := 0

	for characterIndex := 4; characterIndex < len(networkMessage); characterIndex++ {
		if networkMessage[characterIndex] != 'X' {
			continue
		}
		matrixBreakIndex = characterIndex
	}

	orderMatrix := decodeMatrix(networkMessage[4:matrixBreakIndex])
	stateMatrix := decodeVectorMatrix(networkMessage[matrixBreakIndex+1:])

	fmt.Println("Sender ID: ", senderID)
	fmt.Println("Number of floors: ", numFloors)
	fmt.Println("Number of buttons: ", numButtons)

	return orderMatrix, stateMatrix
}

// func getCabOrdersByID

// func main() {
	
// 	matrix1 := [][]int{
// 		{1, 2, 3, 4},
// 		{0, 1, 0, -1},
// 		{0, 1, 0, 1},
// 		{3, 1, 0, 2}}
// 	matrix2 := [][][]int{
// 		{{1}, {2}, {3}, {4}},
// 		{{1, 0, 2}, {2, 1, 1}, {0, 0, 1}, {0, 0, 0}},
// 		{{0, 0, 1}, {1, 3, 1}, {0, 3, 0}, {1, 1, 1}},
// 		{{0, 0, 0}, {0, 0, 0}, {2, 0, 1}, {2, 2, 2}}}

// 	str := encodeData2("01", 4,3, matrix1, matrix2)
// 	orders, states := decodeData2(str)
// 	fmt.Println("Orders")
// 	printMatrix(orders)
// 	fmt.Println("States")
// 	printVectorMatrix(states) 
// }

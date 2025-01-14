package main

import (
	"bufio"
	"container/heap"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	x, y int
}

type PqNode struct {
	point Point
	dist  int
}

type priorityQueue []*PqNode

func (pq *priorityQueue) Len() int {
	return len(*pq)
}

func (pq *priorityQueue) Less(i, j int) bool {
	return (*pq)[i].dist < (*pq)[j].dist
}

func (pq *priorityQueue) Swap(i, j int) {
	(*pq)[i], (*pq)[j] = (*pq)[j], (*pq)[i]
}

func (pq *priorityQueue) Push(x interface{}) {
	item := x.(*PqNode)
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

func readInput() ([][]int, Point, Point, error) {
	scanner := bufio.NewScanner(os.Stdin)

	// Считываем размеры лабиринта
	scanner.Scan()
	size := strings.Fields(scanner.Text())
	if len(size) != 2 {
		return nil, Point{}, Point{}, errors.New("некорректный формат размеров лабиринта")
	}

	rows, err := strconv.Atoi(size[0])
	if err != nil {
		return nil, Point{}, Point{}, errors.New("некорректное значение строк")
	}
	columns, err := strconv.Atoi(size[1])
	if err != nil {
		return nil, Point{}, Point{}, errors.New("некорректное значение столбцов")
	}

	// Считываем лабиринт
	maze := make([][]int, rows)
	for i := 0; i < rows; i++ {
		scanner.Scan()
		line := strings.Fields(scanner.Text())
		if len(line) != columns {
			return nil, Point{}, Point{}, errors.New("некорректное количество столбцов в строке")
		}
		maze[i] = make([]int, columns)
		for j, val := range line {
			num, err := strconv.Atoi(val)
			if err != nil {
				return nil, Point{}, Point{}, errors.New("некорректное значение в лабиринте")
			}
			maze[i][j] = num
		}
	}
	// Считываем стартовую и конечную точку
	scanner.Scan()
	coords := strings.Fields(scanner.Text())
	if len(coords) != 4 {
		return nil, Point{}, Point{}, errors.New("некорректный формат координат")
	}
	startX, err := strconv.Atoi(coords[0])
	if err != nil {
		return nil, Point{}, Point{}, errors.New("некорректное значение X для стартовой точки")
	}
	startY, err := strconv.Atoi(coords[1])
	if err != nil {
		return nil, Point{}, Point{}, errors.New("некорректное значение Y для стартовой точки")
	}
	endX, err := strconv.Atoi(coords[2])
	if err != nil {
		return nil, Point{}, Point{}, errors.New("некорректное значение X для конечной точки")
	}
	endY, err := strconv.Atoi(coords[3])
	if err != nil {
		return nil, Point{}, Point{}, errors.New("некорректное значение Y для конечной точки")
	}

	return maze, Point{startX, startY}, Point{endX, endY}, nil
}

func Dijkstra(grid [][]int, start, end Point) ([]Point, error) {
	rows := len(grid)
	columns := len(grid[0])

	if start.x < 0 || start.x >= rows ||
		start.y < 0 || start.y >= columns {
		return nil, errors.New("стартовая точка вне лабиринта")
	}
	if end.x < 0 || end.x >= rows ||
		end.y < 0 || end.y >= columns {
		return nil, errors.New("финишная точка вне лабиринта")
	}

	if grid[start.x][start.y] == 0 {
		return nil, errors.New("стартовая точка - тупик")
	}
	if grid[end.x][end.y] == 0 {
		return nil, errors.New("финишная точка - тупик")
	}

	dist := make([][]int, rows)
	for i := 0; i < rows; i++ {
		dist[i] = make([]int, columns)
		for j := 0; j < columns; j++ {
			dist[i][j] = 1 << 20
		}
	}
	dist[start.x][start.y] = grid[start.x][start.y]

	parent := make([][]*Point, rows)
	for i := 0; i < rows; i++ {
		parent[i] = make([]*Point, columns)
	}

	pq := make(priorityQueue, 1)
	pq[0] = &PqNode{
		point: start,
		dist:  dist[start.x][start.y],
	}
	heap.Init(&pq)

	directions := []Point{
		{-1, 0}, // вверх
		{1, 0},  // вниз
		{0, -1}, // влево
		{0, 1},  // вправо
	}

	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*PqNode)
		curPoint := item.point
		curDist := item.dist

		if curPoint == end {
			break
		}

		if curDist > dist[curPoint.x][curPoint.y] {
			continue
		}

		for _, d := range directions {
			nextX, nextY := curPoint.x+d.x, curPoint.y+d.y
			if nextX < 0 || nextX >= rows || nextY < 0 || nextY >= columns {
				continue
			}
			if grid[nextX][nextY] == 0 {
				continue
			}

			newDist := dist[curPoint.x][curPoint.y] + grid[nextX][nextY]
			if newDist < dist[nextX][nextY] {
				dist[nextX][nextY] = newDist
				parent[nextX][nextY] = &Point{curPoint.x, curPoint.y}
				heap.Push(&pq, &PqNode{
					point: Point{nextX, nextY},
					dist:  newDist,
				})
			}
		}
	}

	if dist[end.x][end.y] == 1<<20 {
		return nil, errors.New("нет пути до финиша")
	}

	// Восстанавливаем путь
	path := make([]Point, 0)
	cur := end
	for {
		path = append(path, cur)
		if cur == start {
			break
		}
		p := parent[cur.x][cur.y]
		cur = *p
	}

	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return path, nil
}

func main() {
	grid, start, finish, err := readInput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка ввода: %v\n", err)
		os.Exit(1)
	}
	path, err := Dijkstra(grid, start, finish)
	if err != nil {
		if err.Error() == "нет пути до финиша" {
			fmt.Println("Пути не существует")
			return
		} else {
			fmt.Fprintf(os.Stderr, "Ошибка во время расчета пути: %v\n", err)
			os.Exit(1)
		}
	}
	for _, p := range path {
		fmt.Printf("%d %d\n", p.x, p.y)
	}
	fmt.Printf("%s\n", ".")
}

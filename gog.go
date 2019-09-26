package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sync"
	"time"
)

const (
	GREEN = "\033[0;32m"
	RESET = "\x1b[0m"
	GREY  = "\033[0;37m"
)

func clear_terminal() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", name, elapsed)
}

func display(world [][]uint8) {
	clear_terminal()
	for _, world_slice := range world {
		for _, cell := range world_slice {
			if cell == 1 {
				fmt.Printf(" %v%v%v", GREEN, cell, RESET)
			} else {
				fmt.Printf(" %v%v%v", GREY, cell, RESET)
			}
		}
		fmt.Printf("\n")
	}
}

func init_world(size int) [][]uint8 {
	defer timeTrack(time.Now(), "init_world")
	var pop_density float64 = 0.5
	world := make([][]uint8, size)
	for i := range world {
		world[i] = make([]uint8, size)
		for j := range world[i] {
			rand_num := rand.Float64()
			if rand_num < pop_density {
				world[i][j] = 1
			}
		}
	}
	return world
}

func iterate_world_par(world [][]uint8) [][]uint8 {
	size := len(world)
	next_world := init_world(size)
	var wg sync.WaitGroup
	wg.Add(len(next_world))
	start := time.Now()
	for i := range next_world {
		go func(i int) {
			defer wg.Done()
			for j := range next_world[i] {
				next_world[i][j] = update_cell(i, j, world)
			}
		}(i)
	}
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("Parallel Update took %s\n", elapsed)
	return next_world
}

func iterate_world(world [][]uint8) [][]uint8 {
	size := len(world)
	next_world := init_world(size)
	start := time.Now()
	for i := range next_world {
		for j := range next_world[i] {
			next_world[i][j] = update_cell(i, j, world)
		}
	}
	elapsed := time.Since(start)
	fmt.Printf("Serial Update took %s\n", elapsed)
	return next_world
}

func update_cell(i, j int, world [][]uint8) uint8 {
	current_state := world[i][j]
	var new_state uint8 = 9 // if a 9 gets through the logic below is broken
	live_count := live_neighbor_count(i, j, world)
	if current_state == 1 {
		if live_count < 2 {
			new_state = 0
		} else if live_count < 4 {
			new_state = 1
		} else {
			new_state = 0
		}
	} else {
		if live_count == 3 {
			new_state = 1
		} else {
			new_state = 0
		}
	}
	return uint8(new_state)
}

func live_neighbor_count(i, j int, world [][]uint8) int {
	var live_count int = 0
	coords := neighbor_coordinates(i, j, world)
	for _, coord_pair := range coords {
		if world[coord_pair[0]][coord_pair[1]] == 1 {
			live_count++
		}
	}
	return live_count
}

func neighbor_coordinates(i, j int, world [][]uint8) [][2]int {
	coords := make([][2]int, 0, 8) // len 0 cap 8
	size := len(world)
	for n := i - 1; n < i+2; n++ {
		for m := j - 1; m < j+2; m++ {
			if !((n == i) && (m == j)) {
				if (n >= 0) && (m >= 0) {
					if (n < size) && (m < size) {
						new_coord_pair := [2]int{n, m}
						coords = append(coords, new_coord_pair)
					}
				}
			}
		}
	}
	return coords
}

func run(iter_num, size int, parallel_comp, print_world bool) {
	rand.Seed(time.Now().UTC().UnixNano())
	world := init_world(size)
	for i := 0; i < iter_num; i++ {
		if print_world {
			display(world)
		}
		if parallel_comp {
			world = iterate_world_par(world)
		} else {
			world = iterate_world(world)
		}
		if print_world {
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func main() {
	start := time.Now()

	// Variable program parameters here //
	iter_num := 10
	size := 20
	parallel_comp := false
	print_world := true
	/////////////////////////////////////

	if print_world && (size > 50) {
		fmt.Printf("World too big to print! If you want to watch the world in the terminal pick a size of 50 or less.\n")
		return
	}

	run(iter_num, size, parallel_comp, print_world)

	elapsed := time.Since(start)
	fmt.Printf("Whole program took %s\n", elapsed)
}

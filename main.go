package main

import (
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

type Dinosaur struct {
	x, y    int
	jumping bool
}

type Obstacle struct {
	x, y int
}

const DINO = "🦖"

func gameOver() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	msg := "Game Over! Press ESC to exit."
	for i, c := range msg {
		termbox.SetCell(5+i, 5, c, termbox.ColorRed, termbox.ColorDefault)
	}
	termbox.Flush()

	// Wait for ESC key to exit
	for {
		ev := termbox.PollEvent()
		if ev.Key == termbox.KeyEsc {
			break
		}
	}
}

func render(dino Dinosaur, obstacles []Obstacle) {
	// Clear the screen before drawing
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// Draw dinosaur
	termbox.SetCell(dino.x, dino.y, '🦖', termbox.ColorGreen, termbox.ColorDefault)

	// Draw obstacles
	for _, obs := range obstacles {
		termbox.SetCell(obs.x, obs.y, 'O', termbox.ColorRed, termbox.ColorDefault)
	}

	// Refresh the screen
	termbox.Flush()
}

func checkCollision(dino Dinosaur, obstacles []Obstacle) bool {
	for _, obs := range obstacles {
		if obs.x == dino.x && obs.y == dino.y {
			return true // Collision detected
		}
	}
	return false
}

func updateObstacles(obstacles []Obstacle) {
	for i := 0; i < len(obstacles); i++ {
		obstacles[i].x -= 1 // Move each obstacle to the left

		// Remove obstacle if it's out of view (off the left side of the screen)
		if obstacles[i].x < 0 {
			obstacles = append(obstacles[:i], obstacles[i+1:]...) // Remove obstacle
			i--
		}
	}

	// Randomly add a new obstacle with a certain chance
	if rand.Intn(10) == 0 { // 1 in 10 chance to add an obstacle
		width, _ := termbox.Size()
		newObstacle := Obstacle{x: width, y: 10}
		obstacles = append(obstacles, newObstacle)
	}
}

func updateDinosaur(dino *Dinosaur) {
	if dino.jumping {
		// Check if the dinosaur has reached the peak of the jump
		if dino.y > 5 { // Assuming 5 is the highest jump height
			dino.y -= 1 // Move dinosaur up
		} else {
			dino.jumping = false // Start falling back down
		}
	} else if dino.y < 10 { // Assume 10 is the ground level
		dino.y += 1 // Move dinosaur down (gravity effect)
	}
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	// make sure we clear up everything before leaving
	defer termbox.Close()

	dino := Dinosaur{x: 5, y: 10}
	obstacles := []Obstacle{{x: 20, y: 10}}

	gameTick := time.NewTicker(50 * time.Millisecond) // Timer to trigger updates
	defer gameTick.Stop()

	// Channel to receive events
	eventChannel := make(chan termbox.Event)

	// Start a goroutine to poll for events and send them to the channel
	// This is needed so that we can read keyboard event in the `non-blocking` manner
	go func() {
		for {
			eventChannel <- termbox.PollEvent() // Blocking call to PollEvent
		}
	}()

	// Main game loop
	for {
		select {
		case <-gameTick.C: // Game update and render every tick
			updateDinosaur(&dino)
			updateObstacles(obstacles)

			if checkCollision(dino, obstacles) {
				gameOver()
				return
			}

			render(dino, obstacles)

		case ev := <-eventChannel: // Handle keyboard events
			if ev.Type == termbox.EventKey {
				if ev.Key == termbox.KeyEsc {
					// Finish game
					return
				}
				if ev.Key == termbox.KeySpace && !dino.jumping {
					dino.jumping = true
				}
			}
		}
	}
}

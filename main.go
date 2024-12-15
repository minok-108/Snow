package main


import (
  "fmt"
  "math"
  "math/rand/v2"
  "image/color"
  "github.com/hajimehoshi/ebiten/v2"
  "github.com/hajimehoshi/ebiten/v2/vector"
  "github.com/hajimehoshi/ebiten/v2/ebitenutil"

)

/******
Structs
******/

type Vector2 struct {
  X float64
  Y float64
}

type Snowflake struct {
  Position Vector2
  Size float64
  Velocity Vector2
}

type Wind struct {
  Position Vector2
  Size float64
  Force Vector2
}

type Game struct {
  Snowflakes []*Snowflake
  Winds []*Wind
  Map [screenH][screenW]bool
}

/****
Const
****/

const velX = 0.5
const velY = 0.4

const screenW = 240
const screenH = 120

/***
Game
***/

var backgroundImage *ebiten.Image

func init() {

  img, _, err := ebitenutil.NewImageFromFile("background.png")
	if err != nil {
		fmt.Println(err)
	}
	backgroundImage = ebiten.NewImageFromImage(img)
}


func loopX(x int) int {
  x %= screenW

  if x >= 0 {
    return x
  } else {
    return screenW + x
  }

  return x
}


func slip(g *Game, x, y int) (int, int) {
  if y < screenH - 1 {
    if g.Map[y + 1][loopX(x + 1)] == false  {
      x += 1
      y += 1
      x, y = slip(g, x % screenW, y)
    } else if g.Map[y + 1][loopX(x - 1)] == false {
      x -= 1
      y += 1
      x, y = slip(g, loopX(x), y)
    }
  }

  for y + 1 < screenH && g.Map[y + 1][x] == false  {
    y += 1
  }

  return x, y
}


func createSnowflakes(g *Game) {
  for i := 0; i <= 1; i ++ {
    var snowflake Snowflake
    snowflake.Position.X = float64(rand.IntN(screenW))
    snowflake.Size = 1
    snowflake.Velocity = Vector2{
      velX * (rand.Float64() - 0.5),
      velY * (rand.Float64() + 1),
    }

    g.Snowflakes = append(g.Snowflakes, &snowflake)
  }
}


func processSnowflakes(g *Game) {
  for i, s := range g.Snowflakes {
    for _, w := range g.Winds {
      var d = math.Pow(s.Position.X - w.Position.X, 2) + math.Pow(s.Position.Y - w.Position.Y, 2)
      if d < math.Pow(w.Size, 2) {
        s.Velocity.X += w.Force.X
      }
    }

    s.Position.X += s.Velocity.X
    s.Position.Y += s.Velocity.Y

    var mapPosX = int(s.Position.X)
    var mapPosY = int(s.Position.Y)

    if mapPosX >= screenW {
      s.Position.X -= screenW
    }

    var setPixel = false

    if mapPosY >= screenH {
      mapPosY = screenH - 1
      setPixel = true
    }

    for g.Map[mapPosY][loopX(mapPosX)] == true {
      mapPosY -= 1
      setPixel = true
    }

    if setPixel {
      mapPosX, mapPosY = slip(g, loopX(mapPosX), mapPosY)
      g.Map[mapPosY][mapPosX] = true

      backgroundImage.Set(mapPosX, mapPosY, color.RGBA{255, 255, 255, 255})

      g.Snowflakes = append(g.Snowflakes[:i], g.Snowflakes[i+1:]...)
    }
  }
}


func wind(g *Game) {
  var r = rand.IntN(100)

  if len(g.Winds) > 1 && r == 0 {
    var i = rand.IntN(len(g.Winds) - 1)
    g.Winds = append(g.Winds[:i], g.Winds[i+1:]...)
  }

  if len(g.Winds) < 5 && r == 1 {
    var wind Wind

    wind.Position = Vector2{float64(rand.IntN(screenW)), float64(rand.IntN(screenH))}
    wind.Size = 20 + float64(rand.IntN(30))
    wind.Force = Vector2{rand.Float64() / 40 * float64(rand.IntN(2) * 2 - 1), 0}

    g.Winds = append(g.Winds, &wind)
  }
}


func (g *Game) Update() error {
  createSnowflakes(g)
  processSnowflakes(g)

  wind(g)

  return nil
}


func (g *Game) Draw(screen *ebiten.Image) {
  screen.DrawImage(backgroundImage, nil)

  for _, w := range g.Winds {
    vector.DrawFilledCircle(screen, float32(w.Position.X), float32(w.Position.Y), float32(w.Size), color.RGBA{88, 124, 120, 100}, false)
  }

  for _, s := range g.Snowflakes {
    screen.Set(int(s.Position.X), int(s.Position.Y), color.RGBA{255, 255, 255, 255})
  }
}


func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 240, 120
}

/********
Functions
********/

func main() {
  fmt.Println("Hello Snow!")

  g := &Game{}

  for h := 0; h < screenH; h ++ {
    for w := 0; w < screenW; w ++ {
      g.Map[h][w] = false
    }
  }

  ebiten.SetWindowSize(960 * 2, 480 * 2)
	ebiten.SetWindowTitle("Hello Snow!")

	if err := ebiten.RunGame(g); err != nil {
		fmt.Println(err)
	}
}

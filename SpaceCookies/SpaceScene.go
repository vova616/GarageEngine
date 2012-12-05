package SpaceCookies

import (
	"fmt"
	"github.com/vova616/GarageEngine/Engine"
	"github.com/vova616/GarageEngine/Engine/Components"
	_ "image/jpeg"
	_ "image/png"
	//"gl"  
	"strconv"
	"time"
	//"strings"
	//"math"
	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"
	//"image"
	//"image/color"
	"encoding/json"
	"math/rand"
	"os"
)

type GameScene struct {
	*Engine.SceneData
	Layer1 *Engine.GameObject
	Layer2 *Engine.GameObject
	Layer3 *Engine.GameObject
	Layer4 *Engine.GameObject
}

var (
	GameSceneGeneral *GameScene
	cir              *Engine.Texture
	boxt             *Engine.Texture
	cookie           *Engine.GameObject
	defender         *Engine.GameObject

	Player     *Engine.GameObject
	PlayerShip *ShipController

	Explosion *Engine.GameObject
	PowerUpGO *Engine.GameObject

	Wall *Engine.GameObject

	atlas        *Engine.ManagedAtlas
	atlasSpace   *Engine.ManagedAtlas
	atlasPowerUp *Engine.ManagedAtlas
	backgroung   *Engine.Texture
	ArialFont    *Engine.Font
	ArialFont2   *Engine.Font

	queenDead = false
)

const (
	MissleTag = "Missle"
	CookieTag = "Cookie"
)
const SpaceShip_A = 333
const Missle_A = 334
const HP_A = 123
const HPGUI_A = 124
const Queen_A = 666
const Jet_A = 125

func CheckError(err error) bool {
	if err != nil {
		fmt.Println(err)
		return true
	}
	return false
}

func LoadTextures() {
	atlas = Engine.NewManagedAtlas(2048, 1024)
	atlasSpace = Engine.NewManagedAtlas(1024, 1024)
	atlasPowerUp = Engine.NewManagedAtlas(256, 256)

	CheckError(atlas.LoadImage("./data/SpaceCookies/Ship1.png", SpaceShip_A))
	CheckError(atlas.LoadImage("./data/SpaceCookies/missile.png", Missle_A))
	CheckError(atlas.LoadGroupSheet("./data/SpaceCookies/Explosion.png", 128, 128, 6*8))

	CheckError(atlas.LoadImage("./data/SpaceCookies/HealthBar.png", HP_A))
	CheckError(atlas.LoadImage("./data/SpaceCookies/HealthBarGUI.png", HPGUI_A))
	CheckError(atlas.LoadImage("./data/SpaceCookies/Queen.png", Queen_A))
	CheckError(atlas.LoadImage("./data/SpaceCookies/Jet.png", Jet_A))

	atlas.BuildAtlas()
	atlas.BuildMipmaps()
	atlas.SetFiltering(Engine.MipMapLinearNearest, Engine.Nearest)
	atlas.Texture.SetReadOnly()

	var e error
	boxt, e = Engine.LoadTexture("./data/SpaceCookies/wall.png")

	boxt.BuildMipmaps()
	boxt.SetFiltering(Engine.MipMapLinearNearest, Engine.Nearest)

	backgroung, e = Engine.LoadTexture("./data/SpaceCookies/background.png")
	CheckError(e)
	cir, e = Engine.LoadTexture("./data/SpaceCookies/Cookie.png")
	CheckError(e)

	cir.BuildMipmaps()
	cir.SetFiltering(Engine.MipMapLinearNearest, Engine.Nearest)

	backgroung.BuildMipmaps()
	backgroung.SetFiltering(Engine.MipMapLinearNearest, Engine.Nearest)

	CheckError(atlasSpace.LoadGroup("./data/SpaceCookies/Space/"))
	atlasSpace.BuildAtlas()
	atlasSpace.BuildMipmaps()
	atlasSpace.SetFiltering(Engine.MipMapLinearNearest, Engine.Nearest)
	atlasSpace.Texture.SetReadOnly()

	CheckError(atlasPowerUp.LoadGroupSheet("./data/SpaceCookies/powerups.png", 61, 61, 3*4))
	atlasPowerUp.BuildAtlas()
	atlasPowerUp.SetFiltering(Engine.Linear, Engine.Linear)

	ArialFont, e = Engine.NewFont("./data/Fonts/arial.ttf", 48)
	if e != nil {
		panic(e)
	}
	ArialFont.Texture.SetReadOnly()

	ArialFont2, e = Engine.NewFont("./data/Fonts/arial.ttf", 24)
	if e != nil {
		panic(e)
	}
	ArialFont2.Texture.SetReadOnly()
}

func init() {
	Engine.Title = "Space Cookies"
}

func (s *GameScene) Load() {

	LoadTextures()

	queenDead = false

	rand.Seed(time.Now().UnixNano())

	GameSceneGeneral = s

	s.Camera = Engine.NewCamera()

	cam := Engine.NewGameObject("Camera")
	cam.AddComponent(s.Camera)

	cam.Transform().SetScalef(1, 1)

	gui := Engine.NewGameObject("GUI")

	Layer1 := Engine.NewGameObject("Layer1")
	Layer2 := Engine.NewGameObject("Layer2")
	Layer3 := Engine.NewGameObject("Layer3")
	Layer4 := Engine.NewGameObject("Layer3")

	s.Layer1 = Layer1
	s.Layer2 = Layer2
	s.Layer3 = Layer3
	s.Layer4 = Layer4

	mouse := Engine.NewGameObject("Mouse")
	mouse.AddComponent(Engine.NewMouse())
	mouse.AddComponent(NewMouseDebugger())
	mouse.Transform().SetParent2(cam)

	FPSDrawer := Engine.NewGameObject("FPS")
	FPSDrawer.Transform().SetParent2(cam)
	txt := FPSDrawer.AddComponent(Components.NewUIText(ArialFont2, "")).(*Components.UIText)
	fps := FPSDrawer.AddComponent(Engine.NewFPS()).(*Engine.FPS)
	fps.SetAction(func(fps float32) {
		txt.SetString("FPS: " + strconv.FormatFloat(float64(fps), 'f', 2, 32))
	})
	txt.SetAlign(Engine.AlignLeft)

	FPSDrawer.Transform().SetPositionf(20, float32(Engine.Height)-20)
	FPSDrawer.Transform().SetScalef(20, 20)

	label := Engine.NewGameObject("Label")
	label.Transform().SetParent2(cam)
	label.Transform().SetPositionf(20, float32(Engine.Height)-40)
	label.Transform().SetScalef(20, 20)

	txt2 := label.AddComponent(Components.NewUIText(ArialFont2, "Input: ")).(*Components.UIText)
	txt2.SetFocus(true)
	txt2.SetAlign(Engine.AlignLeft)

	//SPACCCEEEEE
	Engine.Space.Gravity.Y = 0
	Engine.Space.Iterations = 10

	Health := Engine.NewGameObject("HP")
	Health.Transform().SetParent2(cam)
	Health.Transform().SetPositionf(150, 50)

	HealthGUI := Engine.NewGameObject("HPGUI")
	HealthGUI.AddComponent(Engine.NewSprite2(atlas.Texture, Engine.IndexUV(atlas, HPGUI_A)))
	HealthGUI.Transform().SetParent2(Health)
	HealthGUI.Transform().SetScalef(50, 50)

	HealthBar := Engine.NewGameObject("HealthBar")
	HealthBar.Transform().SetParent2(Health)
	HealthBar.Transform().SetPositionf(-82, 0)
	HealthBar.Transform().SetScalef(100, 50)

	uvHP := Engine.IndexUV(atlas, HP_A)

	HealthBarGUI := Engine.NewGameObject("HealthBarGUI")
	HealthBarGUI.Transform().SetParent2(HealthBar)
	HealthBarGUI.AddComponent(Engine.NewSprite2(atlas.Texture, uvHP))
	HealthBarGUI.Transform().SetScalef(0.52, 1)
	HealthBarGUI.Transform().SetPositionf((uvHP.Ratio/2)*HealthBarGUI.Transform().Scale().X, 0)

	JetFire := Engine.NewGameObject("Jet")
	JetFire.AddComponent(Engine.NewSprite2(atlas.Texture, Engine.IndexUV(atlas, Jet_A)))

	ship := Engine.NewGameObject("Ship")
	Player = ship
	ship.AddComponent(Engine.NewSprite2(atlas.Texture, Engine.IndexUV(atlas, SpaceShip_A)))
	PlayerShip = ship.AddComponent(NewShipController()).(*ShipController)
	ship.Transform().SetParent2(Layer2)
	ship.Transform().SetPositionf(400, 200)
	ship.Transform().SetScalef(100, 100)
	shipHP := float32(1000)
	PlayerShip.HPBar = HealthBar
	PlayerShip.JetFire = JetFire

	settings := struct {
		Ship                *ShipController
		PowerUpChance       *int
		PowerUpRepairChance *int
		ShipHP              *float32
		Debug               *bool
	}{
		PlayerShip,
		&PowerUpChance,
		&PowerUpRepairChance,
		&shipHP,
		&Engine.Debug,
	}

	f, e := os.Open("./data/SpaceCookies/game.dat")
	if e != nil {
		f, e = os.Create("./data/SpaceCookies/game.dat")
		if e != nil {
			fmt.Println(e)
		}
		defer f.Close()
		encoder := json.NewEncoder(f)
		encoder.Encode(settings)
	} else {
		defer f.Close()
	}
	decoder := json.NewDecoder(f)
	e = decoder.Decode(&settings)
	if e != nil {
		fmt.Println(e)
	}
	ship.AddComponent(NewDestoyable(shipHP, 1))

	uvs, ind := Engine.AnimatedGroupUVs(atlas, "Explosion")
	Explosion = Engine.NewGameObject("Explosion")
	Explosion.AddComponent(Engine.NewSprite3(atlas.Texture, uvs))
	Explosion.Sprite.BindAnimations(ind)
	Explosion.Sprite.AnimationSpeed = 25
	Explosion.Sprite.AnimationEndCallback = func(sprite *Engine.Sprite) {
		sprite.GameObject().Destroy()
	}
	Explosion.Transform().SetScalef(30, 30)

	missle := Engine.NewGameObject("Missle")
	missle.AddComponent(Engine.NewSprite2(atlas.Texture, Engine.IndexUV(atlas, Missle_A)))
	missle.AddComponent(Engine.NewPhysics(false, 10, 10))
	missle.Transform().SetScalef(20, 20)
	missle.AddComponent(NewDamageDealer(50))

	m := NewMissle(30000)
	missle.AddComponent(m)
	PlayerShip.Missle = m
	m.Explosion = Explosion
	ds := NewDestoyable(0, 1)
	ds.SetDestroyTime(1)
	missle.AddComponent(ds)

	cookie = Engine.NewGameObject("Cookie")
	cookie.AddComponent(Engine.NewSprite(cir))
	cookie.AddComponent(NewDestoyable(100, 2))
	cookie.AddComponent(NewDamageDealer(20))
	cookie.AddComponent(NewEnemeyAI(Player, Enemey_Cookie))
	cookie.Transform().SetScalef(50, 50)
	cookie.Transform().SetPositionf(400, 400)
	cookie.AddComponent(Engine.NewPhysics2(false, chipmunk.NewCircle(vect.Vect{0, 0}, 25)))
	cookie.Tag = CookieTag

	defender = Engine.NewGameObject("Box")
	ds = NewDestoyable(30, 3)
	ds.SetDestroyTime(5)
	defender.AddComponent(ds)
	defender.AddComponent(Engine.NewSprite(boxt))
	defender.Tag = CookieTag
	defender.Transform().SetScalef(50, 50)

	phx := defender.AddComponent(Engine.NewPhysics(false, 50, 50)).(*Engine.Physics)
	phx.Body.SetMass(2.5)
	phx.Body.SetMoment(phx.Shape.Moment(2.5))
	phx.Shape.SetFriction(0.5)
	//phx.Shape.Group = 2
	phx.Shape.SetElasticity(0.5)

	QueenCookie := Engine.NewGameObject("Cookie")
	QueenCookie.AddComponent(Engine.NewSprite2(atlas.Texture, Engine.IndexUV(atlas, Queen_A)))
	QueenCookie.AddComponent(NewDestoyable(5000, 2))
	QueenCookie.AddComponent(NewDamageDealer(200))
	QueenCookie.AddComponent(NewEnemeyAI(Player, Enemey_Boss))
	QueenCookie.Transform().SetParent2(Layer2)
	QueenCookie.Transform().SetScalef(300, 300)
	QueenCookie.Transform().SetPositionf(999999, 999999)
	QueenCookie.AddComponent(Engine.NewPhysics2(false, chipmunk.NewCircle(vect.Vect{0, 0}, 25)))
	QueenCookie.Tag = CookieTag

	staticCookie := Engine.NewGameObject("Cookie")
	staticCookie.AddComponent(Engine.NewSprite(cir))
	staticCookie.Transform().SetScalef(400, 400)
	staticCookie.Transform().SetPositionf(400, 400)
	staticCookie.AddComponent(NewDestoyable(float32(Engine.Inf), 2))
	staticCookie.AddComponent(Engine.NewPhysics2(true, chipmunk.NewCircle(vect.Vect{0, 0}, 200)))

	staticCookie.Physics.Shape.SetElasticity(0)
	staticCookie.Physics.Body.SetMass(999999999999)
	staticCookie.Physics.Body.SetMoment(staticCookie.Physics.Shape.Moment(999999999999))
	staticCookie.Tag = CookieTag

	uvs, ind = Engine.AnimatedGroupUVs(atlasSpace, "s")
	Background := Engine.NewGameObject("Background")
	Background.AddComponent(Engine.NewSprite3(atlasSpace.Texture, uvs))
	Background.Sprite.BindAnimations(ind)
	Background.Sprite.SetAnimation("s")
	Background.Sprite.AnimationSpeed = 0
	Background.Transform().SetScalef(50, 50)
	Background.Transform().SetPositionf(400, 400)

	uvs, ind = Engine.AnimatedGroupUVs(atlasPowerUp, "powerups")
	PowerUpGO = Engine.NewGameObject("Background")
	//PowerUpGO.Transform().SetParent2(Layer2)
	PowerUpGO.AddComponent(Engine.NewSprite3(atlasPowerUp.Texture, uvs))
	PowerUpGO.AddComponent(Engine.NewPhysics(false, 61, 61))
	PowerUpGO.Physics.Shape.IsSensor = true
	PowerUpGO.Sprite.BindAnimations(ind)
	PowerUpGO.Sprite.SetAnimation("powerups")
	PowerUpGO.Sprite.AnimationSpeed = 0
	index := (rand.Int() % 6) + 6
	PowerUpGO.Sprite.SetAnimationIndex(int(index))
	PowerUpGO.Transform().SetScalef(61, 61)
	PowerUpGO.Transform().SetPositionf(0, 0)

	background := Engine.NewGameObject("Background")
	background.AddComponent(Engine.NewSprite(backgroung))
	background.AddComponent(NewBackground(background.Sprite))
	background.Sprite.Render = false
	//background.Transform().SetScalef(float32(backgroung.Height()), float32(backgroung.Height()), 1)
	background.Transform().SetScalef(800, 800)
	background.Transform().SetPositionf(0, 0)

	for i := 0; i < 300; i++ {
		c := Background.Clone()
		c.Transform().SetParent2(Layer4)
		size := 20 + rand.Float32()*50
		p := Engine.Vector{(rand.Float32() * 5000) - 1000, (rand.Float32() * 5000) - 1000, 1}

		index := rand.Int() % 7

		Background.Sprite.SetAnimationIndex(int(index))

		c.Transform().SetRotationf(rand.Float32() * 360)

		c.Transform().SetPosition(p)
		c.Transform().SetScalef(size, size)
	}

	for i := 0; i < 600; i++ {
		c := cookie.Clone()
		//c.Tag = CookieTag
		c.Transform().SetParent2(Layer2)
		size := 40 + rand.Float32()*100
		p := Engine.Vector{(rand.Float32() * 4000), (rand.Float32() * 4000), 1}

		if p.X < 1100 && p.Y < 800 {
			p.X += 1100
			p.Y += 800
		}

		c.Transform().SetPosition(p)
		c.Transform().SetScalef(size, size)
	}

	Wall = Engine.NewGameObject("Wall")
	Wall.Transform().SetParent2(Layer2)

	for i := 0; i < (4000/400)+2; i++ {
		c := staticCookie.Clone()
		c.Transform().SetParent2(Wall)
		p := Engine.Vector{float32(i) * 400, -200, 1}
		c.Transform().SetPosition(p)
		c.Transform().SetScalef(400, 400)
	}
	for i := 0; i < (4000/400)+2; i++ {
		c := staticCookie.Clone()
		c.Transform().SetParent2(Wall)
		p := Engine.Vector{float32(i) * 400, 4200, 1}
		c.Transform().SetPosition(p)
		c.Transform().SetScalef(400, 400)
	}
	for i := 0; i < (4000/400)+2; i++ {
		c := staticCookie.Clone()
		c.Transform().SetParent2(Wall)
		p := Engine.Vector{-200, float32(i) * 400, 1}
		c.Transform().SetPosition(p)
		c.Transform().SetScalef(400, 400)
	}
	for i := 0; i < (4000/400)+2; i++ {
		c := staticCookie.Clone()
		c.Transform().SetParent2(Wall)
		p := Engine.Vector{4200, float32(i) * 400, 1}
		c.Transform().SetPosition(p)
		c.Transform().SetScalef(400, 400)
	}

	s.AddGameObject(cam)
	s.AddGameObject(gui)
	s.AddGameObject(Layer1)
	s.AddGameObject(Layer2)
	s.AddGameObject(Layer3)
	s.AddGameObject(Layer4)
	s.AddGameObject(background)
	//s.AddGameObject(shadowShader)

	fmt.Println("Scene loaded")
}

func (s *GameScene) SceneBase() *Engine.SceneData {
	return s.SceneData
}

func (s *GameScene) New() Engine.Scene {
	gs := new(GameScene)
	gs.SceneData = Engine.NewScene("GameScene")
	return gs
}

package game

import (
	"fmt"
	"github.com/vova616/garageEngine/engine"
	"github.com/vova616/garageEngine/engine/components"
	_ "image/jpeg"
	_ "image/png"
	//"gl"  
	"strconv"
	"time"
	//"strings"
	//"math"
	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"
	"github.com/vova616/garageEngine/spaceCookies/server"
	//"image"
	//"image/color"
	"encoding/json"
	"math/rand"
	"os"
)

type GameScene struct {
	*engine.SceneData
	Layer1 *engine.GameObject
	Layer2 *engine.GameObject
	Layer3 *engine.GameObject
	Layer4 *engine.GameObject
}

var (
	GameSceneGeneral *GameScene
	cir              *engine.Texture
	boxt             *engine.Texture
	cookie           *engine.GameObject
	defender         *engine.GameObject
	missle           *Missle

	Player     *engine.GameObject
	PlayerShip *ShipController

	Explosion *engine.GameObject
	PowerUpGO *engine.GameObject

	Wall *engine.GameObject

	atlas        *engine.ManagedAtlas
	atlasSpace   *engine.ManagedAtlas
	atlasPowerUp *engine.ManagedAtlas
	backgroung   *engine.Texture
	ArialFont    *engine.Font
	ArialFont2   *engine.Font

	Players map[server.ID]*engine.GameObject = make(map[server.ID]*engine.GameObject)

	queenDead = false
)

const (
	MissleTag = "Missle"
	CookieTag = "Cookie"
)

var SpaceShip_A = "Ship"
var Explosion_ID engine.ID
var PowerUps_ID engine.ID

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
	atlas = engine.NewManagedAtlas(2048, 1024)
	atlasSpace = engine.NewManagedAtlas(1024, 1024)
	atlasPowerUp = engine.NewManagedAtlas(256, 256)

	var e error

	CheckError(atlas.LoadImageID("./data/spaceCookies/Ship1.png", SpaceShip_A))
	CheckError(atlas.LoadImageID("./data/spaceCookies/missile.png", Missle_A))
	e, Explosion_ID = atlas.LoadGroupSheet("./data/spaceCookies/Explosion.png", 128, 128, 6*8)
	CheckError(e)

	CheckError(atlas.LoadImageID("./data/spaceCookies/HealthBar.png", HP_A))
	CheckError(atlas.LoadImageID("./data/spaceCookies/HealthBarGUI.png", HPGUI_A))
	CheckError(atlas.LoadImageID("./data/spaceCookies/Queen.png", Queen_A))
	CheckError(atlas.LoadImageID("./data/spaceCookies/Jet.png", Jet_A))

	atlas.BuildAtlas()
	atlas.BuildMipmaps()
	atlas.SetFiltering(engine.MipMapLinearNearest, engine.Nearest)
	atlas.Texture.SetReadOnly()

	boxt, e = engine.LoadTexture("./data/spaceCookies/wall.png")

	boxt.BuildMipmaps()
	boxt.SetFiltering(engine.MipMapLinearNearest, engine.Nearest)

	backgroung, e = engine.LoadTexture("./data/spaceCookies/background.png")
	CheckError(e)
	cir, e = engine.LoadTexture("./data/spaceCookies/Cookie.png")
	CheckError(e)

	cir.BuildMipmaps()
	cir.SetFiltering(engine.MipMapLinearNearest, engine.Nearest)

	backgroung.BuildMipmaps()
	backgroung.SetFiltering(engine.MipMapLinearNearest, engine.Nearest)

	CheckError(atlasSpace.LoadGroup("./data/spaceCookies/Space/"))
	atlasSpace.BuildAtlas()
	atlasSpace.BuildMipmaps()
	atlasSpace.SetFiltering(engine.MipMapLinearNearest, engine.Nearest)
	atlasSpace.Texture.SetReadOnly()

	e, PowerUps_ID = atlasPowerUp.LoadGroupSheet("./data/spaceCookies/powerups.png", 61, 61, 3*4)
	CheckError(e)
	atlasPowerUp.BuildAtlas()
	atlasPowerUp.SetFiltering(engine.Linear, engine.Linear)

	ArialFont, e = engine.NewFont("./data/Fonts/arial.ttf", 48)
	if e != nil {
		panic(e)
	}
	ArialFont.Texture.SetReadOnly()

	ArialFont2, e = engine.NewFont("./data/Fonts/arial.ttf", 24)
	if e != nil {
		panic(e)
	}
	ArialFont2.Texture.SetReadOnly()
}

func SpawnMainPlayer(spawnPlayer server.SpawnPlayer) {
	Health := engine.NewGameObject("HP")
	Health.Transform().SetParent2(GameSceneGeneral.Camera.GameObject())
	Health.Transform().SetPositionf(150, 50)

	HealthGUI := engine.NewGameObject("HPGUI")
	HealthGUI.AddComponent(engine.NewSprite2(atlas.Texture, engine.IndexUV(atlas, HPGUI_A)))
	HealthGUI.Transform().SetParent2(Health)
	HealthGUI.Transform().SetPositionf(0, 0)
	HealthGUI.Transform().SetScalef(50, 50)

	HealthBar := engine.NewGameObject("HealthBar")
	HealthBar.Transform().SetParent2(Health)
	HealthBar.Transform().SetPositionf(-82, 0)
	HealthBar.Transform().SetScalef(100, 50)

	uvHP := engine.IndexUV(atlas, HP_A)

	HealthBarGUI := engine.NewGameObject("HealthBarGUI")
	HealthBarGUI.Transform().SetParent2(HealthBar)
	HealthBarGUI.AddComponent(engine.NewSprite2(atlas.Texture, uvHP))
	HealthBarGUI.Transform().SetScalef(0.52, 1)
	HealthBarGUI.Transform().SetPositionf((uvHP.Ratio/2)*HealthBarGUI.Transform().Scale().X, 0)

	JetFire := engine.NewGameObject("Jet")
	JetFire.AddComponent(engine.NewSprite2(atlas.Texture, engine.IndexUV(atlas, Jet_A)))

	Player.AddComponent(engine.NewSprite2(atlas.Texture, engine.IndexUV(atlas, SpaceShip_A)))
	PlayerShip = Player.AddComponent(NewShipController()).(*ShipController)
	Player.Transform().SetWorldPositionf(spawnPlayer.PlayerTransform.X, spawnPlayer.PlayerTransform.Y)
	Player.Transform().SetWorldRotationf(spawnPlayer.PlayerTransform.Rotation)
	Player.Transform().SetWorldScalef(100, 100)
	Player.AddComponent(components.NewSmoothFollow(nil, 2, 200))
	shipHP := float32(1000)
	PlayerShip.HPBar = HealthBar
	PlayerShip.JetFire = JetFire
	PlayerShip.Missle = missle
	Player.AddComponent(NewDestoyable(shipHP, 1))
}

func SpawnPlayer(spawnPlayer server.SpawnPlayer) {
	newPlayer, exists := Players[spawnPlayer.PlayerInfo.PlayerID]
	if !exists {
		newPlayer = engine.NewGameObject(spawnPlayer.PlayerInfo.Name)
	}
	newPlayer.Transform().SetParent2(GameSceneGeneral.Layer2)
	newPlayer.Transform().SetWorldPositionf(spawnPlayer.PlayerTransform.X, spawnPlayer.PlayerTransform.Y)
	newPlayer.Transform().SetWorldRotationf(spawnPlayer.PlayerTransform.Rotation)
	newPlayer.Transform().SetWorldScalef(100, 100)
	newPlayer.AddComponent(engine.NewSprite2(atlas.Texture, engine.IndexUV(atlas, SpaceShip_A)))

	if !exists {
		Players[spawnPlayer.PlayerInfo.PlayerID] = newPlayer
	}
}

func (s *GameScene) Load() {
	Players = make(map[server.ID]*engine.GameObject)
	LoadTextures()
	engine.SetTitle("Space Cookies")
	queenDead = false

	rand.Seed(time.Now().UnixNano())

	GameSceneGeneral = s

	s.Camera = engine.NewCamera()

	cam := engine.NewGameObject("Camera")
	cam.AddComponent(s.Camera)

	cam.Transform().SetScalef(1, 1)

	gui := engine.NewGameObject("GUI")

	Layer1 := engine.NewGameObject("Layer1")
	Layer2 := engine.NewGameObject("Layer2")
	Layer3 := engine.NewGameObject("Layer3")
	Layer4 := engine.NewGameObject("Layer3")

	s.Layer1 = Layer1
	s.Layer2 = Layer2
	s.Layer3 = Layer3
	s.Layer4 = Layer4

	mouse := engine.NewGameObject("Mouse")
	mouse.AddComponent(engine.NewMouse())
	mouse.AddComponent(NewMouseDebugger())
	mouse.Transform().SetParent2(cam)

	FPSDrawer := engine.NewGameObject("FPS")
	FPSDrawer.Transform().SetParent2(cam)
	txt := FPSDrawer.AddComponent(components.NewUIText(ArialFont2, "")).(*components.UIText)
	fps := FPSDrawer.AddComponent(engine.NewFPS()).(*engine.FPS)
	fps.SetAction(func(fps float64) {
		txt.SetString("FPS: " + strconv.FormatFloat(fps, 'f', 2, 32))
	})
	txt.SetAlign(engine.AlignLeft)

	FPSDrawer.Transform().SetPositionf(20, float32(engine.Height)-20)
	FPSDrawer.Transform().SetScalef(20, 20)
	/*
		label := engine.NewGameObject("Label")
		label.Transform().SetParent2(cam)
		label.Transform().SetPositionf(20, float32(engine.Height)-40)
		label.Transform().SetScalef(20, 20)

		txt2 := label.AddComponent(components.NewUIText(ArialFont2, "Input: ")).(*components.UIText)
		txt2.SetFocus(true)
		txt2.SetWritable(true)
		txt2.SetAlign(engine.AlignLeft)
	*/
	//SPACCCEEEEE
	engine.Space.Gravity.Y = 0
	engine.Space.Iterations = 10

	uvs, ind := engine.AnimatedGroupUVs(atlas, Explosion_ID)
	Explosion = engine.NewGameObject("Explosion")
	Explosion.AddComponent(engine.NewSprite3(atlas.Texture, uvs))
	Explosion.Sprite.BindAnimations(ind)
	Explosion.Sprite.AnimationSpeed = 25
	Explosion.Sprite.AnimationEndCallback = func(sprite *engine.Sprite) {
		sprite.GameObject().Destroy()
	}
	Explosion.Transform().SetScalef(30, 30)

	missleGameObject := engine.NewGameObject("Missle")
	missleGameObject.AddComponent(engine.NewSprite2(atlas.Texture, engine.IndexUV(atlas, Missle_A)))
	missleGameObject.AddComponent(engine.NewPhysics(false, 10, 10))
	missleGameObject.Transform().SetScalef(20, 20)
	missleGameObject.AddComponent(NewDamageDealer(50))

	missle = NewMissle(30000)
	missleGameObject.AddComponent(missle)
	missle.Explosion = Explosion
	ds := NewDestoyable(0, 1)
	ds.SetDestroyTime(1)
	missleGameObject.AddComponent(ds)

	ship := engine.NewGameObject("Ship")
	Player = ship
	Player.Transform().SetParent2(Layer2)
	Player.AddComponent(MyClient)
	/*
		Health := engine.NewGameObject("HP")
		Health.Transform().SetParent2(cam)
		Health.Transform().SetPositionf(150, 50)

		HealthGUI := engine.NewGameObject("HPGUI")
		HealthGUI.AddComponent(engine.NewSprite2(atlas.Texture, engine.IndexUV(atlas, HPGUI_A)))
		HealthGUI.Transform().SetParent2(Health)
		HealthGUI.Transform().SetScalef(50, 50)

		HealthBar := engine.NewGameObject("HealthBar")
		HealthBar.Transform().SetParent2(Health)
		HealthBar.Transform().SetPositionf(-82, 0)
		HealthBar.Transform().SetScalef(100, 50)

		uvHP := engine.IndexUV(atlas, HP_A)

		HealthBarGUI := engine.NewGameObject("HealthBarGUI")
		HealthBarGUI.Transform().SetParent2(HealthBar)
		HealthBarGUI.AddComponent(engine.NewSprite2(atlas.Texture, uvHP))
		HealthBarGUI.Transform().SetScalef(0.52, 1)
		HealthBarGUI.Transform().SetPositionf((uvHP.Ratio/2)*HealthBarGUI.Transform().Scale().X, 0)

		JetFire := engine.NewGameObject("Jet")
		JetFire.AddComponent(engine.NewSprite2(atlas.Texture, engine.IndexUV(atlas, Jet_A)))

		ship.AddComponent(engine.NewSprite2(atlas.Texture, engine.IndexUV(atlas, SpaceShip_A)))
		PlayerShip = ship.AddComponent(NewShipController()).(*ShipController)
		ship.Transform().SetParent2(Layer2)
		ship.Transform().SetPositionf(400, 200)
		ship.Transform().SetScalef(100, 100)
		shipHP := float32(1000)
		PlayerShip.HPBar = HealthBar
		PlayerShip.JetFire = JetFire
		PlayerShip.Missle = missle
		ship.AddComponent(NewDestoyable(shipHP, 1))
	*/

	cookie = engine.NewGameObject("Cookie")
	cookie.AddComponent(engine.NewSprite(cir))
	cookie.AddComponent(NewDestoyable(100, 2))
	cookie.AddComponent(NewDamageDealer(20))
	cookie.AddComponent(NewEnemeyAI(Player, Enemey_Cookie))
	cookie.Transform().SetScalef(50, 50)
	cookie.Transform().SetPositionf(400, 400)
	cookie.AddComponent(engine.NewPhysics2(false, chipmunk.NewCircle(vect.Vect{0, 0}, 25)))
	cookie.Tag = CookieTag

	defender = engine.NewGameObject("Box")
	ds = NewDestoyable(30, 3)
	ds.SetDestroyTime(5)
	defender.AddComponent(ds)
	defender.AddComponent(engine.NewSprite(boxt))
	defender.Tag = CookieTag
	defender.Transform().SetScalef(50, 50)

	phx := defender.AddComponent(engine.NewPhysics(false, 50, 50)).(*engine.Physics)
	phx.Body.SetMass(2.5)
	phx.Body.SetMoment(phx.Shape.Moment(2.5))
	phx.Shape.SetFriction(0.5)
	//phx.Shape.Group = 2
	phx.Shape.SetElasticity(0.5)

	QueenCookie := engine.NewGameObject("Cookie")
	QueenCookie.AddComponent(engine.NewSprite2(atlas.Texture, engine.IndexUV(atlas, Queen_A)))
	QueenCookie.AddComponent(NewDestoyable(5000, 2))
	QueenCookie.AddComponent(NewDamageDealer(200))
	//QueenCookie.AddComponent(NewEnemeyAI(Player, Enemey_Boss))
	//QueenCookie.Transform().SetParent2(Layer2)
	QueenCookie.Transform().SetScalef(300, 300)
	QueenCookie.Transform().SetPositionf(999999, 999999)
	QueenCookie.AddComponent(engine.NewPhysics2(false, chipmunk.NewCircle(vect.Vect{0, 0}, 25)))
	QueenCookie.Tag = CookieTag

	staticCookie := engine.NewGameObject("Cookie")
	staticCookie.AddComponent(engine.NewSprite(cir))
	staticCookie.Transform().SetScalef(400, 400)
	staticCookie.Transform().SetPositionf(400, 400)
	staticCookie.AddComponent(NewDestoyable(float32(engine.Inf), 2))
	staticCookie.AddComponent(engine.NewPhysics2(true, chipmunk.NewCircle(vect.Vect{0, 0}, 200)))

	staticCookie.Physics.Shape.SetElasticity(0)
	staticCookie.Physics.Body.SetMass(999999999999)
	staticCookie.Physics.Body.SetMoment(staticCookie.Physics.Shape.Moment(999999999999))
	staticCookie.Tag = CookieTag

	uvs, ind = engine.AnimatedGroupUVs(atlasSpace, "s")
	Background := engine.NewGameObject("Background")
	Background.AddComponent(engine.NewSprite3(atlasSpace.Texture, uvs))
	Background.Sprite.BindAnimations(ind)
	Background.Sprite.SetAnimation("s")
	Background.Sprite.AnimationSpeed = 0
	Background.Transform().SetScalef(50, 50)
	Background.Transform().SetPositionf(400, 400)

	uvs, ind = engine.AnimatedGroupUVs(atlasPowerUp, PowerUps_ID)
	PowerUpGO = engine.NewGameObject("Background")
	//PowerUpGO.Transform().SetParent2(Layer2)
	PowerUpGO.AddComponent(engine.NewSprite3(atlasPowerUp.Texture, uvs))
	PowerUpGO.AddComponent(engine.NewPhysics(false, 61, 61))
	PowerUpGO.Physics.Shape.IsSensor = true
	PowerUpGO.Sprite.BindAnimations(ind)
	PowerUpGO.Sprite.SetAnimation(PowerUps_ID)
	PowerUpGO.Sprite.AnimationSpeed = 0
	index := (rand.Int() % 6) + 6
	PowerUpGO.Sprite.SetAnimationIndex(int(index))
	PowerUpGO.Transform().SetScalef(61, 61)
	PowerUpGO.Transform().SetPositionf(0, 0)

	background := engine.NewGameObject("Background")
	background.AddComponent(engine.NewSprite(backgroung))
	background.AddComponent(NewBackground(background.Sprite))
	background.Sprite.Render = false
	//background.Transform().SetScalef(float32(backgroung.Height()), float32(backgroung.Height()), 1)
	background.Transform().SetScalef(800, 800)
	background.Transform().SetPositionf(0, 0)

	for i := 0; i < 300; i++ {
		c := Background.Clone()
		c.Transform().SetParent2(Layer4)
		size := 20 + rand.Float32()*50
		p := engine.Vector{(rand.Float32() * 5000) - 1000, (rand.Float32() * 5000) - 1000, 1}

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
		p := engine.Vector{(rand.Float32() * 4000), (rand.Float32() * 4000), 1}

		if p.X < 1100 && p.Y < 800 {
			p.X += 1100
			p.Y += 800
		}

		c.Transform().SetPosition(p)
		c.Transform().SetScalef(size, size)
	}

	Wall = engine.NewGameObject("Wall")
	Wall.Transform().SetParent2(Layer2)

	for i := 0; i < (4000/400)+2; i++ {
		c := staticCookie.Clone()
		c.Transform().SetParent2(Wall)
		p := engine.Vector{float32(i) * 400, -200, 1}
		c.Transform().SetPosition(p)
		c.Transform().SetScalef(400, 400)
	}
	for i := 0; i < (4000/400)+2; i++ {
		c := staticCookie.Clone()
		c.Transform().SetParent2(Wall)
		p := engine.Vector{float32(i) * 400, 4200, 1}
		c.Transform().SetPosition(p)
		c.Transform().SetScalef(400, 400)
	}
	for i := 0; i < (4000/400)+2; i++ {
		c := staticCookie.Clone()
		c.Transform().SetParent2(Wall)
		p := engine.Vector{-200, float32(i) * 400, 1}
		c.Transform().SetPosition(p)
		c.Transform().SetScalef(400, 400)
	}
	for i := 0; i < (4000/400)+2; i++ {
		c := staticCookie.Clone()
		c.Transform().SetParent2(Wall)
		p := engine.Vector{4200, float32(i) * 400, 1}
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

	fmt.Println("GameScene loaded")
}

func (s *GameScene) New() engine.Scene {
	gs := new(GameScene)
	gs.SceneData = engine.NewScene("GameScene")
	return gs
}

func (s *GameScene) OldLoad() {

	LoadTextures()

	queenDead = false

	rand.Seed(time.Now().UnixNano())

	GameSceneGeneral = s

	s.Camera = engine.NewCamera()

	cam := engine.NewGameObject("Camera")
	cam.AddComponent(s.Camera)

	cam.Transform().SetScalef(1, 1)

	gui := engine.NewGameObject("GUI")

	Layer1 := engine.NewGameObject("Layer1")
	Layer2 := engine.NewGameObject("Layer2")
	Layer3 := engine.NewGameObject("Layer3")
	Layer4 := engine.NewGameObject("Layer3")

	s.Layer1 = Layer1
	s.Layer2 = Layer2
	s.Layer3 = Layer3
	s.Layer4 = Layer4

	mouse := engine.NewGameObject("Mouse")
	mouse.AddComponent(engine.NewMouse())
	mouse.AddComponent(NewMouseDebugger())
	mouse.Transform().SetParent2(cam)

	FPSDrawer := engine.NewGameObject("FPS")
	FPSDrawer.Transform().SetParent2(cam)
	txt := FPSDrawer.AddComponent(components.NewUIText(ArialFont2, "")).(*components.UIText)
	fps := FPSDrawer.AddComponent(engine.NewFPS()).(*engine.FPS)
	fps.SetAction(func(fps float64) {
		txt.SetString("FPS: " + strconv.FormatFloat(fps, 'f', 2, 32))
	})
	txt.SetAlign(engine.AlignLeft)

	FPSDrawer.Transform().SetPositionf(20, float32(engine.Height)-20)
	FPSDrawer.Transform().SetScalef(20, 20)

	label := engine.NewGameObject("Label")
	label.Transform().SetParent2(cam)
	label.Transform().SetPositionf(20, float32(engine.Height)-40)
	label.Transform().SetScalef(20, 20)

	txt2 := label.AddComponent(components.NewUIText(ArialFont2, "Input: ")).(*components.UIText)
	txt2.SetFocus(true)
	txt2.SetWritable(true)
	txt2.SetAlign(engine.AlignLeft)

	//SPACCCEEEEE
	engine.Space.Gravity.Y = 0
	engine.Space.Iterations = 10

	Health := engine.NewGameObject("HP")
	Health.Transform().SetParent2(cam)
	Health.Transform().SetPositionf(150, 50)

	HealthGUI := engine.NewGameObject("HPGUI")
	HealthGUI.AddComponent(engine.NewSprite2(atlas.Texture, engine.IndexUV(atlas, HPGUI_A)))
	HealthGUI.Transform().SetParent2(Health)
	HealthGUI.Transform().SetScalef(50, 50)

	HealthBar := engine.NewGameObject("HealthBar")
	HealthBar.Transform().SetParent2(Health)
	HealthBar.Transform().SetPositionf(-82, 0)
	HealthBar.Transform().SetScalef(100, 50)

	uvHP := engine.IndexUV(atlas, HP_A)

	HealthBarGUI := engine.NewGameObject("HealthBarGUI")
	HealthBarGUI.Transform().SetParent2(HealthBar)
	HealthBarGUI.AddComponent(engine.NewSprite2(atlas.Texture, uvHP))
	HealthBarGUI.Transform().SetScalef(0.52, 1)
	HealthBarGUI.Transform().SetPositionf((uvHP.Ratio/2)*HealthBarGUI.Transform().Scale().X, 0)

	JetFire := engine.NewGameObject("Jet")
	JetFire.AddComponent(engine.NewSprite2(atlas.Texture, engine.IndexUV(atlas, Jet_A)))

	ship := engine.NewGameObject("Ship")
	Player = ship
	ship.AddComponent(engine.NewSprite2(atlas.Texture, engine.IndexUV(atlas, SpaceShip_A)))
	PlayerShip = ship.AddComponent(NewShipController()).(*ShipController)
	ship.AddComponent(components.NewSmoothFollow(nil, 0, 50))
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
		&engine.Debug,
	}

	f, e := os.Open("./data/spaceCookies/game.dat")
	if e != nil {
		f, e = os.Create("./data/spaceCookies/game.dat")
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

	uvs, ind := engine.AnimatedGroupUVs(atlas, Explosion_ID)
	Explosion = engine.NewGameObject("Explosion")
	Explosion.AddComponent(engine.NewSprite3(atlas.Texture, uvs))
	Explosion.Sprite.BindAnimations(ind)
	Explosion.Sprite.AnimationSpeed = 25
	Explosion.Sprite.AnimationEndCallback = func(sprite *engine.Sprite) {
		sprite.GameObject().Destroy()
	}
	Explosion.Transform().SetScalef(30, 30)

	missle := engine.NewGameObject("Missle")
	missle.AddComponent(engine.NewSprite2(atlas.Texture, engine.IndexUV(atlas, Missle_A)))
	missle.AddComponent(engine.NewPhysics(false, 10, 10))
	missle.Transform().SetScalef(20, 20)
	missle.AddComponent(NewDamageDealer(50))

	m := NewMissle(30000)
	missle.AddComponent(m)
	PlayerShip.Missle = m
	m.Explosion = Explosion
	ds := NewDestoyable(0, 1)
	ds.SetDestroyTime(1)
	missle.AddComponent(ds)

	cookie = engine.NewGameObject("Cookie")
	cookie.AddComponent(engine.NewSprite(cir))
	cookie.AddComponent(NewDestoyable(100, 2))
	cookie.AddComponent(NewDamageDealer(20))
	cookie.AddComponent(NewEnemeyAI(Player, Enemey_Cookie))
	cookie.Transform().SetScalef(50, 50)
	cookie.Transform().SetPositionf(400, 400)
	cookie.AddComponent(engine.NewPhysics2(false, chipmunk.NewCircle(vect.Vect{0, 0}, 25)))
	cookie.Tag = CookieTag

	defender = engine.NewGameObject("Box")
	ds = NewDestoyable(30, 3)
	ds.SetDestroyTime(5)
	defender.AddComponent(ds)
	defender.AddComponent(engine.NewSprite(boxt))
	defender.Tag = CookieTag
	defender.Transform().SetScalef(50, 50)

	phx := defender.AddComponent(engine.NewPhysics(false, 50, 50)).(*engine.Physics)
	phx.Body.SetMass(2.5)
	phx.Body.SetMoment(phx.Shape.Moment(2.5))
	phx.Shape.SetFriction(0.5)
	//phx.Shape.Group = 2
	phx.Shape.SetElasticity(0.5)

	QueenCookie := engine.NewGameObject("Cookie")
	QueenCookie.AddComponent(engine.NewSprite2(atlas.Texture, engine.IndexUV(atlas, Queen_A)))
	QueenCookie.AddComponent(NewDestoyable(5000, 2))
	QueenCookie.AddComponent(NewDamageDealer(200))
	QueenCookie.AddComponent(NewEnemeyAI(Player, Enemey_Boss))
	QueenCookie.Transform().SetParent2(Layer2)
	QueenCookie.Transform().SetScalef(300, 300)
	QueenCookie.Transform().SetPositionf(999999, 999999)
	QueenCookie.AddComponent(engine.NewPhysics2(false, chipmunk.NewCircle(vect.Vect{0, 0}, 25)))
	QueenCookie.Tag = CookieTag

	staticCookie := engine.NewGameObject("Cookie")
	staticCookie.AddComponent(engine.NewSprite(cir))
	staticCookie.Transform().SetScalef(400, 400)
	staticCookie.Transform().SetPositionf(400, 400)
	staticCookie.AddComponent(NewDestoyable(float32(engine.Inf), 2))
	staticCookie.AddComponent(engine.NewPhysics2(true, chipmunk.NewCircle(vect.Vect{0, 0}, 200)))

	staticCookie.Physics.Shape.SetElasticity(0)
	staticCookie.Physics.Body.SetMass(999999999999)
	staticCookie.Physics.Body.SetMoment(staticCookie.Physics.Shape.Moment(999999999999))
	staticCookie.Tag = CookieTag

	uvs, ind = engine.AnimatedGroupUVs(atlasSpace, "s")
	Background := engine.NewGameObject("Background")
	Background.AddComponent(engine.NewSprite3(atlasSpace.Texture, uvs))
	Background.Sprite.BindAnimations(ind)
	Background.Sprite.SetAnimation("s")
	Background.Sprite.AnimationSpeed = 0
	Background.Transform().SetScalef(50, 50)
	Background.Transform().SetPositionf(400, 400)

	uvs, ind = engine.AnimatedGroupUVs(atlasPowerUp, PowerUps_ID)
	PowerUpGO = engine.NewGameObject("Background")
	//PowerUpGO.Transform().SetParent2(Layer2)
	PowerUpGO.AddComponent(engine.NewSprite3(atlasPowerUp.Texture, uvs))
	PowerUpGO.AddComponent(engine.NewPhysics(false, 61, 61))
	PowerUpGO.Physics.Shape.IsSensor = true
	PowerUpGO.Sprite.BindAnimations(ind)
	PowerUpGO.Sprite.SetAnimation(PowerUps_ID)
	PowerUpGO.Sprite.AnimationSpeed = 0
	index := (rand.Int() % 6) + 6
	PowerUpGO.Sprite.SetAnimationIndex(int(index))
	PowerUpGO.Transform().SetScalef(61, 61)
	PowerUpGO.Transform().SetPositionf(0, 0)

	background := engine.NewGameObject("Background")
	background.AddComponent(engine.NewSprite(backgroung))
	background.AddComponent(NewBackground(background.Sprite))
	background.Sprite.Render = false
	//background.Transform().SetScalef(float32(backgroung.Height()), float32(backgroung.Height()), 1)
	background.Transform().SetScalef(800, 800)
	background.Transform().SetPositionf(0, 0)

	for i := 0; i < 300; i++ {
		c := Background.Clone()
		c.Transform().SetParent2(Layer4)
		size := 20 + rand.Float32()*50
		p := engine.Vector{(rand.Float32() * 5000) - 1000, (rand.Float32() * 5000) - 1000, 1}

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
		p := engine.Vector{(rand.Float32() * 4000), (rand.Float32() * 4000), 1}

		if p.X < 1100 && p.Y < 800 {
			p.X += 1100
			p.Y += 800
		}

		c.Transform().SetPosition(p)
		c.Transform().SetScalef(size, size)
	}

	Wall = engine.NewGameObject("Wall")
	Wall.Transform().SetParent2(Layer2)

	for i := 0; i < (4000/400)+2; i++ {
		c := staticCookie.Clone()
		c.Transform().SetParent2(Wall)
		p := engine.Vector{float32(i) * 400, -200, 1}
		c.Transform().SetPosition(p)
		c.Transform().SetScalef(400, 400)
	}
	for i := 0; i < (4000/400)+2; i++ {
		c := staticCookie.Clone()
		c.Transform().SetParent2(Wall)
		p := engine.Vector{float32(i) * 400, 4200, 1}
		c.Transform().SetPosition(p)
		c.Transform().SetScalef(400, 400)
	}
	for i := 0; i < (4000/400)+2; i++ {
		c := staticCookie.Clone()
		c.Transform().SetParent2(Wall)
		p := engine.Vector{-200, float32(i) * 400, 1}
		c.Transform().SetPosition(p)
		c.Transform().SetScalef(400, 400)
	}
	for i := 0; i < (4000/400)+2; i++ {
		c := staticCookie.Clone()
		c.Transform().SetParent2(Wall)
		p := engine.Vector{4200, float32(i) * 400, 1}
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

	fmt.Println("GameScene loaded")
}

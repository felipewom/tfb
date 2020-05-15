package handlers

import (
	"sort"

	"atreugo/src/storage"
	"atreugo/src/templates"

	"github.com/savsgio/atreugo/v11"
)

const helloWorldStr = "Hello, World!"

func queriesParam(ctx *atreugo.RequestCtx) int {
	n := ctx.QueryArgs().GetUintOrZero("queries")
	if n < 1 {
		n = 1
	} else if n > 500 {
		n = 500
	}

	return n
}

// JSONHandler . Test 1: JSON serialization
func JSONHandler(ctx *atreugo.RequestCtx) error {
	message := AcquireMessage()
	message.Message = helloWorldStr
	err := ctx.JSONResponse(message)

	ReleaseMessage(message)

	return err
}

// DBHandler . Test 2: Single database query
func DBHandler(db storage.DB) atreugo.View {
	return func(ctx *atreugo.RequestCtx) error {
		world := storage.AcquireWorld()
		db.GetOneRandomWorld(world)
		err := ctx.JSONResponse(world)

		storage.ReleaseWorld(world)

		return err
	}
}

// QueriesHandler . Test 3: Multiple database queries
func QueriesHandler(db storage.DB) atreugo.View {
	return func(ctx *atreugo.RequestCtx) error {
		queries := queriesParam(ctx)
		worlds := storage.AcquireWorlds()[:queries]

		for i := 0; i < queries; i++ {
			db.GetOneRandomWorld(&worlds[i])
		}

		err := ctx.JSONResponse(worlds)

		storage.ReleaseWorlds(worlds)

		return err
	}
}

// FortuneHandler . Test 4: Fortunes
func FortuneHandler(db storage.DB) atreugo.View {
	return func(ctx *atreugo.RequestCtx) error {
		fortunes, _ := db.GetFortunes()
		newFortune := templates.AcquireFortune()
		newFortune.Message = "Additional fortune added at request time."

		fortunes = append(fortunes, *newFortune)

		sort.Slice(fortunes, func(i, j int) bool {
			return fortunes[i].Message < fortunes[j].Message
		})

		ctx.SetContentType("text/html; charset=utf-8")

		err := templates.FortuneTemplate.Execute(ctx, fortunes)

		templates.ReleaseFortune(newFortune)
		templates.ReleaseFortunes(fortunes)

		return err
	}
}

// FortuneQuickHandler . Test 4: Fortunes
func FortuneQuickHandler(db storage.DB) atreugo.View {
	return func(ctx *atreugo.RequestCtx) error {
		fortunes, _ := db.GetFortunes()
		newFortune := templates.AcquireFortune()
		newFortune.Message = "Additional fortune added at request time."

		fortunes = append(fortunes, *newFortune)

		sort.Slice(fortunes, func(i, j int) bool {
			return fortunes[i].Message < fortunes[j].Message
		})

		ctx.SetContentType("text/html; charset=utf-8")
		templates.WriteFortunePage(ctx, fortunes)

		templates.ReleaseFortune(newFortune)
		templates.ReleaseFortunes(fortunes)

		return nil
	}
}

// UpdateHandler . Test 5: Database updates
func UpdateHandler(db storage.DB) atreugo.View {
	return func(ctx *atreugo.RequestCtx) error {
		queries := queriesParam(ctx)
		worlds := storage.AcquireWorlds()[:queries]

		for i := 0; i < queries; i++ {
			w := &worlds[i]
			db.GetOneRandomWorld(w)
			w.RandomNumber = int32(storage.RandomWorldNum())
		}

		db.UpdateWorlds(worlds)
		err := ctx.JSONResponse(worlds)

		storage.ReleaseWorlds(worlds)

		return err
	}
}

// PlaintextHandler . Test 6: Plaintext
func PlaintextHandler(ctx *atreugo.RequestCtx) error {
	ctx.WriteString(helloWorldStr)
	return nil
}
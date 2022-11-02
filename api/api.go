package api

import (
	"net/http"
	"os"

	"github.com/turbitcat/jsonote/v2/service"
	"github.com/turbitcat/jsonote/v2/wsgo"
)

const paramWriteKey = "writekey"
const paramReadKey = "readkey"
const paramNoteId = "id"

type JsonNoteApi struct {
	server *wsgo.ServerMux
	noter  *service.Noter
}

func (s *JsonNoteApi) Run() error {
	addr := os.Getenv("JSONOTE_ADDR")
	if addr == "" {
		addr = ":8088"
	}
	return s.server.Run(addr)
}

func New() *JsonNoteApi {
	noter := service.NewNoter("./data", "./data/history", true)
	server := defaultServer(noter)
	r := JsonNoteApi{server: server, noter: noter}
	return &r
}

func defaultServer(noter *service.Noter) *wsgo.ServerMux {
	getNoteParamById := func(c *wsgo.Context) {
		p := c.StringParams()
		id, ok := p[paramNoteId]
		if !ok {
			c.Log("empty id")
		}
		note, err := noter.GetNote(id)
		if err != nil {
			c.Log("err :%v", err)
		}
		c.AddParam("__note", note)
	}
	getNote := func(c *wsgo.Context) *service.Note {
		_note, ok := c.Param("__note")
		if !ok {
			getNoteParamById(c)
			_note, ok = c.Param("__note")
			if !ok {
				return nil
			}
		}
		note, ok := _note.(*service.Note)
		if !ok {
			return nil
		}
		return note
	}
	authWriteByNoteId := func(c *wsgo.Context) {
		note := getNote(c)
		if note == nil {
			c.String(http.StatusBadRequest, "Invalid id")
			return
		}
		if note.WriteKey != c.StringParams()[paramWriteKey] {
			c.String(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			c.Log("wrong write key")
			return
		}
		c.Next()
	}
	authReadByNoteId := func(c *wsgo.Context) {
		note := getNote(c)
		if note == nil {
			c.String(http.StatusBadRequest, "Invalid id")
			return
		}
		if note.ReadKey != c.StringParams()[paramReadKey] {
			c.String(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			c.Log("wrong read key")
			return
		}
		c.Next()
	}
	r := wsgo.Default()
	// /new [help writekey readkey] body: content
	r.POST("/new", func(c *wsgo.Context) {
		note := service.NewNote()
		help, ok := c.StringParam("help")
		if ok {
			note.Help = help
		}
		writeK, ok := c.StringParam(paramWriteKey)
		if ok {
			note.WriteKey = writeK
		}
		readK, ok := c.StringParam(paramReadKey)
		if ok {
			note.ReadKey = readK
		}
		err := c.BindJSON(&note.Content)
		if err != nil {
			c.Log("Err while parsing body: %v", err)
			c.String(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}
		if err := noter.SaveNote(note); err != nil {
			c.Log("Err while saving note %v: %v", note, err)
			c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		c.String(http.StatusOK, note.Id)
	})
	// /get [readkey]
	r.GET("/get", authReadByNoteId, func(c *wsgo.Context) {
		note := getNote(c)
		c.Json(http.StatusOK, note.Content)
	})
	// /save [writekey] body: content
	r.POST("/save", authWriteByNoteId, func(c *wsgo.Context) {
		note := getNote(c)
		var newContent any
		err := c.BindJSON(&newContent)
		if err != nil {
			c.Log("Err while parsing body")
			c.String(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}
		note.Content = newContent
		if err := noter.SaveNote(note); err != nil {
			c.Log("Err while saving note %v: %v", note, err)
			c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
	})
	return r
}

package controllers

import (
	"github.com/Cloakaac/cloak/models"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/Cloakaac/cloak/util"
	"github.com/julienschmidt/httprouter"
)

type CommunityController struct {
	*BaseController
}

// CharacterView shows a character
func (base *BaseController) CharacterView(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	name, err := url.QueryUnescape(p.ByName("name"))
	if err != nil {
		base.Error = "Invalid character name"
		return
	}
	player := models.GetPlayerByName(name)
	if player == nil {
		base.Redirect = "/"
		return
	}
	player.GetGuild()
	deaths, err := player.GetDeaths()
	if err != nil {
		base.Error = "Error while getting character deaths"
		return
	}
	base.Data["Info"] = player
	base.Data["Deaths"] = deaths
	base.Template = "character_view.html"
}

// SignatureView shows a signature
func (base *BaseController) SignatureView(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	name, err := url.QueryUnescape(p.ByName("name"))
	if err != nil {
		base.Error = "Invalid character name"
		return
	}
	player := models.GetPlayerByName(name)
	if player == nil {
		base.Error = "Unknown character name"
		return
	}
	signatureFile, err := os.Open(util.Parser.Template + "/public/signatures/" + player.Name + ".png")
	if err != nil { // No signature
		signature, err := util.CreateSignature(player.Name, player.Gender, player.Vocation, player.Level, player.LastLogin)
		if err != nil {
			base.Error = "Error while creating signature"
			return
		}
		w.Header().Set("Content-type", "image/png")
		w.Write(signature)
		return
	}
	defer signatureFile.Close()
	signatureFileStats, err := signatureFile.Stat()
	if err != nil {
		base.Error = "Error while reading signature stats"
		return
	}
	if signatureFileStats.ModTime().Unix()+(1*60) > time.Now().Unix() {
		buffer, err := ioutil.ReadAll(signatureFile)
		if err != nil {
			base.Error = "Error while reading signature bytes"
			return
		}
		w.Header().Set("Content-type", "image/png")
		w.Write(buffer)
		return
	}
	signature, err := util.CreateSignature(player.Name, player.Gender, player.Vocation, player.Level, player.LastLogin)
	if err != nil {
		base.Error = "Error while creating signature"
		return
	}
	w.Header().Set("Content-type", "image/png")
	w.Write(signature)
	return
}

// SearchCharacter searchs for names LIKE
func (base *BaseController) SearchCharacter(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	players, err := models.SearchPlayers(req.FormValue("name"))
	if err != nil {
		base.Error = "Error while searching for players"
		return
	}
	base.Data["Current"] = req.FormValue("name")
	base.Data["Characters"] = players
	base.Template = "character_search.html"
}

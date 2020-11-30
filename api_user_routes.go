package main

import (
	"alphaflow/models"
	"encoding/json"
	"io/ioutil"

	"github.com/attache/attache"
)

func (c *AlphaFlow) GET_APIUserList() {
	w := c.ResponseWriter()
	all, err := c.DB().All(new(models.User))
	if err != nil {
		attache.ErrorFatal(err)
	}
	attache.RenderJSON(w, all)
}

func (c *AlphaFlow) GET_APIUser() {
	w := c.ResponseWriter()
	r := c.Request()
	id := r.FormValue("id")
	var target models.User
	if err := c.DB().Get(&target, id); err != nil {
		if err == attache.ErrRecordNotFound {
			attache.Error(404)
		}
		attache.ErrorFatal(err)
	}
	attache.RenderJSON(w, target)
}

func (c *AlphaFlow) POST_APIUserNew() {
	w := c.ResponseWriter()
	r := c.Request()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		attache.ErrorFatal(err)
	}
	var target models.User
	if err := json.Unmarshal(body, &target); err != nil {
		attache.ErrorFatal(err)
	}
	if err := c.DB().Insert(&target); err != nil {
		attache.ErrorFatal(err)
	}
	w.WriteHeader(200)
}

func (c *AlphaFlow) POST_APIUser() {
	w := c.ResponseWriter()
	r := c.Request()
	id := r.FormValue("id")
	var target models.User
	if err := c.DB().Get(&target, id); err != nil {
		if err == attache.ErrRecordNotFound {
			attache.Error(404)
		}
		attache.ErrorFatal(err)
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		attache.ErrorFatal(err)
	}
	if err := json.Unmarshal(body, &target); err != nil {
		attache.ErrorFatal(err)
	}
	if err := c.DB().Update(&target); err != nil {
		attache.ErrorFatal(err)
	}
	w.WriteHeader(200)
}

func (c *AlphaFlow) DELETE_APIUser() {
	w := c.ResponseWriter()
	r := c.Request()
	id := r.FormValue("id")
	var target models.User
	if err := c.DB().Get(&target, id); err != nil {
		if err == attache.ErrRecordNotFound {
			w.WriteHeader(200)
			return
		}
		attache.ErrorFatal(err)
	}
	if err := c.DB().Delete(&target); err != nil {
		attache.ErrorFatal(err)
	}
	w.WriteHeader(200)
}

// Copyright 2015 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package example

import (
	"github.com/GoogleCloudPlatform/go-endpoints/endpoints"
	"github.com/luci/gae/impl/prod"
	dstore "github.com/luci/gae/service/datastore"
	"golang.org/x/net/context"
)

// CurrentValueReq describes the inputs to the CurrentValueReq RPC.
type CurrentValueReq struct {
	Name string `endpoints:"required"`
}

// CurrentValueRsp describes the outputs of the CurrentValueReq RPC.
type CurrentValueRsp struct {
	Val int64 `json:",string"`
}

// CurrentValue gets the current value of a counter (duh)
func (Example) CurrentValue(c context.Context, r *CurrentValueReq) (rsp *CurrentValueRsp, err error) {
	c = prod.Use(c)
	ds := dstore.Get(c)

	ctr := &Counter{Name: r.Name}
	if err = ds.Get(ctr); err != nil {
		return
	}

	rsp = &CurrentValueRsp{ctr.Val}
	return
}

func init() {
	mi["CurrentValue"] = &endpoints.MethodInfo{
		Path: "counter/{Name}",
		Desc: "Returns the current value held by the named counter",
	}
}

/**

    Plik upload server

The MIT License (MIT)

Copyright (c) <2015>
	- Mathieu Bodjikian <mathieu@bodjikian.fr>
	- Charles-Antoine Mathieu <skatkatt@root.gg>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
**/

package handlers

import (
	"net/http"

	"github.com/root-gg/juliet"
	"github.com/root-gg/plik/server/common"
	"github.com/root-gg/plik/server/data"
	"github.com/root-gg/plik/server/metadata"
)

// RemoveUpload create a new upload
func RemoveUpload(ctx *juliet.Context, resp http.ResponseWriter, req *http.Request) {
	log := common.GetLogger(ctx)

	// Get upload from context
	upload := common.GetUpload(ctx)
	if upload == nil {
		// This should never append
		log.Critical("Missing upload in removeUploadHandler")
		common.Fail(ctx, req, resp, "Internal error", 500)
		return
	}

	// Check authorization
	if !upload.Removable && !upload.IsAdmin {
		log.Warningf("Unable to remove upload : unauthorized")
		common.Fail(ctx, req, resp, "You are not allowed to remove this upload", 403)
		return
	}

	// Remove from data backend
	err := data.GetDataBackend().RemoveUpload(ctx, upload)
	if err != nil {
		log.Warningf("Unable to remove upload data : %s", err)
		common.Fail(ctx, req, resp, "Unable to remove upload", 500)
		return
	}

	// Remove from metadata backend
	err = metadata.GetMetaDataBackend().Remove(ctx, upload)
	if err != nil {
		log.Warningf("Unable to remove upload metadata : %s", err)
		common.Fail(ctx, req, resp, "Unable to remove upload metadata", 500)
	}
}

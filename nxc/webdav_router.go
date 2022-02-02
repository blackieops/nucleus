package nxc

import (
	"encoding/xml"
	"net/http"
	"strings"
	"time"
	"strconv"
	"io/ioutil"

	"com.blackieops.nucleus/data"
	"com.blackieops.nucleus/files"
	"com.blackieops.nucleus/webdav"
	"github.com/gin-gonic/gin"
)

type WebdavRouter struct {
	DBContext *data.Context
	Backend   files.StorageBackend
}

// A "fake" directory to be used as the user's root directory handle, as the
// "root" directory doesn't really exist but we still need to serialize it in
// some places.
var rootDirectory = files.Directory{
	Name:      "",
	FullName:  "",
	CreatedAt: time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
	UpdatedAt: time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
}

func (wr *WebdavRouter) HandlePropfind(c *gin.Context) {
	w := c.Writer
	w.WriteHeader(http.StatusMultiStatus)
	w.Header().Add("content-type", "application/xml; charset=utf-8")

	user, err := CurrentUser(wr.DBContext, c)
	if err != nil {
		panic(err)
	}

	opts, err := wr.buildPropfindOptionsFromRequest(c)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var dir *files.Directory
	dirPath := strings.Trim(c.Params.ByName("filePath"), "/")
	if dirPath == "" {
		// Optimize out query if we know the directory won't exist, as the root
		// path is just a placeholder and not in the DB.
		dir = nil
	} else {
		dir, err = files.FindDirByPath(wr.DBContext, user, dirPath)
		if err != nil {
			panic(err)
		}
	}

	composite, err := files.ListAll(wr.DBContext, user, opts.Depth, dir)
	if err != nil {
		panic(err)
	}

	response := webdav.BuildMultiResponse(user, composite, opts.Properties)

	x, err := xml.Marshal(response)
	if err != nil {
		panic(err)
	}
	w.Write(x)
}

func (wr *WebdavRouter) HandleGet(c *gin.Context) {
	user, err := CurrentUser(wr.DBContext, c)
	if err != nil {
		panic(err)
	}
	file, err := files.FindFileByPath(wr.DBContext, user, c.Params.ByName("filePath")[1:])
	if err != nil {
		panic(err)
	}
	fileBytes, err := wr.Backend.ReadFile(user, file)
	if err != nil {
		panic(err)
	}
	c.Header("etag", file.Digest)
	c.Writer.Write(fileBytes)
}

func (wr *WebdavRouter) buildPropfindOptionsFromRequest(c *gin.Context) (*webdav.PropfindOptions, error) {
	depth, err := strconv.Atoi(c.Request.Header.Get("Depth"))
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return nil, err
	}
	return webdav.BuildPropfindOptions(depth, body)
}

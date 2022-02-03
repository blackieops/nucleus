package nxc

import (
	"fmt"
	"encoding/xml"
	"net/http"
	"strings"
	"time"
	"strconv"
	"io/ioutil"
	"path/filepath"

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
			c.Status(http.StatusNotFound)
			return
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
		c.Status(http.StatusNotFound)
		return
	}
	fileBytes, err := wr.Backend.ReadFile(user, file)
	if err != nil {
		panic(err)
	}
	c.Header("etag", file.Digest)
	c.Writer.Write(fileBytes)
}

func (wr *WebdavRouter) HandlePut(c *gin.Context) {
	user, err := CurrentUser(wr.DBContext, c)
	if err != nil {
		panic(err)
	}
	filePath := c.Params.ByName("filePath")[1:]
	contents, err := ioutil.ReadAll(c.Request.Body)
	fileEntity := &files.File{
		Name:     filepath.Base(filePath),
		FullName: filePath,
		Size:     c.Request.ContentLength,
		User:     *user,
		Digest:   wr.Backend.FileDigest(user, contents),
	}
	if parentDir, err := files.FindDirByPath(wr.DBContext, user, filepath.Dir(filePath)); err == nil {
		fileEntity.Parent = parentDir
	}
	file, err := files.CreateFile(wr.DBContext, fileEntity)
	if err != nil {
		panic(err)
	}
	err = wr.Backend.WriteFile(user, file, contents)
	if err != nil {
		panic(err)
	}
	c.Header("etag", file.Digest)
	c.Status(http.StatusOK)
}

func (wr *WebdavRouter) HandleMkcol(c *gin.Context) {
	user, err := CurrentUser(wr.DBContext, c)
	if err != nil {
		panic(err)
	}
	filePath := c.Params.ByName("filePath")[1:]
	_, err = files.FindDirByPath(wr.DBContext, user, filePath)
	if err == nil {
		// Bail if the directory already exists to make this idempotent.
		c.Status(http.StatusOK)
		return
	}
	dir := &files.Directory{
		Name: filepath.Base(filePath),
		FullName: filePath,
		User: *user,
	}
	var parentDir *files.Directory
	parentDir, err = files.FindDirByPath(wr.DBContext, user, filepath.Dir(filePath))
	if err == nil {
		dir.Parent = parentDir
	}
	err = wr.Backend.CreateDirectory(user, dir)
	if err != nil {
		panic(err)
	}
	_, err = files.CreateDir(wr.DBContext, dir)
	if err != nil {
		c.Status(http.StatusUnprocessableEntity)
		return
	}
	c.Status(http.StatusCreated)
}

func (wr *WebdavRouter) HandleChunkMkcol(c *gin.Context) {
	user, err := CurrentUser(wr.DBContext, c)
	if err != nil {
		panic(err)
	}
	name := filepath.Base(c.Params.ByName("filePath")[1:])
	wr.Backend.CreateChunkDirectory(user, name)
	c.Status(http.StatusCreated)
}

func (wr *WebdavRouter) HandleChunkPut(c *gin.Context) {
	user, err := CurrentUser(wr.DBContext, c)
	if err != nil {
		panic(err)
	}
	contents, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		panic(err)
	}
	err = wr.Backend.WriteChunk(user, c.Params.ByName("filePath")[1:], contents)
	if err != nil {
		c.Status(http.StatusUnprocessableEntity)
		return
	}
}

func (wr *WebdavRouter) HandleChunkMove(c *gin.Context) {
	user, err := CurrentUser(wr.DBContext, c)
	if err != nil {
		panic(err)
	}
	dest := strings.TrimPrefix(
		c.Request.Header.Get("Destination"),
		"/nextcloud/remote.php/dav/files/"+user.Username+"/",
	)
	var dir *files.Directory
	if filepath.Dir(dest) == "." {
		dir = nil
	} else {
		dir, err = files.FindDirByPath(wr.DBContext, user, filepath.Dir(dest))
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
	}
	err = wr.Backend.ReconstructChunks(user, filepath.Dir(c.Params.ByName("filePath")[1:]), dest)
	if err != nil {
		fmt.Printf("Failed to reconstruct chunked upload: %v\n", err)
		c.Status(http.StatusUnprocessableEntity)
		return
	}
	crawler := &files.Crawler{DBContext: wr.DBContext, Backend: wr.Backend}
	file, err := crawler.DiscoverFile(user, dir, dest)
	if err != nil {
		fmt.Printf("Failed to index reconstructed chunked file: %v", err)
		c.Status(http.StatusUnprocessableEntity)
		return
	}
	file, err = files.CreateFile(wr.DBContext, file)
	if err != nil {
		fmt.Printf("Failed to save reconstructed chunked file in database: %v", err)
		c.Status(http.StatusUnprocessableEntity)
		return
	}
	c.Header("etag", file.Digest)
	c.Header("OC-FileID", fmt.Sprint(file.ID))
	c.Header("Location", "/nextcloud/remote.php/dav/files/"+user.Username+"/"+file.FullName)
	c.Status(http.StatusCreated)
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

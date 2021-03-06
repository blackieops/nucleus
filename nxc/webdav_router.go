package nxc

import (
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.b8s.dev/nucleus/data"
	"go.b8s.dev/nucleus/files"
	"go.b8s.dev/nucleus/webdav"
)

type WebdavRouter struct {
	DBContext  *data.Context
	Backend    files.StorageBackend
	Middleware *Middleware
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

func (wr *WebdavRouter) Mount(r *gin.RouterGroup) {
	r.Handle("PROPFIND", "/files/:username",
		wr.Middleware.EnsureAuthorization(), wr.handlePropfind)

	r.Handle("PROPFIND", "/files/:username/*filePath",
		wr.Middleware.EnsureAuthorization(), wr.handlePropfind)

	//r.Handle("PROPPATCH", "/remote.php/dav/files/:username/*filePath", forwardToWebdav)

	r.Handle("GET", "/files/:username/*filePath",
		wr.Middleware.EnsureAuthorization(), wr.handleGet)

	r.Handle("PUT", "/files/:username/*filePath",
		wr.Middleware.EnsureAuthorization(), wr.handlePut)

	r.Handle("MKCOL", "/files/:username/*filePath",
		wr.Middleware.EnsureAuthorization(), wr.handleMkcol)

	r.Handle("DELETE", "/files/:username/*filePath",
		wr.Middleware.EnsureAuthorization(), wr.handleDelete)

	r.Handle("MOVE", "/files/:username/*filePath",
		wr.Middleware.EnsureAuthorization(), wr.handleMove)

	// r.Handle("COPY", "/files/:username/*filePath",
	//	wr.Middleware.EnsureAuthorization(), wr.handleCopy)

	r.Handle("MKCOL", "/uploads/:username/*filePath",
		wr.Middleware.EnsureAuthorization(), wr.handleChunkMkcol)

	r.Handle("PUT", "/uploads/:username/*filePath",
		wr.Middleware.EnsureAuthorization(), wr.handleChunkPut)

	r.Handle("MOVE", "/uploads/:username/*filePath",
		wr.Middleware.EnsureAuthorization(), wr.handleChunkMove)
}

func (wr *WebdavRouter) handlePropfind(c *gin.Context) {
	w := c.Writer
	w.WriteHeader(http.StatusMultiStatus)
	w.Header().Add("content-type", "application/xml; charset=utf-8")

	user := wr.Middleware.GetCurrentUser(c)

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

func (wr *WebdavRouter) handleGet(c *gin.Context) {
	user := wr.Middleware.GetCurrentUser(c)
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

func (wr *WebdavRouter) handlePut(c *gin.Context) {
	user := wr.Middleware.GetCurrentUser(c)
	filePath := c.Params.ByName("filePath")[1:]
	contents, err := ioutil.ReadAll(c.Request.Body)
	digest := sha1.Sum(contents)
	fileEntity := &files.File{
		Name:     filepath.Base(filePath),
		FullName: filePath,
		Size:     c.Request.ContentLength,
		User:     *user,
		Digest:   fmt.Sprintf("%x", digest[:]),
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
	c.Header("oc-fileid", fmt.Sprint(file.ID))
	c.Status(http.StatusOK)
}

func (wr *WebdavRouter) handleMkcol(c *gin.Context) {
	user := wr.Middleware.GetCurrentUser(c)
	filePath := c.Params.ByName("filePath")[1:]
	_, err := files.FindDirByPath(wr.DBContext, user, filePath)
	if err == nil {
		// Bail if the directory already exists to make this idempotent.
		c.Status(http.StatusOK)
		return
	}
	dir := &files.Directory{
		Name:     filepath.Base(filePath),
		FullName: filePath,
		User:     *user,
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

func (wr *WebdavRouter) handleDelete(c *gin.Context) {
	user := wr.Middleware.GetCurrentUser(c)
	path := c.Params.ByName("filePath")[1:]
	err := files.DeletePath(wr.DBContext, user, path)
	if err != nil {
		fmt.Printf("Failed to delete files or directory from database: %v\n", err)
		c.Status(http.StatusUnprocessableEntity)
		return
	}
	err = wr.Backend.DeletePath(user, path)
	if err != nil {
		fmt.Printf("Failed to delete path: %v\n", err)
		c.Status(http.StatusUnprocessableEntity)
		return
	}
	c.Status(http.StatusNoContent)
}

func (wr *WebdavRouter) handleMove(c *gin.Context) {
	user := wr.Middleware.GetCurrentUser(c)
	dest := strings.TrimPrefix(
		c.Request.Header.Get("Destination"),
		"/nextcloud/remote.php/dav/files/"+user.Username+"/",
	)
	dest, err := url.QueryUnescape(dest)
	if err != nil {
		fmt.Printf("Error trying to decode path: %v", err)
		c.Status(http.StatusBadRequest)
		return
	}
	src := c.Params.ByName("filePath")[1:]
	err = wr.Backend.RenamePath(user, src, dest)
	if err != nil {
		fmt.Printf("[Webdav] MOVE error on backend: %v\n", err)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	err = files.RenamePath(wr.DBContext, user, src, dest)
	if err != nil {
		fmt.Printf("[Webdav] MOVE errored at database layer: %v", err)
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}
	// TODO: is it dangerous to interpolate user values here?
	c.Header("Location", "/nextcloud/remote.php/dav/files/"+user.Username+"/"+dest)
	c.Status(http.StatusCreated)
}

func (wr *WebdavRouter) handleChunkMkcol(c *gin.Context) {
	user := wr.Middleware.GetCurrentUser(c)
	name := filepath.Base(c.Params.ByName("filePath")[1:])
	wr.Backend.CreateChunkDirectory(user, name)
	c.Status(http.StatusCreated)
}

func (wr *WebdavRouter) handleChunkPut(c *gin.Context) {
	user := wr.Middleware.GetCurrentUser(c)
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

func (wr *WebdavRouter) handleChunkMove(c *gin.Context) {
	user := wr.Middleware.GetCurrentUser(c)
	dest := strings.TrimPrefix(
		c.Request.Header.Get("Destination"),
		"/nextcloud/remote.php/dav/files/"+user.Username+"/",
	)
	dest, err := url.QueryUnescape(dest)
	if err != nil {
		fmt.Printf("Error trying to decode path: %v", err)
		c.Status(http.StatusUnprocessableEntity)
		return
	}
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
	srcDir := filepath.Dir(c.Params.ByName("filePath")[1:])
	err = wr.Backend.ReconstructChunks(user, srcDir, dest)
	if err != nil {
		fmt.Printf("Failed to reconstruct chunked upload: %v\n", err)
		c.Status(http.StatusUnprocessableEntity)
		return
	}
	crawler := &files.Crawler{DBContext: wr.DBContext, Backend: wr.Backend}
	file, err := crawler.DiscoverFile(user, dir, filepath.Base(dest))
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
	err = wr.Backend.DeleteChunkDirectory(user, srcDir)
	if err != nil {
		fmt.Printf("[Webdav] Failed to remove temporary upload folder: %v\n", err)
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

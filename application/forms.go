package application

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"github.com/mholt/binding"
	"github.com/thoas/gostorages"
	"github.com/thoas/picfit/image"
	"io"
	"mime/multipart"
	"path"
	"strconv"
	"time"
)

type MultipartForm struct {
	Data *multipart.FileHeader `json:"data"`
}

func (f *MultipartForm) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&f.Data: "data",
	}
}

func (f *MultipartForm) Upload(storage gostorages.Storage) (*image.ImageFile, error) {
	var fh io.ReadCloser

	fh, err := f.Data.Open()

	if err != nil {
		return nil, err
	}

	defer fh.Close()

	dataBytes := bytes.Buffer{}

	_, err = dataBytes.ReadFrom(fh)

	if err != nil {
		return nil, err
	}
	//file storage path=>16进制(yMd)/16进制(hhm)
	//if not exists,create one
	now := time.Now()
	dir1, _ := strconv.ParseInt(now.Format("0612"), 10, 64)
	dir2, _ := strconv.ParseInt(now.Format("1504"), 10, 64)
	ext := strconv.FormatInt(dir1, 16) + "/" + strconv.FormatInt(dir2, 16) + "/"

	h := md5.New()
	h.Write(dataBytes.Bytes())
	filename := hex.EncodeToString(h.Sum(nil)) + path.Ext(f.Data.Filename)

	err = storage.Save(ext+filename, gostorages.NewContentFile(dataBytes.Bytes()))

	if err != nil {
		return nil, err
	}

	return &image.ImageFile{
		Filepath: filename,
		Storage:  storage,
	}, nil
}

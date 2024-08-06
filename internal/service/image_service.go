package service

import (
	"bytes"
	"github.com/gophero/goal/aws/s3"
	"github.com/gophero/goal/errorx"
	"github.com/gophero/goal/imagex"
	"golang.org/x/image/bmp"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"path"
	"strings"
)

type ImageService struct {
	*Service
	s3c *s3.Client
}

func NewImageService(s *Service, s3c *s3.Client) *ImageService {
	return &ImageService{
		Service: s,
		s3c:     s3c,
	}
}

func (srv *ImageService) UploadImage(userId uint, r io.Reader, fileName string, ip string, dir string, success func(url string) error, illegal func()) (string, bool, error) {
	if err := srv.s3c.UploadFile(srv.C.AwsS3().Bucket, dir+"/"+fileName, r); err != nil {
		return "", false, err
	}
	//uid := fmt.Sprintf("%v", userId)
	url := srv.C.AwsS3().FormatUrl(dir + "/" + fileName)
	//pass, err := green.ImageCheckerApp.Check(&green.ImageCheckParam{
	//	UserId:   uid,
	//	UserNick: "",
	//	Ip:       ip,
	//	DataId:   "", // 业务关联数据
	//	ImageUrl: url,
	//})
	//if err != nil {
	//	return "", false, err
	//}
	var pass = true
	if !pass {
		// 违规则删除文件
		go func() {
			defer func() {
				if err := recover(); err != nil {
					srv.L.Errorf("handle illegal image asynchronously error: %v", err)
				}
			}()
			illegal()
		}()
		return "", true, nil
	}
	return url, false, success(dir + "/" + fileName) // 只存储相对路径
}

func (srv *ImageService) CompressImage(file io.Reader, size uint, filename string) (r io.Reader) {
	bs, err := io.ReadAll(file)
	fr := bytes.NewReader(bs)
	srcr := bytes.NewReader(bs)
	errorx.Throw(err)

	w := size
	// var dstr io.Reader
	defer func() {
		if err := recover(); err != nil {
			srv.L.Errorf("compress avatar failed, filename: %s, use original image instead", filename)
			r = srcr // 出错了返回原始 reader，不能直接给 dstr 赋值，不能执行到 return
		}
	}()
	ext := strings.ToLower(path.Ext(filename))
	switch ext {
	case ".png":
		img, err := png.Decode(fr)
		errorx.Throw(err)
		img = imagex.Thumbnail(w, w, img)
		var buf bytes.Buffer
		errorx.Throw(png.Encode(&buf, img))
		r = bytes.NewReader(buf.Bytes())
	case ".gif":
		img, err := gif.Decode(fr)
		errorx.Throw(err)
		img = imagex.Thumbnail(w, w, img)
		var buf bytes.Buffer
		errorx.Throw(gif.Encode(&buf, img, nil))
		r = bytes.NewReader(buf.Bytes())
	case ".bmp":
		img, err := bmp.Decode(fr)
		errorx.Throw(err)
		img = imagex.Thumbnail(w, w, img)
		var buf bytes.Buffer
		errorx.Throw(bmp.Encode(&buf, img))
		r = bytes.NewReader(buf.Bytes())
	case ".jpg", ".jpeg": // jpg
		img, err := jpeg.Decode(fr)
		errorx.Throw(err)
		img = imagex.Thumbnail(w, w, img)
		var buf bytes.Buffer
		errorx.Throw(jpeg.Encode(&buf, img, nil))
		r = bytes.NewReader(buf.Bytes())
	default:
		r = srcr // 不支持压缩
	}
	// return dstr
	return r
}

package utils

import (
	"fmt"
	"goserver/pkg/infrastructure/config"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
)

type UploadDir struct {
	Thumb    string `json:"thumb"`
	Original string `json:"original"`
}

func GetUUID() (string, error) {
	uid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	id := strings.ReplaceAll(uid.String(), "-", "")

	return id, nil
}

func GetUploadDir(cfg config.FileConfig) *UploadDir {
	var dir = &UploadDir{}

	now := time.Now()
	folder := fmt.Sprintf("%d-%02d", now.Year(), now.Month())

	dir.Thumb = filepath.Join(folder, cfg.Thumb)
	dir.Original = filepath.Join(folder, cfg.Original)

	return dir
}

func ImageCompression(srcPath, dstPath, fileExt string) error {

	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return err
	}

	src, err := imaging.Open(srcPath, imaging.AutoOrientation(true))

	if err != nil {
		return err
	}

	/*
		滤镜|特点|速度|画质|适用场景
		imaging.Lanczos|重采样滤镜的“画质标杆”|较慢|最高，锐利清晰|对画质要求极高的场景，例如用户点开大图查看的预览。
		imaging.CatmullRom|高画质与速度的平衡点|快|很高，与 Lanczos 非常接近|绝大多数通用场景，如相册列表、缩略图等，是性能和画质的“甜点区”。
		imaging.MitchellNetravali|顺滑的立方滤镜|中等|高，比 CatmullRom 更平滑，振铃效应更少|适合需要柔和过渡，避免边缘过锐的图片。
		imaging.Linear|双线性插值滤镜|快|一般，画面较平滑|对速度有极致要求，且能接受轻微画质损失的场景。
		maging.Box|简单的平均滤镜|非常快|较低，用于缩小图片时效果尚可|追求极致速度，画质要求不高的快速处理。
		imaging.NearestNeighbor|最近邻插值，无抗锯齿|最快|最差，图像锯齿感强|基本不使用，除非是对速度有极端要求的特殊场景。
	*/
	dst := imaging.Resize(src, 1280, 720, imaging.NearestNeighbor)
	log.Println(fileExt)
	switch fileExt {
	case ".jpg", ".jpeg":
		// JPEG 格式：质量 65 获得较小体积（可调整）
		return imaging.Save(dst, dstPath, imaging.JPEGQuality(60))
	case ".png":
		/*
			级别|编码速度|文件大小|适用场景
			png.BestSpeed|最快|最大（压缩最少）|实时生成，对体积不敏感
			png.DefaultCompression|中等|中等|多数场景的默认选择
			png.BestCompression|很慢|最小（压缩最高）|离线优化，对体积要求严格
			png.NoCompression|极快|极大（几乎无压缩）|极少使用
		*/
		return imaging.Save(dst, dstPath, imaging.PNGCompressionLevel(png.DefaultCompression))
	default:
		// 其他格式（如 BMP, GIF）直接保存，无额外压缩选项
		return imaging.Save(dst, dstPath)
	}
}

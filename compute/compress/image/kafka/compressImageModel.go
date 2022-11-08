package kafkaCompressImage

type CompressImageModel struct {
	ImageName   string `valid:"required" json:"imageName"`
	ContentType string `valid:"required" json:"contentType"`
}

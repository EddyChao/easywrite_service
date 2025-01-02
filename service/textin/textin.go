package textin

import "github.com/gin-gonic/gin"

// Dewarp
// @Summary dewarp
// @Description dewarp
// @Tags Textin
// @Produce json
// @Param cookie header string true "Cookie"
// @Success 200
// @Router /textin/ai/service/v1/dewarp [post]
func Dewarp(c *gin.Context) {
	c.String(200, "ok")
}

// CropEnhanceImage
// @Summary Crop enhance image
// @Description Crop enhance image
// @Tags Textin
// @Produce json
// @Param cookie header string true "Cookie"
// @Success 200
// @Router /textin/ai/service/v1/crop_enhance_image [post]
func CropEnhanceImage(c *gin.Context) {
	c.String(200, "ok")
}

// BillsCrop
// @Summary Bills crop
// @Description Bills crop
// @Tags Textin
// @Produce json
// @Param cookie header string true "Cookie"
// @Success 200
// @Router /textin/robot/v1.0/api/bills_crop [post]
func BillsCrop(c *gin.Context) {
	c.String(200, "ok")
}

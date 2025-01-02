package bill

import (
	"easywrite-service/common"
	"easywrite-service/db"
	"easywrite-service/model"
	"easywrite-service/service"
	"easywrite-service/service/account"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/weili71/go-filex"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

var baseUrl = ""

func InitBaseUrl(baseUrl_ string) {
	baseUrl = baseUrl_
}

// GetBillHandler
// @Summary Get bills
// @Description Get a list of bills
// @Tags Bill
// @Accept json
// @Produce json
// @Param cookie header string true "Cookie"
// @Param limit query int false "Limit the number of results"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} model.JsonResponse[[]model.Bill]
// @Failure 400 {object} model.JsonResponse[any]
// @Failure 500 {object} model.JsonResponse[any]
// @Router /bill [get]
func GetBillHandler(c *gin.Context) {
	p := model.GetBillParameters{
		Limit:  -1,
		Offset: -1,
	}
	err := c.BindQuery(&p)
	if err != nil {
		service.HttpParameterError(c)
		return
	}
	logged, username := account.IsLoggedWithResponse(c)
	if !logged {
		return
	}
	var bills = make([]model.Bill, 0)
	stat := db.Mysql.Model(&model.Bill{}).Where("username = ?", username)
	if p.Offset >= 0 {
		stat = stat.Offset(p.Offset)
	}
	if p.Limit >= 0 {
		stat = stat.Limit(p.Limit)
	}
	err = stat.Find(&bills).Error
	if err != nil {
		service.HttpServerInternalError(c)
		return
	}
	for i := 0; i < len(bills); i++ {
		bills[i].ImagesComment = getFullImagesComment(bills[i].ImagesComment, baseUrl)
	}
	c.JSON(200, model.JsonResponse[[]model.Bill]{
		Code: 200,
		Msg:  "ok",
		Data: bills,
	})
}

func getURLPath(rawURL string) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}

	path1 := parsedURL.Path + "?" + parsedURL.RawQuery + "#" + parsedURL.Fragment
	return path1
}

// AddBillHandler
// @Summary Add a bill
// @Description Add a new bill
// @Tags Bill
// @Accept json
// @Produce json
// @Param cookie header string true "Cookie"
// @Security ApiKeyAuth
// @Param bill body model.Bill true "Bill object to be added"
// @Success 200 {object} model.JsonResponse[[]model.AddBillResponse]
// @Failure 400 {object} model.JsonResponse[any]
// @Failure 500 {object} model.JsonResponse[any]
// @Router /bill [post]
func AddBillHandler(c *gin.Context) {
	logged, username := account.IsLoggedWithResponse(c)
	if !logged {
		return
	}
	bill := model.Bill{}
	err := c.Bind(&bill)
	if err != nil {
		service.HttpParameterError(c)
		return
	}
	bill.Username = username
	bill.ID = int64(common.Snowflake.Generate())

	if bill.ThirdPartyID == "" {
		bill.ThirdPartyID = fmt.Sprintf("easywrite-%d", bill.ID)
	}

	for i, value := range bill.ImagesComment.Values {
		if strings.HasPrefix(value, "data:") {
			base64Data := getBase64FromDataURL(value)
			destPath := base64FileHandler(base64Data, "a.jpg")
			bill.ImagesComment.Values[i] = destPath
		} else {
			bill.ImagesComment.Values[i] = getURLPath(value)
		}
	}
	bill.ImagesComment = model.StringList{
		Values: lo.Filter(bill.ImagesComment.Values, func(item string, index int) bool {
			return item != ""
		}),
	}

	err = db.Mysql.Create(&bill).Error
	c.JSON(200, model.JsonResponse[[]model.AddBillResponse]{
		Code: 200,
		Msg:  "ok",
		Data: []model.AddBillResponse{
			{
				ID:            bill.ID,
				ThirdPartyId:  bill.ThirdPartyID,
				ImagesComment: getFullImagesComment(bill.ImagesComment, baseUrl),
			},
		},
	})
}

func getBase64FromDataURL(base64DataURL string) string {
	// 分割 Data URL
	parts := strings.SplitN(base64DataURL, ",", 2)
	if len(parts) != 2 {
		return ""
	}

	// 提取 Base64 部分
	base64Part := parts[1]
	return base64Part
}

func decodeBase64Stream(encodedString string) io.Reader {
	encodedReader := strings.NewReader(encodedString)
	decoder := base64.NewDecoder(base64.StdEncoding, encodedReader)
	return decoder
}

func decodeBase64(encodedString string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(encodedString)
	if err != nil {
		return "", err
	}

	decodedString := string(decodedBytes)
	return decodedString, nil
}

func base64FileHandler(base64Text string, filename string) string {

	relativePath := filepath.Join("/picture", uuid.NewString()+filepath.Ext(filename))

	f := filex.NewFile(filepath.Join(common.UploadSavePath, relativePath))
	dir := f.ParentFile()
	if !dir.IsExist() {
		err := dir.MkdirAll(0777)
		if err != nil {
			log.Panicln(err)
			return ""
		}
	}

	writer, _ := f.OpenFile(os.O_CREATE|os.O_RDWR, 0777)
	defer writer.Close()
	_, err := io.Copy(writer, decodeBase64Stream(base64Text))
	if err != nil {
		fmt.Println(err)
		return ""
	}

	relativePath = strings.ReplaceAll(relativePath, "\\", "/")
	return path.Join(common.UploadDir, relativePath)
}

func fileHandler(c *gin.Context, key string) []string {
	files := c.Request.MultipartForm.File[key]
	var picturePaths []string
	for _, file := range files {
		fmt.Println(file.Filename)
		relativePath := filepath.Join("/picture", uuid.NewString()+filepath.Ext(file.Filename))

		f := filex.NewFile(filepath.Join(common.UploadSavePath, relativePath))
		dir := f.ParentFile()
		if !dir.IsExist() {
			err := dir.MkdirAll(0666)
			if err != nil {
				log.Panicln(err)
				return nil
			}
		}
		destPath := f.Pathname
		if err := c.SaveUploadedFile(file, destPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return nil
		}
		relativePath = strings.ReplaceAll(relativePath, "\\", "/")
		picturePaths = append(picturePaths, path.Join(common.UploadDir, relativePath))
	}
	return picturePaths
}

// AddBillListHandler
// @Summary Add a list of bills
// @Description Add multiple bills at once
// @Tags Bill
// @Accept json
// @Produce json
// @Param cookie header string true "Cookie"
// @Security ApiKeyAuth
// @Param bills body []model.Bill true "List of bills to be added"
// @Success 200 {object} model.JsonResponse[[]model.AddBillResponse]
// @Failure 400 {object} model.JsonResponse[any]
// @Failure 500 {object} model.JsonResponse[any]
// @Router /bill/list [post]
func AddBillListHandler(c *gin.Context) {
	logged, username := account.IsLoggedWithResponse(c)
	if !logged {
		return
	}
	bills := make([]model.Bill, 0)
	err := c.Bind(&bills)
	if err != nil {
		service.HttpParameterError(c)
		return
	}

	for i := range bills {
		bills[i].Username = username
		bills[i].ID = int64(common.Snowflake.Generate())
		for j, value := range bills[i].ImagesComment.Values {
			if strings.HasPrefix(value, "data:") {
				base64Data := getBase64FromDataURL(value)
				destPath := base64FileHandler(base64Data, "a.jpeg")
				bills[i].ImagesComment.Values[j] = destPath
			} else {
				bills[i].ImagesComment.Values[i] = getURLPath(value)
			}
		}
	}

	ids := lo.Map(bills, func(bill model.Bill, index int) model.AddBillResponse {
		return model.AddBillResponse{
			ID:            bill.ID,
			ThirdPartyId:  bill.ThirdPartyID,
			ImagesComment: getFullImagesComment(bill.ImagesComment, baseUrl),
		}
	})

	err = db.Mysql.Create(&bills).Error
	c.JSON(200, model.JsonResponse[[]model.AddBillResponse]{
		Code: 200,
		Msg:  "ok",
		Data: ids,
	})
}

func getFullImagesComment(list model.StringList, baseUrl string) model.StringList {
	return model.StringList{
		Values: lo.Map(list.Values, func(item string, index int) string {
			return baseUrl + item
		}),
	}
}

// DeleteBillHandler
// @Summary Delete a bill
// @Description Delete a bill by ID
// @Tags Bill
// @Accept json
// @Produce json
// @Param cookie header string true "Cookie"
// @Security ApiKeyAuth
// @Param id query int true "ID of the bill to be deleted"
// @Success 200 {object} model.JsonResponse[any]
// @Failure 400 {object} model.JsonResponse[any]
// @Failure 500 {object} model.JsonResponse[any]
// @Router /bill [delete]
func DeleteBillHandler(c *gin.Context) {
	logged, username := account.IsLoggedWithResponse(c)
	if !logged {
		return
	}
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil {
		service.HttpParameterError(c)
		return
	}
	err = db.Mysql.Delete(&model.Bill{}, map[string]any{"username": username, "id": id}).Error
	if err != nil {
		service.HttpServerInternalError(c)
		return
	}
	c.JSON(200, model.JsonResponse[any]{
		Code: 200,
		Msg:  "ok",
		Data: nil,
	})
}

// UpdateBillHandler
// @Summary Update a bill
// @Description Update an existing bill
// @Tags Bill
// @Accept json
// @Produce json
// @Param cookie header string true "Cookie"
// @Security ApiKeyAuth
// @Param bill body model.Bill true "Updated bill object"
// @Success 200 {object} model.JsonResponse[[]model.AddBillResponse]
// @Failure 400 {object} model.JsonResponse[any]
// @Failure 500 {object} model.JsonResponse[any]
// @Router /bill [put]
func UpdateBillHandler(c *gin.Context) {
	logged, username := account.IsLoggedWithResponse(c)
	if !logged {
		return
	}
	bill := model.Bill{}
	err := c.Bind(&bill)
	if err != nil || bill.ID == 0 {
		service.HttpParameterError(c)
		return
	}
	bill.Username = username
	for i, value := range bill.ImagesComment.Values {
		if strings.HasPrefix(value, "data:") {
			base64Data := getBase64FromDataURL(value)
			destPath := base64FileHandler(base64Data, "a.jpeg")
			bill.ImagesComment.Values[i] = destPath
		} else {
			bill.ImagesComment.Values[i] = getURLPath(value)
		}
	}
	err = db.Mysql.Model(&model.Bill{}).Where("username = ? AND id = ?", username, bill.ID).Updates(&bill).Error
	if err != nil {
		service.HttpServerInternalError(c)
		return
	}
	c.JSON(200, model.JsonResponse[[]model.AddBillResponse]{
		Code: 200,
		Msg:  "ok",
		Data: []model.AddBillResponse{
			{
				ID:            bill.ID,
				ThirdPartyId:  bill.ThirdPartyID,
				ImagesComment: getFullImagesComment(bill.ImagesComment, baseUrl),
			},
		},
	})
}

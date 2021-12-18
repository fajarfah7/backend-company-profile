package delivery

import (
	"project_company_profile/models/user"

	"net/http"

	"github.com/gin-gonic/gin"
)

type Delivery struct {
	Ucs user.TextUsecase
}

func NewTextDelivery(ucs user.TextUsecase) user.TextDelivery {
	return &Delivery{
		Ucs: ucs,
	}
}

func (dlv *Delivery) Texts(c *gin.Context) {
	type binderTexts struct {
		TextType string `json:"text_type" form:"text_type"`
	}

	binder := new(binderTexts)
	if binder.TextType == "" {
		messages := []map[string]string{{"api": "Text type can not be empty"}}
		c.JSON(http.StatusBadRequest, gin.H{
			"data": map[string][]map[string]string{
				"messages": messages,
			},
		})
	}
	res := dlv.Ucs.Texts(binder.TextType)
	c.JSON(res.Status, res)
}

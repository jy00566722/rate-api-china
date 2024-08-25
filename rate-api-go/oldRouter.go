package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func LoadV1Router(e *gin.Engine) {
	var aliexpress_element = make([][]string, 1)
	aliexpress_element[0] = []string{".U9mS2 ._2FkhA", "style", "overflow:visible;"}

	var aliexpress_nodes = make([][]string, 20)
	aliexpress_nodes[0] = []string{"span", "class", "SearchProductFeed_Price__titleWrapper__1jg3h"}
	aliexpress_nodes[1] = []string{"div", "class", "_12A8D"}
	aliexpress_nodes[2] = []string{"span", "class", "_1AHOM"}
	aliexpress_nodes[3] = []string{"span", "class", "price-current"}
	aliexpress_nodes[4] = []string{"span", "class", "product-price-value"}
	aliexpress_nodes[5] = []string{"div", "class", "simple-card-price"}
	aliexpress_nodes[6] = []string{"span", "class", "_2FkhA"}
	aliexpress_nodes[7] = []string{"span", "class", "text"}
	aliexpress_nodes[8] = []string{"div", "class", "_335Dg"}
	aliexpress_nodes[9] = []string{"span", "class", "styles_price__1VOmb styles_descriptionItem__eRonF undefined"}
	aliexpress_nodes[10] = []string{"div", "class", "current-price"}
	aliexpress_nodes[11] = []string{"span", "class", "styles_price__xeG3I styles_descriptionItem__39Djx styles_italicPrice__1dKZ4 undefined"}
	aliexpress_nodes[12] = []string{"span", "class", "ali-kit_Base__base__1odrub ali-kit_Base__default__1odrub ali-kit_Base__strong__1odrub price ali-kit_Price__size-xl__12ybyf Product_Price__current__1uqb8 product-price-current"}
	aliexpress_nodes[13] = []string{"span", "class", "uniform-banner-box-price"}
	aliexpress_nodes[14] = []string{"span", "class", "_1PVo-"}
	aliexpress_nodes[15] = []string{"span", "class", "ali-kit_Base__base__1odrub ali-kit_Base__default__1odrub ali-kit_Base__strong__1odrub price ali-kit_Price__size-m__12ybyf"}
	aliexpress_nodes[16] = []string{"span", "class", "ali-kit_Base__base__1odrub ali-kit_Base__default__1odrub ali-kit_Base__strong__1odrub price ali-kit_Price__size-s__12ybyf"}
	aliexpress_nodes[17] = []string{"b", "class", "notranslate"}
	aliexpress_nodes[18] = []string{"span", "class", "notranslate"}
	aliexpress_nodes[19] = []string{"span", "class", "rax-text "}

	var aliexpress_special_price = make([]string, 4)
	aliexpress_special_price[0] = `div[class="_13_ga"`
	aliexpress_special_price[1] = `div[class="_3kmwU"]`
	aliexpress_special_price[2] = `div[class="_13_ga _37W_B"]`
	aliexpress_special_price[3] = `div[class="mGXnE _37W_B"]`

	var amazon_nodes = make([][]string, 13)
	amazon_nodes[0] = []string{"span", "class", "a-color-price _p13n-desktop-sims-fbt_fbt-desktop_total-amount__wLVdU"}
	amazon_nodes[1] = []string{"span", "class", "_p13n-desktop-sims-fbt_price_p13n-sc-price__bCZQt"}
	amazon_nodes[2] = []string{"span", "class", "_p13n-zg-list-carousel-desktop_price_p13n-sc-price__3mJ9Z"}
	amazon_nodes[3] = []string{"span", "aria-hidden", "true"}
	amazon_nodes[4] = []string{"span", "id", "priceblock_ourprice"}
	amazon_nodes[5] = []string{"span", "id", "priceblock_saleprice"}
	amazon_nodes[6] = []string{"span", "id", "priceblock_dealprice"}
	amazon_nodes[7] = []string{"span", "class", "a-size-medium a-color-price"}
	amazon_nodes[8] = []string{"span", "class", "p13n-sc-price"}
	amazon_nodes[9] = []string{"span", "class", "a-color-price a-text-bold"}
	amazon_nodes[10] = []string{"p", "class", "a-spacing-none a-text-left a-size-mini twisterSwatchPrice"}
	amazon_nodes[11] = []string{"span", "class", "netPrice--2L7XO"}
	amazon_nodes[12] = []string{"span", "class", "a-color-price"}
	GetConfig()
	feedback_flag := viper.GetInt("feedback")

	var lazada_main_price_nodes = make([][]string, 1)
	lazada_main_price_nodes[0] = []string{"span[class=\"pdp-price pdp-price_type_normal pdp-price_color_orange pdp-price_size_xl\"]"}

	var lazada_nodes = make([][]string, 21)
	lazada_nodes[0] = []string{"span", "price"}
	lazada_nodes[1] = []string{"span", "super-deals-item-sale-price-value"}
	lazada_nodes[2] = []string{"span", "global-brand-item-price-value"}
	lazada_nodes[3] = []string{"span", "best-seller-item-price-value"}
	lazada_nodes[4] = []string{"div", "sale-price"}
	lazada_nodes[5] = []string{"div", "store-product-price"}
	lazada_nodes[6] = []string{"p", "product-price"}
	lazada_nodes[7] = []string{"div", "product-price"}
	lazada_nodes[8] = []string{"span", "c13VH6"}
	lazada_nodes[9] = []string{"span", "pswt-product-price"}
	lazada_nodes[10] = []string{"div", "pswt-product-price"}
	lazada_nodes[11] = []string{"div", "item-discount-price"}
	lazada_nodes[12] = []string{"div", "delivery-option-item__shipping-fee"}
	lazada_nodes[13] = []string{"div", "discount-price"}
	lazada_nodes[14] = []string{"span", "price-label price-label-prim"}
	lazada_nodes[15] = []string{"div", "p-slider-product-price"}
	lazada_nodes[16] = []string{"span", "hotdeal-product3-item-price-discount"}
	lazada_nodes[17] = []string{"span", "text"}
	lazada_nodes[18] = []string{"div", "voucher-tile-discount-value-text"}
	lazada_nodes[19] = []string{"span", "index__currency___Q78Jz"}
	lazada_nodes[20] = []string{"span", "ooOxS"}

	var ozon_nodes = make([]string, 5)
	ozon_nodes[0] = "span[class=\"ui-p9 ui-q1 ui-q4\"]"
	ozon_nodes[1] = "span[class=\"c2h5 c2h6\"] span"
	ozon_nodes[2] = "span[class=\"c2h5\"] span"
	ozon_nodes[3] = "span[class=\"ui-p9 ui-q1\"]"
	ozon_nodes[4] = "span[class=\"c2h5 c2h7\"] span"

	var qoo10jp_nodes = make([]string, 4)
	qoo10jp_nodes[0] = "strong"
	qoo10jp_nodes[1] = "strong[class=\"price\"]"
	qoo10jp_nodes[2] = "div[class=\"info\"] div[class=\"prc\"]"
	qoo10jp_nodes[3] = "span[class=\"prc\"]"

	var shopee_nodes = make([][]string, 10)
	shopee_nodes[0] = []string{"span", "_2wRB2a"}
	shopee_nodes[1] = []string{"div", "pmmxKx"}
	shopee_nodes[2] = []string{"span", "_3c5u7X"}
	shopee_nodes[3] = []string{"span", "_4XNe7q"}
	shopee_nodes[4] = []string{"span", "j0vBz2"}
	shopee_nodes[5] = []string{"div", "item-card-special__current-price item-card-special__current-price--special"}
	shopee_nodes[6] = []string{"div", "item-card-special__current-price item-card-special__current-price--special item-card-special__current-price--ofs"}
	shopee_nodes[7] = []string{"span", "IL6Kt2 _24JoLh"}
	shopee_nodes[8] = []string{"div", "shopee-item-card__current-price"}
	shopee_nodes[9] = []string{"span", "ZEgDH9"}

	e.GET("/v1/shopee_nodes", func(c *gin.Context) {
		c.JSON(200, gin.H{"code": 20000, "nodes": shopee_nodes})
	})
	e.GET("/v1/amazon_nodes", func(c *gin.Context) {
		c.JSON(200, gin.H{"code": 20000, "nodes": amazon_nodes})
	})
	e.GET("/v1/aliexpress_nodes", func(c *gin.Context) {
		c.JSON(200, gin.H{"code": 20000, "nodes": aliexpress_nodes})
	})
	e.GET("/v1/qoo10jp_nodes", func(c *gin.Context) {
		c.JSON(200, gin.H{"code": 20000, "nodes": qoo10jp_nodes})
	})
	e.GET("/v1/aliexpress_element", func(c *gin.Context) {
		c.JSON(200, gin.H{"code": 20000, "nodes": aliexpress_element})
	})
	e.GET("/v1/aliexpress_special_price", func(c *gin.Context) {
		c.JSON(200, gin.H{"code": 20000, "nodes": aliexpress_special_price})
	})
	e.GET("/v1/lazada_nodes", func(c *gin.Context) {
		c.JSON(200, gin.H{"code": 20000, "nodes": lazada_nodes})
	})
	e.GET("/v1/feedback_flag", func(c *gin.Context) {
		c.JSON(200, gin.H{"code": 20000, "feedback_flag": feedback_flag})
	})
	e.POST("/v1/element_feedback", func(c *gin.Context) {
		c.String(200, "ok")
	})

}

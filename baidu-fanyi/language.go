package main

// LanguageItem represents an item in the language array.
type LanguageItem struct {
	Code         string `json:"code"`
	Bd           string `json:"bd,omitempty"`
	ChineseNawme string `json:"chineseNawme,omitempty"`
}

var languageList = []LanguageItem{
	{Code: "bg", Bd: "bul", ChineseNawme: "保加利亚语"},
	{Code: "bn", Bd: "ben", ChineseNawme: "孟加拉语"},
	{Code: "cs", Bd: "cs", ChineseNawme: "捷克语"},
	{Code: "da", Bd: "dan", ChineseNawme: "丹麦语"},
	{Code: "de", Bd: "de", ChineseNawme: "德语"},
	{Code: "el", Bd: "el", ChineseNawme: "希腊语"},
	{Code: "en", Bd: "en", ChineseNawme: "英语"},
	{Code: "en_AU", Bd: "en", ChineseNawme: "英语（澳大利亚）"},
	{Code: "en_GB", Bd: "en", ChineseNawme: "英语（英国）"},
	{Code: "en_US", Bd: "en", ChineseNawme: "英语（美国）"},
	{Code: "es", Bd: "spa", ChineseNawme: "西班牙语"},
	{Code: "et", Bd: "est", ChineseNawme: "爱沙尼亚语"},
	{Code: "fa", Bd: "per", ChineseNawme: "波斯语"},
	{Code: "fi", Bd: "fin", ChineseNawme: "芬兰语"},
	{Code: "fil", Bd: "fil", ChineseNawme: "菲律宾语"},
	{Code: "fr", Bd: "fra", ChineseNawme: "法语"},
	{Code: "gu", Bd: "guj", ChineseNawme: "古吉拉特语"},
	{Code: "he", Bd: "heb", ChineseNawme: "希伯来语"},
	{Code: "hi", Bd: "hi", ChineseNawme: "印地语"},
	{Code: "hr", Bd: "hrv", ChineseNawme: "克罗地亚语"},
	{Code: "hu", Bd: "hu", ChineseNawme: "匈牙利语"},
	{Code: "id", Bd: "id", ChineseNawme: "印度尼西亚语"},
	{Code: "it", Bd: "it", ChineseNawme: "意大利语"},
	{Code: "ja", Bd: "jp", ChineseNawme: "日语"},
	{Code: "kn", Bd: "kan", ChineseNawme: "卡纳达语"},
	{Code: "ko", Bd: "kor", ChineseNawme: "韩语"},
	{Code: "lt", Bd: "lit", ChineseNawme: "立陶宛语"},
	{Code: "lv", Bd: "lav", ChineseNawme: "拉脱维亚语"},
	{Code: "ms", Bd: "may", ChineseNawme: "马来语"},
	{Code: "nl", Bd: "nl", ChineseNawme: "荷兰语"},
	{Code: "no", Bd: "nor", ChineseNawme: "挪威语"},
	{Code: "pl", Bd: "pl", ChineseNawme: "波兰语"},
	{Code: "pt_PT", Bd: "pt", ChineseNawme: "葡萄牙语"},
	{Code: "ro", Bd: "rom", ChineseNawme: "罗马尼亚语"},
	{Code: "ru", Bd: "ru", ChineseNawme: "俄语"},
	{Code: "sk", Bd: "sk", ChineseNawme: "斯洛伐克语"},
	{Code: "sl", Bd: "slo", ChineseNawme: "斯洛文尼亚语"},
	{Code: "sr", Bd: "srp", ChineseNawme: "塞尔维亚语"},
	{Code: "sv", Bd: "swe", ChineseNawme: "瑞典语"},
	{Code: "sw", Bd: "swa", ChineseNawme: "斯瓦希里语"},
	{Code: "ta", Bd: "tam", ChineseNawme: "泰米尔语"},
	{Code: "th", Bd: "th", ChineseNawme: "泰语"},
	{Code: "tr", Bd: "tr", ChineseNawme: "土耳其语"},
	{Code: "uk", Bd: "ukr", ChineseNawme: "乌克兰语"},
	{Code: "vi", Bd: "vie", ChineseNawme: "越南语"},
	{Code: "zh_CN", Bd: "zh", ChineseNawme: "中文（简体）"},
	{Code: "zh_TW", Bd: "cht", ChineseNawme: "中文（繁体）"},
}

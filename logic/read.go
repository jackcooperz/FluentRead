package logic

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"FluentRead/misc/log"
	"FluentRead/models"
	"FluentRead/repo/cache"
	"FluentRead/repo/db"
	"FluentRead/utils"
)

const (
	pageCacheConst = "page:%s" // page:link
	shortTimeOut   = 30 * time.Second
	longTimeOut    = 72 * time.Hour
)

// PageRead 按域名全文读取
func PageRead(ctx *gin.Context, link string) (transMap map[string]string, err error) {

	link = utils.GetHostByString(link)

	// 1、查缓存，获取失败不影响后续流程
	key := fmt.Sprintf(pageCacheConst, link)
	transString, err := cache.GetKey(ctx, key)
	if err != nil {
		log.Warnf("缓存读取失败：%v", err)
	}
	if len(transString) > 0 {
		err = models.StringToMap(transString, &transMap)
		if err != nil {
			log.Errorf("缓存数据转换失败：%v", err)
			return nil, err
		}
		return transMap, nil
	}

	// 2、查数据库
	transs, err := db.GetAllByPage(ctx, link)
	if err != nil {
		log.Errorf("全文读取失败：%v", err)
		return nil, err
	}
	transMap = models.BatchTransToMap(transs)

	// 3、写缓存（+防缓存穿透）	TODO 12.15 发送到 MQ，统计网站的请求次数，用以判断是否需要翻译
	timeout := shortTimeOut
	if len(transs) > 0 {
		timeout = longTimeOut
	}
	err = cache.SetKeyWithTimeout(ctx, key, timeout, models.MapToString(transMap))
	if err != nil {
		log.Warnf("缓存写入失败：%v", err)
	}

	return transMap, nil
}

// BatchRead 批量读取
func BatchRead(ctx *gin.Context, hashs []string) (map[string]string, error) {

	transModels, err := db.BatchGet(ctx, hashs)
	if err != nil {
		log.Errorf("批量读取失败：%v", err)
		return nil, err
	}
	return models.BatchTransToMap(transModels), nil
}

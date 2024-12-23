package controllers

import (
	"exchangeapp/global"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

// LikeArticle 增加文章的点赞数。
// @Summary 点赞文章
// @Description 根据文章ID，为指定文章增加一次点赞数。
// @Tags 文章操作
// @Param id path string true "文章ID"
// @Router /articles/{id}/like [post]
func LikeArticle(ctx *gin.Context) {
	// 获取文章ID
	articleID := ctx.Param("id")

	// 生成 Redis 中存储点赞数的键名
	likeKey := "article:" + articleID + ":likes"

	// 使用 Redis 的 Incr 方法为点赞数加一
	if err := global.RedisDB.Incr(likeKey).Err(); err != nil {
		// 如果 Redis 操作失败，返回内部服务器错误
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully liked the article"})
}

// GetArticleLikes 获取文章的点赞数。
// @Summary 获取文章点赞数
// @Description 根据文章ID，获取指定文章的点赞数。
// @Tags 文章操作
// @Param id path string true "文章ID"
// @Router /articles/{id}/likes [get]
func GetArticleLikes(ctx *gin.Context) {
	// 获取文章ID
	articleID := ctx.Param("id")

	// 生成 Redis 中存储点赞数的键名
	likeKey := "article:" + articleID + ":likes"

	// 从 Redis 获取点赞数
	likes, err := global.RedisDB.Get(likeKey).Result()

	if err == redis.Nil {
		// 如果键不存在，则设置点赞数为0
		likes = "0"
	} else if err != nil {
		// 如果 Redis 操作失败，返回内部服务器错误
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回点赞数
	ctx.JSON(http.StatusOK, gin.H{"likes": likes})
}

package util

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetParam(ctx *gin.Context, param string) string {
	return ctx.Param(param)
}

func GetParamUint64(ctx *gin.Context, param string) (uint64, error) {
	return strconv.ParseUint(ctx.Param(param), 10, 64)
}

func GetQuery(ctx *gin.Context, query string, defaultValue string) string {
	return ctx.DefaultQuery(query, defaultValue)
}

func GetQueryInt64(ctx *gin.Context, query string, defaultValue int64) (int64, error) {
	str := strconv.FormatInt(defaultValue, 10)
	return strconv.ParseInt(ctx.DefaultQuery(query, str), 10, 64)
}

func GetQueryUint64(ctx *gin.Context, query string, defaultValue uint64) (uint64, error) {
	str := strconv.FormatUint(defaultValue, 10)
	return strconv.ParseUint(ctx.DefaultQuery(query, str), 10, 64)
}

func GetQueryStrings(ctx *gin.Context, query string) []string {
	str := GetQuery(ctx, query, "")
	if str == "" {
		return []string{}
	}
	return strings.Split(str, ",")
}

func GetQueryUint64s(ctx *gin.Context, query string) ([]uint64, error) {
	str := GetQuery(ctx, query, "")
	if str == "" {
		return []uint64{}, nil
	}
	strs := strings.Split(str, ",")

	res := make([]uint64, len(strs))
	for i := range strs {
		num, err := strconv.ParseUint(strs[i], 10, 64)
		if err != nil {
			return nil, err
		}
		res[i] = num
	}
	return res, nil
}

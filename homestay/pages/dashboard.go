package pages

import (
	"fmt"
	"html/template"
	"time"

	"github.com/GoAdminGroup/go-admin/modules/config"
	template2 "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/gin-gonic/gin"
	homestaydb "github.com/moruowait/admin/homestay/db"
)

// 获取民宿统计图表
func GetHomestayDashboard(ctx *gin.Context) (types.Panel, error) {
	components := template2.Get(config.GetTheme())
	colComp := components.Col()
	// 默认查询本月
	now := time.Now()
	startTime := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local).Format(dateFormat)
	endTime := time.Date(now.Year(), now.Month(), 0, 0, 0, 0, 0, time.Local).AddDate(0, 1, 0).Format(dateFormat)

	// 获取收入明细
	incomeDetails, err := homestaydb.GetHomestayDB().GetHomestayIncomeDetailData(startTime, endTime)
	if err != nil {
		return types.Panel{}, err
	}

	// 渠道收入占比图
	incomeChannelBox := colComp.SetSize(types.SizeMD(3)).
		SetContent(
			components.Box().
				WithHeadBorder().SetHeader("渠道收入占比图").
				WithSecondHeadBorder().SetSecondHeader(template2.HTML(fmt.Sprintf("%s至%s明细", startTime, endTime))).
				SetBody(
					generateHomestayChannelPieChart(incomeDetails),
				).GetContent(),
		).GetContent()

	// 图表行
	chartRow := components.Row().SetContent(incomeChannelBox).GetContent()

	return types.Panel{
		Content:     chartRow,
		Title:       "民宿统计",
		Description: "民宿统计图表",
	}, nil
}

// 生成民宿月支出收入统计表
func generateHomestayMonthlySpendIncomeBarChart(room map[int]*homestaydb.HomestayRoom, incomes []*homestaydb.HomestayIncomeDetail) template.HTML {

}

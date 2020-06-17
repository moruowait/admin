package pages

import (
	"fmt"
	"html/template"
	"time"

	"github.com/GoAdminGroup/components/echarts"
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/modules/language"
	template2 "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/icon"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	"github.com/gin-gonic/gin"
	"github.com/go-echarts/go-echarts/charts"
	homestaydb "github.com/moruowait/admin/homestay/db"
)

const dateFormat = "2006-01-02"

func GetHomestaySearchDashboard(ctx *gin.Context) (types.Panel, error) {
	components := template2.Get(config.GetTheme())
	colComp := components.Col()
	// 初始查询数据
	startTime := ctx.Request.Form.Get("time_start__goadmin")
	endTime := ctx.Request.Form.Get("time_end__goadmin")
	if startTime == "" || endTime == "" {
		// 默认查询本月
		now := time.Now()
		startTime = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local).Format(dateFormat)
		endTime = time.Date(now.Year(), now.Month(), 0, 0, 0, 0, 0, time.Local).AddDate(0, 1, 0).Format(dateFormat)
	}

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

	// 房间收入占比图
	incomeRoomBox := colComp.SetSize(types.SizeMD(3)).
		SetContent(
			components.Box().
				WithHeadBorder().SetHeader("房间收入占比图").
				WithSecondHeadBorder().SetSecondHeader(template2.HTML(fmt.Sprintf("%s至%s明细", startTime, endTime))).
				SetBody(
					generateHomestayRoomPieChart(incomeDetails),
				).GetContent(),
		).GetContent()
	// 房间空置率图
	incomeRoomEmptyRateBox := colComp.SetSize(types.SizeMD(3)).
		SetContent(
			components.Box().
				WithHeadBorder().SetHeader("房间空置率图").
				WithSecondHeadBorder().SetSecondHeader(template2.HTML(fmt.Sprintf("%s至%s明细", startTime, endTime))).
				SetBody(
					generateHomestayRoomEmptyRateBarChart(startTime, endTime, incomeDetails),
				).GetContent(),
		).GetContent()
	// 渠道收入柱状图
	incomeBarBox := colComp.SetSize(types.SizeMD(12)).
		SetContent(
			components.Box().
				WithHeadBorder().SetHeader("渠道收入柱状图").
				WithSecondHeadBorder().SetSecondHeader(template2.HTML(fmt.Sprintf("%s至%s明细", startTime, endTime))).
				SetBody(
					generateHomestayChannelBarChart(startTime, endTime, incomeDetails),
				).GetContent(),
		).GetContent()
	// 图表行
	chartRow := components.Row().SetContent(incomeChannelBox + incomeRoomBox + incomeRoomEmptyRateBox).GetContent()
	chartRow2 := components.Row().SetContent(incomeBarBox).GetContent()
	// 搜索按钮
	searchBtn := components.Button().
		SetContent(language.GetFromHtml("search")).
		SetThemePrimary().
		SetOrientationLeft().
		SetLoadingText(icon.Icon("fa-spinner fa-spin", 2) + `Search`).
		GetContent()
	// 重置按钮
	resetBtn := components.Button().SetType("reset").
		SetContent(language.GetFromHtml("Reset")).
		SetOrientationRight().
		SetThemeWarning().
		SetHref("/admin/homestay/search_dashboard").
		GetContent()

	emptyCol := components.Col().SetSize(types.SizeMD(4)).GetContent()                                // 占用符
	btnCol := components.Col().SetSize(types.SizeMD(4)).SetContent(searchBtn + resetBtn).GetContent() // 搜索 + 重置 按钮

	fields, headers := types.NewFormPanel().
		AddField("日期", "time", db.Date, form.DateRange).
		SetTabGroups(types.TabGroups{{"time"}}).
		SetTabHeaders("查询条件").GroupField()
	searchForm := template2.Default().Form().
		SetMethod("get").
		SetTabHeaders(headers).
		SetTabContents(fields).
		SetUrl("/admin/homestay/search_dashboard").
		SetTitle("form").
		SetOperationFooter(emptyCol + btnCol).
		GetContent()
	searchRow := components.Row().SetContent(searchForm).GetContent()
	emptyRow := components.Row().SetContent(`<div style="margin:10px 10px 10px 10px;"></div>`).GetContent()
	return types.Panel{
		Content:     searchRow + emptyRow + chartRow + chartRow2,
		Title:       "民宿统计",
		Description: "民宿统计图表",
	}, nil
}

// 生成民宿收入渠道占比图
func generateHomestayChannelPieChart(details []*homestaydb.HomestayIncomeDetail) template.HTML {
	var data = make(map[string]interface{})
	for _, v := range details {
		// 删除刷单的渠道
		if v.Money < 0 {
			continue
		}
		if _, ok := data[v.Channel]; ok {
			data[v.Channel] = data[v.Channel].(float64) + v.Money
		} else {
			data[v.Channel] = v.Money
		}
	}

	pie := charts.NewPie()
	pie.Add("渠道收入占比", data)
	pie.SetGlobalOptions(charts.TooltipOpts{
		Show:      true,
		Trigger:   "item",
		Formatter: "{a} <br/>{b} : {c} ({d}%)",
	})
	pie.Width = "250px"
	pie.Height = "250px"
	return echarts.NewChart().SetContent(pie).GetContent()
}

// 生成民宿收入房间占比图
func generateHomestayRoomPieChart(details []*homestaydb.HomestayIncomeDetail) template.HTML {
	var data = make(map[string]interface{})
	for _, v := range details {
		// 删除刷单的渠道
		if v.Money < 0 {
			continue
		}
		room := fmt.Sprint(v.Room)
		if _, ok := data[room]; ok {
			data[room] = data[room].(float64) + v.Money
		} else {
			data[room] = v.Money
		}
	}

	pie := charts.NewPie()
	pie.Add("房间收入占比", data)
	pie.SetGlobalOptions(charts.TooltipOpts{
		Show:      true,
		Trigger:   "item",
		Formatter: "{a} <br/>{b} : {c} ({d}%)",
	})
	pie.Width = "250px"
	pie.Height = "250px"
	return echarts.NewChart().SetContent(pie).GetContent()
}

// 生成民宿房间空置率图
func generateHomestayRoomEmptyRateBarChart(startTime, endTime string, details []*homestaydb.HomestayIncomeDetail) template.HTML {
	var count = make(map[string]int)
	st, _ := time.Parse(dateFormat, startTime)
	et, _ := time.Parse(dateFormat, endTime)
	days := int(et.Sub(st).Hours()) / 24

	for _, v := range details {
		// 删除刷单的渠道
		if v.Money < 0 {
			continue
		}
		room := fmt.Sprint(v.Room)
		if _, ok := count[room]; ok {
			count[room] = 1
		} else {
			count[room] += 1
		}
	}

	var xAxis []string
	var yAxis []string
	for key, v := range count {
		xAxis = append(xAxis, key)
		yAxis = append(yAxis, fmt.Sprintf("%.2f", float64(v*100)/float64(days)))
	}

	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.YAxisOpts{AxisLabel: charts.LabelTextOpts{Formatter: "{value}%"}},
	)
	bar.AddXAxis(xAxis)
	bar.AddYAxis("百分比", yAxis, charts.BarOpts{YAxisIndex: 0, BarGap: "50%", BarCategoryGap: "50%"})

	bar.Width = "270px"
	bar.Height = "250px"
	return echarts.NewChart().SetContent(bar).GetContent()
}

// 生成民宿收入价格趋势图
func generateHomestayChannelBarChart(startTime, endTime string, details []*homestaydb.HomestayIncomeDetail) template.HTML {
	var channels = make(map[string][]float64)
	st, _ := time.Parse(dateFormat, startTime)
	et, _ := time.Parse(dateFormat, endTime)
	days := int(et.Sub(st).Hours()) / 24

	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.YAxisOpts{AxisLabel: charts.LabelTextOpts{Formatter: "{value} 元"}},
	)
	// X轴日期
	var xAxis []string
	var date = st
	for i := 1; i <= days; i++ {
		xAxis = append(xAxis, fmt.Sprintf("%d-%d", date.Month(), date.Day()))
		date = date.AddDate(0, 0, 1)
	}
	bar.AddXAxis(xAxis)

	for _, v := range details {
		// 删除刷单
		if v.Money < 0 {
			continue
		}
		if _, ok := channels[v.Channel]; !ok {
			channels[v.Channel] = make([]float64, days)
		}
		channels[v.Channel][int(v.IncomeTime.Sub(st).Hours())/24] = v.Money
	}

	for key, moneys := range channels {
		bar.AddYAxis(key, moneys, charts.BarOpts{YAxisIndex: 0})
	}
	bar.Width = "1800px"
	bar.Height = "250px"

	return echarts.NewChart().SetContent(bar).GetContent()
}

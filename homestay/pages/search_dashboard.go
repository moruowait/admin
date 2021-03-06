package pages

import (
	"fmt"
	"html/template"
	"sort"
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

// 获取民宿查询图表
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
	// 获取支出明细
	spendDetails, err := homestaydb.GetHomestayDB().GetHomestaySpendDetailData(startTime, endTime)
	if err != nil {
		return types.Panel{}, err
	}
	fmt.Println("spendDetails:", spendDetails)
	// 获取房间信息
	roomMap, err := homestaydb.GetHomestayDB().GetHomestayRoomMap()
	if err != nil {
		return types.Panel{}, err
	}
	fmt.Println("roomMap:", roomMap)
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
	incomeRoomEmptyRateBox := colComp.SetSize(types.SizeMD(6)).
		SetContent(
			components.Box().
				WithHeadBorder().SetHeader("房间空置率图").
				WithSecondHeadBorder().SetSecondHeader(template2.HTML(fmt.Sprintf("%s至%s明细", startTime, endTime))).
				SetBody(
					generateHomestayRoomEmptyRateBarChart(startTime, endTime, incomeDetails),
				).GetContent(),
		).GetContent()
	//generateHomestayMonthlyRoomEmptyRateBarChart
	// 房间空置率图
	incomeMonthlyRoomEmptyRateBox := colComp.SetSize(types.SizeMD(6)).
		SetContent(
			components.Box().
				WithHeadBorder().SetHeader("房间空置率图").
				WithSecondHeadBorder().SetSecondHeader(template2.HTML(fmt.Sprintf("%s至%s明细", startTime, endTime))).
				SetBody(
					generateHomestayMonthlyRoomEmptyRateBarChart(startTime, endTime, incomeDetails),
				).GetContent(),
		).GetContent()

	// 收入支出对比图
	incomeSpendMonthlyBox := colComp.SetSize(types.SizeMD(6)).
		SetContent(
			components.Box().
				WithHeadBorder().SetHeader("月收入支出图").
				WithSecondHeadBorder().SetSecondHeader(template2.HTML(fmt.Sprintf("%s至%s明细", startTime, endTime))).
				SetBody(
					generateHomestaySpendIncomeBarChart(startTime, endTime, spendDetails, incomeDetails, roomMap),
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
	chartRow := components.Row().SetContent(incomeChannelBox + incomeRoomBox + incomeRoomEmptyRateBox + incomeMonthlyRoomEmptyRateBox + incomeSpendMonthlyBox).GetContent()
	chartRow2 := components.Row().SetContent(incomeBarBox).GetContent()
	//chartRow3 := components.Row().SetContent(incomeSpendMonthlyBox).GetContent()
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
		yAxis = append(yAxis, fmt.Sprintf("%.2f", 100-float64(v*100)/float64(days)))
	}

	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.YAxisOpts{AxisLabel: charts.LabelTextOpts{Formatter: "{value}%"}},
	)
	bar.AddXAxis(xAxis)
	bar.AddYAxis("百分比", yAxis, charts.BarOpts{YAxisIndex: 0, BarGap: "50%", BarCategoryGap: "50%"})

	bar.Width = "450px"
	bar.Height = "250px"
	return echarts.NewChart().SetContent(bar).GetContent()
}

// 生成民宿房间月空置率图
func generateHomestayMonthlyRoomEmptyRateBarChart(startTime, endTime string, incomes []*homestaydb.HomestayIncomeDetail) template.HTML {
	st, _ := time.Parse(dateFormat, startTime)
	et, _ := time.Parse(dateFormat, endTime)

	var xAxis []string // X轴日期
	for it := st; it.Before(et); it = it.AddDate(0, 1, 0) {
		xAxis = append(xAxis, formatYearMonth(it))
	}
	// 月天数
	dayMap := getMonthDay(st, et)
	var incomeMap = make(map[string]map[string]float64)
	for _, income := range incomes {
		// 删除刷单的渠道
		if income.Money < 0 {
			continue
		}

		room := fmt.Sprint(income.Room)

		if _, ok := incomeMap[room]; !ok {
			incomeMap[room] = map[string]float64{
				formatYearMonth(income.IncomeTime): 1,
			}
		} else {
			if _, ok := incomeMap[room][formatYearMonth(income.IncomeTime)]; !ok {
				incomeMap[room][formatYearMonth(income.IncomeTime)] = 1
			} else {
				incomeMap[room][formatYearMonth(income.IncomeTime)]++
			}

		}
	}

	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.YAxisOpts{AxisLabel: charts.LabelTextOpts{Formatter: "{value}%"}},
	)
	bar.AddXAxis(xAxis)
	for room, months := range incomeMap {
		var yAxis []string
		for month, _ := range months {
			//incomeMap[room][month] = incomeMap[room][month] / float64(dayMap[month])
			yAxis = append(yAxis, fmt.Sprintf("%.2f", 100-incomeMap[room][month]*100/float64(dayMap[month])))
		}
		bar.AddYAxis(room, yAxis, charts.BarOpts{YAxisIndex: 0, BarGap: "50%", BarCategoryGap: "50%"})
	}
	bar.Width = "450px"
	bar.Height = "250px"
	return echarts.NewChart().SetContent(bar).GetContent()
}

// 生成民宿收支对比图
func generateHomestaySpendIncomeBarChart(startTime, endTime string, spends []*homestaydb.HomestaySpendDetail, incomes []*homestaydb.HomestayIncomeDetail, room map[int]*homestaydb.HomestayRoom) template.HTML {
	st, _ := time.Parse(dateFormat, startTime)
	et, _ := time.Parse(dateFormat, endTime)

	var spendMap = make(map[string]float64)
	var incomeMap = make(map[string]float64)
	var makeDealMap = make(map[string]float64) // 刷单
	var xAxis []string                         // X轴日期
	for it := st; it.Before(et); it = it.AddDate(0, 1, 0) {
		key := formatYearMonth(it)
		spendMap[key] = 0
		incomeMap[key] = 0
		xAxis = append(xAxis, key)
	}
	// 支出
	for _, v := range spends {
		spendMap[formatYearMonth(v.Time)] += v.Money
	}
	// 物业费+房租
	for _, r := range room {
		for k, _ := range spendMap {
			spendMap[k] += r.MonthlyManageFee + r.MonthlyRent
		}
	}
	// 收入
	for _, v := range incomes {
		if v.Money < 0 { // 刷单
			makeDealMap[formatYearMonth(v.IncomeTime)] += -v.Money
			continue
		}
		incomeMap[formatYearMonth(v.IncomeTime)] += v.Money
	}

	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.YAxisOpts{AxisLabel: charts.LabelTextOpts{Formatter: "{value} 元"}},
	)
	bar.AddXAxis(xAxis)
	bar.AddYAxis("支出", sortedMap(spendMap), charts.BarOpts{YAxisIndex: 0})
	bar.AddYAxis("收入", sortedMap(incomeMap), charts.BarOpts{YAxisIndex: 0})
	bar.AddYAxis("刷单", sortedMap(makeDealMap), charts.BarOpts{YAxisIndex: 0})
	bar.Width = "650px"
	bar.Height = "250px"

	return echarts.NewChart().SetContent(bar).GetContent()
}

// 获取范围内月天数
func getMonthDay(startTime, endTime time.Time) map[string]int {
	var MDays = make(map[string]int)
	for it := startTime; it.Before(endTime); it = it.AddDate(0, 0, 1) {
		MDays[formatYearMonth(it)]++
	}
	return MDays
}

func formatYearMonth(t time.Time) string {
	return fmt.Sprintf("%d-%d", t.Year(), t.Month())
}

func sortedMap(m map[string]float64) []float64 {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var value []float64
	for _, k := range keys {
		value = append(value, m[k])
	}
	return value
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

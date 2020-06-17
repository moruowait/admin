package homestay

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

func GetHomestayIncomeDetailTable(ctx *context.Context) table.Table {

	homestayIncomeDetailTable := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))

	info := homestayIncomeDetailTable.GetInfo()

	info.AddField("序号", "id", db.Int)
	info.AddField("渠道", "name", db.Varchar).FieldJoin(types.Join{
		Table:     "homestay_channel",
		Field:     "channel_id",
		JoinField: "id",
		BaseTable: "homestay_income_detail",
	}).FieldFilterable(types.FilterType{
		FormType: form.SelectSingle,
	}).FieldFilterOptionsFromTable("homestay_channel", "name", "name")

	info.AddField("房间号", "number", db.Int).FieldJoin(types.Join{
		Table:     "homestay_room",
		Field:     "room_id",
		JoinField: "id",
		BaseTable: "homestay_income_detail",
	}).FieldFilterable(types.FilterType{
		FormType: form.SelectSingle,
	}).FieldFilterOptionsFromTable("homestay_room", "number", "number")

	info.AddField("收入", "money", db.Float)
	info.AddField("收入时间", "income_time", db.Timestamp)
	info.AddField("描述", "desc", db.Varchar)
	info.AddField("更新时间", "update_time", db.Timestamp)

	info.SetTable("homestay_income_detail").SetTitle("收入明细").SetDescription("收入明细")

	formList := homestayIncomeDetailTable.GetForm()

	formList.AddField("Id", "id", db.Int, form.Default).FieldNotAllowAdd()
	formList.AddField("渠道名称", "channel_id", db.Int, form.SelectSingle).FieldOptionsFromTable("homestay_channel", "name", "id")
	formList.AddField("房间号", "room_id", db.Int, form.SelectSingle).FieldOptionsFromTable("homestay_room", "number", "id")
	formList.AddField("收入", "money", db.Float, form.Currency)
	formList.AddField("收入时间", "income_time", db.Timestamp, form.Datetime)
	formList.AddField("描述", "desc", db.Varchar, form.TextArea)

	formList.SetTable("homestay_income_detail").SetTitle("收入明细").SetDescription("收入明细")

	return homestayIncomeDetailTable
}

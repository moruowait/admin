package homestay

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

func GetHomestaySpendDetailTable(ctx *context.Context) table.Table {

	homestaySpendDetailTable := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))

	info := homestaySpendDetailTable.GetInfo()

	info.AddField("序号", "id", db.Int)
	info.AddField("项目", "name", db.Varchar).FieldJoin(types.Join{
		Table:     "homestay_spend_item",
		Field:     "item_id",
		JoinField: "id",
		BaseTable: "homestay_spend_detail",
	}).FieldFilterable(types.FilterType{
		FormType: form.SelectSingle,
	}).FieldFilterOptionsFromTable("homestay_spend_item", "name", "name")

	info.AddField("支出", "money", db.Float)
	info.AddField("支出时间", "time", db.Timestamp)
	info.AddField("描述", "desc", db.Varchar)
	info.AddField("更新时间", "update_time", db.Timestamp)

	info.SetTable("homestay_spend_detail").SetTitle("支出明细").SetDescription("支出明细")

	formList := homestaySpendDetailTable.GetForm()

	formList.AddField("序号", "id", db.Int, form.Default).FieldNotAllowAdd()
	formList.AddField("项目序号", "item_id", db.Int, form.SelectSingle).FieldOptionsFromTable("homestay_spend_item", "name", "id")
	formList.AddField("支出", "money", db.Float, form.Currency)
	formList.AddField("支出时间", "time", db.Timestamp, form.Datetime)
	formList.AddField("描述", "desc", db.Varchar, form.TextArea)

	formList.SetTable("homestay_spend_detail").SetTitle("支出明细").SetDescription("支出明细")

	return homestaySpendDetailTable
}

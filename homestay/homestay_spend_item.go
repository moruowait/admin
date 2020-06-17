package homestay

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

func GetHomestaySpendItemTable(ctx *context.Context) table.Table {

	homestaySpendItemTable := table.NewDefaultTable(table.DefaultConfigWithDriver("mysql"))

	info := homestaySpendItemTable.GetInfo()

	info.AddField("Id", "id", db.Int)
	info.AddField("Name", "name", db.Varchar)
	info.AddField("Desc", "desc", db.Varchar)

	info.SetTable("homestay_spend_item").SetTitle("支出项目").SetDescription("支出项目")

	formList := homestaySpendItemTable.GetForm()

	formList.AddField("Id", "id", db.Int, form.Default).FieldNotAllowAdd()
	formList.AddField("Name", "name", db.Varchar, form.Text)
	formList.AddField("Desc", "desc", db.Varchar, form.Text)

	formList.SetTable("homestay_spend_item").SetTitle("支出项目").SetDescription("支出项目")

	return homestaySpendItemTable
}

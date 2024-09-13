package main

import (
	"github.com/gin-gonic/gin"
)

type JurnalSetting struct {
	IdEntitas   int
	KodeSetting string
	KodeSubakun string
}
type MasterJurnalSetting struct {
	KodeSetting string
	Setting     string
	Keterangan  string
}
type MasterJurnalKategoriAkun struct {
	KodeKategori string
	IdEntitas    int
	Kategori     string
	Tipe         string
}
type MasterJurnalGrupAkun struct {
	KodeGrup     string
	KodeKategori string
	GrupAkun     string
}

func GetMasterJurnalGrupAkun(KodeKategori string) ([]MasterJurnalGrupAkun, error) {
	var jurnalGrupAkuns []MasterJurnalGrupAkun
	tx := db.Debug().Where("kode_kategori = ?", KodeKategori).Find(&jurnalGrupAkuns)

	if tx.Error != nil {
		return nil, tx.Error
	}
	return jurnalGrupAkuns, nil
}

func MasterJurnalSettingRoute(router *gin.Engine) {
	group := router.Group("/keuangan/master-jurnal-setting", corsMiddleware())
	{
		group.GET("/get", func(context *gin.Context) {
			// db_go := db.KoneksiCore()
			var callback = gin.H{}
			status := 200
			identitas := context.Query("identitas")

			result := []map[string]interface{}{}
			sql := "SELECT A.kode_setting, A.setting, A.keterangan, B.kode_subakun, C.subakun, D.kode_akun, D.akun FROM jurnal_master_setting AS A "
			sql += "LEFT JOIN jurnal_setting_subakun AS B ON (A.kode_setting = B.kode_setting AND B.identitas=?) "
			sql += "LEFT JOIN jurnal_master_subakun AS C ON (B.kode_subakun = C.kode_subakun) "
			sql += "LEFT JOIN jurnal_master_akun AS D ON (C.kode_akun = D.kode_akun) "

			db.Raw(sql, identitas).Scan(&result)

			callback["success"] = true
			callback["data"] = result
			DB, _ := db.DB()
			DB.Close()

			context.JSON(status, callback)
		})

		group.GET("/get-akun", func(context *gin.Context) {
			// db_go := db.KoneksiCore()
			var callback = gin.H{}
			status := 200
			identitas := context.Query("identitas")
			kode_akun := context.Query("kode_akun")

			result_akun := []map[string]interface{}{}
			sql := "SELECT A.kode_akun, A.akun FROM master_jurnal_akun AS A "
			sql += "INNER JOIN master_jurnal_group_akun AS B ON (A.kode_group = B.kode_group) "
			sql += "INNER JOIN master_jurnal_kategori_akun AS C ON (B.kode_kategori = C.kode_kategori) "
			sql += "WHERE C.identitas=?"
			db.Raw(sql, identitas).Scan(&result_akun)

			result_subakun := []map[string]interface{}{}
			sql = "SELECT B.akun, A.subakun, A.kode_subakun FROM master_jurnal_subakun AS A "
			sql += "INNER JOIN master_jurnal_akun AS B ON (A.kode_akun = B.kode_akun) "
			sql += "INNER JOIN master_jurnal_group_akun AS C ON (B.kode_group = C.kode_group) "
			sql += "INNER JOIN master_jurnal_kategori_akun AS D ON (C.kode_kategori = D.kode_kategori) "
			sql += "WHERE D.identitas=? AND A.kode_akun=?"
			db.Raw(sql, identitas, kode_akun).Scan(&result_subakun)

			callback["success"] = true
			callback["data"] = map[string]interface{}{
				"akun":     result_akun,
				"sub_akun": result_subakun,
			}
			DB, _ := db.DB()
			DB.Close()

			context.JSON(status, callback)
		})

		group.GET("/get-subakun", func(context *gin.Context) {
			// db_go := db.KoneksiCore()
			var callback = gin.H{}
			status := 200
			kode_akun := context.Query("kode_akun")
			identitas := context.Query("identitas")
			// kode_akun := ("FA.01.02.01")

			result := []map[string]interface{}{}
			sql := "SELECT B.akun, A.subakun, A.kode_subakun FROM master_jurnal_subakun AS A "
			sql += "INNER JOIN master_jurnal_akun AS B ON (A.kode_akun = B.kode_akun) "
			sql += "INNER JOIN master_jurnal_group_akun AS C ON (B.kode_group = C.kode_group) "
			sql += "INNER JOIN master_jurnal_kategori_akun AS D ON (C.kode_kategori = D.kode_kategori) "
			sql += "WHERE D.identitas=? AND A.kode_akun=?"
			db.Raw(sql, identitas, kode_akun).Scan(&result)

			callback["success"] = true
			callback["data"] = result
			DB, _ := db.DB()
			DB.Close()

			context.JSON(status, callback)
		})

		group.POST("/edit", func(context *gin.Context) {
			// db_go := db.KoneksiCore()
			var callback = gin.H{}
			// akun := context.PostForm("akun")
			// subakun := context.PostForm("subakun")
			kode_setting := context.PostForm("kode_setting")
			kode_subakun := context.PostForm("kode_subakun")
			identitas := context.Query("identitas")

			db.Exec("DELETE FROM jurnal_setting WHERE identitas=? AND kode_setting=?", identitas, kode_setting)

			data := map[string]interface{}{
				"identitas":    identitas,
				"kode_setting": kode_setting,
				"kode_subakun": kode_subakun,
			}
			create := db.Table("jurnal_setting").Create(&data)

			if create.Error == nil {
				callback["success"] = true
				callback["msg"] = "Data berhasil diupdate"
			} else {
				callback["success"] = false
				callback["msg"] = "Update Gagal"
			}
			DB, _ := db.DB()
			DB.Close()

			context.JSON(200, callback)
		})
	}
}

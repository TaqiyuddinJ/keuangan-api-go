package keuangan

import (
	"apigo/lib/db"
	"apigo/lib/mainlib"
	"apigo/lib/middleware"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func RekeningSettingRoute(router *gin.Engine) {
	group := router.Group("/keuangan/rekening-setting", middleware.CORSMiddleware())
	{
		group.GET("/get", func(context *gin.Context) {
			db_go := db.KoneksiCore()
			var callback = gin.H{}
			status := 200
			idkoperasi := mainlib.GetKoperasiID(context)
			// idkoperasi := 3

			result := []map[string]interface{}{}
			sql := "SELECT A.bank, A.norek, A.nama_rekening, A.idkoperasi, A.logo, A.topup_aktif, A.kode_subakun, B.subakun, C.kode_akun FROM koperasi_bank AS A "
			sql += "INNER JOIN jurnal_master_subakun AS B ON (A.kode_subakun = B.kode_subakun) "
			sql += "INNER JOIN jurnal_master_akun AS C ON (B.kode_akun = C.kode_akun) "
			sql += "WHERE A.idkoperasi=?"

			db_go.Raw(sql, idkoperasi).Scan(&result)

			callback["success"] = true
			callback["data"] = result

			DB, _ := db_go.DB()
			DB.Close()
			context.JSON(status, callback)
		})
		group.GET("/get-akun", func(context *gin.Context) {
			db_go := db.KoneksiCore()
			var callback = gin.H{}
			status := 200
			idkoperasi := mainlib.GetKoperasiID(context)
			kode_akun := context.Query("kode_akun")
			// idkoperasi := 3

			result_akun := []map[string]interface{}{}
			sql := "SELECT A.kode_akun, A.akun FROM master_jurnal_akun AS A "
			sql += "INNER JOIN master_jurnal_group_akun AS B ON (A.kode_group = B.kode_group) "
			sql += "INNER JOIN master_jurnal_kategori_akun AS C ON (B.kode_kategori = C.kode_kategori) "
			sql += "WHERE C.idkoperasi=?"
			db_go.Raw(sql, idkoperasi).Scan(&result_akun)

			result_subakun := []map[string]interface{}{}
			sql = "SELECT B.akun, A.subakun, A.kode_subakun FROM master_jurnal_subakun AS A "
			sql += "INNER JOIN master_jurnal_akun AS B ON (A.kode_akun = B.kode_akun) "
			sql += "INNER JOIN master_jurnal_group_akun AS C ON (B.kode_group = C.kode_group) "
			sql += "INNER JOIN master_jurnal_kategori_akun AS D ON (C.kode_kategori = D.kode_kategori) "
			sql += "WHERE D.idkoperasi=? AND A.kode_akun=?"
			db_go.Raw(sql, idkoperasi, kode_akun).Scan(&result_subakun)

			callback["success"] = true
			callback["data"] = map[string]interface{}{
				"akun":     result_akun,
				"sub_akun": result_subakun,
			}

			DB, _ := db_go.DB()
			DB.Close()
			context.JSON(status, callback)
		})

		group.GET("/get-subakun", func(context *gin.Context) {
			db_go := db.KoneksiCore()
			var callback = gin.H{}
			status := 200
			kode_akun := context.Query("kode_akun")
			idkoperasi := mainlib.GetKoperasiID(context)
			// idkoperasi := 3
			// kode_akun := ("FA.01.02.01")

			result := []map[string]interface{}{}
			sql := "SELECT B.akun, A.subakun, A.kode_subakun FROM master_jurnal_subakun AS A "
			sql += "INNER JOIN master_jurnal_akun AS B ON (A.kode_akun = B.kode_akun) "
			sql += "INNER JOIN master_jurnal_group_akun AS C ON (B.kode_group = C.kode_group) "
			sql += "INNER JOIN master_jurnal_kategori_akun AS D ON (C.kode_kategori = D.kode_kategori) "
			sql += "WHERE D.idkoperasi=? AND A.kode_akun=?"
			db_go.Raw(sql, idkoperasi, kode_akun).Scan(&result)

			callback["success"] = true
			callback["data"] = result
			DB, _ := db_go.DB()
			DB.Close()

			context.JSON(status, callback)
		})

		group.POST("/add", func(context *gin.Context) {
			db_go := db.KoneksiCore()
			var callback = gin.H{}
			norek := context.PostForm("norek")
			bank := context.PostForm("bank")
			nama_rekening := context.PostForm("nama_rekening")
			kode_subakun := context.PostForm("kode_subakun")
			topup_aktif := context.PostForm("topup_aktif")
			logo, err := context.FormFile("logo")
			idkoperasi := mainlib.GetKoperasiID(context)

			if err != nil {
				callback["logo"] = "No Image Found"
				// return
			}

			extension := filepath.Ext(logo.Filename)
			newFileName := "logoBank-" + time.Now().Format("20060102150405") + extension

			filepath := "../uploads/logoBank/" + newFileName
			if err := context.SaveUploadedFile(logo, filepath); err != nil {
				callback["logo"] = "Unable Save Foto"
				callback["logoz"] = err
			}
			data := map[string]interface{}{
				"idkoperasi":    idkoperasi,
				"norek":         norek,
				"bank":          bank,
				"nama_rekening": nama_rekening,
				"kode_subakun":  kode_subakun,
				"topup_aktif":   topup_aktif,
				"logo":          newFileName,
			}
			create := db_go.Table("koperasi_bank").Create(&data)

			if create.Error == nil {
				callback["success"] = true
				callback["msg"] = "Tambah Data Berhasil"
				callback["z"] = filepath
			} else {
				callback["success"] = false
				callback["msg"] = create.Error
			}
			DB, _ := db_go.DB()
			DB.Close()

			context.JSON(200, callback)
		})

		group.POST("/edit", func(context *gin.Context) {
			db_go := db.KoneksiCore()
			var callback = gin.H{}
			// akun := context.PostForm("akun")
			// subakun := context.PostForm("subakun")
			// kode_setting := context.PostForm("kode_setting")
			// kode_subakun := context.PostForm("kode_subakun")
			norek := context.PostForm("norek")
			norek_data := context.PostForm("norek_data")
			bank := context.PostForm("bank")
			nama_rekening := context.PostForm("nama_rekening")
			kode_subakun := context.PostForm("kode_subakun")
			topup_aktif := context.PostForm("topup_aktif")
			logo, err := context.FormFile("logo")
			// idkoperasi := mainlib.GetKoperasiID(context)

			if err != nil {

				update := db_go.Exec("UPDATE koperasi_bank SET norek=?, bank=?, nama_rekening=?, kode_subakun=?, topup_aktif=? WHERE norek=?", norek, bank, nama_rekening, kode_subakun, topup_aktif, norek_data)

				if update.Error == nil {
					callback["success"] = true
					callback["msg"] = "Data berhasil diupdate"
				} else {
					callback["success"] = false
					callback["msg"] = "Update Gagal"
				}
				callback["foto"] = "No Image Found"
				// return
			} else {
				extension := filepath.Ext(logo.Filename)
				newFileName := "logoBank-" + time.Now().Format("20060102150405") + extension

				filepath := "../uploads/logoBank/" + newFileName
				if err := context.SaveUploadedFile(logo, filepath); err != nil {
					callback["logo"] = "Unable Save Foto"
					callback["logoz"] = err
				}

				update := db_go.Exec("UPDATE koperasi_bank SET norek=?, bank=?, nama_rekening=?, kode_subakun=?, topup_aktif=?, logo=? WHERE norek=?", norek, bank, nama_rekening, kode_subakun, topup_aktif, newFileName, norek_data)

				if update.Error == nil {
					callback["success"] = true
					callback["msg"] = "Data berhasil diupdate"
					callback["z"] = filepath
				} else {
					callback["success"] = false
					callback["msg"] = "Update Gagal"
				}
			}
			DB, _ := db_go.DB()
			DB.Close()

			context.JSON(200, callback)
		})

		group.POST("/delete", func(context *gin.Context) {
			db_go := db.KoneksiCore()
			var callback = gin.H{}
			norek := context.PostForm("norek")
			idkoperasi := mainlib.GetKoperasiID(context)

			result := db_go.Exec("DELETE FROM koperasi_bank WHERE idkoperasi=? AND norek = ?", idkoperasi, norek)
			if result.Error == nil {
				callback["success"] = true
				callback["msg"] = "Data berhasil dihapus"
			} else {
				callback["success"] = false
				callback["msg"] = "Hapus Gagal"
			}

			//CEK PARAMETER POST
			DB, _ := db_go.DB()
			DB.Close()

			context.JSON(200, callback)
		})
	}
}

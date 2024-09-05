package main

import (
	"apigo/lib/db"
	"apigo/lib/mainlib"
	"apigo/lib/middleware"

	"github.com/gin-gonic/gin"
)

func MasterSubAkunRoute(router *gin.Engine) {
	group := router.Group("/keuangan/master-subakun", middleware.CORSMiddleware())
	{
		group.GET("/get", func(context *gin.Context) {
			db_go := db.KoneksiCore()
			var callback = gin.H{}
			status := 200
			idkoperasi := mainlib.GetKoperasiID(context)

			result := []map[string]interface{}{}
			sql := "SELECT A.kode_subakun, A.kode_akun, A.subakun, B.akun, D.kategori, A.keterangan "
			sql += "FROM jurnal_master_subakun AS A "
			sql += "INNER JOIN jurnal_master_akun AS B ON (A.kode_akun = B.kode_akun) "
			sql += "INNER JOIN jurnal_master_kategori AS D ON (B.kode_kategori = D.kode_kategori) "
			sql += "WHERE B.idkoperasi=?"
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

			result := []map[string]interface{}{}
			db_go.Raw("SELECT * FROM jurnal_master_akun AS A INNER JOIN jurnal_master_kategori AS C ON (A.kode_kategori = C.kode_kategori) WHERE idkoperasi=?", idkoperasi).Scan(&result)

			callback["success"] = true
			callback["data"] = result

			DB, _ := db_go.DB()
			DB.Close()

			context.JSON(status, callback)
		})

		group.POST("/add", func(context *gin.Context) {
			db_go := db.KoneksiCore()
			var callback = gin.H{}
			kode_subakun := context.PostForm("kode_subakun")
			kode_akun := context.PostForm("kode_akun")
			subakun := context.PostForm("subakun")
			keterangan := context.PostForm("keterangan")
			//CEK PARAMETER POST
			//callback["master"] = master

			result := map[string]interface{}{}
			check := db_go.Raw("SELECT * FROM jurnal_master_subakun WHERE kode_subakun = ?", kode_subakun).Scan(&result)

			if check.RowsAffected == 0 {
				data := map[string]interface{}{
					"kode_subakun": kode_subakun,
					"kode_akun":    kode_akun,
					"subakun":      subakun,
					"keterangan":   keterangan,
					"tetap":        0,
				}
				create := db_go.Table("jurnal_master_subakun").Create(&data)

				if create.Error == nil {
					callback["success"] = true
					callback["msg"] = "Tambah Data Berhasil"
				} else {
					callback["success"] = false
					callback["msg"] = create.Error
				}
			} else {
				callback["success"] = false
				callback["msg"] = "Kode Sudah Digunakan!"
			}

			DB, _ := db_go.DB()
			DB.Close()

			context.JSON(200, callback)
		})

		group.POST("/edit", func(context *gin.Context) {
			db_go := db.KoneksiCore()
			var callback = gin.H{}
			kode_subakun_lama := context.PostForm("kode_subakun_lama")
			kode_subakun := context.PostForm("kode_subakun")
			kode_akun := context.PostForm("kode_akun")
			subakun := context.PostForm("subakun")
			keterangan := context.PostForm("keterangan")

			result := map[string]interface{}{}
			check := db_go.Raw("SELECT * FROM jurnal_master_subakun WHERE kode_subakun =? AND kode_subakun != ?", kode_subakun, kode_subakun_lama).Scan(&result)

			if check.RowsAffected == 0 {
				update := db_go.Exec("UPDATE jurnal_master_subakun SET kode_subakun=?, kode_akun=?, subakun=?, keterangan=? WHERE kode_subakun=?", kode_subakun, kode_akun, subakun, keterangan, kode_subakun_lama)

				if update.Error == nil {
					callback["success"] = true
					callback["msg"] = "Data berhasil diupdate"
				} else {
					callback["success"] = false
					callback["msg"] = "Update Gagal"
				}
			} else {
				callback["success"] = false
				callback["msg"] = "Kode sudah digunakan!"
			}

			// CEK PARAMETER POST

			DB, _ := db_go.DB()
			DB.Close()

			context.JSON(200, callback)
		})

		group.POST("/delete", func(context *gin.Context) {
			db_go := db.KoneksiCore()
			var callback = gin.H{}
			kode_subakun := context.PostForm("kode_subakun")

			result := db_go.Exec("DELETE FROM jurnal_master_subakun WHERE kode_subakun = ?", kode_subakun)
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

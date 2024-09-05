package keuangan

import (
	"apigo/lib/db"
	"apigo/lib/mainlib"
	"apigo/lib/middleware"

	"github.com/gin-gonic/gin"
)

func MasterAkunRoute(router *gin.Engine) {
	group := router.Group("/keuangan/master-akun", middleware.CORSMiddleware())
	{
		group.GET("/get", func(context *gin.Context) {
			db_go := db.KoneksiCore()
			var callback = gin.H{}
			status := 200
			idkoperasi := mainlib.GetKoperasiID(context)

			result := []map[string]interface{}{}
			sql := "SELECT * FROM jurnal_master_akun AS A "
			sql += "INNER JOIN jurnal_master_kategori AS C ON (A.kode_kategori = C.kode_kategori) "
			sql += "WHERE A.idkoperasi=?"
			db_go.Raw(sql, idkoperasi).Scan(&result)

			callback["success"] = true
			callback["data"] = result
			// callback["idpersonalia"] = idpersonalia
			DB, _ := db_go.DB()
			DB.Close()

			context.JSON(status, callback)
		})

		group.GET("/get-kategori", func(context *gin.Context) {
			db_go := db.KoneksiCore()
			var callback = gin.H{}
			status := 200
			// idkoperasi := mainlib.GetKoperasiID(context)

			result := []map[string]interface{}{}
			db_go.Raw("SELECT * FROM jurnal_master_kategori ").Scan(&result)

			callback["success"] = true
			callback["data"] = result

			DB, _ := db_go.DB()
			DB.Close()

			context.JSON(status, callback)
		})

		group.POST("/add", func(context *gin.Context) {
			db_go := db.KoneksiCore()
			var callback = gin.H{}
			kode_akun := context.PostForm("kode_akun")
			kode_kategori := context.PostForm("kode_kategori")
			akun := context.PostForm("akun")
			keterangan := context.PostForm("keterangan")
			idkoperasi := mainlib.GetKoperasiID(context)
			//CEK PARAMETER POST
			//callback["master"] = master

			//CEK DUPLIKASI
			result := map[string]interface{}{}
			check := db_go.Raw("SELECT * FROM jurnal_master_akun  WHERE kode_akun = ?", kode_akun).Scan(&result)

			if check.RowsAffected == 0 {
				data := map[string]interface{}{
					"kode_akun":     kode_akun,
					"kode_kategori": kode_kategori,
					"akun":          akun,
					"keterangan":    keterangan,
					"idkoperasi":    idkoperasi,
					"tetap":         0,
				}
				create := db_go.Table("jurnal_master_akun").Create(&data)

				if create.Error == nil {
					callback["success"] = true
					callback["msg"] = "Tambah Data Berhasil"
				} else {
					callback["success"] = false
					callback["msg"] = create.Error
				}
			} else {
				callback["success"] = false
				callback["msg"] = "Kode Sudah digunakan!"
			}

			DB, _ := db_go.DB()
			DB.Close()

			context.JSON(200, callback)
		})

		group.POST("/edit", func(context *gin.Context) {
			db_go := db.KoneksiCore()
			var callback = gin.H{}
			kode_akun_lama := context.PostForm("kode_akun_lama")
			kode_akun := context.PostForm("kode_akun")
			kode_kategori := context.PostForm("kode_kategori")
			akun := context.PostForm("akun")
			keterangan := context.PostForm("keterangan")

			result := map[string]interface{}{}
			check := db_go.Raw("SELECT * FROM jurnal_master_akun WHERE kode_akun = ? AND kode_akun != ?", kode_akun, kode_akun_lama).Scan(&result)

			if check.RowsAffected == 0 {
				update := db_go.Exec("UPDATE jurnal_master_akun SET kode_akun=?, kode_kategori=?, akun=?, keterangan=? WHERE kode_akun=?", kode_akun, kode_kategori, akun, keterangan, kode_akun_lama)

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
			kode_akun := context.PostForm("kode_akun")

			result := db_go.Exec("DELETE FROM jurnal_master_akun WHERE kode_akun = ?", kode_akun)
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

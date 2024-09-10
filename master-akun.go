package main

import (
	"github.com/gin-gonic/gin"
)

type MasterJurnalAkun struct {
	KodeAkun string
	KodeGrup string
	Akun     string
}

func GetMasterAkun() ([]MasterJurnalAkun, error) {
	var jurnalAkuns []MasterJurnalAkun
	tx := db.Debug().Find(&jurnalAkuns)

	if tx.Error != nil {
		return nil, tx.Error
	}
	return jurnalAkuns, nil
}
func GetMasterAkunDariEntitas(identitas int) ([]MasterJurnalAkun, error) {
	var jurnalakuns []MasterJurnalAkun
	tx := db.Debug().Joins("INNER JOIN master_jurnal_grup_akun ON master_jurnal_akun.kode_grup = master_jurnal_grup_akun.kode_grup").Joins("INNER JOIN master_jurnal_kategori_akun ON master_jurnal_grup_akun.kode_kategori = master_jurnal_grup_akun.kode_kategori").Where("identitas = ?", identitas).Find(&jurnalakuns)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return jurnalakuns, nil
}

func GetMasterAkunMap(identitas int) ([]map[string]interface{}, error) {
	var jurnalakuns []map[string]interface{}
	sql := "SELECT * FROM jurnal_master_akun AS A "
	sql += "INNER JOIN jurnal_master_kategori AS C ON (A.kode_kategori = C.kode_kategori) "
	sql += "WHERE A.identitas=?"
	tx := db.Debug().Raw(sql, identitas).Scan(&jurnalakuns)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return jurnalakuns, nil
}

func MasterAkunRoute(router *gin.Engine) {
	group := router.Group("/keuangan/master-akun", corsMiddleware())
	{
		// group.GET("/get", func(context *gin.Context) {
		// 	// db_go := db.KoneksiCore()
		// 	var callback = gin.H{}
		// 	status := 200
		// 	identitas := context.Query("identitas")

		// 	result := []map[string]interface{}{}
		// 	sql := "SELECT * FROM jurnal_master_akun AS A "
		// 	sql += "INNER JOIN jurnal_master_kategori AS C ON (A.kode_kategori = C.kode_kategori) "
		// 	sql += "WHERE A.identitas=?"
		// 	db.Raw(sql, identitas).Scan(&result)

		// 	callback["success"] = true
		// 	callback["data"] = result
		// 	// callback["idpersonalia"] = idpersonalia
		// 	DB, _ := db.DB()
		// 	DB.Close()

		// 	context.JSON(status, callback)
		// })

		group.GET("/get-kategori", func(context *gin.Context) {
			// db_go := db.KoneksiCore()
			var callback = gin.H{}
			status := 200
			// identitas := mainlib.GetKoperasiID(context)

			result := []map[string]interface{}{}
			db.Raw("SELECT * FROM jurnal_master_kategori ").Scan(&result)

			callback["success"] = true
			callback["data"] = result

			DB, _ := db.DB()
			DB.Close()

			context.JSON(status, callback)
		})

		group.POST("/add", func(context *gin.Context) {
			// db_go := db.KoneksiCore()
			var callback = gin.H{}
			kode_akun := context.PostForm("kode_akun")
			kode_kategori := context.PostForm("kode_kategori")
			akun := context.PostForm("akun")
			keterangan := context.PostForm("keterangan")
			identitas := context.Query("identitas")

			//CEK PARAMETER POST
			//callback["master"] = master

			//CEK DUPLIKASI
			result := map[string]interface{}{}
			check := db.Raw("SELECT * FROM jurnal_master_akun  WHERE kode_akun = ?", kode_akun).Scan(&result)

			if check.RowsAffected == 0 {
				data := map[string]interface{}{
					"kode_akun":     kode_akun,
					"kode_kategori": kode_kategori,
					"akun":          akun,
					"keterangan":    keterangan,
					"identitas":     identitas,
					"tetap":         0,
				}
				create := db.Table("jurnal_master_akun").Create(&data)

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

			DB, _ := db.DB()
			DB.Close()

			context.JSON(200, callback)
		})

		group.POST("/edit", func(context *gin.Context) {
			// db_go := db.KoneksiCore()
			var callback = gin.H{}
			kode_akun_lama := context.PostForm("kode_akun_lama")
			kode_akun := context.PostForm("kode_akun")
			kode_kategori := context.PostForm("kode_kategori")
			akun := context.PostForm("akun")
			keterangan := context.PostForm("keterangan")

			result := map[string]interface{}{}
			check := db.Raw("SELECT * FROM jurnal_master_akun WHERE kode_akun = ? AND kode_akun != ?", kode_akun, kode_akun_lama).Scan(&result)

			if check.RowsAffected == 0 {
				update := db.Exec("UPDATE jurnal_master_akun SET kode_akun=?, kode_kategori=?, akun=?, keterangan=? WHERE kode_akun=?", kode_akun, kode_kategori, akun, keterangan, kode_akun_lama)

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

			DB, _ := db.DB()
			DB.Close()

			context.JSON(200, callback)
		})

		group.POST("/delete", func(context *gin.Context) {
			// db_go := db.KoneksiCore()
			var callback = gin.H{}
			kode_akun := context.PostForm("kode_akun")

			result := db.Exec("DELETE FROM jurnal_master_akun WHERE kode_akun = ?", kode_akun)
			if result.Error == nil {
				callback["success"] = true
				callback["msg"] = "Data berhasil dihapus"
			} else {
				callback["success"] = false
				callback["msg"] = "Hapus Gagal"
			}

			//CEK PARAMETER POST

			DB, _ := db.DB()
			DB.Close()

			context.JSON(200, callback)
		})
	}
}

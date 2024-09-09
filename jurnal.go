package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Jurnal struct {
	IdJurnal   int
	Noref      string
	Tanggal    time.Time
	Keterangan string
	IdKoperasi int
}

func JurnalRoute(router *gin.Engine) {
	group := router.Group("/keuangan/jurnal")
	{
		group.GET("/umum", func(context *gin.Context) {
			// db_go := db.KoneksiCore()
			var callback = gin.H{}

			result := []map[string]interface{}{}
			sql := "SELECT *,(SELECT SUM(jurnal_trx_subakun.jumlah) FROM jurnal_trx_subakun LEFT JOIN jurnal_trx ON (jurnal_trx_subakun.idjurnal = jurnal_trx.idjurnal) WHERE jurnal_trx_subakun.dk = ? AND jurnal_trx_subakun.idjurnal = jurnal_trx.idjurnal AND jurnal_trx.idjurnal = A.idjurnal) AS debit, (SELECT SUM(jurnal_trx_subakun.jumlah) FROM jurnal_trx_subakun LEFT JOIN jurnal_trx ON (jurnal_trx_subakun.idjurnal = jurnal_trx.idjurnal) WHERE jurnal_trx_subakun.dk = ? AND jurnal_trx_subakun.idjurnal = jurnal_trx.idjurnal AND jurnal_trx.idjurnal = A.idjurnal) AS kredit FROM jurnal_trx AS A ORDER BY A.tanggal DESC, A.idjurnal DESC"
			db.Raw(sql, "D", "K").Scan(&result)

			callback["success"] = true
			callback["data"] = result
			DB, _ := db.DB()
			DB.Close()

			context.JSON(200, callback)
		})

		group.GET("/umum-akun", func(context *gin.Context) {
			// db_go := db.KoneksiCore()
			var callback = gin.H{}

			idjurnal := context.Query("idjurnal")

			result := []map[string]interface{}{}
			sql := "SELECT * FROM jurnal_trx_subakun AS A INNER JOIN jurnal_master_subakun AS B ON (A.kode_subakun = B.kode_subakun) INNER JOIN jurnal_master_akun AS C ON (B.kode_akun = C.kode_akun) WHERE A.idjurnal = ?"
			db.Raw(sql, idjurnal).Scan(&result)

			callback["success"] = true
			callback["data"] = result
			callback["z"] = idjurnal
			DB, _ := db.DB()
			DB.Close()

			context.JSON(200, callback)
		})

		group.POST("/transaksi", func(context *gin.Context) {
			// db_go := db.KoneksiCore()
			var callback = gin.H{}
			// idkoperasi := mainlib.GetKoperasiID(context)
			idkoperasi := context.Query("idkoperasi")

			dataSubakun := context.PostForm("dataSubakun")
			noref := context.PostForm("no_ref")
			tanggal := context.PostForm("tanggal")
			keterangan := context.PostForm("keterangan")

			total_debit := 0
			total_kredit := 0
			// length_subakun := len(dataSubakun)
			// fmt.Println(length_subakun)

			tipe := context.PostForm("type_")
			fmt.Println(tipe)
			dataSubakunint, _ := strconv.Atoi(dataSubakun)

			for i := 0; i < dataSubakunint; i++ {
				tipe = context.PostForm("type_" + strconv.Itoa(i))
				if tipe == "D" {
					nilai_debit, _ := strconv.Atoi(context.PostForm("nilai_" + strconv.Itoa(i)))
					total_debit += nilai_debit
				} else {
					total_debit += 0
				}

				if tipe == "K" {
					nilai_kredit, _ := strconv.Atoi(context.PostForm("nilai_" + strconv.Itoa(i)))
					total_kredit += nilai_kredit
				} else {
					total_kredit += 0
				}
			}

			if total_debit == total_kredit && total_debit > 0 && total_kredit > 0 {
				type Jurnal struct {
					Idjurnal      int `gorm:"primaryKey"` // returm lastinsertID
					Noref         string
					Tanggal       string
					Keterangan    string
					Idkoperasi    string
					Tipe          string
					Posisi        string
					Tanggal_input string
					Tanggal_ubah  string
				}

				jurnal := Jurnal{
					Noref:         noref,
					Tanggal:       tanggal,
					Keterangan:    keterangan,
					Idkoperasi:    idkoperasi,
					Tipe:          "manual",
					Posisi:        "berjalan",
					Tanggal_input: tanggal,
					Tanggal_ubah:  tanggal,
				}

				db.Table("jurnal_trx").Create(&jurnal)

				for i := 0; i < dataSubakunint; i++ {
					jurnal_trx_subakun := map[string]interface{}{
						"idjurnal":     jurnal.Idjurnal,
						"kode_subakun": context.PostForm("kode_subakun_" + strconv.Itoa(i)),
						"dk":           context.PostForm("type_" + strconv.Itoa(i)),
						"jumlah":       context.PostForm("nilai_" + strconv.Itoa(i)),
					}
					db.Table("jurnal_trx_subakun").Create(&jurnal_trx_subakun)
				}

				callback["success"] = true
				callback["msg"] = "Data Berhasil ditambahkan"
			} else {
				callback["success"] = false
				callback["msg"] = "jurnal anda tidak balance, silahkan cek ulang"
			}
			DB, _ := db.DB()
			DB.Close()

			context.JSON(200, callback)
		})
	}
}

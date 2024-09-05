package keuangan

import (
	"apigo/lib/db"

	"github.com/gin-gonic/gin"

	// "strconv"
	// "apigo/lib/mainlib"
	"fmt"
)

func LaporanNeracaRoute(router *gin.Engine) {
	group := router.Group("/keuangan/laporan-neraca")
	{
		group.GET("/get", func(context *gin.Context) {
			db_go := db.KoneksiCore()
			var callback = gin.H{}

			tahun := context.Query("tahun")
			bulan := context.Query("bulan")
			shu_akun := "322"

			// SHU BERJALAN
			phu_ := []map[string]interface{}{}
			sql := "SELECT * FROM jurnal_master_kategori WHERE tipe = ? OR tipe = ? OR tipe = ? "
			db_go.Raw(sql, "laba", "rugi", "pajak").Scan(&phu_)

			where := ""
			fmt.Println(where)
			if bulan != "" {
				where = " AND YEAR(jurnal_trx.tanggal) = " + tahun + " AND MONTH(jurnal_trx.tanggal) <= " + bulan
			} else {
				where = " AND YEAR(jurnal_trx.tanggal) = " + tahun
			}

			var total_pendapatan float64
			var total_beban float64
			var total_pajak float64
			var shu float64

			for _, phu := range phu_ {
				fmt.Println(phu)
				resultShu := []map[string]interface{}{}
				sql := "SELECT *, IFNULL((SELECT SUM(jurnal_trx_subakun.jumlah) FROM jurnal_trx, jurnal_trx_subakun, jurnal_master_subakun WHERE jurnal_trx.idjurnal = jurnal_trx_subakun.idjurnal AND jurnal_master_subakun.kode_subakun = jurnal_trx_subakun.kode_subakun AND jurnal_trx_subakun.dk = 'D' AND jurnal_master_subakun.kode_akun = jurnal_master_akun.kode_akun " + where + "),0) AS debit, "

				sql += "IFNULL((SELECT SUM(jurnal_trx_subakun.jumlah) FROM jurnal_trx, jurnal_trx_subakun, jurnal_master_subakun WHERE jurnal_trx.idjurnal = jurnal_trx_subakun.idjurnal AND jurnal_master_subakun.kode_subakun = jurnal_trx_subakun.kode_subakun AND jurnal_trx_subakun.dk = 'K' AND jurnal_master_subakun.kode_akun = jurnal_master_akun.kode_akun " + where + "),0 ) AS kredit "

				sql += "FROM jurnal_master_akun INNER JOIN jurnal_master_kategori ON (jurnal_master_kategori.kode_kategori = jurnal_master_akun.kode_kategori) WHERE jurnal_master_kategori.kode_kategori = ? "

				db_go.Raw(sql, phu["kode_kategori"].(string)).Scan(&resultShu)

				for _, rshu := range resultShu {
					if rshu["tipe"].(string) == "laba" {
						total_pendapatan += (rshu["kredit"].(float64) - rshu["debit"].(float64))
					}
					if rshu["tipe"].(string) == "rugi" {
						total_beban += (rshu["debit"].(float64) - rshu["kredit"].(float64))
					} else {
						total_pajak += (rshu["debit"].(float64) - rshu["kredit"].(float64))
					}
				}

				shu = (total_pendapatan - total_beban - total_pajak)
				fmt.Println(shu)
			}

			// ASET
			tempAset := []map[string]interface{}{}
			aset_ := []map[string]interface{}{}
			sqltwo := "SELECT * FROM jurnal_master_kategori WHERE tipe = ? OR tipe = ? "
			db_go.Raw(sqltwo, "ASET_LANCAR", "ASET_TDK_LANCAR").Scan(&aset_)

			where_aset := ""
			if bulan != "" {
				where_aset = " AND (YEAR(jurnal_trx.tanggal)) < " + tahun + " OR (YEAR(jurnal_trx.tanggal)) =  " + tahun + " AND MONTH(jurnal_trx.tanggal) <= " + bulan
			} else {
				where_aset = " AND YEAR(jurnal_trx.tanggal) = " + tahun
			}

			for _, aset := range aset_ {
				resultGolongan := []map[string]interface{}{}
				sql := "SELECT *, IFNULL((SELECT SUM(jurnal_trx_subakun.jumlah) FROM jurnal_trx, jurnal_trx_subakun, jurnal_master_subakun WHERE jurnal_trx.idjurnal = jurnal_trx_subakun.idjurnal AND jurnal_master_subakun.kode_subakun = jurnal_trx_subakun.kode_subakun AND jurnal_trx_subakun.dk = 'D' AND jurnal_master_subakun.kode_akun = jurnal_master_akun.kode_akun " + where_aset + "),0) AS debit, "

				sql += "IFNULL((SELECT SUM(jurnal_trx_subakun.jumlah) FROM jurnal_trx, jurnal_trx_subakun, jurnal_master_subakun WHERE jurnal_trx.idjurnal = jurnal_trx_subakun.idjurnal AND jurnal_master_subakun.kode_subakun = jurnal_trx_subakun.kode_subakun AND jurnal_trx_subakun.dk = 'K' AND jurnal_master_subakun.kode_akun = jurnal_master_akun.kode_akun " + where_aset + "),0 ) AS kredit "

				sql += "FROM jurnal_master_akun INNER JOIN jurnal_master_kategori ON (jurnal_master_kategori.kode_kategori = jurnal_master_akun.kode_kategori) WHERE jurnal_master_kategori.kode_kategori = ? "

				db_go.Raw(sql, aset["kode_kategori"].(string)).Scan(&resultGolongan)

				temp := []map[string]interface{}{}
				var total_aset float64
				fmt.Println(total_aset)
				for _, gol := range resultGolongan {
					gol["saldo"] = (gol["debit"].(float64) - gol["kredit"].(float64))
					total_aset += (gol["debit"].(float64) - gol["kredit"].(float64))

					temp = append(temp, gol)
				}

				aset["akun"] = temp
				aset["saldo"] = total_aset

				tempAset = append(tempAset, aset)
			}

			// EKUITAS
			tempEkuitas := []map[string]interface{}{}
			ekuitas_ := []map[string]interface{}{}
			sqlthree := "SELECT * FROM jurnal_master_kategori WHERE tipe = ? OR tipe = ?"
			db_go.Raw(sqlthree, "EKUITAS", "KEWAJIBAN").Scan(&ekuitas_)

			for _, ekuitas := range ekuitas_ {
				resultGolongan := []map[string]interface{}{}
				sql := "SELECT *, IFNULL((SELECT SUM(jurnal_trx_subakun.jumlah) FROM jurnal_trx, jurnal_trx_subakun, jurnal_master_subakun WHERE jurnal_trx.idjurnal = jurnal_trx_subakun.idjurnal AND jurnal_master_subakun.kode_subakun = jurnal_trx_subakun.kode_subakun AND jurnal_trx_subakun.dk = 'D' AND jurnal_master_subakun.kode_akun = jurnal_master_akun.kode_akun " + where_aset + "),0) AS debit, "

				sql += "IFNULL((SELECT SUM(jurnal_trx_subakun.jumlah) FROM jurnal_trx, jurnal_trx_subakun, jurnal_master_subakun WHERE jurnal_trx.idjurnal = jurnal_trx_subakun.idjurnal AND jurnal_master_subakun.kode_subakun = jurnal_trx_subakun.kode_subakun AND jurnal_trx_subakun.dk = 'K' AND jurnal_master_subakun.kode_akun = jurnal_master_akun.kode_akun " + where_aset + "),0 ) AS kredit "

				sql += "FROM jurnal_master_akun INNER JOIN jurnal_master_kategori ON (jurnal_master_kategori.kode_kategori = jurnal_master_akun.kode_kategori) WHERE jurnal_master_kategori.kode_kategori = ? "

				db_go.Raw(sql, ekuitas["kode_kategori"].(string)).Scan(&resultGolongan)

				temp := []map[string]interface{}{}
				var total_ekuitas float64

				for _, gol := range resultGolongan {
					var gol_saldo float64
					if gol["kode_akun"].(string) == shu_akun {
						gol_saldo = (gol["kredit"].(float64) - gol["debit"].(float64)) + shu
					} else {
						gol_saldo = (gol["kredit"].(float64) - gol["debit"].(float64))
					}

					total_ekuitas += gol_saldo
					gol["saldo"] = gol_saldo
					temp = append(temp, gol)
				}

				ekuitas["akun"] = temp
				ekuitas["saldo"] = total_ekuitas
				tempEkuitas = append(tempEkuitas, ekuitas)
			}

			datas := map[string]interface{}{
				"aset":    tempAset,
				"ekuitas": tempEkuitas,
			}

			callback["success"] = true
			callback["data"] = datas
			DB, _ := db_go.DB()
			DB.Close()

			context.JSON(200, callback)
		})
	}
}

package main

import (
	"github.com/gin-gonic/gin"

	"fmt"
)

func NeracaSaldoRoute(router *gin.Engine) {
	group := router.Group("/keuangan/neraca-saldo")
	{
		group.GET("/get", func(context *gin.Context) {
			// db_go := db.KoneksiCore()
			var callback = gin.H{}

			tahun := context.Query("tahun")
			bulan := context.Query("bulan")

			where := ""
			where_awal := ""

			if bulan != "" {
				where = " AND MONTH(jurnal_trx.tanggal) = " + bulan

				where_awal = " AND jurnal_trx.tanggal < " + tahun + "-" + bulan + "-01"
			}

			dataKategori := []map[string]interface{}{}
			sql := "SELECT * FROM jurnal_master_kategori WHERE tipe = ? OR tipe = ? OR tipe = ? OR tipe = ?"
			db.Raw(sql, "ASET_LANCAR", "ASET_TDK_LANCAR", "KEWAJIBAN", "EKUITAS").Scan(&dataKategori)

			tempKel := []map[string]interface{}{}
			for _, kategori := range dataKategori {
				resultGroup := []map[string]interface{}{}
				fmt.Println(kategori["kode_kategori"].(string))
				sql := "SELECT *, "
				sql += "IFNULL((SELECT SUM(jurnal_trx_subakun.jumlah) FROM jurnal_trx, jurnal_trx_subakun, jurnal_master_subakun WHERE jurnal_trx.idjurnal = jurnal_trx_subakun.idjurnal AND jurnal_master_subakun.kode_subakun = jurnal_trx_subakun.kode_subakun AND jurnal_trx_subakun.dk = 'D' AND jurnal_master_subakun.kode_akun = jurnal_master_akun.kode_akun AND YEAR(jurnal_trx.tanggal) = ? " + where + "),0) AS debit_saldo, "
				sql += "IFNULL((SELECT SUM(jurnal_trx_subakun.jumlah) FROM jurnal_trx, jurnal_trx_subakun, jurnal_master_subakun WHERE jurnal_trx.idjurnal = jurnal_trx_subakun.idjurnal AND jurnal_master_subakun.kode_subakun = jurnal_trx_subakun.kode_subakun AND jurnal_trx_subakun.dk = 'K' AND jurnal_master_subakun.kode_akun = jurnal_master_akun.kode_akun AND YEAR(jurnal_trx.tanggal) = ? " + where + "),0) AS kredit_saldo, "
				sql += "IFNULL((SELECT SUM(jurnal_trx_subakun.jumlah) FROM jurnal_trx, jurnal_trx_subakun, jurnal_master_subakun WHERE jurnal_trx.idjurnal = jurnal_trx_subakun.idjurnal AND jurnal_master_subakun.kode_subakun = jurnal_trx_subakun.kode_subakun AND jurnal_trx_subakun.dk = 'D' AND jurnal_master_subakun.kode_akun = jurnal_master_akun.kode_akun " + where_awal + "),0) AS debit_saldoawal, "
				sql += "IFNULL((SELECT SUM(jurnal_trx_subakun.jumlah) FROM jurnal_trx, jurnal_trx_subakun, jurnal_master_subakun WHERE jurnal_trx.idjurnal = jurnal_trx_subakun.idjurnal AND jurnal_master_subakun.kode_subakun = jurnal_trx_subakun.kode_subakun AND jurnal_trx_subakun.dk = 'K' AND jurnal_master_subakun.kode_akun = jurnal_master_akun.kode_akun " + where_awal + "),0) AS kredit_saldoawal "
				sql += "FROM jurnal_master_akun INNER JOIN jurnal_master_kategori ON (jurnal_master_kategori.kode_kategori = jurnal_master_akun.kode_kategori) WHERE jurnal_master_kategori.kode_kategori = ? "
				db.Raw(sql, tahun, tahun, kategori["kode_kategori"].(string)).Scan(&resultGroup)

				temp := []map[string]interface{}{}
				var kategori_saldo float64
				var kategori_saldo_awal float64

				for _, gol := range resultGroup {
					var group_saldo_awal float64
					fmt.Println(group_saldo_awal)
					var saldo float64
					fmt.Println(saldo)

					if kategori["tipe"].(string) == "ASET_LANCAR" {
						group_saldo_awal = (gol["debit_saldoawal"].(float64) - gol["kredit_saldoawal"].(float64))
						saldo = (gol["debit_saldo"].(float64) - gol["kredit_saldo"].(float64))
					} else {
						group_saldo_awal = (gol["kredit_saldoawal"].(float64) - gol["debit_saldoawal"].(float64))
						saldo = (gol["kredit_saldo"].(float64) - gol["debit_saldo"].(float64))
					}

					gol["saldo"] = (group_saldo_awal + saldo)
					// fmt.Println(group_saldo)
					kategori_saldo_awal += group_saldo_awal
					kategori_saldo += (group_saldo_awal + saldo)

					temp = append(temp, gol)
				}
				kategori["akun"] = temp
				kategori["saldo"] = kategori_saldo
				kategori["saldo_awal"] = kategori_saldo_awal

				tempKel = append(tempKel, kategori)
			}

			callback["success"] = true
			callback["data"] = tempKel

			DB, _ := db.DB()
			DB.Close()

			context.JSON(200, callback)
		})
	}
}

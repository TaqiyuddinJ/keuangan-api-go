package main

import (
	"github.com/gin-gonic/gin"

	"fmt"
)

func LaporanPhuRoute(router *gin.Engine) {
	group := router.Group("/keuangan/laporan-phu")
	{
		group.GET("/get", func(context *gin.Context) {
			// db_go := db.KoneksiCore()
			var callback = gin.H{}

			tahun := context.Query("tahun")
			bulan := context.Query("bulan")

			where := ""
			fmt.Println(where)
			if bulan != "" {
				where = " AND MONTH(jurnal_trx.tanggal) = " + bulan
			}

			//Kategori
			tempKel := []map[string]interface{}{}
			dataKelompok := []map[string]interface{}{}
			sql := "SELECT * FROM jurnal_master_kategori WHERE tipe = ? OR tipe = ? "
			db.Raw(sql, "PENDAPATAN", "BEBAN").Scan(&dataKelompok)

			// fmt.Println(dataKelompok)
			for _, each := range dataKelompok {
				fmt.Println(each)
				resultGolongan := []map[string]interface{}{}
				sql := "SELECT *, IFNULL((SELECT SUM(jurnal_trx_subakun.jumlah) FROM jurnal_trx,jurnal_trx_subakun,jurnal_master_subakun WHERE jurnal_trx.idjurnal = jurnal_trx_subakun.idjurnal AND jurnal_master_subakun.kode_subakun = jurnal_trx_subakun.kode_subakun AND jurnal_trx_subakun.dk = 'D' AND jurnal_master_subakun.kode_akun = jurnal_master_akun.kode_akun AND YEAR(jurnal_trx.tanggal) = ? AND MONTH(jurnal_trx.tanggal) >= 1 AND MONTH(jurnal_trx.tanggal) <= ?),0) AS debit_akumulasi, "

				sql += "IFNULL((SELECT SUM(jurnal_trx_subakun.jumlah) FROM jurnal_trx,jurnal_trx_subakun,jurnal_master_subakun WHERE jurnal_trx.idjurnal = jurnal_trx_subakun.idjurnal AND jurnal_master_subakun.kode_subakun = jurnal_trx_subakun.kode_subakun AND jurnal_trx_subakun.dk = 'K' AND jurnal_master_subakun.kode_akun = jurnal_master_akun.kode_akun AND YEAR(jurnal_trx.tanggal) = ? AND MONTH(jurnal_trx.tanggal) >= 1 AND MONTH(jurnal_trx.tanggal) <= ?),0) AS kredit_akumulasi, "

				sql += "IFNULL((SELECT SUM(jurnal_trx_subakun.jumlah) FROM jurnal_trx, jurnal_trx_subakun, jurnal_master_subakun WHERE jurnal_trx.idjurnal = jurnal_trx_subakun.idjurnal AND jurnal_master_subakun.kode_subakun = jurnal_trx_subakun.kode_subakun AND jurnal_trx_subakun.dk = 'D' AND jurnal_master_subakun.kode_akun = jurnal_master_akun.kode_akun AND YEAR(jurnal_trx.tanggal) = ? " + where + "),0) AS debit_bulan, "

				sql += "IFNULL((SELECT SUM(jurnal_trx_subakun.jumlah) FROM jurnal_trx, jurnal_trx_subakun, jurnal_master_subakun WHERE jurnal_trx.idjurnal = jurnal_trx_subakun.idjurnal AND jurnal_master_subakun.kode_subakun = jurnal_trx_subakun.kode_subakun AND jurnal_trx_subakun.dk = 'K' AND jurnal_master_subakun.kode_akun = jurnal_master_akun.kode_akun AND YEAR(jurnal_trx.tanggal) = ? " + where + "),0 ) AS kredit_bulan "

				sql += "FROM jurnal_master_akun INNER JOIN jurnal_master_kategori ON (jurnal_master_kategori.kode_kategori = jurnal_master_akun.kode_kategori) WHERE jurnal_master_kategori.kode_kategori = ? "

				db.Raw(sql, tahun, bulan, tahun, bulan, tahun, tahun, each["kode_kategori"].(string)).Scan(&resultGolongan)

				temp := []map[string]interface{}{}
				var totalakumulasi_laba float64
				var totalbulan_laba float64

				for _, element := range resultGolongan {
					// var saldo_akumulasi_laba float64
					// var saldo_bulan_laba float64

					if element["tipe"] == "laba" {
						element["saldo_akumulasi_laba"] = (element["kredit_akumulasi"].(float64) - element["debit_akumulasi"].(float64))
						element["saldo_bulan_laba"] = (element["kredit_bulan"].(float64) - element["debit_bulan"].(float64))
						totalakumulasi_laba += (element["kredit_akumulasi"].(float64) - element["debit_akumulasi"].(float64))
						totalbulan_laba += (element["kredit_bulan"].(float64) - element["debit_bulan"].(float64))
						// element["saldo_akumulasi"] =
					} else {
						element["saldo_akumulasi_laba"] = (element["debit_akumulasi"].(float64) - element["kredit_akumulasi"].(float64))
						element["saldo_bulan_laba"] = (element["debit_bulan"].(float64) - element["kredit_bulan"].(float64))
						totalakumulasi_laba += (element["debit_akumulasi"].(float64) - element["kredit_akumulasi"].(float64))
						totalbulan_laba += (element["debit_bulan"].(float64) - element["kredit_bulan"].(float64))
					}
					// fmt.Println(element["saldo_akumulasi_laba"])
					// fmt.Println(element["saldo_bulan_laba"])
					temp = append(temp, element)
				}

				each["akun"] = temp
				each["total_akumulasi"] = totalakumulasi_laba
				each["total_bulan"] = totalbulan_laba

				tempKel = append(tempKel, each)
			}

			//Pajak
			tempPajakGlobal := []map[string]interface{}{}
			dataPajak := []map[string]interface{}{}
			sqlpajak := "SELECT * FROM jurnal_master_kategori WHERE tipe = ?"
			db.Raw(sqlpajak, "PAJAK").Scan(&dataPajak)

			for _, eachPajak := range dataPajak {
				fmt.Println(eachPajak)
				resultGolonganPajak := []map[string]interface{}{}
				sqlpajaktwo := "SELECT *, IFNULL((SELECT SUM(jurnal_trx_subakun.jumlah) FROM jurnal_trx,jurnal_trx_subakun,jurnal_master_subakun WHERE jurnal_trx.idjurnal = jurnal_trx_subakun.idjurnal AND jurnal_master_subakun.kode_subakun = jurnal_trx_subakun.kode_subakun AND jurnal_trx_subakun.dk = 'D' AND jurnal_master_subakun.kode_akun = jurnal_master_akun.kode_akun AND YEAR(jurnal_trx.tanggal) = ? AND MONTH(jurnal_trx.tanggal) >= 1 AND MONTH(jurnal_trx.tanggal) <= ?),0) AS debit_akumulasi, "

				sqlpajaktwo += "IFNULL((SELECT SUM(jurnal_trx_subakun.jumlah) FROM jurnal_trx,jurnal_trx_subakun,jurnal_master_subakun WHERE jurnal_trx.idjurnal = jurnal_trx_subakun.idjurnal AND jurnal_master_subakun.kode_subakun = jurnal_trx_subakun.kode_subakun AND jurnal_trx_subakun.dk = 'K' AND jurnal_master_subakun.kode_akun = jurnal_master_akun.kode_akun AND YEAR(jurnal_trx.tanggal) = ? AND MONTH(jurnal_trx.tanggal) >= 1 AND MONTH(jurnal_trx.tanggal) <= ?),0) AS kredit_akumulasi, "

				sqlpajaktwo += "IFNULL((SELECT SUM(jurnal_trx_subakun.jumlah) FROM jurnal_trx, jurnal_trx_subakun, jurnal_master_subakun WHERE jurnal_trx.idjurnal = jurnal_trx_subakun.idjurnal AND jurnal_master_subakun.kode_subakun = jurnal_trx_subakun.kode_subakun AND jurnal_trx_subakun.dk = 'D' AND jurnal_master_subakun.kode_akun = jurnal_master_akun.kode_akun AND YEAR(jurnal_trx.tanggal) = ? " + where + "),0) AS debit_bulan, "

				sqlpajaktwo += "IFNULL((SELECT SUM(jurnal_trx_subakun.jumlah) FROM jurnal_trx, jurnal_trx_subakun, jurnal_master_subakun WHERE jurnal_trx.idjurnal = jurnal_trx_subakun.idjurnal AND jurnal_master_subakun.kode_subakun = jurnal_trx_subakun.kode_subakun AND jurnal_trx_subakun.dk = 'K' AND jurnal_master_subakun.kode_akun = jurnal_master_akun.kode_akun AND YEAR(jurnal_trx.tanggal) = ? " + where + "),0 ) AS kredit_bulan "

				sqlpajaktwo += "FROM jurnal_master_akun INNER JOIN jurnal_master_kategori ON (jurnal_master_kategori.kode_kategori = jurnal_master_akun.kode_kategori) WHERE jurnal_master_kategori.kode_kategori = ? "

				db.Raw(sqlpajaktwo, tahun, bulan, tahun, bulan, tahun, tahun, eachPajak["kode_kategori"].(string)).Scan(&resultGolonganPajak)

				tempPajak := []map[string]interface{}{}
				var totalakumulasi_pajak float64 = 0
				var totalbulan_pajak float64 = 0

				for _, elementPajak := range resultGolonganPajak {

					// var saldo_akumulasi_pajak float64
					// var saldo_bulan_pajak float64
					if elementPajak["tipe"] == "pajak" {
						elementPajak["saldo_akumulasi_pajak"] = (elementPajak["kredit_akumulasi"].(float64) - elementPajak["debit_akumulasi"].(float64))
						elementPajak["saldo_bulan_pajak"] = (elementPajak["kredit_bulan"].(float64) - elementPajak["debit_bulan"].(float64))
						totalakumulasi_pajak += (elementPajak["kredit_akumulasi"].(float64) - elementPajak["debit_akumulasi"].(float64))
						totalbulan_pajak += (elementPajak["kredit_bulan"].(float64) - elementPajak["debit_bulan"].(float64))
						// elementPajak["saldo_akumulasi"] =
					} else {
						elementPajak["saldo_akumulasi_pajak"] = (elementPajak["debit_akumulasi"].(float64) - elementPajak["kredit_akumulasi"].(float64))
						elementPajak["saldo_bulan_pajak"] = (elementPajak["debit_bulan"].(float64) - elementPajak["kredit_bulan"].(float64))
						totalakumulasi_pajak += (elementPajak["debit_akumulasi"].(float64) - elementPajak["kredit_akumulasi"].(float64))
						totalbulan_pajak += (elementPajak["debit_bulan"].(float64) - elementPajak["kredit_bulan"].(float64))
					}

					tempPajak = append(tempPajak, elementPajak)
				}

				eachPajak["akun"] = tempPajak
				eachPajak["total_akumulasi"] = totalakumulasi_pajak
				eachPajak["total_bulan"] = totalbulan_pajak

				tempPajakGlobal = append(tempPajakGlobal, eachPajak)
			}

			callback["success"] = true
			callback["data"] = tempKel
			callback["pajak"] = tempPajakGlobal
			DB, _ := db.DB()
			DB.Close()

			context.JSON(200, callback)
		})
	}
}

package keuangan

import (
	"apigo/lib/db"

	"github.com/gin-gonic/gin"

	// "strconv"
	// "apigo/lib/mainlib"
	"fmt"
	"sort"

	"github.com/xuri/excelize/v2"
)

func BukuBesarRoute(router *gin.Engine) {
	group := router.Group("/keuangan/buku-besar")
	{
		group.GET("/get", func(context *gin.Context) {
			db_go := db.KoneksiCore()
			var callback = gin.H{}

			start := context.Query("start")
			end := context.Query("end")

			if start == "null" && end == "null" {

				result := []map[string]interface{}{}
				sql := "SELECT jurnal_master_subakun.*, jurnal_master_akun.akun, "
				sql += "IFNULL((SELECT SUM(jurnal_trx_subakun.nilai) FROM jurnal_trx_subakun WHERE jurnal_trx_subakun.kode_subakun=jurnal_master_subakun.kode_subakun AND jurnal_trx_subakun.dk='D'),0) AS debit, "
				sql += "IFNULL((SELECT SUM(jurnal_trx_subakun.nilai) FROM jurnal_trx_subakun WHERE jurnal_trx_subakun.kode_subakun=jurnal_master_subakun.kode_subakun AND jurnal_trx_subakun.dk='K'),0) AS kredit "
				sql += "FROM jurnal_master_subakun "
				sql += "INNER JOIN jurnal_master_akun ON (jurnal_master_akun.kode_akun=jurnal_master_subakun.kode_akun)"
				db_go.Raw(sql).Scan(&result)

				callback["success"] = true
				callback["data"] = result
			} else {

				result := []map[string]interface{}{}
				sql := "SELECT *, "
				sql += "IFNULL((SELECT SUM(jurnal_trx_subakun.jumlah) FROM jurnal_trx_subakun "
				sql += "LEFT JOIN jurnal_trx ON (jurnal_trx_subakun.idjurnal = jurnal_trx.idjurnal) "
				sql += "WHERE jurnal_trx_subakun.dk = 'D' AND jurnal_trx_subakun.kode_subakun = A.kode_subakun), 0) AS debit, "
				sql += "IFNULL((SELECT SUM(jurnal_trx_subakun.jumlah) FROM jurnal_trx_subakun "
				sql += "LEFT JOIN jurnal_trx ON (jurnal_trx_subakun.idjurnal = jurnal_trx.idjurnal) "
				sql += "WHERE jurnal_trx_subakun.dk = 'K' "
				sql += "AND jurnal_trx_subakun.kode_subakun = A.kode_subakun), 0) AS kredit "
				sql += "FROM jurnal_trx_subakun AS A INNER JOIN jurnal_master_subakun AS B ON (A.kode_subakun = B.kode_subakun) INNER JOIN jurnal_master_akun AS C ON (B.kode_akun = C.kode_akun) INNER JOIN jurnal_trx AS D ON (A.idjurnal = D.idjurnal) "
				sql += "WHERE D.tanggal >= ? AND  D.tanggal <= ? GROUP BY A.kode_subakun"
				db_go.Raw(sql, start, end).Scan(&result)

				callback["success"] = true
				callback["data"] = result
			}
			DB, _ := db_go.DB()
			DB.Close()

			context.JSON(200, callback)
		})

		group.GET("/detail", func(context *gin.Context) {
			db_go := db.KoneksiCore()
			var callback = gin.H{}

			start := context.Query("start") + " 00:00:00"
			end := context.Query("end") + " 23:59:59"
			kode_subakun := context.Query("kode_subakun")

			if start == "null" && end == "null" {

				result := []map[string]interface{}{}
				sql := "SELECT * FROM jurnal_trx_subakun AS A INNER JOIN jurnal_master_subakun AS B ON (A.kode_subakun = B.kode_subakun) "
				sql += "LEFT JOIN jurnal_trx AS C ON (A.idjurnal = C.idjurnal) "
				sql += "WHERE A.kode_subakun = ? ORDER BY C.idjurnal ASC "
				db_go.Raw(sql, kode_subakun).Scan(&result)

				dk := "D"
				saldo := 0
				temp := []map[string]interface{}{}

				for _, element := range result {
					jumlah := int(element["jumlah"].(float64))
					if element["dk"] == dk {
						saldo += jumlah
					} else {
						if saldo >= jumlah {
							saldo -= jumlah
						} else {
							saldo = (jumlah - saldo)
						}
					}

					element["saldo"] = saldo
					temp = append(temp, element)
					fmt.Println(temp)
				}

				temp_awal := map[string]interface{}{
					"approval":     "0",
					"dk":           "D",
					"idtransaksi":  "0",
					"keterangan":   "Saldo Awal",
					"kode_akun":    "0",
					"kode_subakun": "0",
					"nilai":        "0",
					"no_ref":       "null",
					"saldo":        "0",
					"subakun":      "-",
					"tanggal":      "00-00-00 00:00:00",
				}

				//sort array
				sort.Slice(temp, func(i, j int) bool {
					return true
				})

				temp = append(temp, temp_awal)

				callback["success"] = true
				callback["data"] = temp

			} else {

				dataAwal := map[string]interface{}{}
				sql := "SELECT *, "
				sql += "IFNULL((SELECT SUM(jurnal_trx_subakun.jumlah) FROM jurnal_trx_subakun LEFT JOIN jurnal_trx ON (jurnal_trx_subakun.idjurnal = jurnal_trx.idjurnal) "
				sql += "WHERE jurnal_trx_subakun.kode_subakun=jurnal_master_subakun.kode_subakun AND jurnal_trx_subakun.dk='D' AND jurnal_trx.tanggal < ?),0) AS debit, "
				sql += "IFNULL((SELECT SUM(jurnal_trx_subakun.jumlah) FROM jurnal_trx_subakun LEFT JOIN jurnal_trx ON (jurnal_trx_subakun.idjurnal = jurnal_trx.idjurnal) "
				sql += "WHERE jurnal_trx_subakun.kode_subakun=jurnal_master_subakun.kode_subakun AND jurnal_trx_subakun.dk='K' AND jurnal_trx.tanggal < ?),0) AS kredit "
				sql += "FROM jurnal_master_subakun WHERE jurnal_master_subakun.kode_subakun = ? "
				db_go.Raw(sql, start, start, kode_subakun).Scan(&dataAwal)

				dataSubakun := map[string]interface{}{}
				sqltwo := "SELECT * FROM jurnal_master_subakun AS A WHERE A.kode_subakun = ?"
				db_go.Raw(sqltwo, kode_subakun).Scan(&dataSubakun)

				result := []map[string]interface{}{}
				sqlthree := "SELECT * FROM jurnal_trx_subakun AS A INNER JOIN jurnal_master_subakun AS B ON (A.kode_subakun = B.kode_subakun) "
				sqlthree += "LEFT JOIN jurnal_trx AS C ON (A.idjurnal = C.idjurnal) "
				sqlthree += "WHERE A.kode_subakun = ? AND C.tanggal >= ? AND C.tanggal <= ? ORDER BY C.idjurnal ASC "
				db_go.Raw(sqlthree, kode_subakun, start, end).Scan(&result)

				dk := ""
				fmt.Println(dk)
				saldo := 0
				if int(dataAwal["debit"].(float64)) > int(dataAwal["kredit"].(float64)) {
					dk = "D"
				} else {
					dk = "K"
				}

				if int(dataAwal["debit"].(float64)) > int(dataAwal["kredit"].(float64)) {
					saldo = (int(dataAwal["debit"].(float64)) - int(dataAwal["kredit"].(float64)))
				} else {
					saldo = (int(dataAwal["kredit"].(float64)) - int(dataAwal["debit"].(float64)))
				}

				temp := []map[string]interface{}{}

				for _, element := range result {
					jumlah := int(element["jumlah"].(float64))
					if element["dk"] == dk {
						saldo += jumlah
					} else {
						if saldo >= jumlah {
							saldo -= jumlah
						} else {
							saldo = (jumlah - saldo)
						}
					}

					element["saldo"] = saldo
					temp = append(temp, element)
					fmt.Println(temp)
				}

				temp_awal := map[string]interface{}{
					"approval":     0,
					"debit":        int(dataAwal["debit"].(float64)),
					"dk":           dk,
					"idtransaksi":  "0",
					"keterangan":   "Saldo Awal",
					"kode_akun":    dataAwal["kode_akun"].(string),
					"kode_subakun": dataAwal["kode_subakun"].(string),
					"kredit":       int(dataAwal["kredit"].(float64)),
					"nilai":        saldo,
					"saldo":        saldo,
					"no_ref":       "null",
					"subakun":      dataAwal["subakun"].(string),
					"tanggal":      "00-00-00 00:00:00",
				}

				//sort array
				sort.Slice(temp, func(i, j int) bool {
					return true
				})

				temp = append(temp, temp_awal)

				callback["success"] = true
				callback["data"] = temp
				callback["subakun"] = dataSubakun
			}
			DB, _ := db_go.DB()
			DB.Close()

			context.JSON(200, callback)
		})

		group.GET("/detail/export", func(context *gin.Context) {
			db_go := db.KoneksiCore()
			var callback = gin.H{}

			start := context.Query("start") + " 00:00:00"
			end := context.Query("end") + " 23:59:59"
			kode_subakun := context.Query("kode_subakun")

			dataAwal := map[string]interface{}{}
			sql := "SELECT *, "
			sql += "IFNULL((SELECT SUM(jurnal_trx_subakun.jumlah) FROM jurnal_trx_subakun LEFT JOIN jurnal_trx ON (jurnal_trx_subakun.idjurnal = jurnal_trx.idjurnal) "
			sql += "WHERE jurnal_trx_subakun.kode_subakun=jurnal_master_subakun.kode_subakun AND jurnal_trx_subakun.dk='D' AND jurnal_trx.tanggal < ?),0) AS debit, "
			sql += "IFNULL((SELECT SUM(jurnal_trx_subakun.jumlah) FROM jurnal_trx_subakun LEFT JOIN jurnal_trx ON (jurnal_trx_subakun.idjurnal = jurnal_trx.idjurnal) "
			sql += "WHERE jurnal_trx_subakun.kode_subakun=jurnal_master_subakun.kode_subakun AND jurnal_trx_subakun.dk='K' AND jurnal_trx.tanggal < ?),0) AS kredit "
			sql += "FROM jurnal_master_subakun WHERE jurnal_master_subakun.kode_subakun = ? "
			db_go.Raw(sql, start, start, kode_subakun).Scan(&dataAwal)

			dataSubakun := map[string]interface{}{}
			sqltwo := "SELECT * FROM jurnal_master_subakun AS A WHERE A.kode_subakun = ?"
			db_go.Raw(sqltwo, kode_subakun).Scan(&dataSubakun)

			result := []map[string]interface{}{}
			sqlthree := "SELECT * FROM jurnal_trx_subakun AS A INNER JOIN jurnal_master_subakun AS B ON (A.kode_subakun = B.kode_subakun) "
			sqlthree += "LEFT JOIN jurnal_trx AS C ON (A.idjurnal = C.idjurnal) "
			sqlthree += "WHERE A.kode_subakun = ? AND C.tanggal >= ? AND C.tanggal <= ? ORDER BY C.idjurnal ASC "
			db_go.Raw(sqlthree, kode_subakun, start, end).Scan(&result)

			dk := ""
			fmt.Println(dk)
			saldo := 0
			if int(dataAwal["debit"].(float64)) > int(dataAwal["kredit"].(float64)) {
				dk = "D"
			} else {
				dk = "K"
			}

			if int(dataAwal["debit"].(float64)) > int(dataAwal["kredit"].(float64)) {
				saldo = (int(dataAwal["debit"].(float64)) - int(dataAwal["kredit"].(float64)))
			} else {
				saldo = (int(dataAwal["kredit"].(float64)) - int(dataAwal["debit"].(float64)))
			}

			temp := []map[string]interface{}{}

			for _, element := range result {
				jumlah := int(element["jumlah"].(float64))
				if element["dk"] == dk {
					saldo += jumlah
				} else {
					if saldo >= jumlah {
						saldo -= jumlah
					} else {
						saldo = (jumlah - saldo)
					}
				}

				element["saldo"] = saldo
				temp = append(temp, element)
				fmt.Println(temp)
			}

			temp_awal := map[string]interface{}{
				"approval":     0,
				"debit":        int(dataAwal["debit"].(float64)),
				"dk":           dk,
				"idtransaksi":  "0",
				"keterangan":   "Saldo Awal",
				"kode_akun":    dataAwal["kode_akun"].(string),
				"kode_subakun": dataAwal["kode_subakun"].(string),
				"kredit":       int(dataAwal["kredit"].(float64)),
				"jumlah":       saldo,
				"saldo":        saldo,
				"no_ref":       "null",
				"subakun":      dataAwal["subakun"].(string),
				"tanggal":      "00-00-00 00:00:00",
			}

			//sort array
			sort.Slice(temp, func(i, j int) bool {
				return true
			})

			temp = append(temp, temp_awal)

			f, err := excelize.OpenFile("../template/TemplateBukuBesar.xlsx")
			if err != nil {
				fmt.Println(err)
				return
			}
			// Get value from cell by given worksheet name and axis.
			cell, err := f.GetCellValue("Sheet1", "B2")
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(cell)
			// Get all the rows in the Sheet1.
			// rows, err := f.GetRows("Sheet1")
			// for _, row := range rows {
			//     for _, colCell := range row {
			//         fmt.Print(colCell, "\t")
			//     }
			//     fmt.Println()
			// }

			sheet1Name := "Sheet One"
			f.SetSheetName(f.GetSheetName(1), sheet1Name)

			// f.SetCellValue("Sheet1", "A3", "No")
			// f.SetCellValue("Sheet1", "B3", "Tanggal")
			// f.SetCellValue("Sheet1", "C3", "Keterangan")
			// f.SetCellValue("Sheet1", "D3", "Debit")
			// f.SetCellValue("Sheet1", "E3", "Kredit")
			// f.SetCellValue("Sheet1", "F3", "Saldo")

			//isi dengan foreach
			for i, each := range temp {
				fmt.Println(each["jumlah"])
				var debit float64
				fmt.Println(debit)
				var kredit float64
				fmt.Println(kredit)
				if each["dk"].(string) == "D" {
					// debit = int(each["jumlah"].(float64))
					debit = each["jumlah"].(float64)
				}
				if each["dk"].(string) == "K" {
					// kredit = int(each["jumlah"].(float64))
					kredit = each["jumlah"].(float64)
				}

				f.SetCellValue(sheet1Name, fmt.Sprintf("A%d", i+2), i+1)
				f.SetCellValue(sheet1Name, fmt.Sprintf("B%d", i+2), each["tanggal"])
				f.SetCellValue(sheet1Name, fmt.Sprintf("D%d", i+2), each["keterangan"])
				f.SetCellValue(sheet1Name, fmt.Sprintf("E%d", i+2), debit)
				f.SetCellValue(sheet1Name, fmt.Sprintf("F%d", i+2), kredit)
				f.SetCellValue(sheet1Name, fmt.Sprintf("F%d", i+2), saldo)
			}

			// f.SetActiveSheet(cell)
			// Save spreadsheet by the given path.
			if err := f.SaveAs("../download/excel/Book1.xlsx"); err != nil {
				fmt.Println(err)
			}

			// Buat File
			// f := excelize.NewFile()
			// // Create a new sheet.
			// index := f.NewSheet("Sheet1")
			// // Set value of a cell.
			// f.SetCellValue("Sheet1", "A2", "No")
			// f.SetCellValue("Sheet1", "B2", "Tanggal")
			// f.SetCellValue("Sheet1", "C2", "Keterangan")
			// f.SetCellValue("Sheet1", "D2", "Debit")
			// f.SetCellValue("Sheet1", "E2", "Kredit")
			// f.SetCellValue("Sheet1", "F2", "Saldo")
			// // Set active sheet of the workbook.
			// f.SetActiveSheet(index)
			// // Save spreadsheet by the given path.
			// if err := f.SaveAs("Book1.xlsx"); err != nil {
			//     fmt.Println(err)
			// }

			callback["success"] = true
			callback["msg"] = "Sukses Download"
			DB, _ := db_go.DB()
			DB.Close()

			context.JSON(200, callback)
		})
	}
}

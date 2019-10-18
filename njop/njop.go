package main

import (
	"database/sql"
	"github.com/360EntSecGroup-Skylar/excelize"
	_ "gopkg.in/goracle.v2"
)
import ."fmt"

var Nkec, Nkel, Nblok, Nurut, Lbumi, Njopbumi string
func main() {
	getExcelData()
	//getNJOP(Nkec, Nkel, Nblok, Nurut)
}

func getExcelData(){
	f, err := excelize.OpenFile("documents/Rekap jalan.xlsx")
	if err != nil {
		Println(err)
		return
	}

	//get value from cell given worksheet name and axis
	// cell, err := f.GetCellValue("Sheet1", "A2")
	// if err != nil{
	//     Println(err)
	//     return
	// }
	// Println(cell)

	//get all the rows in the sheet1
	//rows, err := f.GetRows("Sheet1")
	//for _, row := range rows {
	//	for _, colCell := range row {
	//		Print(colCell, "\t")
	//	}
	//	Println()
	//}

	// get one column
	// rows, err := f.Rows("Sheet1")
	// if err != nil {
	//     Println(err)
	// }
	// for rows.Next() {
	//     row, err := rows.Columns()
	//     if err != nil {
	//         Println(err)
	//     }
	//     Printf("%s\t\n", row[0]) // Print values column A
	// }

	n := 40
	for i := 2; i < n; i++ {
		a, err := f.GetCellValue("Sheet1", Sprintf("A%d", i))
		if err != nil {
			Println(err)
		}
		//b, err := f.GetCellValue("Sheet1", Sprintf("B%d", i))
		// c, err := f.GetCellValue("Sheet1", "A2")
		if err != nil {
			Println(err)
		}
		//Printf("%s\t%s\n", a, b) //print values colum A
		Nkec = string(a[6:9])
		Nkel = string(a[10:13])
		Nblok = string(a[14:17])
		Nurut = string(a[18:22])
		//Printf("Kecamatan: %s, Kelurahan: %s, Blok: %s, Urut: %s\n\n", Nkec, Nkel, Nblok, Nurut)
		getNJOP(Nkec, Nkel, Nblok, Nurut)
		Printf("Kecamatan: %s, Kelurahan: %s, Blok: %s, Urut: %s, Bumi: %s, NJOP: %s\n\n", Nkec, Nkel, Nblok,Nurut,Lbumi,Njopbumi)
		f.SetCellValue("Sheet1",Sprintf("B%d", i), Sprintf("%s", Njopbumi))
		err = f.Save()
		if err != nil {
			println(err)
		}
	}
}

func getNJOP(vkdKec, vkdKel, vkdBlok, vkdUrut string) {
	db, err := konek()
	if err != nil {
		Println(err)
		return
	}
	defer db.Close()
	rows, err := db.Query("select * from (select s.luas_bumi_sppt, s.njop_bumi_sppt from SPPT s where " +
		"s.KD_KECAMATAN = '" + vkdKec + "' and s.KD_KELURAHAN = '" + vkdKel + "' and s.KD_BLOK = '" + vkdBlok + "' and " +
		"s.NO_URUT = '" + vkdUrut + "' and s.THN_PAJAK_SPPT < 2019 order by s.THN_PAJAK_SPPT desc) where ROWNUM =1")
	if err != nil {
		Println("Error running query")
		Printf("%v\n", rows)
		Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&Lbumi, &Njopbumi)
		if err != nil {
			Println("Error when scanning")
			Println("NOP not found")
			Println(err)
		}
	}
}

func konek() (*sql.DB, error) {
	connString := "pbb/PBB@103.76.175.175:26/ORCL"
	db, err := sql.Open("goracle", connString)
	if err != nil {
		return nil, err
	}
	return db, nil
}

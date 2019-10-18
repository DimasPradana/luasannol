package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	_ "gopkg.in/goracle.v2"
)

//type getsppt struct {
//	xkdProp, xkdDati, xkdKec, xkdKel, xkdBlok, xkdUrut, xkdJns, xThn string
//	xLuasBumi, xLuasBng int16
//}
//
//type getdatobyekpajak struct {
//	tBumi, tBng int16
//	xKdZnt string
//}

var (
	Kel      string
	Blok     string
	Urut     string
	JnsOp    string
	LBumi    string
	LBng     string
	KdZnt    string
	Form     string
	JnsTanah string
	NoBng    string
	Kec      string
)

const DayaListrik string = "750"

func main() {
	/*
		TODO
		- tabel yang dipakai SPPT, DAT_OP_BANGUNAN, DAT_FASILITAS_BANGUNAN, DAT_OP_BUMI, DAT_OBJEK_PAJAK
		- jika luas bangunan di SPPT 0 maka no bangunan =0 dan jenis tanah = 3
	*/
	kecPtr := flag.String("kec", "000", "Kode Kecamatan")
	kelPtr := flag.String("kel", "000", "Kode Kelurahan")
	blokPtr := flag.String("blok", "000", "Kode Blok")
	urutPtr := flag.String("urut", "000", "Nomor Urut")
	formPtr := flag.String("form", "1000", "Nomor Formulir")
	flag.Parse()
	getSPPT(*kecPtr, *kelPtr, *blokPtr, *urutPtr)
	getDatObyekPajak(*kecPtr, *kelPtr, *blokPtr, *urutPtr)
	Form = "2019900" + *formPtr
	fmt.Printf("NOP main\t: 35-12-%s-%s-%s-%s-%s\n", Kec, Kel, Blok, Urut, JnsOp)
	fmt.Printf("Formulir main\t: %s\n", Form)
	fmt.Printf("Bumi main\t: %s\n", LBumi)
	fmt.Printf("Bangunan main\t: %s\n", LBng)
	fmt.Printf("Zona Tanah\t: %s\n", KdZnt)
	fmt.Printf("Jenis Tanah\t: %s\n", JnsTanah)
	fmt.Printf("Nomor Bangunan\t: %s\n", NoBng)
	fmt.Printf("Daya Listrik\t: %s\n\n", DayaListrik)
	if JnsTanah == "3" {
		insertDATOPBumi(Kec, Kel, Blok, Urut, KdZnt, LBumi, JnsTanah)
		updateDATObjekPajak(Form, LBumi, LBng, Kec, Kel, Blok, Urut)
	} else if JnsTanah == "1" {
		insertDATOPBangunan(Kec, Kel, Blok, Urut, NoBng, Form, LBng)
		insertDATFasilitasBangunan(Kec, Kel, Blok, Urut, NoBng, DayaListrik)
		insertDATOPBumi(Kec, Kel, Blok, Urut, KdZnt, LBumi, JnsTanah)
		updateDATObjekPajak(Form, LBumi, LBng, Kec, Kel, Blok, Urut)
	}
}

func getSPPT(vkdKec, vkdKel, vkdBlok, vkdUrut string) {
	db, err := connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	var pr, dati, kec, kel, blok, urut, jns, thn, lbumi, lbng string
	rows, err := db.Query("select * from (select s.kd_propinsi, s.kd_dati2, s.kd_kecamatan, s.kd_kelurahan, s.kd_blok, s.no_urut, " +
		"s.kd_jns_op, s.thn_pajak_sppt, s.luas_bumi_sppt, s.luas_bng_sppt from SPPT s where " +
		"s.KD_KECAMATAN = '" + vkdKec + "' and s.KD_KELURAHAN = '" + vkdKel + "' and s.KD_BLOK = '" + vkdBlok + "' and " +
		"s.NO_URUT = '" + vkdUrut + "' and s.THN_PAJAK_SPPT < 2019 order by s.THN_PAJAK_SPPT desc) where ROWNUM =1")
	if err != nil {
		fmt.Println("Error running query")
		fmt.Printf("%v\n", rows)
		fmt.Println(err)
		log.Fatal(err)
		//return
	}
	defer rows.Close()
	for rows.Next() {
		//rows.Scan(&pr, &dati, &kec, &kel, &blok, &urut, &jns, &thn, &tagihan, &lunas)
		err := rows.Scan(&pr, &dati, &kec, &kel, &blok, &urut, &jns, &thn, &lbumi, &lbng)
		if err != nil {
			fmt.Println("Error when scanning")
			fmt.Println("NOP not found")
			fmt.Println(err)
			log.Fatal(err)
			//return
		}
		fmt.Println("\ntable SPPT")
		fmt.Printf("NOP\t\t: %s-%s-%s-%s-%s-%s-%s\n", pr, dati, kec, kel, blok, urut, jns)
		fmt.Printf("Tahun\t\t: %s\n", thn)
		fmt.Printf("Luas Bumi\t: %s\n", lbumi)
		fmt.Printf("Luas Bangunan\t: %s\n", lbng)
		LBumi = lbumi
		LBng = lbng
		if lbng == "0" {
			JnsTanah = "3"
			NoBng = "0"
		} else {
			JnsTanah = "1"
			NoBng = "1"
		}
	}
}

func getDatObyekPajak(vkdKec, vkdKel, vkdBlok, vkdUrut string) {
	db, err := connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	rows, err := db.Query("select s.kd_propinsi, s.kd_dati2, s.kd_kecamatan, s.kd_kelurahan, s.kd_blok, s.no_urut," +
		" s.kd_jns_op, s.total_luas_bumi, s.total_luas_bng, d.kd_znt, e.jns_bumi from DAT_OBJEK_PAJAK s left join DAT_PETA_ZNT d " +
		"on s.KD_PROPINSI = d.KD_PROPINSI and s.KD_DATI2 = d.KD_DATI2 and s.KD_KECAMATAN = d.KD_KECAMATAN and " +
		"s.KD_KELURAHAN = d.KD_KELURAHAN and s.KD_BLOK = d.KD_BLOK left join DAT_OP_BUMI e " +
		"on s.KD_PROPINSI = e.KD_PROPINSI and s.KD_DATI2 = e.KD_DATI2 and s.KD_KECAMATAN = e.KD_KECAMATAN and " +
		"s.KD_KELURAHAN = e.KD_KELURAHAN and s.KD_BLOK = e.KD_BLOK where s.KD_KECAMATAN='" + vkdKec + "' and " +
		"s.KD_KELURAHAN='" + vkdKel + "' and s.KD_BLOK='" + vkdBlok + "' and s.NO_URUT='" + vkdUrut + "' and ROWNUM =1")
	if err != nil {
		fmt.Println("Error running query")
		fmt.Printf("%v\n", rows)
		fmt.Println(err)
		return
	}
	defer rows.Close()
	var pr, dati, kec, kel, blok, urut, jns, znt, jnsbumi string
	var tbumi, tbng int16
	for rows.Next() {
		//rows.Scan(&pr, &dati, &kec, &kel, &blok, &urut, &jns, &thn, &tagihan, &lunas)
		err := rows.Scan(&pr, &dati, &kec, &kel, &blok, &urut, &jns, &tbumi, &tbng, &znt, &jnsbumi)
		if err != nil {
			fmt.Println("Error when scanning")
			fmt.Println("NOP not found")
			fmt.Println(err)
			return
		}
		fmt.Println("============================")
		fmt.Println("table DAT_OBYEK_PAJAK")
		fmt.Printf("NOP\t\t: %s-%s-%s-%s-%s-%s-%s\n", pr, dati, kec, kel, blok, urut, jns)
		fmt.Printf("Luas Bumi\t: %d\n", tbumi)
		fmt.Printf("Luas Bangunan\t: %d\n", tbng)
		fmt.Printf("ZNT\t\t: %s\n", znt)
		fmt.Printf("Jenis Bumi\t: %s\n\n", jnsbumi)
		Kec = kec
		Kel = kel
		Blok = blok
		Urut = urut
		JnsOp = jns
		KdZnt = znt
	}
}

func connect() (*sql.DB, error) {
	connString := "pbb/PBB@103.76.175.175:26/ORCL"
	db, err := sql.Open("goracle", connString)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func insertDATOPBangunan(vkec, vkel, vblok, vurut, vnbng, vform, vlbng string) {
	/*
		INSERT INTO PBB.DAT_OP_BANGUNAN (KD_PROPINSI, KD_DATI2, KD_KECAMATAN, KD_KELURAHAN, KD_BLOK, NO_URUT, KD_JNS_OP, NO_BNG,
		                                 KD_JPB, NO_FORMULIR_LSPOP, THN_DIBANGUN_BNG, THN_RENOVASI_BNG, LUAS_BNG,
		                                 JML_LANTAI_BNG, KONDISI_BNG, JNS_KONSTRUKSI_BNG, JNS_ATAP_BNG, KD_DINDING,
		                                 KD_LANTAI, KD_LANGIT_LANGIT, NILAI_SISTEM_BNG, JNS_TRANSAKSI_BNG,
		                                 TGL_PENDATAAN_BNG, NIP_PENDATA_BNG, TGL_PEMERIKSAAN_BNG, NIP_PEMERIKSA_BNG,
		                                 TGL_PEREKAMAN_BNG, NIP_PEREKAM_BNG)
		                VALUES ('35', '12', '$(kec)', '$(kel)', '$(blok)', '$(urut)', '9', $(no_bangunan),
		                        '01', '$(no_formulir)', '2000', null, $(luas_bangunan),
		                        1, '2', '3', '3', '3', '2', '2', 0, '1',
		                TO_DATE(CURRENT_DATE, 'DD-MM-YYYY HH24:MI:SS'), '196408142', TO_DATE(CURRENT_DATE, 'DD-MM-YYYY HH24:MI:SS'), '197505192',
		                                        TO_DATE(CURRENT_DATE, 'DD-MM-YYYY HH24:MI:SS'), '060000000');
	*/
	db, err := connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	rows, err := db.Exec("INSERT INTO PBB.DAT_OP_BANGUNAN (KD_PROPINSI, KD_DATI2, KD_KECAMATAN, KD_KELURAHAN, " +
		"KD_BLOK, NO_URUT, KD_JNS_OP, NO_BNG, KD_JPB, NO_FORMULIR_LSPOP, THN_DIBANGUN_BNG, " +
		"THN_RENOVASI_BNG, LUAS_BNG, JML_LANTAI_BNG, KONDISI_BNG, JNS_KONSTRUKSI_BNG, JNS_ATAP_BNG, KD_DINDING," +
		"KD_LANTAI, KD_LANGIT_LANGIT, NILAI_SISTEM_BNG, JNS_TRANSAKSI_BNG, TGL_PENDATAAN_BNG, NIP_PENDATA_BNG," +
		" TGL_PEMERIKSAAN_BNG, NIP_PEMERIKSA_BNG, TGL_PEREKAMAN_BNG, NIP_PEREKAM_BNG) VALUES ('35', '12', '" + vkec + "', " +
		"'" + vkel + "', '" + vblok + "', '" + vurut + "', '9', " + vnbng + ", '01', '" + vform + "', '2000', null, " + vlbng + ", " +
		"1, '2', '3', '3', '3', '2', '2', 0, '1', SYSDATE, " +
		"'196408142', SYSDATE, '197505192', SYSDATE, '060000000')")
	if err != nil {
		fmt.Println("Error running query")
		fmt.Printf("%v\n", rows)
		fmt.Println(err)
		return
	}
	fmt.Println(&rows)
	fmt.Println("masuk insertDatOpBangunan")
}
func insertDATFasilitasBangunan(vkec, vkel, vblok, vurut, vnbng, vdayalistrik string) {
	/*
		INSERT INTO PBB.DAT_FASILITAS_BANGUNAN (KD_PROPINSI, KD_DATI2, KD_KECAMATAN, KD_KELURAHAN, KD_BLOK, NO_URUT, KD_JNS_OP,
		                                        NO_BNG, KD_FASILITAS, JML_SATUAN)
		                        VALUES ('35', '12', '$(kec)', '$(kel)', '$(blok)', '$(urut)', '9',
		                                $(no_bangunan), '44', $(daya_listrik));
	*/
	db, err := connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	rows, err := db.Exec("INSERT INTO PBB.DAT_FASILITAS_BANGUNAN (KD_PROPINSI, KD_DATI2, KD_KECAMATAN," +
		"KD_KELURAHAN, KD_BLOK, NO_URUT, KD_JNS_OP, NO_BNG, KD_FASILITAS, JML_SATUAN) VALUES ('35', '12', " +
		"'" + vkec + "', '" + vkel + "', '" + vblok + "', '" + vurut + "', '9', " + vnbng + ", '44', " + vdayalistrik + ")")
	if err != nil {
		fmt.Println("Error running query")
		fmt.Printf("%v\n", rows)
		fmt.Println(err)
		return
	}
	fmt.Println(&rows)
	fmt.Println("masuk insertDatFasBangunan")
}
func insertDATOPBumi(vkec, vkel, vblok, vurut, vkdznt, vlbumi, vjnstnh string) {
	/*
		INSERT INTO PBB.DAT_OP_BUMI (KD_PROPINSI, KD_DATI2, KD_KECAMATAN, KD_KELURAHAN, KD_BLOK, NO_URUT, KD_JNS_OP, NO_BUMI,KD_ZNT, LUAS_BUMI, JNS_BUMI, NILAI_SISTEM_BUMI)
		                VALUES ('35', '12', '$(kec)', '$(kel)', '$(blok)', '$(urut)', '9', 1,'$(kd_znt)', $(luas_bumi), $(jenis_tanah), 0);

	*/
	db, err := connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	rows, err := db.Exec("INSERT INTO PBB.DAT_OP_BUMI (KD_PROPINSI, KD_DATI2, KD_KECAMATAN, KD_KELURAHAN," +
		" KD_BLOK, NO_URUT, KD_JNS_OP, NO_BUMI,KD_ZNT, LUAS_BUMI, JNS_BUMI, NILAI_SISTEM_BUMI) VALUES ('35', '12', " +
		"'" + vkec + "', '" + vkel + "', '" + vblok + "', '" + vurut + "', '9', 1,'" + vkdznt + "', " + vlbumi + "," + vjnstnh + " , 0)")
	if err != nil {
		fmt.Println("Error running query")
		fmt.Printf("%v\n", rows)
		fmt.Println(err)
		return
	}
	fmt.Println(&rows)
	fmt.Println("masuk insertDatOpBumi")
}
func updateDATObjekPajak(vform, vlbumi, vlbng, vkec, vkel, vblok, vurut string) {
	/*
		UPDATE PBB.DAT_OBJEK_PAJAK
		            SET NO_FORMULIR_SPOP = '$(no_formulir)',
		            NO_PERSIL = null, BLOK_KAV_NO_OP = null, RW_OP = '00', RT_OP = '000',
		            KD_STATUS_CABANG = '0', KD_STATUS_WP = '1',
		            TOTAL_LUAS_BUMI = $(luas_bumi), TOTAL_LUAS_BNG = $(luas_bangunan), NJOP_BUMI = 0, NJOP_BNG = 0, STATUS_PETA_OP = '1',
		            JNS_TRANSAKSI_OP = '2', TGL_PENDATAAN_OP = TO_DATE(CURRENT_DATE, 'DD-MM-YYYY HH24:MI:SS'),
		            NIP_PENDATA = '196408142', TGL_PEMERIKSAAN_OP = TO_DATE(CURRENT_DATE, 'DD-MM-YYYY HH24:MI:SS'),
		            NIP_PEMERIKSA_OP = '197505192', TGL_PEREKAMAN_OP = TO_DATE(CURRENT_DATE, 'DD-MM-YYYY HH24:MI:SS'),
		            NIP_PEREKAM_OP = '060000000'
		        WHERE KD_PROPINSI = '35' AND KD_DATI2 = '12' AND KD_KECAMATAN = '$(kec)' AND KD_KELURAHAN = '$(kel)' AND KD_BLOK = '$(blok)' AND NO_URUT = '$(urut)' AND KD_JNS_OP = '9';

	*/
	db, err := connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	rows, err := db.Exec("UPDATE PBB.DAT_OBJEK_PAJAK SET NO_FORMULIR_SPOP = '" + vform + "', " +
		"NO_PERSIL = null, BLOK_KAV_NO_OP = null, RW_OP = '00', RT_OP = '000', KD_STATUS_CABANG = '0', " +
		"KD_STATUS_WP = '1', TOTAL_LUAS_BUMI = " + vlbumi + ", TOTAL_LUAS_BNG = " + vlbng + ", " +
		"NJOP_BUMI = 0, NJOP_BNG = 0, STATUS_PETA_OP = '1', JNS_TRANSAKSI_OP = '2', " +
		//"TGL_PENDATAAN_OP = TO_DATE(CURRENT_DATE, 'DD-MM-YYYY HH24:MI:SS'),	NIP_PENDATA = '196408142', " +
		"TGL_PENDATAAN_OP = SYSDATE, NIP_PENDATA = '196408142', " +
		"TGL_PEMERIKSAAN_OP = SYSDATE, NIP_PEMERIKSA_OP = '197505192', " +
		"TGL_PEREKAMAN_OP = SYSDATE, NIP_PEREKAM_OP = '060000000' " +
		"WHERE KD_PROPINSI = '35' AND KD_DATI2 = '12' AND KD_KECAMATAN = '" + vkec + "' AND KD_KELURAHAN = '" + vkel + "' " +
		"AND KD_BLOK = '" + vblok + "' AND NO_URUT = '" + vurut + "' AND KD_JNS_OP = '9'")
	if err != nil {
		fmt.Println("Error running query")
		fmt.Printf("%v\n", rows)
		fmt.Println(err)
		return
	}
	fmt.Println(&rows)
	fmt.Println("masuk updateDatObjekPajak")
}

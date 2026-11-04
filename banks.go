package main

import "strings"

// Bank catalogs the Vietnamese banks supported by the VietQR NAPAS 247 network.
// Code is the identifier expected by github.com/sunary/vietqr.
type Bank struct {
	Code       string `json:"code"`
	Bin        string `json:"bin"`
	VietQRCode string `json:"vietqrCode,omitempty"`
	ShortName  string `json:"shortName"`
	Name       string `json:"name"`
	Logo       string `json:"logo,omitempty"`
}

var banks = []Bank{
	{
		Code:       "VIETINBANK",
		Bin:        "970415",
		VietQRCode: "ICB",
		ShortName:  "VietinBank",
		Name:       "Ngân hàng TMCP Công thương Việt Nam",
		Logo:       "/assets/logos/VIETINBANK.png",
	},
	{
		Code:       "VIETCOMBANK",
		Bin:        "970436",
		VietQRCode: "VCB",
		ShortName:  "Vietcombank",
		Name:       "Ngân hàng TMCP Ngoại Thương Việt Nam",
		Logo:       "/assets/logos/VIETCOMBANK.png",
	},
	{
		Code:       "BIDV",
		Bin:        "970418",
		VietQRCode: "BIDV",
		ShortName:  "BIDV",
		Name:       "Ngân hàng TMCP Đầu tư và Phát triển Việt Nam",
		Logo:       "/assets/logos/BIDV.png",
	},
	{
		Code:       "AGRIBANK",
		Bin:        "970405",
		VietQRCode: "VBA",
		ShortName:  "Agribank",
		Name:       "Ngân hàng Nông nghiệp và Phát triển Nông thôn Việt Nam",
		Logo:       "/assets/logos/AGRIBANK.png",
	},
	{
		Code:       "OCB",
		Bin:        "970448",
		VietQRCode: "OCB",
		ShortName:  "OCB",
		Name:       "Ngân hàng TMCP Phương Đông",
		Logo:       "/assets/logos/OCB.png",
	},
	{
		Code:       "MBBANK",
		Bin:        "970422",
		VietQRCode: "MB",
		ShortName:  "MBBank",
		Name:       "Ngân hàng TMCP Quân đội",
		Logo:       "/assets/logos/MBBANK.png",
	},
	{
		Code:       "TECHCOMBANK",
		Bin:        "970407",
		VietQRCode: "TCB",
		ShortName:  "Techcombank",
		Name:       "Ngân hàng TMCP Kỹ thương Việt Nam",
		Logo:       "/assets/logos/TECHCOMBANK.png",
	},
	{
		Code:       "ACB",
		Bin:        "970416",
		VietQRCode: "ACB",
		ShortName:  "ACB",
		Name:       "Ngân hàng TMCP Á Châu",
		Logo:       "/assets/logos/ACB.png",
	},
	{
		Code:       "VPBANK",
		Bin:        "970432",
		VietQRCode: "VPB",
		ShortName:  "VPBank",
		Name:       "Ngân hàng TMCP Việt Nam Thịnh Vượng",
		Logo:       "/assets/logos/VPBANK.png",
	},
	{
		Code:       "TPBANK",
		Bin:        "970423",
		VietQRCode: "TPB",
		ShortName:  "TPBank",
		Name:       "Ngân hàng TMCP Tiên Phong",
		Logo:       "/assets/logos/TPBANK.png",
	},
	{
		Code:       "SACOMBANK",
		Bin:        "970403",
		VietQRCode: "STB",
		ShortName:  "Sacombank",
		Name:       "Ngân hàng TMCP Sài Gòn Tài Lộc",
		Logo:       "/assets/logos/SACOMBANK.png",
	},
	{
		Code:       "HDBANK",
		Bin:        "970437",
		VietQRCode: "HDB",
		ShortName:  "HDBank",
		Name:       "Ngân hàng TMCP Phát triển Thành phố Hồ Chí Minh",
		Logo:       "/assets/logos/HDBANK.png",
	},
	{
		Code:       "BANVIET",
		Bin:        "970454",
		VietQRCode: "VCCB",
		ShortName:  "VietCapitalBank",
		Name:       "Ngân hàng TMCP Bản Việt",
		Logo:       "/assets/logos/BANVIET.png",
	},
	{
		Code:       "SCB",
		Bin:        "970429",
		VietQRCode: "SCB",
		ShortName:  "SCB",
		Name:       "Ngân hàng TMCP Sài Gòn",
		Logo:       "/assets/logos/SCB.png",
	},
	{
		Code:       "VIB",
		Bin:        "970441",
		VietQRCode: "VIB",
		ShortName:  "VIB",
		Name:       "Ngân hàng TMCP Quốc tế Việt Nam",
		Logo:       "/assets/logos/VIB.png",
	},
	{
		Code:       "SHB",
		Bin:        "970443",
		VietQRCode: "SHB",
		ShortName:  "SHB",
		Name:       "Ngân hàng TMCP Sài Gòn - Hà Nội",
		Logo:       "/assets/logos/SHB.png",
	},
	{
		Code:       "EXIMBANK",
		Bin:        "970431",
		VietQRCode: "EIB",
		ShortName:  "Eximbank",
		Name:       "Ngân hàng TMCP Xuất Nhập khẩu Việt Nam",
		Logo:       "/assets/logos/EXIMBANK.png",
	},
	{
		Code:       "MSB",
		Bin:        "970426",
		VietQRCode: "MSB",
		ShortName:  "MSB",
		Name:       "Ngân hàng TMCP Hàng Hải",
		Logo:       "/assets/logos/MSB.png",
	},
	{
		Code:       "CAKE",
		Bin:        "546034",
		VietQRCode: "CAKE",
		ShortName:  "CAKE",
		Name:       "TMCP Việt Nam Thịnh Vượng - Ngân hàng số CAKE by VPBank",
		Logo:       "/assets/logos/CAKE.png",
	},
	{
		Code:       "UBANK",
		Bin:        "546035",
		VietQRCode: "Ubank",
		ShortName:  "Ubank",
		Name:       "TMCP Việt Nam Thịnh Vượng - Ngân hàng số Ubank by VPBank",
		Logo:       "/assets/logos/UBANK.png",
	},
	{
		Code:       "TIMO",
		Bin:        "963388",
		VietQRCode: "TIMO",
		ShortName:  "Timo",
		Name:       "Ngân hàng số Timo by Ban Viet Bank (Timo by Ban Viet Bank)",
		Logo:       "/assets/logos/TIMO.png",
	},
	{
		Code:       "SAIGONBANK",
		Bin:        "970400",
		VietQRCode: "SGICB",
		ShortName:  "SaigonBank",
		Name:       "Ngân hàng TMCP Sài Gòn Công Thương",
		Logo:       "/assets/logos/SAIGONBANK.png",
	},
	{
		Code:       "BAC_A_BANK",
		Bin:        "970409",
		VietQRCode: "BAB",
		ShortName:  "BacABank",
		Name:       "Ngân hàng TMCP Bắc Á",
		Logo:       "/assets/logos/BAC_A_BANK.png",
	},
	{
		Code:       "PVCOM_BANK",
		Bin:        "970412",
		VietQRCode: "PVCB",
		ShortName:  "PVcomBank",
		Name:       "Ngân hàng TMCP Đại Chúng Việt Nam",
		Logo:       "/assets/logos/PVCOM_BANK.png",
	},
	{
		Code:       "OCEANBANK",
		Bin:        "970414",
		VietQRCode: "MBV",
		ShortName:  "MBV",
		Name:       "Ngân hàng TNHH MTV Việt Nam Hiện Đại",
		Logo:       "/assets/logos/OCEANBANK.png",
	},
	{
		Code:       "NCB",
		Bin:        "970419",
		VietQRCode: "NCB",
		ShortName:  "NCB",
		Name:       "Ngân hàng TMCP Quốc Dân",
		Logo:       "/assets/logos/NCB.png",
	},
	{
		Code:       "SHINHAN_BANK",
		Bin:        "970424",
		VietQRCode: "SHBVN",
		ShortName:  "ShinhanBank",
		Name:       "Ngân hàng TNHH MTV Shinhan Việt Nam",
		Logo:       "/assets/logos/SHINHAN_BANK.png",
	},
	{
		Code:       "ABBANK",
		Bin:        "970425",
		VietQRCode: "ABB",
		ShortName:  "ABBANK",
		Name:       "Ngân hàng TMCP An Bình",
		Logo:       "/assets/logos/ABBANK.png",
	},
	{
		Code:       "VIET_A_BANK",
		Bin:        "970427",
		VietQRCode: "VAB",
		ShortName:  "VietABank",
		Name:       "Ngân hàng TMCP Việt Á",
		Logo:       "/assets/logos/VIET_A_BANK.png",
	},
	{
		Code:       "NAM_A_BANK",
		Bin:        "970428",
		VietQRCode: "NAB",
		ShortName:  "NamABank",
		Name:       "Ngân hàng TMCP Nam Á",
		Logo:       "/assets/logos/NAM_A_BANK.png",
	},
	{
		Code:       "PGBANK",
		Bin:        "970430",
		VietQRCode: "PGB",
		ShortName:  "PGBank",
		Name:       "Ngân hàng TMCP Thịnh vượng và Phát triển",
		Logo:       "/assets/logos/PGBANK.png",
	},
	{
		Code:       "VIET_BANK",
		Bin:        "970433",
		VietQRCode: "VIETBANK",
		ShortName:  "VietBank",
		Name:       "Ngân hàng TMCP Việt Nam Thương Tín",
		Logo:       "/assets/logos/VIET_BANK.png",
	},
	{
		Code:       "BAOVIET_BANK",
		Bin:        "970438",
		VietQRCode: "BVB",
		ShortName:  "BaoVietBank",
		Name:       "Ngân hàng TMCP Bảo Việt",
		Logo:       "/assets/logos/BAOVIET_BANK.png",
	},
	{
		Code:       "SEA_BANK",
		Bin:        "970440",
		VietQRCode: "SEAB",
		ShortName:  "SeABank",
		Name:       "Ngân hàng TMCP Đông Nam Á",
		Logo:       "/assets/logos/SEA_BANK.png",
	},
	{
		Code:       "COOP_BANK",
		Bin:        "970446",
		VietQRCode: "COOPBANK",
		ShortName:  "COOPBANK",
		Name:       "Ngân hàng Hợp tác xã Việt Nam",
		Logo:       "/assets/logos/COOP_BANK.png",
	},
	{
		Code:       "LIENVIETPOST_BANK",
		Bin:        "970449",
		VietQRCode: "LPB",
		ShortName:  "LPBank",
		Name:       "Ngân hàng TMCP Lộc Phát Việt Nam",
		Logo:       "/assets/logos/LIENVIETPOST_BANK.png",
	},
	{
		Code:       "KIENLONG_BANK",
		Bin:        "970452",
		VietQRCode: "KLB",
		ShortName:  "KienLongBank",
		Name:       "Ngân hàng TMCP Kiên Long",
		Logo:       "/assets/logos/KIENLONG_BANK.png",
	},
	{
		Code:       "KASIKORN_BANK",
		Bin:        "668888",
		VietQRCode: "KBank",
		ShortName:  "KBank",
		Name:       "Ngân hàng Đại chúng TNHH Kasikornbank",
		Logo:       "/assets/logos/KASIKORN_BANK.png",
	},
	{
		Code:       "KEB_HANA_BANK_HCM",
		Bin:        "970466",
		VietQRCode: "KEBHANAHCM",
		ShortName:  "KEBHanaHCM",
		Name:       "Ngân hàng KEB Hana – Chi nhánh Thành phố Hồ Chí Minh",
		Logo:       "/assets/logos/KEB_HANA_BANK_HCM.png",
	},
	{
		Code:       "KEB_HANA_BANK_HN",
		Bin:        "970467",
		VietQRCode: "KEBHANAHN",
		ShortName:  "KEBHANAHN",
		Name:       "Ngân hàng KEB Hana – Chi nhánh Hà Nội",
		Logo:       "/assets/logos/KEB_HANA_BANK_HN.png",
	},
	{
		Code:       "SINOPAC_BANK_HCM",
		Bin:        "970465",
		VietQRCode: "BSP",
		ShortName:  "BSP",
		Name:       "Ngân hàng SinoPac - Chi nhánh Thành phố Hồ Chí Minh",
		Logo:       "/assets/logos/SINOPAC_BANK_HCM.png",
	},
	{
		Code:       "HSBC",
		Bin:        "458761",
		VietQRCode: "HSBC",
		ShortName:  "HSBC",
		Name:       "Ngân hàng TNHH MTV HSBC (Việt Nam)",
		Logo:       "/assets/logos/HSBC.png",
	},
	{
		Code:       "HONGLEONG_BANK",
		Bin:        "970442",
		VietQRCode: "HLBVN",
		ShortName:  "HongLeong",
		Name:       "Ngân hàng TNHH MTV Hong Leong Việt Nam",
		Logo:       "/assets/logos/HONGLEONG_BANK.png",
	},
	{
		Code:       "IBK_HN",
		Bin:        "970455",
		VietQRCode: "IBK - HN",
		ShortName:  "IBKHN",
		Name:       "Ngân hàng Công nghiệp Hàn Quốc - Chi nhánh Hà Nội",
		Logo:       "",
	},
	{
		Code:       "IBK_HCM",
		Bin:        "970456",
		VietQRCode: "IBK - HCM",
		ShortName:  "IBKHCM",
		Name:       "Ngân hàng Công nghiệp Hàn Quốc - Chi nhánh TP. Hồ Chí Minh",
		Logo:       "",
	},
	{
		Code:       "INDOVINA_BANK",
		Bin:        "970434",
		VietQRCode: "IVB",
		ShortName:  "IndovinaBank",
		Name:       "Ngân hàng TNHH Indovina",
		Logo:       "/assets/logos/INDOVINA_BANK.png",
	},
	{
		Code:       "DONG_A_BANK",
		Bin:        "970406",
		VietQRCode: "Vikki",
		ShortName:  "Vikki",
		Name:       "Ngân hàng TNHH Một thành viên số Vikki",
		Logo:       "/assets/logos/DONG_A_BANK.png",
	},
	{
		Code:       "NONGHYUP_BANK_HN",
		Bin:        "801011",
		VietQRCode: "NHB HN",
		ShortName:  "Nonghyup",
		Name:       "Ngân hàng Nonghyup - Chi nhánh Hà Nội",
		Logo:       "/assets/logos/NONGHYUP_BANK_HN.png",
	},
	{
		Code:       "DBS_BANK",
		Bin:        "796500",
		VietQRCode: "DBS",
		ShortName:  "DBSBank",
		Name:       "DBS Bank Ltd - Chi nhánh Thành phố Hồ Chí Minh",
		Logo:       "/assets/logos/DBS_BANK.png",
	},
	{
		Code:       "PUBLIC_BANK",
		Bin:        "970439",
		VietQRCode: "PBVN",
		ShortName:  "PublicBank",
		Name:       "Ngân hàng TNHH MTV Public Việt Nam",
		Logo:       "/assets/logos/PUBLIC_BANK.png",
	},
	{
		Code:       "STANDARD_CHARTERED_BANK",
		Bin:        "970410",
		VietQRCode: "SCVN",
		ShortName:  "StandardChartered",
		Name:       "Ngân hàng TNHH MTV Standard Chartered Bank Việt Nam",
		Logo:       "/assets/logos/STANDARD_CHARTERED_BANK.png",
	},
	{
		Code:       "UNITED_OVERSEAS_BANK",
		Bin:        "970458",
		VietQRCode: "UOB",
		ShortName:  "UnitedOverseas",
		Name:       "Ngân hàng United Overseas - Chi nhánh TP. Hồ Chí Minh",
		Logo:       "/assets/logos/UNITED_OVERSEAS_BANK.png",
	},
	{
		Code:       "CBBANK",
		Bin:        "970444",
		VietQRCode: "CBB",
		ShortName:  "CBBank",
		Name:       "Ngân hàng Thương mại trách nhiệm hữu hạn một thành viên Ngoại thương Công nghệ số (Vietcombank Neo Limited)",
		Logo:       "/assets/logos/CBBANK.png",
	},
	{
		Code:       "GPBANK",
		Bin:        "970408",
		VietQRCode: "GPB",
		ShortName:  "GPBank",
		Name:       "Ngân hàng Thương mại TNHH MTV Dầu Khí Toàn Cầu",
		Logo:       "/assets/logos/GPBANK.png",
	},
	{
		Code:       "VRB",
		Bin:        "970421",
		VietQRCode: "VRB",
		ShortName:  "VRB",
		Name:       "Ngân hàng Liên doanh Việt - Nga",
		Logo:       "/assets/logos/VRB.png",
	},
	{
		Code:       "WOORI_BANK",
		Bin:        "970457",
		VietQRCode: "WVN",
		ShortName:  "Woori",
		Name:       "Ngân hàng TNHH MTV Woori Việt Nam",
		Logo:       "/assets/logos/WOORI_BANK.png",
	},
	{
		Code:       "KOOKMIN_BANK_HN",
		Bin:        "970462",
		VietQRCode: "KBHN",
		ShortName:  "KookminHN",
		Name:       "Ngân hàng Kookmin - Chi nhánh Hà Nội",
		Logo:       "/assets/logos/KOOKMIN_BANK_HN.png",
	},
	{
		Code:       "KOOKMIN_BANK_HCM",
		Bin:        "970463",
		VietQRCode: "KBHCM",
		ShortName:  "KookminHCM",
		Name:       "Ngân hàng Kookmin - Chi nhánh Thành phố Hồ Chí Minh",
		Logo:       "/assets/logos/KOOKMIN_BANK_HCM.png",
	},
	{
		Code:       "CIMB",
		Bin:        "970459",
		VietQRCode: "",
		ShortName:  "CIMB",
		Name:       "Ngân hàng CIMB Việt Nam",
		Logo:       "/assets/logos/CIMB.png",
	},
	{
		Code:       "VIET_CREDIT",
		Bin:        "970460",
		VietQRCode: "",
		ShortName:  "VietCredit",
		Name:       "Công ty Tài chính VietCredit",
		Logo:       "/assets/logos/VIET_CREDIT.png",
	},
	{
		Code:       "TNEX",
		Bin:        "970464",
		VietQRCode: "",
		ShortName:  "TNEX",
		Name:       "Ngân hàng số TNEX",
		Logo:       "/assets/logos/TNEX.png",
	},
	{
		Code:       "MIRAE_ASSET",
		Bin:        "970468",
		VietQRCode: "",
		ShortName:  "Mirae Asset",
		Name:       "Công ty Tài chính Mirae Asset",
		Logo:       "",
	},
	{
		Code:       "MB_SHINSEI",
		Bin:        "970470",
		VietQRCode: "",
		ShortName:  "MB Shinsei",
		Name:       "Ngân hàng MB Shinsei",
		Logo:       "/assets/logos/MB_SHINSEI.png",
	},
}

var bankIndex = func() map[string]*Bank {
	index := make(map[string]*Bank, len(banks)*4)
	add := func(alias string, bank *Bank) {
		if key := normalizeLookup(alias); key != "" {
			if _, exists := index[key]; !exists {
				index[key] = bank
			}
		}
	}
	for i := range banks {
		bank := &banks[i]
		add(bank.Code, bank)
		add(bank.Bin, bank)
		add(bank.VietQRCode, bank)
		add(bank.ShortName, bank)
		add(bank.Name, bank)
	}
	return index
}()

func normalizeLookup(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(strings.TrimSpace(s)) {
		switch r {
		case ' ', '.', '-', '_', '/':
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

// findBank resolves a bank by library code, BIN, VietQR code, short name or full name.
func findBank(query string) (*Bank, bool) {
	bank, ok := bankIndex[normalizeLookup(query)]
	return bank, ok
}

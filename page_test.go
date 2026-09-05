package pdf

import (
	"testing"
)

func TestUcs2Encoder(t *testing.T) {
	// raw content in PDF (hexed for pretty)
	// 6c5f82cf94f6884c005c284ea4661362636b3e56de5355005c29

	// raw bytes in showEncodedText (hexed for pretty)
	// 6c5f82cf94f6884c00284ea4661362636b3e56de53550029
	data := []byte{
		0x6C, 0x5F, // 江
		0x82, 0xCF, // 苏
		0x94, 0xF6, // 银
		0x88, 0x4C, // 行
		0x00, 0x28, // (
		0x4E, 0xA4, // 交
		0x66, 0x13, // 易
		0x62, 0x63, // 扣
		0x6B, 0x3E, // 款
		0x56, 0xDE, // 回
		0x53, 0x55, // 单
		0x00, 0x29, // )
	}

	var e ucs2Encoder
	text := e.Decode(string(data))

	const want = "江苏银行(交易扣款回单)"
	if text != want {
		t.Errorf("got %q, want %q", text, want)
	}
}

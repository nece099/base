package ipip

import "testing"

func TestLoad(t *testing.T) {
	var err = Load("C:/lovebizhi/tiantexin/17mon/mydata4vipday2.datx")
	if err != nil {
		t.Error(err)
	}
}

func TestFind(t *testing.T) {

}

func TestFind2(t *testing.T) {

}

func BenchmarkFind(b *testing.B) {
	for i := 0; i < 1000000; i++ {
		Find("118.28.8.8")
	}
}
func BenchmarkFind2(b *testing.B) {
	for i := 0; i < 1000000; i++ {
		Find2("118.28.8.8")
	}
}
func BenchmarkFind3(b *testing.B) {
	for i := 0; i < 1000000; i++ {
		Find3("118.28.8.8")
	}
}

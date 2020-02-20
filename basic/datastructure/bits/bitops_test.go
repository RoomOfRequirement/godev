package bits

import (
	"testing"
)

func TestSign(t *testing.T) {
	if Sign(-53) != -1 || Sign(62) != 1 {
		t.Fail()
	}

	if Sign(0) != 0 || Sign(-0) != 0 {
		t.Fatal(Sign(0), Sign(-0))
	}
}

func TestIsOppositeSign(t *testing.T) {
	if IsOppositeSign(-5, -6) || IsOppositeSign(2, 3) || !IsOppositeSign(-3, 8) || !IsOppositeSign(7, -9) {
		t.Fail()
	}

	if !IsOppositeSign(-1, 0) || IsOppositeSign(1, 0) || !IsOppositeSign(-1, -0) {
		t.Fatal(-1^0, 1^0, -1^-0, 1^-0)
	}
}

func TestAbs(t *testing.T) {
	if Abs(-5) != 5 || Abs(5) != 5 || Abs(0) != 0 {
		t.Fail()
	}
}

func TestMin(t *testing.T) {
	if Min(-10, 5) != -10 || Min(9, 6) != 6 {
		t.Fail()
	}

	// 1 << shiftToSign - 1 (max int)
	if Min(1<<shiftToSign-1, -2) == -2 {
		t.Fatal(Min(1<<shiftToSign-1, -2))
	}
}

func TestMax(t *testing.T) {
	if Max(-10, 5) != 5 || Max(9, 6) != 9 {
		t.Fail()
	}

	// 1 >> shiftToSign - 1 (min int)
	if Max(1>>shiftToSign-1, 2) == -2 {
		t.Fatal(Max(1>>shiftToSign-1, -2))
	}
}

func TestIsPowerOfTwo(t *testing.T) {
	if !IsPowerOfTwo(8) || IsPowerOfTwo(0) || IsPowerOfTwo(3) {
		t.Fail()
	}
}

func TestSetBit(t *testing.T) {
	// set 0
	if SetBit(3, 2, 0) != 1 || SetBit(7, 4, 0) != 3 {
		t.Fail()
	}

	// set 1
	if SetBit(5, 2, 1) != 7 || SetBit(7, 8, 1) != 15 {
		t.Fail()
	}
}

func TestNegateIf(t *testing.T) {
	if NegateIf(5, 1) != -5 || NegateIf(5, 0) != 5 {
		t.Fail()
	}

	if NegateIf(0, 1) != 0 || NegateIf(0, 0) != 0 {
		t.Fail()
	}
}

func TestMerge(t *testing.T) {
	if Merge(7, 0, 2) != 5 || Merge(2, 5, 5) != 7 {
		t.Fail()
	}
}

func TestCountBitsOne(t *testing.T) {
	if CountBitsOne(7) != 3 || CountBitsOne(5) != 2 || CountBitsOne(0) != 0 {
		t.Fail()
	}
}

func TestCountBitsOneUntil(t *testing.T) {
	if CountBitsOneUntil(15, 2) != 2 {
		t.Fail()
	}
}

func TestModulusDivide(t *testing.T) {
	if ModulusDivide(17, 5) != 2 || ModulusDivide(19, 2) != 1 {
		t.Fail()
	}
}

func TestCeilToPowerOfTwo(t *testing.T) {
	if CeilToPowerOfTwo(15) != 16 || CeilToPowerOfTwo(1) != 1 {
		t.Fatal(CeilToPowerOfTwo(0))
	}
}

func TestFloorToPowerOfTwo(t *testing.T) {
	if FloorToPowerOfTwo(15) != 8 || FloorToPowerOfTwo(1) != 1 {
		t.Fatal(FloorToPowerOfTwo(0))
	}
}

func TestCountBitsZeroTailing(t *testing.T) {
	if CountBitsZeroTailing(6) != 1 || CountBitsZeroTailing(8) != 3 {
		t.Fail()
	}
}

func TestIsOdd(t *testing.T) {
	if IsOdd(1) != true || IsOdd(-1) != true || IsOdd(2) != false {
		t.Fail()
	}
}

func TestIsEven(t *testing.T) {
	if IsEven(1) != false || IsEven(-1) != false || IsEven(2) != true {
		t.Fail()
	}
}

func BenchmarkIsOdd(b *testing.B) {
	f1 := func(x int) bool {
		return (x & 1) == 1
	}
	f2 := func(x int) bool {
		return x%2 != 0
	}
	f3 := func(x int) bool {
		return int(^uint(0)>>1) == (x | (int(^uint(0)>>1) - 1))
	}
	b.Run("IsOdd and ops", func(b *testing.B) {
		b.ReportAllocs()
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			_ = f1(i)
		}
	})
	b.Run("IsOdd modulo ops", func(b *testing.B) {
		b.ReportAllocs()
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			_ = f2(i)
		}
	})
	b.Run("IsOdd or ops", func(b *testing.B) {
		b.ReportAllocs()
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			_ = f3(i)
		}
	})
}

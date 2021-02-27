package util

func Lerp(t, start, end float64) float64 {
	return (1-t)*start + t*end
}

func LerpClamp(t, start, end float64) float64 {
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}
	return (1-t)*start + t*end
}

func InvLerp(v, start, end float64) float64 {
	return (v - start) / (end - start)
}

func InvLerpClamp(v, start, end float64) float64 {
	t := (v - start) / (end - start)
	if t < 0 {
		return 0
	}
	if t > 1 {
		return 1
	}
	return t
}

func Remap(v, istart, iend, ostart, oend float64) float64 {
	return Lerp(InvLerp(v, istart, iend), ostart, oend)
}

func RemapClamp(v, istart, iend, ostart, oend float64) float64 {
	return LerpClamp(InvLerpClamp(v, istart, iend), ostart, oend)
}

// Float32 versions for Path and x/image/vector

func Lerp32(t, start, end float32) float32 {
	return (1-t)*start + t*end
}

func LerpClamp32(t, start, end float32) float32 {
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}
	return (1-t)*start + t*end
}

func InvLerp32(v, start, end float32) float32 {
	return (v - start) / (end - start)
}

func InvLerpClamp32(v, start, end float32) float32 {
	t := (v - start) / (end - start)
	if t < 0 {
		return 0
	}
	if t > 1 {
		return 1
	}
	return t
}

func Remap32(v, istart, iend, ostart, oend float32) float32 {
	return Lerp32(InvLerp32(v, istart, iend), ostart, oend)
}

func RemapClamp32(v, istart, iend, ostart, oend float32) float32 {
	return LerpClamp32(InvLerpClamp32(v, istart, iend), ostart, oend)
}

package util

// Linear interpolation.

// Lerp returns the value (1-t)*start + t*end.
func Lerp(t, start, end float64) float64 {
	return (1-t)*start + t*end
}

// LerpClamp is a clamped [0,1] version of Lerp.
func LerpClamp(t, start, end float64) float64 {
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}
	return (1-t)*start + t*end
}

// InvLerp performs the inverse of Lerp and returns the value of t for a value v.
func InvLerp(v, start, end float64) float64 {
	return (v - start) / (end - start)
}

// InvLerpClamp is a clamped [start, end] version of InvLerp.
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

// Remap converts v from one space to another by applying InvLerp to find t in the initial range, and
// then using t to find v' in the new range.
func Remap(v, istart, iend, ostart, oend float64) float64 {
	return Lerp(InvLerp(v, istart, iend), ostart, oend)
}

// RemapClamp is a clamped version of Remap.
func RemapClamp(v, istart, iend, ostart, oend float64) float64 {
	return LerpClamp(InvLerpClamp(v, istart, iend), ostart, oend)
}

// Float32 versions for Path and x/image/vector

// Lerp32 is a float32 version of Lerp.
func Lerp32(t, start, end float32) float32 {
	return (1-t)*start + t*end
}

// LerpClamp32 is a float32 version of LerpClamp.
func LerpClamp32(t, start, end float32) float32 {
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}
	return (1-t)*start + t*end
}

// InvLerp32 is a float32 version of InvLerp.
func InvLerp32(v, start, end float32) float32 {
	return (v - start) / (end - start)
}

// InvLerpClamp32 is a float32 version of InvLerpClamp.
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

// Remap32 is a float32 version of Remap.
func Remap32(v, istart, iend, ostart, oend float32) float32 {
	return Lerp32(InvLerp32(v, istart, iend), ostart, oend)
}

// RemapClamp32 is a float32 version of RemapClamp.
func RemapClamp32(v, istart, iend, ostart, oend float32) float32 {
	return LerpClamp32(InvLerpClamp32(v, istart, iend), ostart, oend)
}

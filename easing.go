package particles

// EasingFn takes a value between 0 and 1 and should return a value between 0 and 1
type EasingFn func(t float32) float32

var (
	// Linear easing
	Linear EasingFn = linearEasing

	// InQuad is an increasing quadratic easing function
	InQuad EasingFn = inQuadEasing

	// OutQuad is an decreasing quadratic easing function
	OutQuad EasingFn = outQuadEasing
)

func linearEasing(t float32) float32 {
	return t
}

func inQuadEasing(t float32) float32 {
	return t * t
}

func outQuadEasing(t float32) float32 {
	return -t * (t - 2)
}

package blend

import "github.com/Nadim147c/goyou/color"

func Blend(bg, fg color.ARGB, alpha float64) color.ARGB {
	_, fgR, fgG, fgB := fg.Values()
	lr1, lg1, lb1 := color.Linearized3(fgR, fgG, fgB)

	_, bgR, bgG, bgB := bg.Values()
	lr2, lg2, lb2 := color.Linearized3(bgR, bgG, bgB)

	return color.ARGBFromLinRGB(
		lr1*alpha+lr2*(1-alpha),
		lg1*alpha+lg2*(1-alpha),
		lb1*alpha+lb2*(1-alpha),
	)
}

/*
Package image contains functions that mostly operate on image.Gray.

Logical Operations

	And
	Or
	Xor
	Not
	Sub
	Equal
	Copy

Morphological Operations

	Dilate
	Erode
	Open
	Close
	TopHat
	BotHat
	HitOrMiss
	Thin
	Skeleton
	LJSkeleton
	LJReconstitute

Remap Operations

	ColorConvert
	RemapGray
	RemapRGBA

Channel Operations

	ExtractChannel
	ReplaceChannel
	SwitchChannels
	CombineChannels

Various metrics

	Histogram
	CDF
	Variance

Image types:

	Patch - replicates a patch of colors across the plane like Uniform does for a single color
	Tile - replicates an image across the plane
*/
package image

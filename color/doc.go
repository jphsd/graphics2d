/*
Package color contains types and functions for color management.

A new color type and model for [HSL] along with functions:

  - [Complement] - return the opposite color on the color wheel
  - [Monochrome] - create a monochrome palatte for a color
  - [Analogous] - returns the analogous (adjacent) colors from the color wheel
  - [Triad] - returns the other two colors in the color wheel triad
  - [Tetrad] - returns the other three colors in the color wheel tetrad
  - [Warmer] - moves a color towards red by 10%
  - [Cooler] - moves a color towards cyan by 10%
  - [Tint] - adds 10% of white to a color
  - [Shade] - adds 10% of black to a color
  - [Boost] - increases saturation by 10%
  - [Tone] - adds 10% of gray to a color
  - [Compound] - returns the analogous colors of the color's complement

An embedded list of [CSS] color names and their colors, [CSSNamedRGBs].

An embedded list of [popular] color names and their colors, [BestNamedRGBs].

Lerping functions for RGB and HSL.

[CSS]: https://www.w3.org/wiki/CSS/Properties/color/keywords
[popular]: https://github.com/meodai/color-names
*/
package color

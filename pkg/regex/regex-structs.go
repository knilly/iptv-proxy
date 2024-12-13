// Package filter provides an option to use regex filter for your content
package regex

import "regexp"

var filter *Regex

type RegexSettings struct {
	M3uGroup		 string
	M3uChannel		 string
	CategoryLive 	 string
	CategoryVod	 	 string
	CategorySeries   string
	StreamsLive	 	 string
	StreamsVod		 string
	StreamsSeries	 string
}

type Regex struct {
	FilterOn		 bool

	M3uGroup		 *regexp.Regexp
	M3uChannel		 *regexp.Regexp
	CategoryLive 	 *regexp.Regexp
	CategoryVod	 	 *regexp.Regexp
	CategorySeries   *regexp.Regexp
	StreamsLive	 	 *regexp.Regexp
	StreamsVod		 *regexp.Regexp
	StreamsSeries	 *regexp.Regexp
}
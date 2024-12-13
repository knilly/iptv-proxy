package regex

import (
	"fmt"
	"log"
	"regexp"
	"runtime"
	"sync"

	x "github.com/tellytv/go.xtream-codes"
)

func NewFilter(cfg *RegexSettings) (*Regex, error) {
	var r Regex

	if cfg.M3uGroup != "" || cfg.M3uChannel != "" {
		r.FilterOn = true
	} else if cfg.CategoryLive != "" || cfg.CategoryVod != "" || cfg.CategorySeries != "" {
		r.FilterOn = true
	} else if cfg.StreamsLive != "" || cfg.StreamsVod != "" || cfg.StreamsSeries != "" {
		r.FilterOn = true
	} else {
		r.FilterOn = false
		return &r, nil
	}

	if r.FilterOn {
		var err error
		
		r.M3uGroup, err = regexp.Compile(cfg.M3uGroup)
		if err != nil {
			return nil, err
		}
		r.M3uChannel, err = regexp.Compile(cfg.M3uChannel)
		if err != nil {
			return nil, err
		}
		r.CategoryLive, err = regexp.Compile(cfg.CategoryLive)
		if err != nil {
			return nil, err
		}
		r.CategoryVod, err = regexp.Compile(cfg.CategoryVod)
		if err != nil {
			return nil, err
		}
		r.CategorySeries, err = regexp.Compile(cfg.CategorySeries)
		if err != nil {
			return nil, err
		}
		r.StreamsLive, err = regexp.Compile(cfg.StreamsLive)
		if err != nil {
			return nil, err
		}
		r.StreamsVod, err = regexp.Compile(cfg.StreamsVod)
		if err != nil {
			return nil, err
		}
		r.StreamsSeries, err = regexp.Compile(cfg.StreamsSeries)
		if err != nil {
			return nil, err
		}
	}

	return &r, nil
}

// DeleteUnmatched provides a method for in-place deletion of elements of a slice 
// For large data we can save disk space and it is much faster than copying slices
func deleteUnmatched[Slice ~[]E, E any](s Slice, index int) (Slice, error) {
	if index < 0 || index > len(s) - 1 {
		return nil, fmt.Errorf("regex.DeleteUnmatched: index: %d out of range", index)
	}

	s[index] = s[len(s) - 1]
	
	var zero E
	s[len(s) - 1] = zero

	s = s[:len(s) - 1]

	return s, nil
}

func (r *Regex) filterCategory(categories []x.Category, err error, filter *regexp.Regexp) ([]x.Category, error) {
    if !r.FilterOn || err != nil {
        return categories, err
    }

    numWorkers := runtime.NumCPU()
	log.Printf("Number of worker CPU: %d\n", numWorkers)

	if numWorkers == 1 {
		return r.filterCategoryInline(categories, err, filter)
	} else if numWorkers > len(categories) {
        numWorkers = len(categories)
    }

    chunkSize := (len(categories) + numWorkers - 1) / numWorkers

    var wg sync.WaitGroup
    resultChan := make(chan []x.Category, numWorkers)

    for i := 0; i < numWorkers; i++ {
        start := i * chunkSize
        end := start + chunkSize
        if end > len(categories) {
            end = len(categories)
        }

        wg.Add(1)
        go func(part []x.Category) {
            defer wg.Done()
            var filtered []x.Category
            for _, item := range part {
                if filter.MatchString(item.Name) {
                    filtered = append(filtered, item)
                }
            }
            resultChan <- filtered
        }(categories[start:end])
    }

    wg.Wait()
    close(resultChan)

    var filteredCategories []x.Category
    for partial := range resultChan {
        filteredCategories = append(filteredCategories, partial...)
    }

    return filteredCategories, nil
}

func (r *Regex) filterCategoryInline(categories []x.Category, err error, filter *regexp.Regexp) ([]x.Category, error) {
	for i := len(categories) - 1; i >= 0; i-- {
		item := categories[i]

		if ! filter.MatchString(item.Name) {
			categories, err = deleteUnmatched(categories, i)
			if err != nil {
				log.Printf("ERROR: categories -> %s\n", err)
			}
		}
	}

	return categories, err
}

func (r *Regex) filterStream(streams []x.Stream, err error, filter *regexp.Regexp) ([]x.Stream, error) {
    if !r.FilterOn || err != nil {
        return streams, err
    }

    numWorkers := runtime.NumCPU()
	log.Printf("Number of worker CPU: %d\n", numWorkers)

	if numWorkers == 1 {
		return r.filterStreamInline(streams, err, filter)
	} else if numWorkers > len(streams) {
        numWorkers = len(streams)
    }

    chunkSize := (len(streams) + numWorkers - 1) / numWorkers

    var wg sync.WaitGroup
    resultChan := make(chan []x.Stream, numWorkers)

    for i := 0; i < numWorkers; i++ {
        start := i * chunkSize
        end := start + chunkSize
        if end > len(streams) {
            end = len(streams)
        }

        wg.Add(1)
        go func(part []x.Stream) {
            defer wg.Done()
            var filtered []x.Stream
            for _, item := range part {
                if filter.MatchString(item.Name) {
                    filtered = append(filtered, item)
                }
            }
            resultChan <- filtered
        }(streams[start:end])
    }

    wg.Wait()
    close(resultChan)

    var filteredStreams []x.Stream
    for partial := range resultChan {
        filteredStreams = append(filteredStreams, partial...)
    }

    return filteredStreams, nil
}

func (r *Regex) filterStreamInline(streams []x.Stream, err error, filter *regexp.Regexp) ([]x.Stream, error) {
	for i := len(streams) - 1; i >= 0; i-- {
		item := streams[i]

		if ! filter.MatchString(item.Name) {
			streams, err = deleteUnmatched(streams, i)
			if err != nil {
				log.Printf("ERROR: streams -> %s\n", err)
			}
		}
	}

	return streams, err
}

func (r *Regex) filterSeries(series []x.SeriesInfo, err error, filter *regexp.Regexp) ([]x.SeriesInfo, error) {
    if !r.FilterOn || err != nil {
        return series, err
    }

    numWorkers := runtime.NumCPU()
	log.Printf("Number of worker CPU: %d\n", numWorkers)

	if numWorkers == 1 {
		return r.filterSeriesInline(series, err, filter)
	} else if numWorkers > len(series) {
        numWorkers = len(series)
    }

    chunkSize := (len(series) + numWorkers - 1) / numWorkers

    var wg sync.WaitGroup
    resultChan := make(chan []x.SeriesInfo, numWorkers)

    for i := 0; i < numWorkers; i++ {
        start := i * chunkSize
        end := start + chunkSize
        if end > len(series) {
            end = len(series)
        }

        wg.Add(1)
        go func(part []x.SeriesInfo) {
            defer wg.Done()
            var filtered []x.SeriesInfo
            for _, item := range part {
                if filter.MatchString(item.Name) {
                    filtered = append(filtered, item)
                }
            }
            resultChan <- filtered
        }(series[start:end])
    }

    wg.Wait()
    close(resultChan)

    var filteredSeries []x.SeriesInfo
    for partial := range resultChan {
        filteredSeries = append(filteredSeries, partial...)
    }

    return filteredSeries, nil
}

func (r *Regex) filterSeriesInline(series []x.SeriesInfo, err error, filter *regexp.Regexp) ([]x.SeriesInfo, error) {
	for i := len(series) - 1; i >= 0; i-- {
		item := series[i]
		
		if ! filter.MatchString(item.Name) {
			series, err = deleteUnmatched(series, i)
			if err != nil {
				log.Printf("ERROR: series -> %s\n", err)
			}
		}
	}

	return series, err
}

func (r *Regex) FilterLiveCategory(categories []x.Category, err error) ([]x.Category, error) {
	return r.filterCategory(categories, err, r.CategoryLive)
}

func (r *Regex) FilterVodCategory(categories []x.Category, err error) ([]x.Category, error) {
	return r.filterCategory(categories, err, r.CategoryVod)
}

func (r *Regex) FilterSeriesCategory(categories []x.Category, err error) ([]x.Category, error) {
	return r.filterCategory(categories, err, r.CategorySeries)
}

func (r *Regex) FilterLiveStreams(streams []x.Stream, err error) ([]x.Stream, error) {
	return r.filterStream(streams, err, r.StreamsLive)
}

func (r *Regex) FilterVodStreams(streams []x.Stream, err error) ([]x.Stream, error) {
	return r.filterStream(streams, err, r.StreamsVod)
}

func (r *Regex) FilterSeries(series []x.SeriesInfo, err error) ([]x.SeriesInfo, error) {
	return r.filterSeries(series, err, r.StreamsSeries)
}

package main

import (
	"flag"
	"fmt"

	pr "github.com/Moleus/os-page-replacement/pkg/page-replacement"
)

var (
	framesCount  = flag.Int("frames", 5, "Number of frames in main memory")
	totalPages   = flag.Int("pages", 22, "Number of pages in virtual memory")
	replacer     = flag.String("replacer", "fifo", "Replacer algorithm: fifo, lru, opt")
	bruteForce   = flag.Bool("brute", false, "Brute force optimal frames count")
	brutePercent = flag.Float64("brute-percent", 0.05, "Brute force optimal frames count: required page faults percentage")
)

func main() {
	flag.Parse()

	pagesAccesses := []int{2, 15, 20, 17, 21, 19, 14, 3, 9, 8, 15, 10, 20, 2, 16, 18, 14, 19, 18, 7, 12, 1, 13, 20, 11, 20, 14, 17, 13, 6, 13, 15, 11, 2, 10}

	if *bruteForce {
		bruteForceOptimal(*brutePercent, pagesAccesses)
	} else {
		normalRun(pagesAccesses)
	}
}

func selectReplacer(replacer string, framesCount int) (pr.Replacer, pr.AccessNotifier) {
	switch replacer {
	case "fifo":
		return pr.NewFIFO(framesCount), &pr.NoopNotifier{}
	case "lru":
		lru := pr.NewLRU(*totalPages)
		return lru, lru
	case "opt":
		return pr.NewOPT(*totalPages), &pr.NoopNotifier{}
	default:
		panic(fmt.Sprint("Unknown replacer: ", replacer))
	}
}

func normalRun(pagesAccesses []int) {
	optimal, notifier := selectReplacer("opt", *framesCount)
	optimalWrapper := pr.NewBasicPageReplacerWrapper(optimal, *framesCount, *totalPages, pagesAccesses, notifier)
	optimalWrapper.Run(false)
	optimalFaults := optimalWrapper.GetPageFaults()

	fmt.Printf("Using '%s' page replacement algorithm\n", *replacer)
	replacer, notifier := selectReplacer(*replacer, *framesCount)
	wrapper := pr.NewBasicPageReplacerWrapper(replacer, *framesCount, *totalPages, pagesAccesses, notifier)
	wrapper.Run(true)
	faults := wrapper.GetPageFaults()

	fmt.Printf("Total page faults / optimal page faults: %d/%d\n", faults, optimalFaults)
}

func bruteForceOptimal(requiredFaultsPercentage float64, pagesAccesses []int) {
	maxFrames := len(pagesAccesses)
	for framesCount := 1; framesCount < maxFrames; framesCount++ {
		optimal, notifier := selectReplacer("opt", framesCount)
		optimalWrapper := pr.NewBasicPageReplacerWrapper(optimal, framesCount, *totalPages, pagesAccesses, notifier)
		optimalWrapper.Run(false)
		faultsCount := optimalWrapper.GetPageFaults()
		faultsPercentage := float64(faultsCount) / float64(len(pagesAccesses))
		fmt.Printf("Frames: %d. Page faults percentage: %f. (%d/%d)\n", framesCount, faultsPercentage, faultsCount, len(pagesAccesses))
		if faultsPercentage < requiredFaultsPercentage {
			return
		}
	}
	panic(fmt.Sprintf("Can't find frames count with page faults percentage more than %f", requiredFaultsPercentage))
}

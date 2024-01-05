package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	pr "github.com/Moleus/os-page-replacement/pkg/page-replacement"
)

var (
	framesCount  = flag.Int("frames", 5, "Number of frames in main memory")
	totalPages   = flag.Int("pages", 22, "Number of pages in virtual memory")
	replacer     = flag.String("replacer", "fifo", "Replacer algorithm: fifo, lru, opt")
	bruteForce   = flag.Bool("brute", false, "Brute force optimal frames count")
	brutePercent = flag.Float64("brute-percent", 0.05, "Brute force optimal frames count: required page faults percentage")
  pagesAccessesInput = flag.String("accesses", "7, 0, 1, 2, 0, 3, 0, 4, 2, 3, 0, 3, 2", "Comma separated list of pages accesses")
  isEmptyPageFault = flag.Bool("empty-is-fault", false, "Non initialized pages (-1) as fault")
)

func main() {
	flag.Parse()

  pagesAccesses := []int{}
  for _, pageAccess := range strings.Split(*pagesAccessesInput, ",") {
    numStr := strings.TrimSpace(pageAccess)
    if numStr == "" {
      continue
    }
    pageAccessInt, err := strconv.Atoi(numStr)
    if err != nil {
      panic(err)
    }
    pagesAccesses = append(pagesAccesses, pageAccessInt)
  }

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
  optimalWrapper.Run(false, *isEmptyPageFault)
	optimalFaults := optimalWrapper.GetPageFaults()

	fmt.Printf("Using '%s' page replacement algorithm\n", *replacer)
	replacer, notifier := selectReplacer(*replacer, *framesCount)
	wrapper := pr.NewBasicPageReplacerWrapper(replacer, *framesCount, *totalPages, pagesAccesses, notifier)
	wrapper.Run(true, *isEmptyPageFault)
	faults := wrapper.GetPageFaults()

  fmt.Println("Statistics:")
  fmt.Printf("faults/non-faults/total accesses: %d/%d/%d\n", faults, len(pagesAccesses) - faults, len(pagesAccesses))
	fmt.Printf("Total page faults %d. Page faults with optimal algo: %d\n", faults, optimalFaults)
}

func bruteForceOptimal(requiredFaultsPercentage float64, pagesAccesses []int) {
	maxFrames := len(pagesAccesses)
	for framesCount := 1; framesCount < maxFrames; framesCount++ {
		optimal, notifier := selectReplacer("opt", framesCount)
		optimalWrapper := pr.NewBasicPageReplacerWrapper(optimal, framesCount, *totalPages, pagesAccesses, notifier)
		optimalWrapper.Run(false, *isEmptyPageFault)
		faultsCount := optimalWrapper.GetPageFaults()
		faultsPercentage := float64(faultsCount) / float64(len(pagesAccesses))
		fmt.Printf("Frames: %d. Page faults percentage: %f. (%d/%d)\n", framesCount, faultsPercentage, faultsCount, len(pagesAccesses))
		if faultsPercentage < requiredFaultsPercentage {
			return
		}
	}
	panic(fmt.Sprintf("Can't find frames count with page faults percentage more than %f", requiredFaultsPercentage))
}

package main

import (
	"flag"
	"fmt"
	"slices"
)

var (
	framesCount = flag.Int("frames", 5, "Number of frames in main memory")
	totalPages  = flag.Int("pages", 22, "Number of pages in virtual memory")
	replacer    = flag.String("replacer", "fifo", "Replacer algorithm: fifo, lru, opt")
)

// there are three replacement algorithms: FIFO, LRU, and OPT

type BasicPageReplacerWrapper struct {
	Replacer       Replacer
	AccessNotifier AccessNotifier
	frames         []int
	framesCount    int
	totalPages     int
	pagesAccesses  []int

	pageFaults int
}

type Replacer interface {
	ChoosePageIdxToReplace(currentIndex int, pagesAccesses []int, frames []int) int
}

type AccessNotifier interface {
	Notify(page int, currentIndex int)
}

func NewBasicPageReplacerWrapper(replacer Replacer, framesCount, totalPages int, pagesAccesses []int, notifier AccessNotifier) *BasicPageReplacerWrapper {
	frames := make([]int, framesCount)
	for i := 0; i < framesCount; i++ {
		frames[i] = -1
	}
	return &BasicPageReplacerWrapper{
		Replacer:       replacer,
		frames:         frames,
		framesCount:    framesCount,
		totalPages:     totalPages,
		pagesAccesses:  pagesAccesses,
		AccessNotifier: notifier,
	}
}

func (b *BasicPageReplacerWrapper) Run(verbose bool) {
	if verbose {
		b.printHeading()
	}
	for i := 0; i < len(b.pagesAccesses); i++ {
		pageToAccess := b.pagesAccesses[i]
		b.AccessNotifier.Notify(pageToAccess, i)
		isFault := !b.isPageInFrames(pageToAccess)
		filled := !slices.Contains(b.frames, -1)
		if isFault {
			pageIndex := getFreeFrame(b.frames)
			if pageIndex == -1 {
				pageIndex = b.Replacer.ChoosePageIdxToReplace(i, b.pagesAccesses, b.frames)
			}
			b.frames[pageIndex] = pageToAccess
			b.pageFaults++
		}
		if !verbose {
			continue
		}
		b.Print(pageToAccess, isFault && filled)
	}
}

func (b *BasicPageReplacerWrapper) GetPageFaults() int {
	return b.pageFaults
}

func getFreeFrame(frames []int) int {
	for i := 0; i < len(frames); i++ {
		if frames[i] == -1 {
			return i
		}
	}
	return -1
}

func (b *BasicPageReplacerWrapper) isPageInFrames(page int) bool {
	for i := 0; i < b.framesCount; i++ {
		if b.frames[i] == page {
			return true
		}
	}
	return false
}

func (b *BasicPageReplacerWrapper) printHeading() {
	print(" P | ")
	for i := 0; i < b.framesCount; i++ {
		fmt.Printf("f%d ", i+1)
	}
	print(" | fault?")
	println()
	fmt.Println("---+-----------+---------")
}

func (b *BasicPageReplacerWrapper) Print(pageToAccess int, isFault bool) {
	fmt.Printf("%2d | ", pageToAccess)
	for i := 0; i < b.framesCount; i++ {
		fmt.Printf("%2d ", b.frames[i])
	}
	print(" | ")
	if isFault {
		print("Page fault")
	}
	println()
}

type FIFO struct {
	framesCount    int
	indexToReplace int
}

func NewFIFO(framesCount int) *FIFO {
	return &FIFO{
		framesCount:    framesCount,
		indexToReplace: 0,
	}
}

func (f *FIFO) ChoosePageIdxToReplace(int, []int, []int) int {
	lastIndex := f.indexToReplace
	f.indexToReplace++
	if f.indexToReplace == f.framesCount {
		f.indexToReplace = 0
	}
	return lastIndex
}

type LRU struct {
	lastAccessesPageToTime map[int]int
}

func NewLRU(totalPagesCount int) *LRU {
	return &LRU{
		lastAccessesPageToTime: make(map[int]int, totalPagesCount),
	}
}

func (l *LRU) Notify(page int, currentIndex int) {
	l.lastAccessesPageToTime[page] = currentIndex
}

func (l *LRU) ChoosePageIdxToReplace(currentIndex int, pagesAccesses []int, frames []int) int {
	minTime := currentIndex
	pageValue := -1
	for i := 0; i < len(pagesAccesses); i++ {
		page := pagesAccesses[i]
		if l.lastAccessesPageToTime[page] < minTime && slices.Contains(frames, page) {
			minTime = l.lastAccessesPageToTime[page]
			pageValue = page
		}
	}
	for i := 0; i < len(frames); i++ {
		if frames[i] == pageValue {
			return i
		}
	}
	panic(fmt.Sprintf("Page %d not found in frames", pageValue))
}

type NoopNotifier struct{}

func (n *NoopNotifier) Notify(int, int) {}

type OPT struct {
	totalPages int
}

func NewOPT(totalPages int) *OPT {
	return &OPT{
		totalPages: totalPages,
	}
}

func (o *OPT) ChoosePageIdxToReplace(currentIndex int, pagesAccesses []int, frames []int) int {
	maxDistance := 0
	pageValue := frames[0]
	for i := 0; i < len(pagesAccesses); i++ {
		page := pagesAccesses[i]
		if o.distanceToNextReference(page, currentIndex, pagesAccesses) > maxDistance && slices.Contains(frames, page) {
			maxDistance = o.distanceToNextReference(page, currentIndex, pagesAccesses)
			pageValue = page
		}
	}
	for i := 0; i < len(frames); i++ {
		if frames[i] == pageValue {
			return i
		}
	}
	panic(fmt.Sprintf("Page %d not found in frames", pageValue))
}

func (o *OPT) distanceToNextReference(page, currentIndex int, pagesAccesses []int) int {
	for i := currentIndex; i < len(pagesAccesses); i++ {
		if pagesAccesses[i] == page {
			return i - currentIndex
		}
	}
	return o.totalPages
}

func selectReplacer(replacer string, framesCount int) (Replacer, AccessNotifier) {
	switch replacer {
	case "fifo":
		return NewFIFO(framesCount), &NoopNotifier{}
	case "lru":
		lru := NewLRU(*totalPages)
		return lru, lru
	case "opt":
		return NewOPT(*totalPages), &NoopNotifier{}
	default:
		panic(fmt.Sprint("Unknown replacer: ", replacer))
	}
}

func main() {
	flag.Parse()

	pagesAccesses := []int{2, 15, 20, 17, 21, 19, 14, 3, 9, 8, 15, 10, 20, 2, 16, 18, 14, 19, 18, 7, 12, 1, 13, 20, 11, 20, 14, 17, 13, 6, 13, 15, 11, 2, 10}
	//pagesAccesses := []int{2, 3, 2, 1, 5, 2, 4, 5, 3, 2, 5, 2}

	optimal, notifier := selectReplacer("opt", *framesCount)
	optimalWrapper := NewBasicPageReplacerWrapper(optimal, *framesCount, *totalPages, pagesAccesses, notifier)
	optimalWrapper.Run(false)
	optimalFaults := optimalWrapper.GetPageFaults()

	fmt.Printf("Using '%s' page replacement algorithm\n", *replacer)
	replacer, notifier := selectReplacer(*replacer, *framesCount)
	wrapper := NewBasicPageReplacerWrapper(replacer, *framesCount, *totalPages, pagesAccesses, notifier)
	wrapper.Run(true)
	faults := wrapper.GetPageFaults()

	fmt.Printf("Total page faults / optimal page faults: %d/%d\n", faults, optimalFaults)
}

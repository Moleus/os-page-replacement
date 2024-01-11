package page_replacement

import (
	"fmt"
	"slices"
)

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

func (b *BasicPageReplacerWrapper) Run(verbose bool, isEmptyPageFault bool) {
	if verbose {
		b.printHeading()
	}
	for i := 0; i < len(b.pagesAccesses); i++ {
		pageToAccess := b.pagesAccesses[i]
		b.AccessNotifier.Notify(pageToAccess, i)
		isFault := !b.isPageInFrames(pageToAccess)
		isEmpty := slices.Contains(b.frames, -1)
    isShowFault := (!isEmpty || isEmptyPageFault) && isFault
		if isFault {
			pageIndex := getFreeFrame(b.frames)
			if pageIndex == -1 {
				pageIndex = b.Replacer.ChoosePageIdxToReplace(i, b.pagesAccesses, b.frames)
			}
			b.frames[pageIndex] = pageToAccess
			if isShowFault {
				b.pageFaults++
			}
		}
		if !verbose {
			continue
		}
		b.Print(pageToAccess, isShowFault)
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
	fmt.Print(" P | ")
	for i := 0; i < b.framesCount; i++ {
		fmt.Printf("f%d ", i+1)
	}
	fmt.Print(" | fault?")
	fmt.Println()
	fmt.Println("---+-----------+---------")
}

func (b *BasicPageReplacerWrapper) Print(pageToAccess int, isFault bool) {
	fmt.Printf("%2d | ", pageToAccess)
	for i := 0; i < b.framesCount; i++ {
		fmt.Printf("%2d ", b.frames[i])
	}
	fmt.Print(" | ")
	if isFault {
		fmt.Print("Page fault")
	}
	fmt.Println()
}

type FIFO struct {
	framesCount    int
	indexToReplace int
}

func NewFIFO(framesCount int) Replacer {
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
	return len(pagesAccesses)
}

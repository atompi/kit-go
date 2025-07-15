package bar

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

type Bar struct {
	mu      sync.Mutex
	graph   string    // 显示符号
	rate    string    // 进度条
	percent int       // 百分比
	current int64     // 当前进度位置
	total   int64     // 总进度
	start   time.Time // 开始时间
}

func (bar *Bar) getPercent() int {
	return int((float64(bar.current) / float64(bar.total)) * 100)
}

func (bar *Bar) getTime() (s, ls string) {
	u := time.Since(bar.start).Seconds()
	h := int(u) / 3600
	m := int(u) % 3600 / 60
	if h > 0 {
		s += strconv.Itoa(h) + "h "
	}
	if h > 0 || m > 0 {
		s += strconv.Itoa(m) + "m "
	}
	s += strconv.Itoa(int(u)%60) + "s"

	l := (float64(bar.total) / float64(bar.current)) * u
	lh := int(l) / 3600
	lm := int(l) % 3600 / 60
	if lh > 0 {
		ls += strconv.Itoa(lh) + "h "
	}
	if lh > 0 || lm > 0 {
		ls += strconv.Itoa(lm) + "m "
	}
	ls += strconv.Itoa(int(l)%60) + "s"
	return
}

func (bar *Bar) load() {
	last := bar.percent
	bar.percent = bar.getPercent()
	if bar.percent != last && bar.percent%2 == 0 {
		bar.rate += bar.graph
	}
	pasted, left := bar.getTime()
	fmt.Fprintf(os.Stderr, "\r[%-50s]% 3d%%    %2s/%2s   %d/%d", bar.rate, bar.percent, pasted, left, bar.current, bar.total)
}

func (bar *Bar) Reset(current int64) {
	bar.mu.Lock()
	defer bar.mu.Unlock()
	bar.current = current
	bar.load()
}

func (bar *Bar) Add(i int64) {
	bar.mu.Lock()
	defer bar.mu.Unlock()
	bar.current += i
	bar.load()
}

func NewBar(current, total int64) *Bar {
	bar := new(Bar)
	bar.current = current
	bar.total = total
	bar.start = time.Now()
	if bar.graph == "" {
		bar.graph = "█"
	}
	bar.percent = bar.getPercent()
	for i := 0; i < bar.percent; i += 2 {
		bar.rate += bar.graph // 初始化进度条位置
	}
	return bar
}

func NewBarWithGraph(start, total int64, graph string) *Bar {
	bar := NewBar(start, total)
	bar.graph = graph
	return bar
}

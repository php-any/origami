package fyne

import (
	"image/color"
	"sync"

	fyneLib "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

type TabItem struct {
	Title    string
	IconText string
	OnTap    func()
}

// BottomTabBar iOS 风格底部标签栏 (Fyne 自定义 Widget)
type BottomTabBar struct {
	widget.BaseWidget
	mu       sync.RWMutex
	Items    []TabItem
	Selected int

	bgColor       color.Color
	activeColor   color.Color
	inactiveColor color.Color
	dividerColor  color.Color
}

func NewBottomTabBar() *BottomTabBar {
	b := &BottomTabBar{
		Items:         make([]TabItem, 0),
		Selected:      0,
		bgColor:       color.NRGBA{R: 246, G: 246, B: 248, A: 255},
		activeColor:   color.NRGBA{R: 0, G: 122, B: 255, A: 255},
		inactiveColor: color.NRGBA{R: 142, G: 142, B: 147, A: 255},
		dividerColor:  color.NRGBA{R: 0, G: 0, B: 0, A: 30},
	}
	b.ExtendBaseWidget(b)
	return b
}

func (b *BottomTabBar) Append(title, iconText string, onTap func()) {
	b.mu.Lock()
	b.Items = append(b.Items, TabItem{Title: title, IconText: iconText, OnTap: onTap})
	b.mu.Unlock()
	b.Refresh()
}

func (b *BottomTabBar) SetSelected(idx int) {
	b.mu.Lock()
	if idx >= 0 && idx < len(b.Items) {
		b.Selected = idx
	}
	b.mu.Unlock()
	b.Refresh()
}

// Tapped 处理点击
func (b *BottomTabBar) Tapped(ev *fyneLib.PointEvent) {
	b.mu.RLock()
	n := len(b.Items)
	if n == 0 {
		b.mu.RUnlock()
		return
	}
	tabW := b.Size().Width / float32(n)
	idx := int(ev.Position.X / tabW)
	if idx < 0 {
		idx = 0
	}
	if idx >= n {
		idx = n - 1
	}
	item := b.Items[idx]
	b.mu.RUnlock()

	b.SetSelected(idx)
	if item.OnTap != nil {
		item.OnTap()
	}
}

var _ fyneLib.Tappable = (*BottomTabBar)(nil)

// ── 渲染器 ──

type bottomTabBarRenderer struct {
	bar       *BottomTabBar
	bg        *canvas.Rectangle
	divider   *canvas.Rectangle
	indicator *canvas.Rectangle
	icons     []*canvas.Text
	labels    []*canvas.Text
	lastN     int
}

func (b *BottomTabBar) CreateRenderer() fyneLib.WidgetRenderer {
	r := &bottomTabBarRenderer{bar: b}
	r.bg = canvas.NewRectangle(b.bgColor)
	r.divider = canvas.NewRectangle(b.dividerColor)
	r.indicator = canvas.NewRectangle(b.activeColor)
	r.indicator.CornerRadius = 1.5
	// 直接构建带正确选中状态的图标/标签
	r.rebuildWithState()
	return r
}

// rebuildWithState 根据当前 Items 和 Selected 状态创建图标/标签对象
func (r *bottomTabBarRenderer) rebuildWithState() {
	b := r.bar
	b.mu.RLock()
	n := len(b.Items)
	sel := b.Selected
	b.mu.RUnlock()

	if n == r.lastN && len(r.icons) == n {
		return
	}
	r.lastN = n

	r.icons = make([]*canvas.Text, n)
	r.labels = make([]*canvas.Text, n)
	for i, item := range b.Items {
		isSel := i == sel
		iconColor := b.inactiveColor
		labelColor := b.inactiveColor
		labelBold := fyneLib.TextStyle{}
		if isSel {
			iconColor = b.activeColor
			labelColor = b.activeColor
			labelBold = fyneLib.TextStyle{Bold: true}
		}

		icon := canvas.NewText(item.IconText, iconColor)
		icon.Alignment = fyneLib.TextAlignCenter
		icon.TextSize = 24
		r.icons[i] = icon

		label := canvas.NewText(item.Title, labelColor)
		label.Alignment = fyneLib.TextAlignCenter
		label.TextSize = 10
		label.TextStyle = labelBold
		r.labels[i] = label
	}
}

func (r *bottomTabBarRenderer) Layout(size fyneLib.Size) {
	r.bg.Resize(size)
	r.layoutTabs(size)
}

func (r *bottomTabBarRenderer) layoutTabs(size fyneLib.Size) {
	b := r.bar
	b.mu.RLock()
	n := len(b.Items)
	sel := b.Selected
	b.mu.RUnlock()

	if n == 0 {
		return
	}

	divH := float32(0.5)
	barH := size.Height
	bodyH := barH - divH
	tabW := size.Width / float32(n)

	// divider
	r.divider.Move(fyneLib.NewPos(0, 0))
	r.divider.Resize(fyneLib.NewSize(size.Width, divH))

	// indicator — small bar above selected tab
	indW := tabW * 0.35
	indX := float32(sel)*tabW + (tabW-indW)/2
	r.indicator.Move(fyneLib.NewPos(indX, divH+2))
	r.indicator.Resize(fyneLib.NewSize(indW, 3))

	// icon + label positions
	iconH := bodyH * 0.44
	labelH := bodyH * 0.24
	iconY := divH + bodyH*0.10
	labelY := iconY + iconH + bodyH*0.02

	for i := 0; i < n && i < len(r.icons); i++ {
		x := float32(i) * tabW
		r.icons[i].Move(fyneLib.NewPos(x, iconY))
		r.icons[i].Resize(fyneLib.NewSize(tabW, iconH))
		r.labels[i].Move(fyneLib.NewPos(x, labelY))
		r.labels[i].Resize(fyneLib.NewSize(tabW, labelH))
	}
}

func (r *bottomTabBarRenderer) MinSize() fyneLib.Size {
	return fyneLib.NewSize(320, 52)
}

func (r *bottomTabBarRenderer) Refresh() {
	b := r.bar
	b.mu.RLock()
	n := len(b.Items)
	sel := b.Selected
	b.mu.RUnlock()

	r.bg.FillColor = b.bgColor
	r.divider.FillColor = b.dividerColor
	r.indicator.FillColor = b.activeColor

	// 如果对象数量变化，需要重建
	if n != len(r.icons) {
		r.rebuildWithState()
	}

	// 原地更新已有对象的属性（Fyne 要求更新属性，不创建新对象）
	for i, item := range b.Items {
		if i >= len(r.icons) {
			break
		}
		isSel := i == sel

		iconColor := b.inactiveColor
		labelColor := b.inactiveColor
		labelBold := fyneLib.TextStyle{}
		if isSel {
			iconColor = b.activeColor
			labelColor = b.activeColor
			labelBold = fyneLib.TextStyle{Bold: true}
		}

		r.icons[i].Text = item.IconText
		r.icons[i].Color = iconColor
		r.icons[i].Refresh()

		r.labels[i].Text = item.Title
		r.labels[i].Color = labelColor
		r.labels[i].TextStyle = labelBold
		r.labels[i].Refresh()
	}

	// Update indicator position & visibility
	size := b.Size()
	if size.Width > 0 && n > 0 {
		divH := float32(0.5)
		tabW := size.Width / float32(n)
		indW := tabW * 0.35
		indX := float32(sel)*tabW + (tabW-indW)/2
		r.indicator.Move(fyneLib.NewPos(indX, divH+2))
		r.indicator.Resize(fyneLib.NewSize(indW, 3))
		r.indicator.Show()
	} else {
		r.indicator.Hide()
	}

	// 重排位置
	if size.Width > 0 {
		r.layoutTabs(size)
	}
}

func (r *bottomTabBarRenderer) Objects() []fyneLib.CanvasObject {
	objs := []fyneLib.CanvasObject{r.bg, r.divider, r.indicator}
	for _, icon := range r.icons {
		objs = append(objs, icon)
	}
	for _, label := range r.labels {
		objs = append(objs, label)
	}
	return objs
}

func (r *bottomTabBarRenderer) Destroy() {}

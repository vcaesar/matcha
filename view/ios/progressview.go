// Package ios implements a native ios views.
package ios

import (
	"image/color"

	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/internal"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/paint"
	pb "gomatcha.io/matcha/proto"
	pbios "gomatcha.io/matcha/proto/view/ios"
	"gomatcha.io/matcha/view"
)

// ProgressView implements a progess view.
type ProgressView struct {
	view.Embed
	Progress         float64
	ProgressNotifier comm.Float64Notifier
	ProgressColor    color.Color
	PaintStyle       *paint.Style
	progressNotifier comm.Float64Notifier
}

// NewProgressView returns a new view.
func NewProgressView() *ProgressView {
	return &ProgressView{}
}

// Lifecycle implements the view.View interface.
func (v *ProgressView) Lifecycle(from, to view.Stage) {
	if view.ExitsStage(from, to, view.StageMounted) {
		if v.progressNotifier != nil {
			v.Unsubscribe(v.progressNotifier)
		}
	}
}

// Build implements the view.View interface.
func (v *ProgressView) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}
	l.Solve(func(s *constraint.Solver) {
		s.Height(2) // 2.5 if its a bar
		s.WidthEqual(l.MinGuide().Width())
		s.TopEqual(l.MaxGuide().Top())
		s.LeftEqual(l.MaxGuide().Left())
	})

	if v.ProgressNotifier != v.progressNotifier {
		if v.progressNotifier != nil {
			v.Unsubscribe(v.progressNotifier)
		}
		if v.ProgressNotifier != nil {
			v.Subscribe(v.ProgressNotifier)
		}
		v.progressNotifier = v.ProgressNotifier
	}

	val := v.Progress
	if v.ProgressNotifier != nil {
		val = v.ProgressNotifier.Value()
	}

	painter := paint.Painter(nil)
	if v.PaintStyle != nil {
		painter = v.PaintStyle
	}
	return view.Model{
		Painter:        painter,
		Layouter:       l,
		NativeViewName: "gomatcha.io/matcha/view/progressview",
		NativeViewState: internal.MarshalProtobuf(&pbios.ProgressView{
			Progress:      val,
			ProgressColor: pb.ColorEncode(v.ProgressColor),
		}),
	}
}

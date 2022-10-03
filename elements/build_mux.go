package elements

import (
	"github.com/tinyzimmer/go-gst/gst"
)

func buildMux(pipeline *gst.Pipeline, name string) (*gst.Element, error) {
	if name == "mp4mux" {
		mux, err := gst.NewElement("mp4mux")
		if err != nil {
			return nil, err
		}
		pipeline.Add(mux)
		return mux, nil
	}

	mux, err := gst.NewElement("flvmux")
	if err != nil {
		return nil, err
	}
	pipeline.Add(mux)
	return mux, nil
}

func muxRequestPads(mux *gst.Element) (*gst.Pad, *gst.Pad) {
	audioPad := mux.GetRequestPad("audio_%u")
	if audioPad == nil {
		audioPad = mux.GetRequestPad("audio")
	}
	videoPad := mux.GetRequestPad("video_%u")
	if videoPad == nil {
		videoPad = mux.GetRequestPad("video")
	}

	return audioPad, videoPad

}

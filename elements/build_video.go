package elements

import (
	"github.com/tinyzimmer/go-gst/gst"
)

func buildVideoElements(pipeline *gst.Pipeline) ([]*gst.Element, error) {
	elementsForVideo, err := gst.NewElementMany("v4l2src", "queue", "videoconvert", "videorate", "videoscale", "capsfilter", "queue", "x264enc", "h264parse", "capsfilter", "queue", "tee")
	if err != nil {
		return nil, err
	}

	//Setting properties and caps
	elementsForVideo[3].Set("silent", false)
	if err := elementsForVideo[5].SetProperty("caps", gst.NewCapsFromString(
		"video/x-raw, width=1280, height=720, framerate=30/1",
	)); err != nil {
		return nil, err
	}

	if err := elementsForVideo[9].SetProperty("caps", gst.NewCapsFromString(
		"video/x-h264, profile=high",
	)); err != nil {
		return nil, err
	}

	elementsForVideo[7].Set("speed-preset", 3)
	elementsForVideo[7].Set("tune", "zerolatency")
	elementsForVideo[7].Set("bitrate", 2500)
	elementsForVideo[7].Set("key-int-max", 100)

	pipeline.AddMany(elementsForVideo...)
	//linking video elements
	gst.ElementLinkMany(elementsForVideo...)

	return elementsForVideo, nil

}

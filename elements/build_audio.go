package elements

import "github.com/tinyzimmer/go-gst/gst"

func buildAudioElements(pipeline *gst.Pipeline) ([]*gst.Element, error) {
	elementsForAudio, err := gst.NewElementMany("openalsrc", "queue", "audioconvert", "audioresample", "audiorate", "capsfilter", "queue", "fdkaacenc", "queue", "tee")
	if err != nil {
		return nil, err
	}

	//Setting properties and caps
	if err := elementsForAudio[5].SetProperty("caps", gst.NewCapsFromString(
		"audio/x-raw, rate=48000, channels=2",
	)); err != nil {
		return nil, err
	}
	elementsForAudio[7].Set("bitrate", 128000)

	pipeline.AddMany(elementsForAudio...)
	//linking audio elements
	gst.ElementLinkMany(elementsForAudio...)

	return elementsForAudio, nil
}

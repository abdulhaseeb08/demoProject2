package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/tinyzimmer/go-glib/glib"
	"github.com/tinyzimmer/go-gst/examples"
	"github.com/tinyzimmer/go-gst/gst"
	"github.com/tinyzimmer/go-gst/gst/app"
)

func buildPipeline() (*gst.Pipeline, error) {
	//initialize gstreamer
	gst.Init(nil)

	//create a new pipeline
	pipeline, err := gst.NewPipeline("")
	if err != nil {
		return nil, err
	}

	//Build the video elements and add them to the pipeline and also link
	elementsForVideo, err := buildVideoElements(pipeline)
	if err != nil {
		return nil, err
	}
	videotee := elementsForVideo[len(elementsForVideo)-1]

	//Build the audio elements and add them to the pipeline and also link
	elementsForAudio, err := buildAudioElements(pipeline)
	if err != nil {
		return nil, err
	}
	audiotee := elementsForAudio[len(elementsForAudio)-1]

	//Build both the muxes (one for file, one for streaming)
	muxFile, err := buildMux(pipeline, "mp4mux")
	if err != nil {
		return nil, err
	}

	muxStream, err := buildMux(pipeline, "flvmux")
	if err != nil {
		return nil, err
	}

	//requesting mux sink pads
	muxFileAudioSink, muxFileVideoSink := muxRequestPads(muxFile)
	muxStreamAudioSink, muxStreamVideoSink := muxRequestPads(muxStream)

	//creating queues for mux, we will link the sink pads of these queues with the audio and video tee elements
	muxQueues, err := gst.NewElementMany("queue", "queue", "queue", "queue")
	if err != nil {
		return nil, err
	}
	muxFileQueueAudio := muxQueues[0]
	muxFileQueueVideo := muxQueues[1]
	muxStreamQueueAudio := muxQueues[2]
	muxStreamQueueVideo := muxQueues[3]

	//link the queues with the FileMux
	muxFileQueueAudio.GetStaticPad("src").Link(muxFileAudioSink)
	muxFileQueueVideo.GetStaticPad("src").Link(muxFileVideoSink)

	//Link the queues with the StreamMux
	muxStreamQueueAudio.GetStaticPad("src").Link(muxStreamAudioSink)
	muxStreamQueueVideo.GetStaticPad("src").Link(muxStreamVideoSink)

	//Requesting the source pads of tee
	teesrcFileAudio := audiotee.GetRequestPad("src_%u")
	teesrcFileVideo := videotee.GetRequestPad("src_%u")
	teesrcStreamAudio := audiotee.GetRequestPad("src_%u")
	teesrcStreamVideo := videotee.GetRequestPad("src_%u")

	//Link the queue sinks with the tee element (file)
	teesrcFileAudio.Link(muxFileQueueAudio.GetStaticPad("sink"))
	teesrcFileVideo.Link(muxFileQueueVideo.GetStaticPad("sink"))

	//Link the queue sinks with the tee element (stream)
	teesrcStreamAudio.Link(muxStreamQueueAudio.GetStaticPad("sink"))
	teesrcStreamVideo.Link(muxStreamQueueVideo.GetStaticPad("sink"))

	//Creating filesink, adding it to the pipline and linking to the mux
	filesink, err := gst.NewElement("filesink")
	if err != nil {
		return nil, err
	}
	filesink.Set("location", "file.mp4")
	pipeline.Add(filesink)
	muxFile.Link(filesink)

	//Creating rtmp2sink, adding it to pipeline and linking to mux
	rtmpsink, err := gst.NewElement("rtmp2sink")
	if err != nil {
		return nil, err
	}
	rtmpsink.Set("location", "rtmps://live-api-s.facebook.com:443/rtmp/FB-178268061401609-0-Aby0txp8ZaVqImP9")
	pipeline.Add(rtmpsink)
	muxStream.Link(rtmpsink)

	//Sending EOS event
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		for sig := range ch {
			switch sig {
			case syscall.SIGINT:
				fmt.Println("Sending EOS")
				pipeline.SendEvent(gst.NewEOSEvent())
				return
			}
		}
	}()

	return pipeline, nil

}

func handleMessage(msg *gst.Message) error {
	switch msg.Type() {
	case gst.MessageEOS:
		return app.ErrEOS
	case gst.MessageError:
		gerr := msg.ParseError()
		if debug := gerr.DebugString(); debug != "" {
			fmt.Println(debug)
		}
		return gerr
	}
	return nil
}

func mainLoop(loop *glib.MainLoop, pipeline *gst.Pipeline) error {
	// Start the pipeline

	pipeline.Ref()
	defer pipeline.Unref()

	pipeline.SetState(gst.StatePlaying)

	// Retrieve the bus from the pipeline and add a watch function
	pipeline.GetPipelineBus().AddWatch(func(msg *gst.Message) bool {
		if err := handleMessage(msg); err != nil {
			fmt.Println(err)
			loop.Quit()
			return false
		}
		return true
	})

	loop.Run()

	return nil
}

func main() {
	examples.RunLoop(func(loop *glib.MainLoop) error {
		var pipeline *gst.Pipeline
		var err error
		if pipeline, err = buildPipeline(); err != nil {
			return err
		}
		return mainLoop(loop, pipeline)
	})
}

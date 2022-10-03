package main

import (
	"fmt"

	"github.com/tinyzimmer/go-gst/examples"
	"github.com/tinyzimmer/go-gst/gst"
	"github.com/ziutek/glib"
	"github.com/abdulhaseeb08/demoProject2"
)

func buildPipeline() (*gst.Pipeline, error) {
	//initialize gstreamer
	gst.Init(nil)

	//create a new pipeline
	pipeline, err := gst.NewPipeline("")
	if err != nil {
		return nil, err
	}

	elementsForVideo, err := buildVideoElements(pipeline)
	if err != nil {
		return nil, err
	}
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
